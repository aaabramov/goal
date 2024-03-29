package lib

import (
	"github.com/google/go-cmp/cmp"
	"reflect"
	"testing"
)

func TestParseCommands(t *testing.T) {

	type args struct {
		bytes []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *Goals
		wantErr bool
	}{
		{
			name: "Simple alias",
			args: args{bytes: []byte(`
workspace:
  desc: Current terraform workspace
  cmd: terraform
  args:
    - workspace
    - show
`)},
			want: &Goals{
				Commands: []Goal{
					{
						Name:   "workspace",
						Cmd:    "terraform",
						Args:   []string{"workspace", "show"},
						Assert: nil,
						Env:    "",
						Desc:   "Current terraform workspace",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Multiple environments",
			args: args{bytes: []byte(`
workspace:
  envs:
    dev:
      desc: tf apply dev
      cmd: terraform
      args:
        - apply
        - -var-file
        - dev.tfvars
    stage:
      desc: tf apply stage
      cmd: terraform
      args:
        - apply
        - -var-file
        - stage.tfvars
`)},
			want: &Goals{
				Commands: sortCommands([]Goal{
					{
						Name:   "workspace",
						Cmd:    "terraform",
						Args:   []string{"apply", "-var-file", "dev.tfvars"},
						Assert: nil,
						Env:    "dev",
						Desc:   "tf apply dev",
					},
					{
						Name:   "workspace",
						Cmd:    "terraform",
						Args:   []string{"apply", "-var-file", "stage.tfvars"},
						Assert: nil,
						Env:    "stage",
						Desc:   "tf apply stage",
					},
				}),
			},
			wantErr: false,
		},
		{
			name: "With assertion",
			args: args{bytes: []byte(`
apply:
  desc: tf apply dev
  assert:
    - desc: Check if on dev workspace
      ref: workspace
      expect: dev
      fix: terraform workspace select dev
  cmd: terraform
  args:
    - apply
    - -var-file
    - dev.tfvars
`)},
			want: &Goals{
				Commands: []Goal{
					{
						Name: "apply",
						Cmd:  "terraform",
						Args: []string{"apply", "-var-file", "dev.tfvars"},
						Assert: []Assertion{
							RefAssertion{
								Desc:   "Check if on dev workspace",
								Ref:    "workspace",
								Expect: "dev",
								Fix:    "terraform workspace select dev",
							},
						},
						Env:  "",
						Desc: "tf apply dev",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Env with assertion",
			args: args{bytes: []byte(`
apply:
  envs:
    dev:
      desc: tf apply dev
      assert:
        - desc: Check if on dev workspace
          ref: workspace
          expect: dev
          fix: terraform workspace select dev
      cmd: terraform
      args:
        - apply
        - -var-file
        - dev.tfvars
`)},
			want: &Goals{
				Commands: []Goal{
					{
						Name: "apply",
						Cmd:  "terraform",
						Args: []string{"apply", "-var-file", "dev.tfvars"},
						Assert: []Assertion{
							RefAssertion{
								Desc:   "Check if on dev workspace",
								Ref:    "workspace",
								Expect: "dev",
								Fix:    "terraform workspace select dev",
							}},
						Env:  "dev",
						Desc: "tf apply dev",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseCommands(tt.args.bytes)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseCommands() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want) {
				t.Errorf("ParseCommands() \n\tgot  = %v, \n\twant = %v\n\n%s", got, tt.want, cmp.Diff(got, tt.want))
			}
		})
	}
}

func Test_normalizeArgs(t *testing.T) {
	type args struct {
		args []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{name: "nil", args: args{args: nil}, want: []string{}},
		{name: "empty array", args: args{args: []string{}}, want: []string{}},
		{name: "non-empty array", args: args{args: []string{"1", "2"}}, want: []string{"1", "2"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := normalizeArgs(tt.args.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("normalizeArgs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCommand_Cli(t *testing.T) {
	type fields struct {
		Name   string
		Cmd    string
		Args   []string
		Assert []Assertion
		Env    string
		Desc   string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{name: "no argument", fields: fields{Cmd: "echo", Args: []string{}}, want: "echo"},
		{name: "single argument", fields: fields{Cmd: "echo", Args: []string{"123"}}, want: "echo 123"},
		{name: "single argument + flag", fields: fields{Cmd: "echo", Args: []string{"-n", "123"}}, want: "echo -n 123"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Goal{
				Name:   tt.fields.Name,
				Cmd:    tt.fields.Cmd,
				Args:   tt.fields.Args,
				Assert: tt.fields.Assert,
				Env:    tt.fields.Env,
				Desc:   tt.fields.Desc,
			}
			if got := c.Cli(); got != tt.want {
				t.Errorf("Cli() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCommands_GetWithEnv(t *testing.T) {
	type fields struct {
		Commands []Goal
	}
	type args struct {
		name string
		env  string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Goal
		exists bool
	}{
		{
			name: "single no env, empty env",
			fields: fields{Commands: []Goal{
				{Name: "test", Env: ""},
			}},
			args: args{
				name: "test",
				env:  "",
			},
			want:   &Goal{Name: "test", Env: ""},
			exists: true,
		},
		{
			name: "single no env, empty env, wrong command",
			fields: fields{Commands: []Goal{
				{Name: "test", Env: ""},
			}},
			args: args{
				name: "wrong",
				env:  "",
			},
			want:   nil,
			exists: false,
		},
		{
			name: "single no env, dev env, wrong command",
			fields: fields{Commands: []Goal{
				{Name: "test", Env: ""},
			}},
			args: args{
				name: "wrong",
				env:  "dev",
			},
			want:   nil,
			exists: false,
		},
		{
			name: "single with env, dev env, wrong command",
			fields: fields{Commands: []Goal{
				{Name: "test", Env: "dev"},
			}},
			args: args{
				name: "wrong",
				env:  "dev",
			},
			want:   nil,
			exists: false,
		},
		{
			name: "single with env, empty env",
			fields: fields{Commands: []Goal{
				{Name: "test", Env: "dev"},
			}},
			args: args{
				name: "test",
				env:  "",
			},
			want:   nil,
			exists: false,
		},
		{
			name: "single with env, dev env",
			fields: fields{Commands: []Goal{
				{Name: "test", Env: "dev"},
			}},
			args: args{
				name: "test",
				env:  "dev",
			},
			want:   &Goal{Name: "test", Env: "dev"},
			exists: true,
		},
		{
			name: "multiple with env, no env",
			fields: fields{Commands: []Goal{
				{Name: "test", Env: "dev"},
				{Name: "test", Env: "stage"},
			}},
			args: args{
				name: "test",
				env:  "",
			},
			want:   nil,
			exists: false,
		},
		{
			name: "multiple with env, dev env",
			fields: fields{Commands: []Goal{
				{Name: "test", Env: "dev"},
				{Name: "test", Env: "stage"},
			}},
			args: args{
				name: "test",
				env:  "dev",
			},
			want:   &Goal{Name: "test", Env: "dev"},
			exists: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Goals{
				Commands: tt.fields.Commands,
			}
			got, got1 := c.GetWithEnv(tt.args.name, tt.args.env)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetWithEnv() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.exists {
				t.Errorf("GetWithEnv() got1 = %v, want %v", got1, tt.exists)
			}
		})
	}
}

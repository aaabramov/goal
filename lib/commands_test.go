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
		want    *Commands
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
			want: &Commands{
				Commands: []Command{
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
			want: &Commands{
				Commands: []Command{
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
				},
			},
			wantErr: false,
		},
		{
			name: "With assertion",
			args: args{bytes: []byte(`
apply:
  desc: tf apply dev
  assert:
    desc: Check if on dev workspace
    ref: workspace
    expect: dev
    fix: terraform workspace select dev
  cmd: terraform
  args:
    - apply
    - -var-file
    - dev.tfvars
`)},
			want: &Commands{
				Commands: []Command{
					{
						Name: "apply",
						Cmd:  "terraform",
						Args: []string{"apply", "-var-file", "dev.tfvars"},
						Assert: &Assert{
							Desc:   "Check if on dev workspace",
							Ref:    "workspace",
							Expect: "dev",
							Fix:    "terraform workspace select dev",
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
        desc: Check if on dev workspace
        ref: workspace
        expect: dev
        fix: terraform workspace select dev
      cmd: terraform
      args:
        - apply
        - -var-file
        - dev.tfvars
`)},
			want: &Commands{
				Commands: []Command{
					{
						Name: "apply",
						Cmd:  "terraform",
						Args: []string{"apply", "-var-file", "dev.tfvars"},
						Assert: &Assert{
							Desc:   "Check if on dev workspace",
							Ref:    "workspace",
							Expect: "dev",
							Fix:    "terraform workspace select dev",
						},
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
		Assert *Assert
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
			c := Command{
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

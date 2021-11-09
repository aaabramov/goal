package cmd

import (
	"github.com/aaabramov/goal/lib"
	"io/ioutil"
	"os"
	"testing"
)

func Test_initGoals(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "generate file", args: args{filename: "goal.test.yaml"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initGoals(tt.args.filename)
			bytes, _ := ioutil.ReadFile(tt.args.filename)
			commands, _ := lib.ParseCommands(bytes)
			if len(commands.Commands) != 10 {
				t.Errorf("expected %d commands to be generated, got: %d", 10, len(commands.Commands))
			}
			_ = os.Remove(tt.args.filename)
		})
	}
}

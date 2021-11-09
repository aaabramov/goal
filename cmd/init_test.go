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
			parseCommands, _ := lib.ParseCommands(bytes)
			if len(parseCommands.Commands) != 3 {
				t.Errorf("expected %d commands to be generated", 3)
			}
			_ = os.Remove(tt.args.filename)
		})
	}
}

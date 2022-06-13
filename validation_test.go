package aform_test

import (
	"github.com/roleupjobboard/aform"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExportBuildValidationChoicesRule(t *testing.T) {
	type args struct {
		choices []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "no choice",
			args: args{
				choices: []string{},
			},
			want: "oneof=",
		},
		{
			name: "multiple choices",
			args: args{
				choices: []string{"one", "two", "three"},
			},
			want: "oneof='one' 'two' 'three'",
		},
		{
			name: "multiple choices with spaces",
			args: args{
				choices: []string{"one monkey", "two monkeys", "three monkeys"},
			},
			want: "oneof='one monkey' 'two monkeys' 'three monkeys'",
		},
		{
			name: "same choice multiple time",
			args: args{
				choices: []string{"one", "one"},
			},
			want: "oneof='one' 'one'",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, aform.ExportBuildValidationChoicesRule(tt.args.choices...), "buildValidationChoicesRule(%v)", tt.name)
		})
	}
}

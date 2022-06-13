package aform_test

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExportValidAutoID(t *testing.T) {
	type args struct {
		autoID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "empty string to deactivate auto ID is valid",
			args: args{
				autoID: "",
			},
			wantErr: false,
		},
		{
			name: "string with one %s verb is valid",
			args: args{
				autoID: "pre_%s_post",
			},
			wantErr: false,
		},
		{
			name: "string with two %s verbs is invalid",
			args: args{
				autoID: "pre_%s_%s_post",
			},
			wantErr: true,
		},
		{
			name: "string with another verb is invalid",
			args: args{
				autoID: "pre_%v_post",
			},
			wantErr: true,
		},
		{
			name: "non empty string without %s verb is invalid",
			args: args{
				autoID: "i",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := aform.ExportValidAutoID(tt.args.autoID); (err != nil) != tt.wantErr {
				t.Errorf("validAutoID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestExportNormalizedName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "already normalized",
			args: args{
				name: "_test_test_",
			},
			want: "_test_test_",
		},
		{
			name: "double space",
			args: args{
				name: "_test  test_",
			},
			want: "_test_test_",
		},
		{
			name: "leading trailing double space",
			args: args{
				name: "  test_test  ",
			},
			want: "_test_test_",
		},
		{
			name: "with new line",
			args: args{
				name: `_test
test_`,
			},
			want: "_test_test_",
		},
		{
			name: "with blank everywhere",
			args: args{
				name: ` test 
 test  `,
			},
			want: "_test_test_",
		},
		{
			name: "with uppercase",
			args: args{
				name: "_tEsT_TEst_",
			},
			want: "_test_test_",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, aform.ExportNormalizedName(tt.args.name), "normalizedName(%v)", tt.args.name)
		})
	}
}

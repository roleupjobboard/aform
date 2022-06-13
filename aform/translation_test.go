package aform_test

import (
	"github.com/roleupjobboard/aform"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
	"testing"
)

func TestExportSelectLanguage(t *testing.T) {
	type args struct {
		availableLanguages  []language.Tag
		matchingLangStrings []string
	}
	tests := []struct {
		name string
		args args
		want language.Tag
	}{
		{
			name: "Test en",
			args: args{
				availableLanguages:  aform.ExportLanguages,
				matchingLangStrings: []string{"en"},
			},
			want: language.English,
		},
		{
			name: "Test EN",
			args: args{
				availableLanguages:  aform.ExportLanguages,
				matchingLangStrings: []string{"EN"},
			},
			want: language.English,
		},
		{
			name: "Test fr",
			args: args{
				availableLanguages:  aform.ExportLanguages,
				matchingLangStrings: []string{"fr"},
			},
			want: language.French,
		},
		{
			name: "No string",
			args: args{
				availableLanguages:  aform.ExportLanguages,
				matchingLangStrings: nil,
			},
			want: language.English,
		},
		{
			name: "Empty string",
			args: args{
				availableLanguages:  aform.ExportLanguages,
				matchingLangStrings: []string{""},
			},
			want: language.English,
		},
		{
			name: "Accept-Language: *",
			args: args{
				availableLanguages:  aform.ExportLanguages,
				matchingLangStrings: []string{"*"},
			},
			want: language.English,
		},
		{
			name: "Accept-Language: fr,en;q=0.9,ru;q=0.8",
			args: args{
				availableLanguages:  aform.ExportLanguages,
				matchingLangStrings: []string{"fr,en;q=0.9,ru;q=0.8"},
			},
			want: language.French,
		},
		{
			name: "Accept-Language: fr-CH, fr;q=0.9, en;q=0.8, de;q=0.7, *;q=0.5",
			args: args{
				availableLanguages:  aform.ExportLanguages,
				matchingLangStrings: []string{"fr-CH, fr;q=0.9, en;q=0.8, de;q=0.7, *;q=0.5"},
			},
			want: language.French,
		},
		{
			name: "en than Accept-Language: fr-CH, fr;q=0.9, en;q=0.8, de;q=0.7, *;q=0.5",
			args: args{
				availableLanguages:  aform.ExportLanguages,
				matchingLangStrings: []string{"en", "fr-CH, fr;q=0.9, en;q=0.8, de;q=0.7, *;q=0.5"},
			},
			want: language.English,
		},
		{
			name: "Accept-Language: de-CH",
			args: args{
				availableLanguages:  aform.ExportLanguages,
				matchingLangStrings: []string{"de-CH"},
			},
			want: language.English,
		},
		{
			name: "Accept-Language: de-CH than fr",
			args: args{
				availableLanguages:  aform.ExportLanguages,
				matchingLangStrings: []string{"de-CH", "fr"},
			},
			want: language.French,
		},
		{
			name: "No languages and no string",
			args: args{
				availableLanguages:  []language.Tag{},
				matchingLangStrings: []string{},
			},
			want: language.English,
		},
		{
			name: "No languages and empty string",
			args: args{
				availableLanguages:  []language.Tag{},
				matchingLangStrings: []string{""},
			},
			want: language.English,
		},
		{
			name: "No languages and fr",
			args: args{
				availableLanguages:  []language.Tag{},
				matchingLangStrings: []string{"fr"},
			},
			want: language.English,
		},
		{
			name: "Not supported language",
			args: args{
				availableLanguages:  []language.Tag{language.Danish},
				matchingLangStrings: []string{language.Danish.String()},
			},
			want: language.English,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, aform.ExportSelectLanguage(tt.args.availableLanguages, tt.args.matchingLangStrings...), "selectLanguage(%v, %v)", tt.args.availableLanguages, tt.args.matchingLangStrings)
		})
	}
}

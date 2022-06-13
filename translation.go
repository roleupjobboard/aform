package aform

import (
	"fmt"
	english "github.com/go-playground/locales/en"
	french "github.com/go-playground/locales/fr"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	fr_translations "github.com/go-playground/validator/v10/translations/fr"
	"golang.org/x/exp/slices"
	"golang.org/x/text/language"
)

// English error messages of the available validations.
const (
	BooleanErrorMessageEn   = "Enter a valid boolean"
	EmailErrorMessageEn     = "Enter a valid email address"
	ChoiceErrorMessageEn    = "Invalid choice"
	MinLengthErrorMessageEn = "Ensure this value has at least {0} characters"
	MaxLengthErrorMessageEn = "Ensure this value has at most {0} characters"
	RequiredErrorMessageEn  = "This field is required"
	URLErrorMessageEn       = "Enter a valid URL"
)

// French error messages of the available validations.
const (
	BooleanErrorMessageFr   = "Entrez un booléen valide"
	EmailErrorMessageFr     = "Entrez une adresse e-mail valide"
	ChoiceErrorMessageFr    = "Choix invalide"
	MinLengthErrorMessageFr = "Assurez-vous que cette valeur fait au minimum {0} caractères"
	MaxLengthErrorMessageFr = "Assurez-vous que cette valeur fait au maximum {0} caractères"
	RequiredErrorMessageFr  = "Ce champ est obligatoire"
	URLErrorMessageFr       = "Entrez une URL valide"
)

var (
	languages = []language.Tag{language.English, language.French}
)

var defaultLanguage = language.English

var universalTranslator *ut.UniversalTranslator

func setValidationTranslations(validate *validator.Validate) {
	en := english.New()
	fr := french.New()
	universalTranslator = ut.New(en, en, fr)
	enTrans, _ := universalTranslator.GetTranslator(language.English.String())
	setEnValidationTranslations(validate, enTrans)
	frTrans, _ := universalTranslator.GetTranslator(language.French.String())
	setFrValidationTranslations(validate, frTrans)
}

func setEnValidationTranslations(validate *validator.Validate, trans ut.Translator) {
	if err := en_translations.RegisterDefaultTranslations(validate, trans); err != nil {
		panic(fmt.Sprintf("fail to load en form error trans %s", err.Error()))
	}
	_ = validate.RegisterTranslation(BooleanErrorCode, trans, func(ut ut.Translator) error {
		return ut.Add(BooleanErrorCode, BooleanErrorMessageEn, true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T(BooleanErrorCode)
		return t
	})
	_ = validate.RegisterTranslation(EmailErrorCode, trans, func(ut ut.Translator) error {
		return ut.Add(EmailErrorCode, EmailErrorMessageEn, true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T(EmailErrorCode)
		return t
	})
	_ = validate.RegisterTranslation(ChoiceErrorCode, trans, func(ut ut.Translator) error {
		return ut.Add(ChoiceErrorCode, ChoiceErrorMessageEn, true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T(ChoiceErrorCode)
		return t
	})
	_ = validate.RegisterTranslation(MinLengthErrorCode, trans, func(ut ut.Translator) error {
		return ut.Add(MinLengthErrorCode, MinLengthErrorMessageEn, true)
	}, func(ut ut.Translator, fe validator.FieldError) string {

		t, _ := ut.T(MinLengthErrorCode, fe.Param())
		return t
	})
	_ = validate.RegisterTranslation(MaxLengthErrorCode, trans, func(ut ut.Translator) error {
		return ut.Add(MaxLengthErrorCode, MaxLengthErrorMessageEn, true)
	}, func(ut ut.Translator, fe validator.FieldError) string {

		t, _ := ut.T(MaxLengthErrorCode, fe.Param())
		return t
	})
	_ = validate.RegisterTranslation(RequiredErrorCode, trans, func(ut ut.Translator) error {
		return ut.Add(RequiredErrorCode, RequiredErrorMessageEn, true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T(RequiredErrorCode)
		return t
	})
	_ = validate.RegisterTranslation(URLErrorCode, trans, func(ut ut.Translator) error {
		return ut.Add(URLErrorCode, URLErrorMessageEn, true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T(URLErrorCode)
		return t
	})
}

func setFrValidationTranslations(validate *validator.Validate, trans ut.Translator) {
	if err := fr_translations.RegisterDefaultTranslations(validate, trans); err != nil {
		panic(fmt.Sprintf("fail to load en form error trans %s", err.Error()))
	}
	_ = validate.RegisterTranslation(BooleanErrorCode, trans, func(ut ut.Translator) error {
		return ut.Add(BooleanErrorCode, BooleanErrorMessageFr, true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T(BooleanErrorCode)
		return t
	})
	_ = validate.RegisterTranslation(EmailErrorCode, trans, func(ut ut.Translator) error {
		return ut.Add(EmailErrorCode, EmailErrorMessageFr, true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T(EmailErrorCode)
		return t
	})
	_ = validate.RegisterTranslation(ChoiceErrorCode, trans, func(ut ut.Translator) error {
		return ut.Add(ChoiceErrorCode, ChoiceErrorMessageFr, true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T(ChoiceErrorCode)
		return t
	})
	_ = validate.RegisterTranslation(MinLengthErrorCode, trans, func(ut ut.Translator) error {
		return ut.Add(MinLengthErrorCode, MinLengthErrorMessageFr, true)
	}, func(ut ut.Translator, fe validator.FieldError) string {

		t, _ := ut.T(MinLengthErrorCode, fe.Param())
		return t
	})
	_ = validate.RegisterTranslation(MaxLengthErrorCode, trans, func(ut ut.Translator) error {
		return ut.Add(MaxLengthErrorCode, MaxLengthErrorMessageFr, true)
	}, func(ut ut.Translator, fe validator.FieldError) string {

		t, _ := ut.T(MaxLengthErrorCode, fe.Param())
		return t
	})
	_ = validate.RegisterTranslation(RequiredErrorCode, trans, func(ut ut.Translator) error {
		return ut.Add(RequiredErrorCode, RequiredErrorMessageFr, true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T(RequiredErrorCode)
		return t
	})
	_ = validate.RegisterTranslation(URLErrorCode, trans, func(ut ut.Translator) error {
		return ut.Add(URLErrorCode, URLErrorMessageFr, true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T(URLErrorCode)
		return t
	})
}

func selectLanguage(availableLanguages []language.Tag, matchingLangStrings ...string) language.Tag {
	matcher := language.NewMatcher(availableLanguages)
	l, _ := language.MatchStrings(matcher, matchingLangStrings...)
	if l == language.Und {
		return defaultLanguage
	}
	if !slices.Contains(languageStrings(languages), l.String()) {
		return defaultLanguage
	}
	return l
}

func languageStrings(langs []language.Tag) []string {
	s := make([]string, len(langs))
	for i, l := range langs {
		s[i] = l.String()
	}
	return s
}

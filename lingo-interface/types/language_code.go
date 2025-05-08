package types

import (
	"github.com/Mad-Pixels/applingo-api/pkg/utils"
	"github.com/pkg/errors"
)

// Language represents a language with a code and a name.
type Language struct {
	// Code is the ISO 639-1 language code.
	Code string
	// Name is the English name of the language.
	Name string
}

// NewLanguage creates a new Language instance.
func NewLanguage(code, name string) Language {
	return Language{
		Code: code,
		Name: name,
	}
}

// NewLanguageFromCode creates a new Language instance from a language code.
func NewLanguageFromCode(code string) (Language, error) {
	langCode, err := ParseLanguageCode(code)
	if err != nil {
		return Language{}, errors.Wrap(err, "invalid language code")
	}
	return Language{
		Code: langCode.String(),
		Name: langCode.Name(),
	}, nil
}

// NewLanguageFromName creates a new Language instance from a language name.
func NewLanguageFromName(name string) (Language, error) {
	for _, code := range AllLanguageCodes() {
		if code.Name() == name {
			return Language{
				Code: code.String(),
				Name: name,
			}, nil
		}
	}
	return Language{}, errors.New("unknown language name")
}

// ParseLanguageString parses a language string and returns a Language instance.
func ParseLanguageString(s string) (Language, error) {
	if len(s) == 2 {
		return NewLanguageFromCode(s)
	}
	return NewLanguageFromName(s)
}

// GetRandomLanguage returns a random language.
func GetRandomLanguage() (Language, error) {
	langCode, err := GetRandomLanguageCode()
	if err != nil {
		return Language{}, err
	}
	return Language{
		Code: langCode.String(),
		Name: langCode.Name(),
	}, nil
}

// GetRandomLanguageExcept returns a random language except the given language.
func GetRandomLanguageExcept(except Language) (Language, error) {
	allLangs := AllLanguageCodes()
	validLangs := make([]LanguageCode, 0, len(allLangs)-1)

	for _, code := range allLangs {
		if code.String() != except.Code {
			validLangs = append(validLangs, code)
		}
	}
	if len(validLangs) == 0 {
		return Language{}, errors.New("no valid languages available")
	}

	idx, err := utils.RandomInt(0, len(validLangs)-1)
	if err != nil {
		return Language{}, err
	}
	randomLang := validLangs[idx]
	return Language{
		Code: randomLang.String(),
		Name: randomLang.Name(),
	}, nil
}

// LanguageCode represents ISO 639-1 language code
type LanguageCode int

const (
	// LangEN represents English ("en").
	LangEN LanguageCode = iota

	// LangES represents Spanish ("es").
	LangES

	// LangFR represents French ("fr").
	// LangFR LanguageCode
	// LangNL represents Dutch ("nl").
	// LangNL LanguageCode
	// LangPL represents Polish ("pl").
	// LangPL LanguageCode
	// LangCS represents Czech ("cs").
	// LangCS LanguageCode
	// LangSV represents Swedish ("sv").
	// LangSV LanguageCode
	// LangDA represents Danish ("da").
	// LangDA LanguageCode
	// LangFI represents Finnish ("fi").
	// LangFI LanguageCode
	// LangNO represents Norwegian ("no").
	// LangNO LanguageCode
	// LangHI represents Hindi ("hi").
	// LangHI LanguageCode
	// LangBN represents Bengali ("bn").
	// LangBN LanguageCode
	// LangID represents Indonesian ("id").
	// LangID LanguageCode
	// LangAR represents Arabic ("ar").
	// LangAR LanguageCode
	// LangFA represents Persian ("fa").
	// LangFA LanguageCode

	// LangDE represents German ("de").
	LangDE

	// LangIT represents Italian ("it").
	LangIT

	// LangPT represents Portuguese ("pt").
	LangPT

	// LangRU represents Russian ("ru").
	LangRU

	// LangHE represents Hebrew ("he").
	LangHE
)

// String returns the ISO 639-1 code
func (l LanguageCode) String() string {
	switch l {
	case LangEN:
		return "en"
	case LangES:
		return "es"
	// case LangFR:
	// 	return "fr"
	case LangDE:
		return "de"
	case LangIT:
		return "it"
	case LangPT:
		return "pt"
	case LangRU:
		return "ru"
	// case LangNL:
	// 	return "nl"
	// case LangPL:
	// 	return "pl"
	// case LangCS:
	// 	return "cs"
	// case LangSV:
	// 	return "sv"
	// case LangDA:
	// 	return "da"
	// case LangFI:
	// 	return "fi"
	// case LangNO:
	// 	return "no"
	// case LangHI:
	// 	return "hi"
	// case LangBN:
	// 	return "bn"
	// case LangID:
	// 	return "id"
	// case LangAR:
	// 	return "ar"
	case LangHE:
		return "he"
	// case LangFA:
	// 	return "fa"
	default:
		return "unknown"
	}
}

// Name returns the English name of the language
func (l LanguageCode) Name() string {
	switch l {
	case LangEN:
		return "English"
	case LangES:
		return "Spanish"
	// case LangFR:
	// 	return "French"
	case LangDE:
		return "German"
	case LangIT:
		return "Italian"
	case LangPT:
		return "Portuguese"
	case LangRU:
		return "Russian"
	// case LangNL:
	// 	return "Dutch"
	// case LangPL:
	// 	return "Polish"
	// case LangCS:
	// 	return "Czech"
	// case LangSV:
	// 	return "Swedish"
	// case LangDA:
	// 	return "Danish"
	// case LangFI:
	// 	return "Finnish"
	// case LangNO:
	// 	return "Norwegian"
	// case LangHI:
	// 	return "Hindi"
	// case LangBN:
	// 	return "Bengali"
	// case LangID:
	// 	return "Indonesian"
	// case LangAR:
	// 	return "Arabic"
	case LangHE:
		return "Hebrew"
	// case LangFA:
	// 	return "Persian"
	default:
		return "Unknown"
	}
}

// AllLanguageCodes returns a slice of all available language codes
func AllLanguageCodes() []LanguageCode {
	return []LanguageCode{
		LangEN, LangES, LangDE, LangIT, LangPT, LangRU, LangHE,
		// LangFR, LangNL, LangPL, LangCS, LangSV, LangDA, LangFI,
		// LangHI, LangBN, LangID, LangAR, LangFA, LangNO,
	}
}

// GetRandomLanguageCode returns a random language code
func GetRandomLanguageCode() (LanguageCode, error) {
	codes := AllLanguageCodes()
	idx, err := utils.RandomInt(0, len(codes)-1)
	if err != nil {
		return 0, errors.Wrap(err, "failed to get random language code")
	}
	return codes[idx], nil
}

// ParseLanguageCode converts string to LanguageCode
func ParseLanguageCode(s string) (LanguageCode, error) {
	switch s {
	case "en":
		return LangEN, nil
	case "es":
		return LangES, nil
	// case "fr":
	// 	return LangFR, nil
	case "de":
		return LangDE, nil
	case "it":
		return LangIT, nil
	case "pt":
		return LangPT, nil
	case "ru":
		return LangRU, nil
	// case "nl":
	// 	return LangNL, nil
	// case "pl":
	// 	return LangPL, nil
	// case "cs":
	// 	return LangCS, nil
	// case "sv":
	// 	return LangSV, nil
	// case "da":
	// 	return LangDA, nil
	// case "fi":
	// 	return LangFI, nil
	// case "no":
	// 	return LangNO, nil
	// case "hi":
	// 	return LangHI, nil
	// case "bn":
	// 	return LangBN, nil
	// case "id":
	// 	return LangID, nil
	// case "ar":
	// 	return LangAR, nil
	case "he":
		return LangHE, nil
	// case "fa":
	// 	return LangFA, nil
	default:
		return 0, errors.New("invalid language code")
	}
}

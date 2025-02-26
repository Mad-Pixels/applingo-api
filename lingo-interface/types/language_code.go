package types

import (
	"github.com/Mad-Pixels/applingo-api/pkg/utils"
	"github.com/pkg/errors"
)

// LanguageCode represents ISO 639-1 language code
type LanguageCode int

const (
	LangEN LanguageCode = iota // English
	LangES                     // Spanish
	LangFR                     // French
	LangDE                     // German
	LangIT                     // Italian
	LangPT                     // Portuguese
	LangRU                     // Russian
	LangNL                     // Dutch
	LangPL                     // Polish
	LangCS                     // Czech
	LangSV                     // Swedish
	LangDA                     // Danish
	LangFI                     // Finnish
	LangNO                     // Norwegian
	LangHI                     // Hindi
	LangBN                     // Bengali
	LangID                     // Indonesian
	LangAR                     // Arabic
	LangHE                     // Hebrew
	LangFA                     // Persian
)

// String returns the ISO 639-1 code
func (l LanguageCode) String() string {
	switch l {
	case LangEN:
		return "en"
	case LangES:
		return "es"
	case LangFR:
		return "fr"
	case LangDE:
		return "de"
	case LangIT:
		return "it"
	case LangPT:
		return "pt"
	case LangRU:
		return "ru"
	case LangNL:
		return "nl"
	case LangPL:
		return "pl"
	case LangCS:
		return "cs"
	case LangSV:
		return "sv"
	case LangDA:
		return "da"
	case LangFI:
		return "fi"
	case LangNO:
		return "no"
	case LangHI:
		return "hi"
	case LangBN:
		return "bn"
	case LangID:
		return "id"
	case LangAR:
		return "ar"
	case LangHE:
		return "he"
	case LangFA:
		return "fa"
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
	case LangFR:
		return "French"
	case LangDE:
		return "German"
	case LangIT:
		return "Italian"
	case LangPT:
		return "Portuguese"
	case LangRU:
		return "Russian"
	case LangNL:
		return "Dutch"
	case LangPL:
		return "Polish"
	case LangCS:
		return "Czech"
	case LangSV:
		return "Swedish"
	case LangDA:
		return "Danish"
	case LangFI:
		return "Finnish"
	case LangNO:
		return "Norwegian"
	case LangHI:
		return "Hindi"
	case LangBN:
		return "Bengali"
	case LangID:
		return "Indonesian"
	case LangAR:
		return "Arabic"
	case LangHE:
		return "Hebrew"
	case LangFA:
		return "Persian"
	default:
		return "Unknown"
	}
}

// AllLanguageCodes returns a slice of all available language codes
func AllLanguageCodes() []LanguageCode {
	return []LanguageCode{
		LangEN, LangES, LangFR, LangDE, LangIT, LangPT, LangRU,
		LangNL, LangPL, LangCS, LangSV, LangDA, LangFI, LangNO,
		LangHI, LangBN, LangID, LangAR, LangHE, LangFA,
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
	case "fr":
		return LangFR, nil
	case "de":
		return LangDE, nil
	case "it":
		return LangIT, nil
	case "pt":
		return LangPT, nil
	case "ru":
		return LangRU, nil
	case "nl":
		return LangNL, nil
	case "pl":
		return LangPL, nil
	case "cs":
		return LangCS, nil
	case "sv":
		return LangSV, nil
	case "da":
		return LangDA, nil
	case "fi":
		return LangFI, nil
	case "no":
		return LangNO, nil
	case "hi":
		return LangHI, nil
	case "bn":
		return LangBN, nil
	case "id":
		return LangID, nil
	case "ar":
		return LangAR, nil
	case "he":
		return LangHE, nil
	case "fa":
		return LangFA, nil
	default:
		return 0, errors.New("invalid language code")
	}
}

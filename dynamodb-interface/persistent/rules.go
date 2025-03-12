package persistent

import (
	"errors"
	"strings"

	"slices"

	"github.com/Mad-Pixels/applingo-api/lingo-interface/types"
	"github.com/Mad-Pixels/applingo-api/pkg/chatgpt"
	"github.com/Mad-Pixels/applingo-api/pkg/utils"
)

func FieldLevelOrError(v string) error {
	_, err := types.ParseLanguageLevel(v)
	return err
}

func FieldLevelOrDefault(v, baseline string) string {
	lvl, err := types.ParseLanguageLevel(v)
	if err != nil {
		return baseline
	}
	return lvl.String()
}

func FieldDescriptionOrError(v string) error {
	if v == "" {
		return errors.New("is empty")
	}
	return nil
}

func FieldDescriptionOrDefault(_, baseline string) string {
	return baseline
}

func FieldSubcategoryOrError(v string) error {
	subcategory := strings.Split(v, "-")
	if len(subcategory) != 2 {
		return errors.New("incorrect format")
	}

	if _, err := types.ParseLanguageCode(subcategory[0]); err != nil {
		return err
	}
	if _, err := types.ParseLanguageCode(subcategory[1]); err != nil {
		return err
	}
	return nil
}

func FieldSubcategoryOrDefault(v, baseline string) string {
	if err := FieldSubcategoryOrError(v); err != nil {
		return baseline
	}
	return v
}

func FieldOverviewOrError(v string) error {
	if v == "" {
		return errors.New("is empty")
	}
	return nil
}

func FieldOverviewOrDefault(v, baseline string) string {
	if err := FieldOverviewOrError(v); err != nil {
		return baseline
	}
	return v
}

func FieldTopicOrError(v string) error {
	if v == "" {
		return errors.New("is empty")
	}
	return nil
}

func FieldTopicOrDefault(v, baseline string) string {
	if err := FieldTopicOrError(v); err != nil {
		return baseline
	}
	return v
}

func FieldLanguagesOrError(v string) error {
	languages := utils.SplitValues(v)
	if len(languages) != 2 {
		return errors.New("incorrect format")
	}

	if _, err := types.ParseLanguageString(languages[0]); err != nil {
		return err
	}
	if _, err := types.ParseLanguageString(languages[1]); err != nil {
		return err
	}
	return nil
}

func FieldLanguagesOrDefault(v, baseline string) string {
	if err := FieldLanguagesOrError(v); err != nil {
		return baseline
	}
	return v
}

func FieldWordsOrError(v int) error {
	if v < 1 {
		return errors.New("incorrect value")
	}
	return nil
}

func FieldWordsOrDefault(v, baseline int) int {
	if err := FieldWordsOrError(v); err != nil {
		return baseline
	}
	return v
}

func FieldPromptOrError(v string) error {
	prompt := utils.SplitValues(v)
	if len(prompt) != 2 {
		return errors.New("incorrect format")
	}
	if prompt[0] == "" {
		return errors.New("incorrect value")
	}

	modelCheck := slices.Contains(chatgpt.AvailableModels(), chatgpt.OpenAIModel(prompt[2]))
	if !modelCheck {
		return errors.New("incorrect value")
	}
	return nil
}

func FieldPromtOrDefault(_, baseline string) string {
	return baseline
}

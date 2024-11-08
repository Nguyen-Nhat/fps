package i18n

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
)

func Test_Get_Set_LanguageWithContext(t *testing.T) {
	ctx := context.Background()

	// 1. Get language when haven't set language yet
	// 1.1. Without default language in input
	gotLangWithoutDefault1 := GetLanguageFromContext(ctx)
	assert.Equalf(t, language.English, gotLangWithoutDefault1, "1.1. GetLanguageFromContext() = %v, want %v", gotLangWithoutDefault1, language.English)
	// 1.2. With default language in input
	gotLangWithDefault1 := GetLanguageFromContext(ctx)
	assert.Equalf(t, language.English, gotLangWithoutDefault1, "1.2. GetLanguageFromContext() = %v, want %v", gotLangWithDefault1, language.Italian)

	// 2. Set language
	lang := language.French
	ctxWithLang := SetLanguageToContext(ctx, language.French)

	// 3. Get language when already set language
	// 3.1. Without default language in input
	gotLangWithoutDefault3 := GetLanguageFromContext(ctxWithLang)
	assert.Equalf(t, lang, gotLangWithoutDefault3, "3.1. GetLanguageFromContext() = %v, want %v", gotLangWithoutDefault3, lang)
	// 3.2. With default language in input
	gotLangWithDefault3 := GetLanguageFromContext(ctxWithLang)
	assert.Equalf(t, lang, gotLangWithoutDefault3, "3.1. GetLanguageFromContext() = %v, want %v", gotLangWithDefault3, lang)
}

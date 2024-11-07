package middleware

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
)

func Test_getAcceptLanguage(t *testing.T) {
	tests := []struct {
		name             string
		acceptLangHeader string
		want             language.Tag
	}{
		{"when no accept-language => return empty", "", language.Vietnamese},
		{"when accept-language is wrong => return empty", "abc", language.English},
		{"when accept-language is not support => return empty", "fr", language.English},

		{"when accept-language is en => return en", "en", language.English},
		{"when accept-language is en-US => return en", "en-US", language.English},
		{"when accept-language is en-GB => return en", "en-gb", language.English},
		{"when accept-language is vi => return vi", "vi", language.Vietnamese},

		{"when accept-language is vi, en-gb;q=0.8, en;q=0.7 => return vi", "vi, en-gb;q=0.8, en;q=0.7", language.Vietnamese},
		{"when accept-language is en-gb;q=0.8, vi;q=0.7 => return en", "en-gb;q=0.8, en;q=0.7", language.English},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getAcceptLanguage(tt.acceptLangHeader)
			assert.Equalf(t, tt.want, got, "getAcceptLanguage() = %v, want %v", got, tt.want)
		})
	}
}

func Test_Get_Set_LanguageWithContext(t *testing.T) {
	ctx := context.Background()

	// 1. Get language when haven't set language yet
	// 1.1. Without default language in input
	gotLangWithoutDefault1 := GetLanguageFromContext(ctx)
	assert.Equalf(t, language.English, gotLangWithoutDefault1, "1.1. GetLanguageFromContext() = %v, want %v", gotLangWithoutDefault1, language.English)
	// 1.2. With default language in input
	gotLangWithDefault1 := GetLanguageFromContext(ctx, language.Italian)
	assert.Equalf(t, language.English, gotLangWithoutDefault1, "1.2. GetLanguageFromContext() = %v, want %v", gotLangWithDefault1, language.Italian)

	// 2. Set language
	lang := language.French
	ctxWithLang := SetLanguageToContext(ctx, language.French)

	// 3. Get language when already set language
	// 3.1. Without default language in input
	gotLangWithoutDefault3 := GetLanguageFromContext(ctxWithLang)
	assert.Equalf(t, lang, gotLangWithoutDefault3, "3.1. GetLanguageFromContext() = %v, want %v", gotLangWithoutDefault3, lang)
	// 3.2. With default language in input
	gotLangWithDefault3 := GetLanguageFromContext(ctxWithLang, language.Italian)
	assert.Equalf(t, lang, gotLangWithoutDefault3, "3.1. GetLanguageFromContext() = %v, want %v", gotLangWithDefault3, lang)
}

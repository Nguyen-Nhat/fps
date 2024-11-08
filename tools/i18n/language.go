package i18n

import (
	"context"

	"golang.org/x/text/language"
)

type ctxKey string

var langKey ctxKey = "language"

// SetLanguageToContext ...
func SetLanguageToContext(ctx context.Context, lang language.Tag) context.Context {
	return context.WithValue(ctx, langKey, lang)
}

// GetLanguageFromContext ... get language from context
func GetLanguageFromContext(ctx context.Context) language.Tag {
	// 1. Get from context
	if lang, ok := ctx.Value(langKey).(language.Tag); ok {
		return lang
	}

	// 2. if no language in context => return default
	return language.English
}

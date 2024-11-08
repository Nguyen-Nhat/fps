package middleware

import (
	"net/http"

	"golang.org/x/text/language"

	"git.teko.vn/loyalty-system/loyalty-file-processing/tools/i18n"
)

const (
	headerAcceptLanguage = "Accept-Language"
)

var supportedLanguagesMatcher = language.NewMatcher([]language.Tag{
	language.Vietnamese, // "vi"
	language.English,    // "en"
})

// LanguageMW ...
func LanguageMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get header
		ctx := r.Context()
		acceptLanguage := r.Header.Get(headerAcceptLanguage)
		tag := getAcceptLanguage(acceptLanguage)

		// Set language to context
		ctx = i18n.SetLanguageToContext(ctx, tag)
		r = r.WithContext(ctx)

		// Next
		next.ServeHTTP(w, r)
	})
}

// getAcceptLanguage ...
func getAcceptLanguage(acceptLangHeader string) language.Tag {
	// 1. If no language in header => return Vietnamese, because previous clients are receiving Vietnamese message (backward compatibility)
	if acceptLangHeader == "" {
		return language.Vietnamese
	}

	// 2. Parse language, if FE send wrong language => return default is English
	tags, _, err := language.ParseAcceptLanguage(acceptLangHeader)
	if err != nil {
		return language.English
	}

	// 3. Match request language
	code, _, confidence := supportedLanguagesMatcher.Match(tags...)
	if confidence == language.No { // if not support => return default English
		return language.English
	}

	// 4. Convert to Base, `en-US` -> `en`
	if !code.IsRoot() {
		base, _ := code.Base()
		code = language.Make(base.String())
	}
	return code
}

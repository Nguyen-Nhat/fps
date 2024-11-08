package middleware

import (
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

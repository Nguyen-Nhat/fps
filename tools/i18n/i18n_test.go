package i18n

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
)

func TestGetMessage(t *testing.T) {
	_, _ = LoadI18n("../../resources/messages")
	type args struct {
		lang string
		key  string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"say hello by English", args{"en", "hello"}, "Hello!"},
		{"say hello by Vietnamese", args{"vi", "hello"}, "Xin chào!"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, GetMessage(tt.args.lang, tt.args.key), "GetMessage(%v, %v)", tt.args.lang, tt.args.key)
		})
	}
}

func TestGetMessageCtx(t *testing.T) {
	_, _ = LoadI18n("../../resources/messages")
	ctx := context.Background()
	type args struct {
		lang         *language.Tag
		key          string
		dataKeyValue []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"no language -> message EN", args{nil, "hello", nil}, "Hello!"},
		// test dataKeyValue
		{"EN, has dataKeyValue (miss data) -> message EN", args{&language.English, "helloSomeone", []string{"name"}}, "Hello <no value>!"},
		{"VI, has dataKeyValue (excess data) -> message EN", args{&language.Vietnamese, "helloSomeone", []string{"name", "quy", "quy"}}, "Xin chào quy!"},
		{"EN, has dataKeyValue (enough) -> message EN", args{&language.English, "helloSomeone", []string{"name", "quy"}}, "Hello quy!"},
		{"VI, has dataKeyValue (wrong key) -> message EN", args{&language.Vietnamese, "helloSomeone", []string{"namesssss", "quy"}}, "Xin chào <no value>!"},
		// test English
		{"EN, no dataKeyValue -> message EN", args{&language.English, "hello", nil}, "Hello!"},
		{"EN, has dataKeyValue -> message EN", args{&language.English, "helloSomeone", []string{"name", "quy"}}, "Hello quy!"},
		// test Vietnamese
		{"VI, no dataKeyValue -> message VI", args{&language.Vietnamese, "hello", nil}, "Xin chào!"},
		{"VI, has dataKeyValue -> message VI", args{&language.Vietnamese, "helloSomeone", []string{"name", "quy"}}, "Xin chào quy!"},
		// test French
		{"FR, no dataKeyValue -> message default EN", args{&language.French, "hello", nil}, "Hello!"},
		{"FR, has dataKeyValue -> message default EN", args{&language.French, "helloSomeone", []string{"name", "quy"}}, "Hello quy!"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctxLang := ctx
			if tt.args.lang != nil {
				ctxLang = SetLanguageToContext(ctx, *tt.args.lang)
			}
			got := GetMessageCtx(ctxLang, tt.args.key, tt.args.dataKeyValue...)
			assert.Equalf(t, tt.want, got,
				"GetMessageCtx(%v, %v, %v) = %v | but expectation is %v", tt.args.lang, tt.args.key, tt.args.dataKeyValue, got, tt.want)
		})
	}
}

func TestGetMessageD(t *testing.T) {
	_, _ = LoadI18n("../../resources/messages")
	type args struct {
		lang         language.Tag
		key          string
		dataKeyValue []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// test English
		{"EN, no dataKeyValue -> message EN", args{language.English, "hello", nil}, "Hello!"},
		{"EN, has dataKeyValue -> message EN", args{language.English, "helloSomeone", []string{"name", "quy"}}, "Hello quy!"},
		// test Vietnamese
		{"VI, no dataKeyValue -> message VI", args{language.Vietnamese, "hello", nil}, "Xin chào!"},
		{"VI, has dataKeyValue -> message VI", args{language.Vietnamese, "helloSomeone", []string{"name", "quy"}}, "Xin chào quy!"},
		// test French
		{"FR, no dataKeyValue -> message default EN", args{language.French, "hello", nil}, "Hello!"},
		{"FR, has dataKeyValue -> message default EN", args{language.French, "helloSomeone", []string{"name", "quy"}}, "Hello quy!"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetMessageD(tt.args.lang.String(), tt.args.key, tt.args.dataKeyValue...)
			assert.Equalf(t, tt.want, got,
				"GetMessageCtx(%v, %v, %v) = %v | but expectation is %v", tt.args.lang, tt.args.key, tt.args.dataKeyValue, got, tt.want)
		})
	}
}

func Test_toTemplateDataMap(t *testing.T) {
	emptyMap := map[string]string{}
	tests := []struct {
		name         string
		dataKeyValue []string
		want         map[string]string
	}{
		{"nil data -> empty", nil, emptyMap},
		{"no data -> empty", []string{}, emptyMap},
		{"1 element -> empty", []string{"a"}, emptyMap},
		{"2 element -> map has 1 ke-value pair", []string{"a", "b"}, map[string]string{"a": "b"}},
		{"3 element -> map has 1 ke-value pair", []string{"a", "b", "c"}, map[string]string{"a": "b"}},
		{"4 element -> map has 2 ke-value pair", []string{"a", "b", "c", "d"}, map[string]string{"a": "b", "c": "d"}},
		{"5 element -> map has 2 ke-value pair", []string{"a", "b", "c", "d", "e"}, map[string]string{"a": "b", "c": "d"}},
		{"6 element -> map has 3 ke-value pair", []string{"a", "b", "c", "d", "e", "f"}, map[string]string{"a": "b", "c": "d", "e": "f"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, toTemplateDataMap(tt.dataKeyValue), "toTemplateDataMap(%v)", tt.dataKeyValue)
		})
	}
}

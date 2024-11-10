package i18n

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var once sync.Once
var i18nInstance *I18n

type I18n struct {
	bundle *i18n.Bundle
}

func LoadI18n(messageFolder string) (*I18n, error) {
	var err error
	if i18nInstance == nil {
		once.Do(func() {
			bundle := i18n.NewBundle(language.English) // Default language is English

			bundle.RegisterUnmarshalFunc("json", json.Unmarshal)                                 // Register JSON unmarshal function
			msgFileEn, errEn := bundle.LoadMessageFile(fmt.Sprintf("%s/en.json", messageFolder)) // Load English messages
			msgFileVi, errVi := bundle.LoadMessageFile(fmt.Sprintf("%s/vi.json", messageFolder)) // Load Vietnamese messages

			if msgFileEn != nil && msgFileVi != nil {
				if len(msgFileEn.Messages) != len(msgFileVi.Messages) { // only warning in this case
					logger.Warnf("language files don't have the same total messages:\n"+
						"\t- en: %d messages\n"+
						"\t- vi: %d messages\n",
						len(msgFileEn.Messages), len(msgFileVi.Messages))
				}
			} else { // when has error, we will return it then stop application
				err = fmt.Errorf("failed to load language file:\n"+
					"\t- en: %v\n"+
					"\t- vi: %v\n",
					errEn, errVi)
			}

			i18nInstance = &I18n{bundle: bundle}
		})
	}

	return i18nInstance, err
}

func GetMessage(lang, key string) string {
	return getMessage(lang, key, i18n.LocalizeConfig{MessageID: key})
}

func GetMessageCtx(ctx context.Context, key string, dataKeyValue ...string) string {
	lang := GetLanguageFromContext(ctx)
	data := toTemplateDataMap(dataKeyValue)

	return getMessage(lang.String(), key, i18n.LocalizeConfig{MessageID: key, TemplateData: data})
}

// GetMessageD ... Get message with Data
func GetMessageD(lang, key string, dataKeyValue ...string) string {
	data := toTemplateDataMap(dataKeyValue)
	return getMessage(lang, key, i18n.LocalizeConfig{MessageID: key, TemplateData: data})
}

// GetMessageDT ... Get message with Data and Data is translated
func GetMessageDT(lang, key string, data map[string]string) string {
	// todo implement Data Translation
	return getMessage(lang, key, i18n.LocalizeConfig{MessageID: key, TemplateData: data})
}

// private function ----------------------------------------------------------------------------------------------------

func getMessage(lang string, key string, localizeConfig i18n.LocalizeConfig) string {
	if message, err := i18n.NewLocalizer(i18nInstance.bundle, lang).Localize(&localizeConfig); err == nil {
		return message
	} else {
		return key // if not found message key => return default is key
	}
}

// toTemplateDataMap ... convert to get TemplateData
func toTemplateDataMap(dataKeyValue []string) map[string]string {
	data := map[string]string{}
	redundant := len(dataKeyValue) % 2 // case total element is odd, need to exclude the last element
	for i := 0; i < len(dataKeyValue)-redundant; i += 2 {
		data[dataKeyValue[i]] = dataKeyValue[i+1]
	}
	return data
}

package translator

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/lantonster/askme/pkg/log"
	myTran "github.com/segmentfault/pacman/contrib/i18n"
	"github.com/segmentfault/pacman/i18n"
	"gopkg.in/yaml.v3"
)

var GlobalTrans i18n.Translator

// LangOption language option
type LangOption struct {
	Label string `json:"label"`
	Value string `json:"value"`
	// Translation completion percentage
	Progress int `json:"progress"`
}

// DefaultLangOption default language option. If user config the language is default, the language option is admin choose.
const DefaultLangOption = "Default"

var (
	// LanguageOptions language
	LanguageOptions []*LangOption
)

// NewTranslator new a translator
func init() {
	dir := "i18n"
	entries, err := os.ReadDir(dir)
	if err != nil {
		log.WithContext(context.Background()).Fatalf("read dir failed: %s", err)
	}

	// read the Bundle resources file from entries
	for _, file := range entries {
		// ignore directory
		if file.IsDir() {
			continue
		}
		// ignore non-YAML file
		if filepath.Ext(file.Name()) != ".yaml" && file.Name() != "i18n.yaml" {
			continue
		}

		buf, err := os.ReadFile(filepath.Join(dir, file.Name()))
		if err != nil {
			log.WithContext(context.Background()).Fatalf("read file failed: %s %s", file.Name(), err)
		}

		// parse the backend translation
		originalTr := struct {
			Backend map[string]map[string]any `yaml:"backend"`
			UI      map[string]any            `yaml:"ui"`
			Plugin  map[string]any            `yaml:"plugin"`
		}{}
		if err = yaml.Unmarshal(buf, &originalTr); err != nil {
			log.WithContext(context.Background()).Fatalf("unmarshal file failed: %s %s", file.Name(), err)
		}
		translation := make(map[string]any, 0)
		for k, v := range originalTr.Backend {
			translation[k] = v
		}
		translation["backend"] = originalTr.Backend
		translation["ui"] = originalTr.UI
		translation["plugin"] = originalTr.Plugin

		content, err := yaml.Marshal(translation)
		if err != nil {
			log.WithContext(context.Background()).Debugf("marshal translation content failed: %s %s", file.Name(), err)
			continue
		}

		// add translator use backend translation
		if err = myTran.AddTranslator(content, file.Name()); err != nil {
			log.WithContext(context.Background()).Debugf("add translator failed: %s %s", file.Name(), err)
			continue
		}
	}
	GlobalTrans = myTran.GlobalTrans

	i18nFile, err := os.ReadFile(filepath.Join(dir, "i18n.yaml"))
	if err != nil {
		log.WithContext(context.Background()).Fatalf("read i18n file failed: %s", err)
	}

	s := struct {
		LangOption []*LangOption `yaml:"language_options"`
	}{}
	err = yaml.Unmarshal(i18nFile, &s)
	if err != nil {
		log.WithContext(context.Background()).Fatalf("unmarshal i18n file failed: %s", err)
	}
	LanguageOptions = s.LangOption
	for _, option := range LanguageOptions {
		option.Label = fmt.Sprintf("%s (%d%%)", option.Label, option.Progress)
	}
}

// CheckLanguageIsValid check user input language is valid
func CheckLanguageIsValid(lang string) bool {
	if lang == DefaultLangOption {
		return true
	}
	for _, option := range LanguageOptions {
		if option.Value == lang {
			return true
		}
	}
	return false
}

// Tr use language to translate data. If this language translation is not available, return default english translation.
func Tr(lang i18n.Language, data string) string {
	if GlobalTrans == nil {
		return data
	}
	translation := GlobalTrans.Tr(lang, data)
	if translation == data {
		return GlobalTrans.Tr(i18n.DefaultLanguage, data)
	}
	return translation
}

// TrWithData translate key with template data, it will replace the template data {{ .PlaceHolder }} in the translation.
func TrWithData(lang i18n.Language, key string, templateData any) string {
	if GlobalTrans == nil {
		return key
	}
	translation := GlobalTrans.TrWithData(lang, key, templateData)
	if translation == key {
		return GlobalTrans.TrWithData(i18n.DefaultLanguage, key, templateData)
	}
	return translation
}

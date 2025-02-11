package i18n

import _ "embed"

type Language string

const (
	LanguageChinese            Language = "zh_CN"
	LanguageChineseTraditional Language = "zh_TW"
	LanguageEnglish            Language = "en_US"
	LanguageGerman             Language = "de_DE"
	LanguageSpanish            Language = "es_ES"
	LanguageFrench             Language = "fr_FR"
	LanguageItalian            Language = "it_IT"
	LanguageJapanese           Language = "ja_JP"
	LanguageKorean             Language = "ko_KR"
	LanguagePortuguese         Language = "pt_PT"
	LanguageRussian            Language = "ru_RU"
	LanguageVietnamese         Language = "vi_VN"

	DefaultLanguage = LanguageEnglish
)

//go:embed  en_us.yaml
var en_us []byte

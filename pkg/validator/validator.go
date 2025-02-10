package validator

import (
	"bytes"
	"context"
	"fmt"
	"reflect"
	"strings"
	"unicode"

	"github.com/go-playground/locales"
	german "github.com/go-playground/locales/de"
	english "github.com/go-playground/locales/en"
	spanish "github.com/go-playground/locales/es"
	french "github.com/go-playground/locales/fr"
	italian "github.com/go-playground/locales/it"
	japanese "github.com/go-playground/locales/ja"
	korean "github.com/go-playground/locales/ko"
	portuguese "github.com/go-playground/locales/pt"
	russian "github.com/go-playground/locales/ru"
	vietnamese "github.com/go-playground/locales/vi"
	chinese "github.com/go-playground/locales/zh"
	chineseTraditional "github.com/go-playground/locales/zh_Hant_TW"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/translations/en"
	"github.com/go-playground/validator/v10/translations/es"
	"github.com/go-playground/validator/v10/translations/fr"
	"github.com/go-playground/validator/v10/translations/it"
	"github.com/go-playground/validator/v10/translations/ja"
	"github.com/go-playground/validator/v10/translations/pt"
	"github.com/go-playground/validator/v10/translations/ru"
	"github.com/go-playground/validator/v10/translations/vi"
	"github.com/go-playground/validator/v10/translations/zh"
	"github.com/go-playground/validator/v10/translations/zh_tw"
	"github.com/lantonster/askme/pkg/errors"
	"github.com/lantonster/askme/pkg/log"
	"github.com/lantonster/askme/pkg/reason"
	"github.com/lantonster/askme/pkg/translator"
	"github.com/microcosm-cc/bluemonday"
	"github.com/segmentfault/pacman/i18n"
)

type TranslatorLocal struct {
	La           i18n.Language
	Lo           locales.Translator
	RegisterFunc func(v *validator.Validate, trans ut.Translator) (err error)
}

type Validator struct {
	Validate *validator.Validate
	Tran     ut.Translator
	Lang     i18n.Language
}

var (
	allLanguageTranslators = []*TranslatorLocal{
		{La: i18n.LanguageChinese, Lo: chinese.New(), RegisterFunc: zh.RegisterDefaultTranslations},
		{La: i18n.LanguageChineseTraditional, Lo: chineseTraditional.New(), RegisterFunc: zh_tw.RegisterDefaultTranslations},
		{La: i18n.LanguageEnglish, Lo: english.New(), RegisterFunc: en.RegisterDefaultTranslations},
		{La: i18n.LanguageGerman, Lo: german.New(), RegisterFunc: nil},
		{La: i18n.LanguageSpanish, Lo: spanish.New(), RegisterFunc: es.RegisterDefaultTranslations},
		{La: i18n.LanguageFrench, Lo: french.New(), RegisterFunc: fr.RegisterDefaultTranslations},
		{La: i18n.LanguageItalian, Lo: italian.New(), RegisterFunc: it.RegisterDefaultTranslations},
		{La: i18n.LanguageJapanese, Lo: japanese.New(), RegisterFunc: ja.RegisterDefaultTranslations},
		{La: i18n.LanguageKorean, Lo: korean.New(), RegisterFunc: nil},
		{La: i18n.LanguagePortuguese, Lo: portuguese.New(), RegisterFunc: pt.RegisterDefaultTranslations},
		{La: i18n.LanguageRussian, Lo: russian.New(), RegisterFunc: ru.RegisterDefaultTranslations},
		{La: i18n.LanguageVietnamese, Lo: vietnamese.New(), RegisterFunc: vi.RegisterDefaultTranslations},
	}
)

var GlobalValidatorMapping = make(map[i18n.Language]*Validator, 0)

func init() {
	for _, t := range allLanguageTranslators {
		tran, val := getTran(t.Lo), createDefaultValidator(t.La)
		if t.RegisterFunc != nil {
			if err := t.RegisterFunc(val, tran); err != nil {
				panic(err)
			}
		}
		GlobalValidatorMapping[t.La] = &Validator{Validate: val, Tran: tran, Lang: t.La}
	}
}

func getTran(lo locales.Translator) ut.Translator {
	tran, ok := ut.New(lo, lo).GetTranslator(lo.Locale())
	if !ok {
		panic(fmt.Sprintf("not found translator %s", lo.Locale()))
	}
	return tran
}

func createDefaultValidator(la i18n.Language) *validator.Validate {
	// 创建一个新的验证器
	validate := validator.New()

	// 注册自定义的验证函数 "notblank"
	validate.RegisterValidation("notblank", NotBlank)
	// 注册自定义的验证函数 "sanitizer"
	validate.RegisterValidation("sanitizer", Sanitizer)
	// 注册用于获取字段标签名的函数
	validate.RegisterTagNameFunc(func(fld reflect.StructField) (res string) {
		defer func() {
			if len(res) > 0 {
				res = translator.Tr(la, res)
			}
		}()
		if jsonTag := fld.Tag.Get("json"); len(jsonTag) > 0 {
			if jsonTag == "-" {
				return ""
			}
			return jsonTag
		}
		if formTag := fld.Tag.Get("form"); len(formTag) > 0 {
			return formTag
		}
		return fld.Name
	})

	return validate
}

// NotBlank 函数用于检查字段是否不为空白。
//
// 参数:
//   - fl: 字段级别信息，包含要检查的字段
//
// 返回:
//   - bool: 如果字段不为空白则为 true，否则为 false
func NotBlank(fl validator.FieldLevel) bool {
	// 获取要检查的字段
	field := fl.Field()

	// 根据字段的类型进行不同的检查
	switch field.Kind() {
	case reflect.String:
		// 去除字符串两端的空格
		trimSpace := strings.TrimSpace(field.String())
		// 如果字符串为空，将其设置为去除空格后的结果
		if len(trimSpace) == 0 {
			field.SetString(trimSpace)
		}
		return true
	case reflect.Chan, reflect.Map, reflect.Slice, reflect.Array:
		// 对于通道、映射、切片和数组，检查长度是否大于 0
		return field.Len() > 0
	case reflect.Ptr, reflect.Interface, reflect.Func:
		// 对于指针、接口和函数，检查是否不为空
		return !field.IsNil()
	default:
		// 对于其他类型，检查是否有效且不等于零值
		return field.IsValid() && field.Interface() != reflect.Zero(field.Type()).Interface()
	}
}

// Sanitizer 函数用于对字段进行清理操作。
//
// 参数:
//   - fl: 字段级别信息，包含要处理的字段
//
// 返回:
//   - bool: 始终返回 true
func Sanitizer(fl validator.FieldLevel) bool {
	// 获取要处理的字段
	field := fl.Field()

	// 根据字段的类型进行不同的操作
	switch field.Kind() {
	case reflect.String:
		// 使用 bluemonday 的 UGCPolicy 进行清理
		filter := bluemonday.UGCPolicy()
		// 清理并替换特定字符
		content := strings.Replace(filter.Sanitize(field.String()), "&amp;", "&", -1)
		// 设置清理后的结果到字段
		field.SetString(content)
		return true
	case reflect.Chan, reflect.Map, reflect.Slice, reflect.Array:
		// 对于通道、映射、切片和数组，检查长度是否大于 0
		return field.Len() > 0
	case reflect.Ptr, reflect.Interface, reflect.Func:
		// 对于指针、接口和函数，检查是否不为空
		return !field.IsNil()
	default:
		// 对于其他类型，检查是否有效且不等于零值
		return field.IsValid() && field.Interface() != reflect.Zero(field.Type()).Interface()
	}
}

type Checker interface {
	Check() (fields []*FieldError, err error)
}

type FieldError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

type ValidationErrors []*FieldError

func (v ValidationErrors) Error() string {
	buff := bytes.NewBufferString("")
	for i := 0; i < len(v); i++ {
		buff.WriteString(v[i].Error)
		buff.WriteString("\n")
	}
	return strings.TrimSpace(buff.String())
}

func GetValidatorByLang(lang i18n.Language) *Validator {
	if GlobalValidatorMapping[lang] != nil {
		return GlobalValidatorMapping[lang]
	}
	return GlobalValidatorMapping[i18n.DefaultLanguage]
}

// Check 函数用于对给定的值进行校验，并处理校验结果。
//
// 参数:
//   - c: 上下文
//   - value: 要校验的值
//
// 返回:
//   - []*FieldError: 包含校验错误信息的字段数组
//   - error: 可能出现的错误
func (v *Validator) Check(c context.Context, value any) (fields []*FieldError, err error) {
	// 对校验错误字段进行一些处理：
	//   - 如果字段名的第一个字符是字母且为拉丁字符，则将其首字母转换为大写，并在末尾添加句号。
	defer func() {
		if len(fields) == 0 {
			return
		}

		for _, field := range fields {
			if len(field.Field) == 0 {
				continue
			}

			firstRune := []rune(field.Error)[0]
			if !unicode.IsLetter(firstRune) || !unicode.Is(unicode.Latin, firstRune) {
				continue
			}

			upper := unicode.ToUpper(firstRune)
			field.Error = string(upper) + field.Error[1:]
			if !strings.HasSuffix(field.Error, ".") {
				field.Error += "."
			}
		}
	}()

	// 进行默认的结构体校验
	if err = v.Validate.Struct(value); err != nil {
		var errs validator.ValidationErrors
		if !errors.As(err, &errs) {
			// 记录校验出错时的错误信息
			log.WithContext(c).Errorf("validate.Struct(%v) error: %v", value, err)
			return nil, fmt.Errorf("validate check exception: %w", err)
		}

		for _, field := range errs {
			fieldErr := &FieldError{
				Field: field.Field(),
				Error: field.Translate(v.Tran),
			}

			structNamespace := field.StructNamespace()
			if _, filedName, found := strings.Cut(structNamespace, "."); found {
				if originalTag := getObjectTagByFieldName(c, value, filedName); len(originalTag) > 0 {
					fieldErr.Field = originalTag
				}
			}

			fields = append(fields, fieldErr)
		}

		if len(fields) > 0 {
			return fields, errors.BadRequest(reason.RequestFormatError).WithMsg("%s", ValidationErrors(fields).Error())
		}
	}

	// 如果值实现了 Checker 接口，进行额外的校验
	if value, ok := value.(Checker); ok {
		if fields, err = value.Check(); err == nil {
			return nil, nil
		}

		return fields, errors.BadRequest(reason.RequestFormatError).WithMsg("%s", ValidationErrors(fields).Error())
	}

	return nil, nil
}

// getObjectTagByFieldName 函数用于根据给定的对象和字段名获取对应的标签值。
//
// 参数:
//   - c: 上下文
//   - obj: 要检查的对象
//   - fieldName: 字段名
//
// 返回: 字段对应的标签值，如果未找到则返回空字符串
func getObjectTagByFieldName(c context.Context, obj any, fieldName string) (tag string) {
	defer func() {
		if err := recover(); err != nil {
			log.WithContext(c).Errorf("getObjectTagByFieldName(%v, %s) error: %v", obj, fieldName, err)
		}
	}()

	// 获取对象的类型
	objT := reflect.TypeOf(obj)
	// 获取指针指向的实际类型
	objT = objT.Elem()

	// 根据字段名查找字段，如果字段不存在，返回空字符串
	structField, exists := objT.FieldByName(fieldName)
	if !exists {
		return ""
	}

	// 尝试获取 "json" 标签的值，如果有且不为空则返回
	if tag = structField.Tag.Get("json"); len(tag) > 0 {
		return tag
	}
	// 否则返回 "form" 标签的值
	return structField.Tag.Get("form")
}

package i18n

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	goI18n "github.com/LinkinStars/go-i18n/v2/i18n"
	"github.com/lantonster/askme/pkg/utils"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v2"
)

var (
	bundle    = goI18n.NewBundle(language.English)
	localizes = make(map[Language]*goI18n.Localizer)
	jsonData  = make(map[Language]any)
)

func init() {
	bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)

	// 添加默认的翻译器
	addTranslator(string(DefaultLanguage), en_us)
}

type Config struct {
	BundleDir string `yaml:"bundle_dir" mapstructure:"bundle_dir"`
}

// SetTranslator 函数根据给定的配置设置翻译器。
//
// 参数:
//   - config: 配置结构体
//
// 返回: 可能返回设置过程中出现的错误
func SetTranslator(config *Config) error {
	// 检查捆绑目录是否存在
	stat, err := os.Stat(config.BundleDir)
	if err != nil {
		return fmt.Errorf("stat bundle dir: %w", err)
	} else if !stat.IsDir() {
		// 如果捆绑目录不是一个目录，返回错误
		return fmt.Errorf("bundle dir [%s] is not a directory", config.BundleDir)
	}

	// 读取捆绑目录下的所有文件
	entries, err := os.ReadDir(config.BundleDir)
	if err != nil {
		return fmt.Errorf("read bundle dir: %w", err)
	}

	for _, file := range entries {
		// 如果条目是目录则跳过
		if file.IsDir() {
			continue
		}

		// 获取文件扩展名
		ext := filepath.Ext(file.Name())
		// 如果扩展名不是 ".yaml" 则跳过
		if ext != ".yaml" {
			continue
		}

		// 读取文件内容
		buffer, err := os.ReadFile(filepath.Join(config.BundleDir, file.Name()))
		if err != nil {
			return fmt.Errorf("read bundle file [%s]: %w", file.Name(), err)
		}

		// 提取语言名称并添加翻译器
		languageName := strings.TrimSuffix(file.Name(), ext)
		addTranslator(languageName, buffer)
	}

	return nil
}

func addTranslator(lang string, data []byte) error {
	// 解析消息文件
	if _, err := bundle.ParseMessageFileBytes(data, lang); err != nil {
		return fmt.Errorf("parse language message file [%s] failed: %w", lang, err)
	}

	// 创建并设置本地化对象
	localizes[Language(lang)] = goI18n.NewLocalizer(bundle, lang)

	// 将 YAML 内容转换为 JSON 格式并保存
	jsonBytes, err := utils.YamlToJson(data)
	if err != nil {
		return fmt.Errorf("parse language message file [%s] to json failed: %w", lang, err)
	}
	jsonData[Language(lang)] = jsonBytes

	return nil
}

// Dump 函数将指定语言的 JSON 数据进行序列化。
//
// 参数:
//   - lang: 语言类型
//
// 返回:
//   - []byte: 序列化后的字节数组
//   - error: 序列化过程中可能出现的错误
func Dump(lang Language) ([]byte, error) {
	// 对指定语言的 JSON 数据进行 JSON 序列化
	return json.Marshal(jsonData[lang])
}

// Tr 函数根据指定语言和键获取翻译后的字符串，不使用模板数据。
//
// 参数:
//   - lang: 语言类型
//   - key: 翻译键
//
// 返回: 翻译后的字符串
func Tr(lang Language, key string) string {
	// 调用带有模板数据的 TrWithData 函数，传入空的模板数据
	return TrWithData(lang, key, nil)
}

// TrWithData 函数根据指定语言、键和模板数据获取翻译后的字符串。
//
// 参数:
//   - lang: 语言类型
//   - key: 翻译键
//   - templateData: 模板数据
//
// 返回: 翻译后的字符串
func TrWithData(lang Language, key string, templateData any) string {
	// 获取指定语言的本地化对象
	localize, ok := localizes[lang]
	if !ok {
		localize = localizes[DefaultLanguage]
	}

	// 进行本地化翻译
	translation, err := localize.Localize(&goI18n.LocalizeConfig{MessageID: key, TemplateData: templateData})
	if err != nil {
		// 如果翻译出错且获取消息模板也出错，返回原始键
		if _, tmpl, err := localize.GetMessageTemplate(key, nil); err != nil {
			return key
		} else {
			// 否则返回模板的其他默认值
			return tmpl.Other
		}
	}
	// 返回翻译结果
	return translation
}

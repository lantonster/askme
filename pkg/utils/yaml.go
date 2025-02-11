package utils

import (
	"encoding/json"
	"fmt"

	"sigs.k8s.io/yaml"
)

// YamlToJson 函数用于将 YAML 格式的字节数据转换为 JSON 格式的映射。
//
// 参数:
//   - buf: YAML 格式的字节数据
//
// 返回:
//   - map[string]any: 转换后的 JSON 映射
//   - error: 转换过程中可能出现的错误
func YamlToJson(buf []byte) (map[string]any, error) {
	// 将 YAML 数据转换为 JSON 字节数据
	jsonBytes, err := yaml.YAMLToJSON(buf)
	if err != nil {
		return nil, fmt.Errorf("yaml to json: %w", err)
	}

	m := make(map[string]any)
	// 将 JSON 字节数据解析到创建的映射中
	if err = json.Unmarshal(jsonBytes, &m); err != nil {
		return nil, fmt.Errorf("unmarshal json: %w", err)
	}
	return m, nil
}

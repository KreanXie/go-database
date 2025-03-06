package table

import (
	"encoding/json"
	"fmt"
)

type Row struct {
	Values any
}

// NewRow 创建一个新的行
func NewRow(values any) Row {
	return Row{
		Values: values,
	}
}

// GetValue 获取指定列的值
func (r Row) GetValue(columnName string) (interface{}, error) {
	// 如果 Values 是 map[string]interface{} 类型
	if valuesMap, ok := r.Values.(map[string]interface{}); ok {
		if value, exists := valuesMap[columnName]; exists {
			return value, nil
		}
		return nil, fmt.Errorf("column %s not found in row", columnName)
	}
	
	// 如果 Values 是其他类型，可能需要其他方式获取值
	return nil, fmt.Errorf("unsupported row values type: %T", r.Values)
}

// SetValue 设置指定列的值
func (r *Row) SetValue(columnName string, value interface{}) error {
	// 如果 Values 是 map[string]interface{} 类型
	if valuesMap, ok := r.Values.(map[string]interface{}); ok {
		valuesMap[columnName] = value
		r.Values = valuesMap
		return nil
	}
	
	// 如果 Values 是其他类型，可能需要其他方式设置值
	return fmt.Errorf("unsupported row values type: %T", r.Values)
}

// ToJSON 将行数据转换为JSON字符串
func (r Row) ToJSON() (string, error) {
	bytes, err := json.Marshal(r.Values)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// FromJSON 从JSON字符串解析行数据
func FromJSON(jsonStr string) (Row, error) {
	var values map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &values)
	if err != nil {
		return Row{}, err
	}
	return Row{Values: values}, nil
}

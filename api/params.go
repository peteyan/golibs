package api

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

// StructToMap 将结构体转换为map[string]string，其中key名称使用json标签
func StructToMap[T any](obj T) map[string]string {
	result := make(map[string]string)
	objValue := reflect.ValueOf(obj)
	objType := reflect.TypeOf(obj)
	if objValue.Kind() == reflect.Ptr {
		objValue = objValue.Elem()
		objType = objType.Elem()
	}
	for i := 0; i < objValue.NumField(); i++ {
		field := objValue.Field(i)
		typeField := objType.Field(i)
		jsonTag := typeField.Tag.Get("json")
		if jsonTag == "" {
			jsonTag = typeField.Name
		} else {
			// 处理逗号分割的标签，例如 "json:"name,omitempty"
			jsonTag = strings.Split(jsonTag, ",")[0]
		}
		if jsonTag != "-" {
			// 将字段值转换为字符串
			var fieldValue string
			switch field.Kind() {
			case reflect.String:
				fieldValue = field.String()
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				fieldValue = strconv.FormatInt(field.Int(), 10)
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				fieldValue = strconv.FormatUint(field.Uint(), 10)
			case reflect.Float32, reflect.Float64:
				fieldValue = strconv.FormatFloat(field.Float(), 'f', -1, 64)
			case reflect.Bool:
				fieldValue = strconv.FormatBool(field.Bool())
			default:
				fieldValue = fmt.Sprintf("%v", field.Interface())
			}
			result[jsonTag] = fieldValue
		}
	}
	return result
}

// MapToSortedString 将 `map[string]string` 按键排序并拼接成字符串
func MapToSortedString(m map[string]string) string {
	keys := make([]string, 0, len(m))
	for key := range m {
		if key == "apiSign" {
			continue
		}
		keys = append(keys, key)
	}
	sort.Strings(keys)
	var sb strings.Builder
	for i, key := range keys {
		sb.WriteString(fmt.Sprintf("%s=%s", key, m[key]))
		if i < len(keys)-1 {
			sb.WriteString("&")
		}
	}
	return sb.String()
}

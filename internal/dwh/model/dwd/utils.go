package dwd

import (
	"fmt"
	"go-cs/internal/dwh/pkg/model"
	"go-cs/pkg/stream"
	"reflect"
	"slices"
)

var modelFieldNames []string
var dwdWitemFieldNames []string

func init() {
	modelFieldNames = []string{getTypeName(model.DimModel{}), getTypeName(model.DwdModel{})}
	dwdWitemFieldNames = stream.Diff(structFieldNames(DwdWitem{}), modelFieldNames)
}

// 获取变量的类型名
func getTypeName(i interface{}) string {
	t := reflect.TypeOf(i)
	if t.Kind() == reflect.Ptr {
		return "*" + t.Elem().Name() // 指针类型的名称
	}
	return t.Name() // 非指针类型的名称
}

func structFieldNames(i any) []string {
	var fieldNames []string

	v := reflect.ValueOf(i)

	// 如果是结构体指针，需要先解引用
	if v.Kind() == reflect.Ptr {
		v = v.Elem() // 获取指针指向的实际值
	}

	// 确保传入的是结构体类型
	if v.Kind() == reflect.Struct {
		// 获取结构体的类型信息
		t := v.Type()

		// 遍历结构体的所有字段
		for i := 0; i < v.NumField(); i++ {
			field := t.Field(i) // 获取字段的元信息
			fieldNames = append(fieldNames, field.Name)
		}
	}

	return fieldNames
}

// / 比较两个结构体的字段值是否相等，忽略指定字段
func compareStructFields(a, b interface{}, ignoreFields []string) bool {
	valA := reflect.ValueOf(a)
	valB := reflect.ValueOf(b)

	if valA.Kind() == reflect.Pointer {
		valA = valA.Elem() // 解引用指针
	}
	if valB.Kind() == reflect.Pointer {
		valB = valB.Elem() // 解引用指针
	}

	// 确保传入的是结构体类型
	if valA.Kind() != reflect.Struct || valB.Kind() != reflect.Struct {
		fmt.Println("One of the values is not a struct")
		return false
	}

	// 获取结构体的类型
	typA := valA.Type()
	typB := valB.Type()

	// 如果结构体类型不同，则直接返回false
	if typA != typB {
		return false
	}

	// 遍历结构体的字段
	for i := 0; i < valA.NumField(); i++ {
		// 获取字段的值
		fieldA := valA.Field(i)
		fieldB := valB.Field(i)

		// 检查字段是否在忽略列表中
		if slices.Contains(ignoreFields, typA.Field(i).Name) {
			continue
		}

		// 比较字段值
		if !reflect.DeepEqual(fieldA.Interface(), fieldB.Interface()) {
			return false
		}
	}

	return true
}

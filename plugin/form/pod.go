package form

import (
	"reflect"
	"strings"
	"tini-paas/api/podapi/proto/podApi"
	"tini-paas/pkg/common"
)

// FromToPodStruct 根据结构体中name标签映射数据到结构体中并且转换类型
func FromToPodStruct(data map[string]*podApi.Pair, obj interface{}) {
	objValue := reflect.ValueOf(obj).Elem()
	for i := 0; i < objValue.NumField(); i++ {
		//获取sql对应的值
		dataTag := strings.Replace(objValue.Type().Field(i).Tag.Get("json"), ",omitempty", "", -1)
		dataSlice, ok := data[dataTag]
		if !ok {
			continue
		}
		valueSlice := dataSlice.Values
		if len(valueSlice) <= 0 {
			continue
		}
		//排除port和env
		if dataTag == "pod_port" || dataTag == "pod_env" {
			continue
		}
		value := valueSlice[0]
		//端口，环境变量的单独处理
		//获取对应字段的名称
		name := objValue.Type().Field(i).Name
		//获取对应字段类型
		structFieldType := objValue.Field(i).Type()
		//获取变量类型，也可以直接写"string类型"
		val := reflect.ValueOf(value)
		var err error
		if structFieldType != val.Type() {
			//类型转换
			val, err = TypeConversion(value, structFieldType.Name()) //类型转换
			if err != nil {
				common.Error(err)
			}
		}
		//设置类型值
		objValue.FieldByName(name).Set(val)
	}
}

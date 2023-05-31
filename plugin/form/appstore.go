package form

import (
	"reflect"
	"sort"
	"strings"
	"tini-paas/api/appstoreapi/proto/appstoreApi"
	"tini-paas/pkg/common"
)

func FormToAppStoreStruct(data map[string]*appstoreApi.Pair, obj interface{}) {
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

		//排除列表
		tagList := []string{"app_image", "app_pod", "app_middle", "app_volume", "app_comment"}
		//进行排除
		if isIn(dataTag, tagList) {
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

func isIn(target string, strArray []string) bool {
	sort.Strings(strArray)
	index := sort.SearchStrings(strArray, target)
	if index < len(strArray) && strArray[index] == target {
		return true
	}
	return false
}

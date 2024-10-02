package luatool

import (
	"bytes"
	"fmt"
	lua "github.com/yuin/gopher-lua"
	"go.mongodb.org/mongo-driver/bson/primitive"
	luar "layeh.com/gopher-luar"
	"reflect"
	"strconv"
	"strings"
	"unicode"
)

const convertCamelKeys = "ConvertCamelKeys"

func interfaceToString(str interface{}) string {
	var ret string
	if reflect.TypeOf(str).Kind() == reflect.String {
		ret = str.(string)
	} else {
		ret = fmt.Sprintf("%v", str)
	}
	return ret
}

func convertLuaRetToGo(src interface{}, camelKey map[string]interface{}) interface{} {
	if reflect.ValueOf(src).Kind() == reflect.Map {
		dst := make(map[string]interface{})
		srcVal := reflect.ValueOf(src)
		for _, key := range srcVal.MapKeys() {
			if interfaceToString(key.Interface()) == convertCamelKeys {
				continue
			} else if camelKey == nil || camelKey[strings.ToLower(interfaceToString(key.Interface()))] == nil {
				dst[camel2Case(key.Interface())] = convertLuaRetToGo(srcVal.MapIndex(key).Interface(), camelKey)
			} else {
				dst[camel(interfaceToString(key.Interface()))] = convertLuaRetToGo(srcVal.MapIndex(key).Interface(), camelKey)
			}
		}
		// 业务一般不会返回空对象，此处强制转换为空数组
		if len(dst) == 0 {
			return []interface{}{}
		}
		return dst
	} else if reflect.ValueOf(src).Kind() == reflect.Slice {
		var dst []interface{}
		srcVal := reflect.ValueOf(src)
		for i := 0; i < srcVal.Len(); i++ {
			dst = append(dst, convertLuaRetToGo(srcVal.Index(i).Interface(), camelKey))
		}
		return dst
	}
	return src
}

func ConvertLuaData(src interface{}) interface{} {
	var camelKey map[string]interface{}
	if reflect.ValueOf(src).Kind() == reflect.Map {
		srcVal := reflect.ValueOf(src)
		for _, key := range srcVal.MapKeys() {
			if strings.ToLower(interfaceToString(key.Interface())) == strings.ToLower(convertCamelKeys) {
				camelKeySlice := srcVal.MapIndex(key)
				camelKeySliceVal := reflect.ValueOf(camelKeySlice.Interface())
				if camelKeySliceVal.Kind() == reflect.Slice {
					camelKey = make(map[string]interface{})
					for i := 0; i < camelKeySliceVal.Len(); i++ {
						camelKey[strings.ToLower(interfaceToString(camelKeySliceVal.Index(i).Interface()))] = 1
					}
				}

				break
			}
		}
	}
	return convertLuaRetToGo(src, camelKey)
}

func camel(rawname string) string {
	if rawname == "" {
		return rawname
	}
	return strings.ToLower(rawname[0:1]) + rawname[1:]
}

func camel2Case(rawname interface{}) string {
	buffer := NewBuffer()
	var name = interfaceToString(rawname)
	for i, r := range name {
		if unicode.IsUpper(r) {
			if i != 0 {
				buffer.Append('_')
			}
			buffer.Append(unicode.ToLower(r))
		} else {
			buffer.Append(r)
		}
	}
	return buffer.String()
}

type Buffer struct {
	*bytes.Buffer
}

func NewBuffer() *Buffer {
	return &Buffer{Buffer: new(bytes.Buffer)}
}

func (b *Buffer) Append(i interface{}) *Buffer {
	switch val := i.(type) {
	case int:
		b.append(strconv.Itoa(val))
	case int64:
		b.append(strconv.FormatInt(val, 10))
	case uint:
		b.append(strconv.FormatUint(uint64(val), 10))
	case uint64:
		b.append(strconv.FormatUint(val, 10))
	case string:
		b.append(val)
	case []byte:
		b.Write(val)
	case rune:
		b.WriteRune(val)
	}
	return b
}

func (b *Buffer) append(s string) *Buffer {
	b.WriteString(s)
	return b
}

func ConvertToTable(state *lua.LState, src interface{}) lua.LValue {
	if reflect.ValueOf(src).Kind() == reflect.Map {
		dst := &lua.LTable{}
		if reflect.TypeOf(src) == reflect.TypeOf(primitive.M{}) {
			for k, v := range src.(primitive.M) {
				dst.RawSet(lua.LString(k), ConvertToTable(state, v))
			}
		} else {
			srcVal := reflect.ValueOf(src)
			for _, key := range srcVal.MapKeys() {
				dst.RawSet(luar.New(state, key.Interface()), ConvertToTable(state, srcVal.MapIndex(key).Interface()))
			}
		}
		return dst
	} else if reflect.ValueOf(src).Kind() == reflect.Slice {
		dst := &lua.LTable{}
		if reflect.TypeOf(src) == reflect.TypeOf(primitive.A{}) {
			for i, v := range src.(primitive.A) {
				dst.RawSetInt(i+1, ConvertToTable(state, v))
			}
		} else {
			srcVal := reflect.ValueOf(src)
			for i := 0; i < srcVal.Len(); i++ {
				dst.Append(ConvertToTable(state, srcVal.Index(i).Interface()))
			}
		}
		return dst
	} else {
		return luar.New(state, src)
	}
}

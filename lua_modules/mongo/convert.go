package mongo

import (
	"fmt"
	lua "github.com/yuin/gopher-lua"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"reflect"
	"time"
)

const (
	convertISODateKeys = "convert_isodate_keys"
	convertOIDKeys     = "convert_oid_keys"
)

func fixOid(bm bson.M) {
	for k, v := range bm {
		if reflect.TypeOf(v) == reflect.TypeOf(primitive.M{}) {
			if v.(primitive.M)["env"] != nil && v.(primitive.M)["value"] != nil {
				bm[k] = v.(primitive.M)["value"]
			} else {
				fixOid(v.(primitive.M))
			}
		}
	}
}

func convertSpecialKeys(bm bson.M, keyName string) {
	if specialKeys, ok := bm[keyName]; ok {
		if reflect.TypeOf(specialKeys) == reflect.TypeOf(primitive.A{}) {
			for _, v := range specialKeys.(primitive.A) {
				if reflect.TypeOf(v).Kind() == reflect.String {
					key := v.(string)
					if bm[key] != nil {
						if reflect.TypeOf(bm[key]) == reflect.TypeOf(primitive.M{}) {
							for k2, v2 := range bm[key].(primitive.M) {
								if keyName == convertISODateKeys {
									if reflect.TypeOf(v2).Kind() == reflect.Float64 {
										bm[key].(primitive.M)[k2] = time.UnixMilli(int64(v2.(float64)))
									}
								} else if keyName == convertOIDKeys {
									if reflect.TypeOf(v2).Kind() == reflect.String {
										oid, err := primitive.ObjectIDFromHex(v2.(string))
										if err == nil {
											bm[key].(primitive.M)[k2] = oid
										}
									} else if reflect.TypeOf(v2).Kind() == reflect.Slice {
										for i, v3 := range v2.(primitive.A) {
											if reflect.TypeOf(v3).Kind() == reflect.String {
												oid, err := primitive.ObjectIDFromHex(v3.(string))
												if err == nil {
													bm[key].(primitive.M)[k2].(primitive.A)[i] = oid
												}
											}
										}
									}
								}
							}
						} else if reflect.TypeOf(bm[key]).Kind() == reflect.String {
							if keyName == convertOIDKeys {
								oid, err := primitive.ObjectIDFromHex(bm[key].(string))
								if err == nil {
									bm[key] = oid
								}
							}
						} else if reflect.TypeOf(bm[key]).Kind() == reflect.Float64 {
							if keyName == convertISODateKeys {
								bm[key] = time.UnixMilli(int64(bm[key].(float64)))
							}
						}
					}
				}
			}
		}
		delete(bm, keyName)
	}
}

func ToGoValue(lv lua.LValue, nameFunc func(s string) string) interface{} {
	switch v := lv.(type) {
	case *lua.LNilType:
		return nil
	case lua.LBool:
		return bool(v)
	case lua.LString:
		return string(v)
	case lua.LNumber:
		return float64(v)
	case *lua.LTable:
		maxn := v.MaxN()
		if maxn == 0 { // table
			ret := make(map[interface{}]interface{})
			v.ForEach(func(key, value lua.LValue) {
				keystr := fmt.Sprint(ToGoValue(key, nameFunc))
				ret[nameFunc(keystr)] = ToGoValue(value, nameFunc)
			})
			return ret
		} else { // array
			ret := make([]interface{}, 0, maxn)
			for i := 1; i <= maxn; i++ {
				ret = append(ret, ToGoValue(v.RawGetInt(i), nameFunc))
			}
			return ret
		}
	case *lua.LUserData:
		return v.Value
	default:
		return v
	}
}

func convertToBson(src interface{}) (interface{}, error) {
	t, dataBytes, err := bson.MarshalValue(src)
	if err != nil {
		return nil, err
	}

	var bs interface{}
	err = bson.UnmarshalValue(t, dataBytes, &bs)
	if err != nil {
		return nil, err
	}

	if bsm, ok := bs.(bson.M); ok {
		// convert to isoDate
		convertSpecialKeys(bsm, convertISODateKeys)

		// convert to oid
		convertSpecialKeys(bsm, convertOIDKeys)
	}

	return bs, nil
}

func convertFromBsonRaw(src bson.Raw) (map[string]interface{}, error) {
	dataBytes, err := bson.Marshal(src)
	if err != nil {
		return nil, err
	}
	dst := make(map[string]interface{})
	err = bson.Unmarshal(dataBytes, &dst)
	if err != nil {
		return nil, err
	}
	return dst, nil
}

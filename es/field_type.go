package es

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type IndexSetting struct {
	NumberOfShards   int `json:"number_of_shards"`
	NumberOfReplicas int `json:"number_of_replicas"`
}

type MappingOptions interface {
	Options() interface{}
}

func (m *Mappings) Options() interface{} {
	return m
}

type PropertiesType map[string]interface{}

type Mappings struct {
	Properties interface{} `json:"properties"`
}

type Field struct {
	Type string `json:"type"`
}

type Index struct {
	Settings *IndexSetting `json:"settings"`
	Mappings *Mappings     `json:"mappings"`
}

type SimpleType struct {
	Type string `json:"type"`
}

type ExtendType struct {
	Type       string         `json:"type"`
	Properties PropertiesType `json:"properties"`
}

type TextType struct {
	Field
	Analyzer string `json:"analyzer"`
}

func Completion() *Field {
	return &Field{Type: "completion"}
}

//func Array() *Field {
//
//}

func Nested(propertiesType PropertiesType) *ExtendType {
	return &ExtendType{
		Type:       "nested",
		Properties: propertiesType,
	}
}

const (
	ES_TYPE_TEXT             = "text"
	ES_TYPE_KEYWORD          = "keyword"
	ES_TYPE_CONSTANT_KEYWORD = "constant_keyword"
)

func textType(tags map[string]string) (resMap map[string]interface{}) {
	resMap = map[string]interface{}{}
	tp, ok := tags["type"]
	if !ok || (tp != ES_TYPE_TEXT && tp != ES_TYPE_KEYWORD && tp != ES_TYPE_CONSTANT_KEYWORD) {
		tp = "text"
	}
	switch tp {
	case ES_TYPE_CONSTANT_KEYWORD:
		if val, ok := tags["value"]; ok {
			if reflect.TypeOf(val).Kind() != reflect.String {
				resMap["value"] = val
			}
		}
	case ES_TYPE_TEXT:
		k1 := "analyzer"
		if val, ok := tags[k1]; ok {
			resMap[k1] = val
		}
	}
	resMap["type"] = tp
	return resMap
}

const (
	ES_TYPE_LONG          = "long"
	ES_TYPE_UNSIGNED_LONG = "unsigned_long"
	ES_TYPE_INTEGER       = "integer"
	ES_TYPE_SHORT         = "short"
	ES_TYPE_BYTE          = "byte"
	ES_TYPE_DOUBLE        = "double"
	ES_TYPE_FLOAT         = "float"
	ES_TYPE_HALF_FLOAT    = "half_float"
)

func numericType(kind reflect.Kind, tags map[string]string) (resMap map[string]interface{}) {
	resMap = map[string]interface{}{}
	tp, ok := tags["type"]
	if ok {
		switch kind {
		case reflect.Int:
			if tp != ES_TYPE_LONG && tp != ES_TYPE_INTEGER && tp != ES_TYPE_SHORT {
				panic(fmt.Sprintf("type error %v", tp))
			}
		case reflect.Int64:
			if tp != ES_TYPE_LONG {
				panic(fmt.Sprintf("type error %v", tp))
			}
		case reflect.Int32:
			if tp != ES_TYPE_INTEGER {
				panic(fmt.Sprintf("type error %v", tp))
			}
		case reflect.Int16:
			if tp != ES_TYPE_SHORT {
				panic(fmt.Sprintf("type error %v", tp))
			}
		case reflect.Float32:
			if tp != ES_TYPE_FLOAT {
				panic(fmt.Sprintf("type error %v", tp))
			}
		case reflect.Float64:
			if tp != ES_TYPE_DOUBLE {
				panic(fmt.Sprintf("type error %v", tp))
			}
		case reflect.Uint, reflect.Uint32, reflect.Uint16, reflect.Uint64:
			if tp != ES_TYPE_UNSIGNED_LONG {
				panic(fmt.Sprintf("type error %v", tp))
			}
		case reflect.Uint8:
			if tp != ES_TYPE_BYTE {
				panic(fmt.Sprintf("type error %v", tp))
			}
		default:
			panic(fmt.Sprintf("unknown reflect type %v", kind))
		}
	} else {
		switch kind {
		case reflect.Int:
			tp = ES_TYPE_LONG
		case reflect.Int64:
			tp = ES_TYPE_LONG
		case reflect.Int32:
			tp = ES_TYPE_INTEGER
		case reflect.Int16:
			tp = ES_TYPE_SHORT
		case reflect.Float32:
			tp = ES_TYPE_FLOAT
		case reflect.Float64:
			tp = ES_TYPE_DOUBLE
		case reflect.Uint, reflect.Uint32, reflect.Uint16, reflect.Uint64:
			tp = ES_TYPE_UNSIGNED_LONG
		case reflect.Uint8:
			tp = ES_TYPE_BYTE
		default:
			panic(fmt.Sprintf("unknown reflect type %v", kind))
		}
	}
	resMap["type"] = tp
	return
}

func baseType(kind reflect.Kind, tags map[string]string) (ifc interface{}) {
	resMap := map[string]interface{}{}
	switch kind {
	case reflect.String:
		return textType(tags)
	case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Float32, reflect.Float64, reflect.Uint, reflect.Uint32, reflect.Uint16, reflect.Uint64, reflect.Uint8:
		return numericType(kind, tags)
	case reflect.Bool:
		return map[string]interface{}{
			"type": "boolean",
		}
	default:
		panic(errors.New(fmt.Sprintf("unknown base type %s", kind.String())))
	}
	return resMap
}

var tagSep = uint8(';')
var kvSep = uint8(':')

var pairMap = map[uint8]uint8{
	'\'': '\'',
	'"':  '"',
}

func parseItems(tagStr string) (items []string) {
	items = []string{}
	patten := tagSep
	mustMatch := []uint8{}
	for i := 0; i < len(tagStr); i++ {
		l := len(mustMatch)
		if l > 0 && tagStr[i] == mustMatch[l-1] {
			mustMatch = mustMatch[:l-1]
		}
		if tagStr[i] == patten && len(mustMatch) == 0 {
			items = append(items, tagStr[:i])
			if i+1 < len(items) {
				items = append(items, parseItems(tagStr[i+1:])...)
			} else {
				break
			}
		}
		if next, ok := pairMap[tagStr[i]]; ok {
			mustMatch = append(mustMatch, next)
		}
	}
	if len(mustMatch) > 0 {
		panic(errors.New(fmt.Sprintf("must match error ! %s", string(mustMatch))))
	}
	return items
}

func parseTag(tagStr string) (tags map[string]string) {
	if tagStr == "-" {
		return map[string]string{}
	}
	if tagStr != "" && tagStr[len(tagStr)-1] != ';' {
		tagStr += string(tagSep)
	}
	tags = map[string]string{}
	tagItems := parseItems(tagStr)
	for _, itemStr := range tagItems {
		if itemStr == "" {
			continue
		}
		items := strings.Split(itemStr, string(kvSep))
		if len(items) < 2 {
			panic(fmt.Sprintf("error format tags %v", itemStr))
		}
		key := items[0]
		value := strings.Join(items[1:], string(kvSep))
		tags[key] = value
	}
	return tags
}

func mappingReflect(inIfc reflect.Type, tags map[string]string, deep int) (out interface{}) {
	t := inIfc
	//fmt.Printf("%s %s\n", inIfc, t.Kind())
	//time.Sleep(time.Second * 3)
	switch t.Kind() {
	case reflect.Pointer:
		//t = reflect.TypeOf(*inIfc)
		return mappingReflect(t.Elem(), tags, deep)
	case reflect.Struct:
		out = map[string]interface{}{}
		newNested := ExtendType{
			Type:       "nested",
			Properties: PropertiesType{},
		}
		for idx := 0; idx < t.NumField(); idx += 1 {
			field := t.Field(idx)
			tagStr := field.Tag.Get("xes")
			if tagStr == "-" {
				continue
			}
			fieldTags := parseTag(tagStr)
			name := field.Name
			tagJsonStr := field.Tag.Get("json")
			if jName, _, _ := strings.Cut(tagJsonStr, ","); jName != "" {
				name = jName
			}
			//if tName, ok := fieldTags["name"]; ok {
			//	name = tName
			//}
			//fmt.Println("struct", name, field.Type.Kind())
			switch field.Type.Kind() {
			case reflect.Ptr, reflect.Array, reflect.Slice:
				//reflect.ValueOf(idx).Elem()
				newNested.Properties[name] = mappingReflect(field.Type.Elem(), fieldTags, deep+1)
			case reflect.Struct:
				newNested.Properties[name] = mappingReflect(field.Type, fieldTags, deep+1)
			default:
				newNested.Properties[name] = baseType(field.Type.Kind(), fieldTags)
			}
		}
		if deep == 0 {
			return newNested.Properties
		} else {
			return newNested
		}
	case reflect.Slice, reflect.Array:
		//newT := t.Elem()
		//switch newT.Kind() {
		//case reflect.Ptr, reflect.Struct, reflect.Array, reflect.Slice:
		//	newNested := ExtendType{
		//		Type:       "nested",
		//		Properties: PropertiesType{},
		//	}
		//	newNested.Properties[name] = mappingReflect(f.Type)
		//default:
		//	newNested.Properties[name] = baseType(f.Type.Kind(), tags)
		//}
		if t.Kind() != reflect.String {
			return mappingReflect(t.Elem(), tags, deep+1)
		}
		fallthrough
	default:
		return baseType(t.Kind(), tags)
	}
}

func MappingIndex(inIfc interface{}) (index *Index) {
	index = &Index{
		Settings: &IndexSetting{
			NumberOfShards:   1,
			NumberOfReplicas: 1,
		},
		Mappings: &Mappings{mappingReflect(reflect.TypeOf(inIfc), nil, 0)},
	}
	return index
}

/*
author:nash
date:2022/9/23
comment: a Go package to interact with arbitrary JSON
*/
package gojson

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strconv"
	"time"

	"github.com/jinzhu/copier"
)

// root 和  parent/key 是互斥的关系
type Json struct {
	//根对象
	root *interface{}

	//用于表示【键值对】对象
	parent map[string]interface{}
	key    string
}

func (t *Json) Load(txt string) (err error) {
	var obj interface{}
	d := json.NewDecoder(bytes.NewBuffer(([]byte)(txt)))
	d.UseNumber()
	err = d.Decode(&obj)
	if err != nil {
		return
	}
	if t.parent != nil {
		t.parent[t.key] = obj
	} else {
		t.root = &obj
	}

	return
}

func (t *Json) LoadString(src string) (err error) {
	var v interface{}
	err = json.Unmarshal([]byte(src), &v)
	if err != nil {
		return
	}
	if t.parent != nil {
		//覆盖当前key的value
		t.parent[t.key] = v
	} else {
		t.root = &v
	}
	return
}

func (t *Json) EncodePretty() (s []byte, err error) {
	if t.root != nil {
		s, err = json.MarshalIndent(t.root, "", "  ")
		if err != nil {
			return
		}

	} else if t.parent != nil {
		s, err = json.MarshalIndent(t.parent[t.key], "", "  ")
		if err != nil {
			return
		}
	}
	return
}
func (t *Json) Encode() (s []byte, err error) {

	if t.root != nil {
		s, err = json.Marshal(t.root)
		if err != nil {
			return
		}
	} else if t.parent != nil {
		s, err = json.Marshal(t.parent[t.key])
		if err != nil {
			return
		}
	}
	return
}

func (t *Json) Set(key string) *Json {

	var child Json
	if t.parent != nil {
		m, ok := t.parent[t.key].(map[string]interface{})
		if ok {
			child.parent = m
			child.key = key
		} else {
			x := make(map[string]interface{})
			child.parent = x
			child.key = key
			t.parent[t.key] = x
		}
	} else {
		if t.root == nil {
			t.root = new(interface{})
		}
		x, ok := (*t.root).(map[string]interface{})
		if !ok {
			x = make(map[string]interface{})
			*t.root = x
		}
		child.parent = x
		child.key = key
	}
	return &child
}

func (t *Json) Del(key string) {
	if t.root != nil {
		x, ok := (*t.root).(map[string]interface{})
		if ok {
			delete(x, key)
		}
	} else if t.parent != nil {
		x, ok := t.parent[t.key].(map[string]interface{})
		if ok {
			delete(x, key)
		}
	}
}

// 设定任意值（和外部的变量脱离关系）
func (t *Json) SetAny(i any) {
	copyI := reflect.New(reflect.TypeOf(i)).Elem().Interface()
	copier.Copy(&copyI, i)

	if t.parent != nil { //当前是子对象
		t.parent[t.key] = copyI //map
	} else {
		if t.root == nil {
			t.root = new(interface{})
		}
		*t.root = copyI
	}
}

func (t *Json) Nil() bool {
	if t.root == nil {
		if t.parent == nil {
			return true
		} else {
			_, ok := t.parent[t.key]
			if !ok {
				return true
			}
		}
	}
	return false
}

func (t *Json) Get(key string) *Json {
	var obj Json
	if t.parent != nil {
		x, ok := t.parent[t.key].(map[string]interface{})
		if ok && x != nil {
			obj.parent = x
			obj.key = key
		}
	} else {
		if t.root == nil {
			return &obj
		}
		x, ok := (*t.root).(map[string]interface{})
		if ok && x != nil {
			obj.parent = x
			obj.key = key
		}
	}
	return &obj
}

// 对数组进行操作
func (t *Json) GetIndex(i int) *Json {
	var obj Json
	if t.parent != nil {
		x, ok := t.parent[t.key].([]interface{})
		if !ok {
			return &obj
		}
		if x != nil && i < len(x) {
			obj.root = &x[i]
		} else {
			return &obj
		}
	}
	if t.root == nil {
		return &obj
	}
	x, ok := (*t.root).([]interface{})
	if !ok {
		return &obj
	}
	if x != nil && i < len(x) {
		obj.root = &x[i]
	}
	return &obj
}

func (t *Json) AddIndex() *Json { //对数组增加一个元素

	var obj Json
	if t.parent != nil {
		x, ok := t.parent[t.key].([]interface{})
		if ok {
			var data interface{}
			x = append(x, data)
		} else {
			x = make([]interface{}, 1)
		}
		t.parent[t.key] = x
		obj.root = &x[len(x)-1]
	} else {
		if t.root == nil {
			t.root = new(interface{})
		}
		x, ok := (*t.root).([]interface{})
		if ok {
			var data interface{}
			x = append(x, data)
		} else {
			x = make([]interface{}, 1)
		}
		*t.root = x
		obj.root = &x[len(x)-1]
	}
	return &obj
}

func (t *Json) ArrayLen() (count int) {
	if t.parent != nil {
		x, ok := t.parent[t.key].([]interface{})
		if ok {
			return len(x)
		}
	} else {
		if t.root != nil {
			x, ok := (*t.root).([]interface{})
			if ok {
				return len(x)
			}
		}
	}
	return 0
}

////////////////////////////////// 返回内部数据/////////////////////
//将json对象 作为 interface{}

func (t *Json) Interface() interface{} {
	if t.parent != nil {
		return t.parent[t.key]
	} else {
		if t.root != nil {
			return *t.root
		}
	}

	return nil
}

//将json内部数据按照Map返回

func (t *Json) Map() map[string]interface{} {
	if t.parent != nil {
		x, ok := t.parent[t.key].(map[string]interface{})
		if ok {
			return x
		}
	} else {
		if t.root != nil {
			x, ok := (*t.root).(map[string]interface{})
			if ok {
				return x
			}
		}
	}
	return nil
}

//将json内部数据按照 Array 返回

func (t *Json) Array() []interface{} {
	if t.parent != nil {
		x, ok := t.parent[t.key].([]interface{})
		if ok && x != nil {
			return x
		}
	} else {
		if t.root != nil {
			x, ok := (*t.root).([]interface{})
			if ok && x != nil {
				return x
			}
		}
	}
	return nil
}

// 作为字符串显示
func (t *Json) String() string {
	if t.parent != nil {
		v := reflect.ValueOf(t.parent[t.key])
		switch v.Kind() {
		case reflect.String:
			return v.String()
		case reflect.Float32, reflect.Float64:
			return strconv.FormatFloat(float64(v.Float()), 'f', -1, 64)
		case reflect.Int, reflect.Int64:
			return strconv.FormatInt(int64(v.Int()), 10)
		case reflect.Uint64, reflect.Uint:
			return strconv.FormatUint(uint64(v.Uint()), 10)
		case reflect.Bool:
			return strconv.FormatBool(v.Bool())
		default:

			switch x := t.parent[t.key].(type) {
			case time.Time:
				return x.Format("2006-01-02 15:04:05")
			}

		}
	} else if t.root != nil {
		v := reflect.ValueOf(*t.root)
		switch v.Kind() {
		case reflect.String:
			return v.String()
		case reflect.Float32, reflect.Float64:
			return strconv.FormatFloat(float64(v.Float()), 'f', -1, 64)
		case reflect.Int, reflect.Int64:
			return strconv.FormatInt(int64(v.Int()), 10)
		case reflect.Uint64, reflect.Uint:
			return strconv.FormatUint(uint64(v.Uint()), 10)
		case reflect.Bool:
			return strconv.FormatBool(v.Bool())
		default:
			switch x := t.parent[t.key].(type) {
			case time.Time:
				return x.Format("2006-01-02 15:04:05")
			}
		}
	}

	return ""
}

// 作为 Int
func (t *Json) Int() int {
	return int(t.Int64())
}

// 作为 Int64
func (t *Json) Int64() int64 {
	if t.parent != nil {
		v := reflect.ValueOf(t.parent[t.key])
		switch v.Kind() {
		case reflect.String:
			x, err := strconv.ParseInt(v.String(), 0, 64)
			if err == nil {
				return x
			}
		case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int8, reflect.Int16:
			return v.Int()
		case reflect.Uint, reflect.Uint16, reflect.Uint64, reflect.Uint8, reflect.Uint32:
			return int64(v.Uint())
		case reflect.Float32, reflect.Float64:
			return int64(v.Float())
		case reflect.Bool:
			if v.Bool() {
				return 1
			} else {
				return 0
			}
		}
	} else if t.root != nil {
		v := reflect.ValueOf(*t.root)
		switch v.Kind() {
		case reflect.String:
			x, err := strconv.ParseInt(v.String(), 0, 64)
			if err == nil {
				return x
			}
		case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int8, reflect.Int16:
			return v.Int()
		case reflect.Uint, reflect.Uint16, reflect.Uint64, reflect.Uint8, reflect.Uint32:
			return int64(v.Uint())
		case reflect.Float32, reflect.Float64:
			return int64(v.Float())
		case reflect.Bool:
			if v.Bool() {
				return 1
			} else {
				return 0
			}
		}
	}
	return 0
}

// 作为 Uint64
func (t *Json) Uint64() uint64 {
	if t.parent != nil {
		v := reflect.ValueOf(t.parent[t.key])
		switch v.Kind() {
		case reflect.String:
			x, err := strconv.ParseUint(v.String(), 0, 64)
			if err == nil {
				return x
			}
		case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int8, reflect.Int16:
			return uint64(v.Int())
		case reflect.Uint, reflect.Uint16, reflect.Uint64, reflect.Uint8, reflect.Uint32:
			return v.Uint()
		case reflect.Float32, reflect.Float64:
			return uint64(v.Float())
		case reflect.Bool:
			if v.Bool() {
				return 1
			} else {
				return 0
			}
		}
	} else if t.root != nil {
		v := reflect.ValueOf(*t.root)
		switch v.Kind() {
		case reflect.String:
			x, err := strconv.ParseUint(v.String(), 0, 64)
			if err == nil {
				return x
			}
		case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int8, reflect.Int16:
			return uint64(v.Int())
		case reflect.Uint, reflect.Uint16, reflect.Uint64, reflect.Uint8, reflect.Uint32:
			return v.Uint()
		case reflect.Float32, reflect.Float64:
			return uint64(v.Float())
		case reflect.Bool:
			if v.Bool() {
				return 1
			} else {
				return 0
			}
		}
	}
	return 0
}

//作为 Float64

func (t *Json) Float64() float64 {
	if t.parent != nil {
		v := reflect.ValueOf(t.parent[t.key])
		switch v.Kind() {
		case reflect.String:
			x, err := strconv.ParseFloat(v.String(), 64)
			if err == nil {
				return x
			}
		case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int8, reflect.Int16:
			return float64(v.Int())
		case reflect.Uint, reflect.Uint16, reflect.Uint64, reflect.Uint8, reflect.Uint32:
			return float64(v.Uint())
		case reflect.Float32, reflect.Float64:
			return v.Float()
		case reflect.Bool:
			if v.Bool() {
				return 1.0
			} else {
				return 0.0
			}
		}
	} else if t.root != nil {
		v := reflect.ValueOf(*t.root)
		switch v.Kind() {
		case reflect.String:
			x, err := strconv.ParseFloat(v.String(), 64)
			if err == nil {
				return x
			}
		case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int8, reflect.Int16:
			return float64(v.Int())
		case reflect.Uint, reflect.Uint16, reflect.Uint64, reflect.Uint8, reflect.Uint32:
			return float64(v.Uint())
		case reflect.Float32, reflect.Float64:
			return v.Float()
		case reflect.Bool:
			if v.Bool() {
				return 1.0
			} else {
				return 0.0
			}
		}
	}
	return 0.0
}

// 作为 Bool 类型
func (t *Json) Bool() bool {
	if t.parent != nil {
		v := reflect.ValueOf(t.parent[t.key])
		switch v.Kind() {
		case reflect.Bool:
			return v.Bool()
		case reflect.String:
			x, err := strconv.ParseBool(v.String())
			if err == nil {
				return x
			}
		case reflect.Int, reflect.Int64:
			if v.Int() != 0 {
				return true
			}
		case reflect.Uint, reflect.Uint64:
			if v.Uint() != 0 {
				return true
			}
		case reflect.Float64, reflect.Float32:
			if v.Float() <= 0.000001 && v.Float() >= -0.000001 {
				return false
			} else {
				return true
			}
		}
	}
	return false
}

func (t *Json) Time() (res time.Time, err error) {
	if t.parent != nil {
		k, ok := (t.parent[t.key]).(time.Time)
		if ok {
			return k, nil
		}

	}

	str := t.String()
	res, err = time.ParseInLocation("2006-01-02", str, time.Local)
	if err != nil {
		res, err = time.ParseInLocation("2006-01-02 15:04:05", str, time.Local)
		if err != nil {
			res, err = time.ParseInLocation("2006-01-02T15:04:05.999999999+08:00", str, time.Local)
		}
	}
	return
}

func (t *Json) Kind() reflect.Kind {
	if t.parent != nil {
		return reflect.ValueOf(t.parent[t.key]).Kind()
	} else {
		if t.root != nil {
			return reflect.ValueOf(*t.root).Kind()
		}
	}

	return reflect.Invalid
}

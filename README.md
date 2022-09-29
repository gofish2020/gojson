# gojson
a Go package to interact with arbitrary JSON


1. Construct json object by API;
2. Decode json string, then Get the target value;
3. Function is simple and easy;
```go
package main

import (
	"fmt"
	"time"

	"github.com/gofish2020/gojson"
)

func main() {

    //construct json
	var js gojson.Json
	js.Set("ints").SetAny(100)                           //整数
	js.Set("strs").SetAny("T")                           //字符串
	js.Set("floats").SetAny(3.1415926)                   //浮点数
	js.Set("timestamp").SetAny(time.Now())               //时间戳
	js.Set("timestampstr").SetAny("2009-01-02 15:02:01") //字符串表示时间
	js.Set("boolen").SetAny(true)                        //布尔
	js.Set("array").AddIndex().SetAny(1)
	js.Set("array").AddIndex().SetAny("2")
	fmt.Println("数组长度:", js.Get("array").ArrayLen())
	//遍历数组
	for i := 0; i < js.Get("array").ArrayLen(); i++ {
		fmt.Println(js.Get("array").GetIndex(i).String())
	}
    //load json string  and set loadstr
	js.Set("loadstr").LoadString(`{"nash":111,"ss":1.2,"bb":"23"}`) //用字符串设定loadstr的值
	fmt.Println("///////////////////////")

    //one json object set other json object
	var newJs gojson.Json
	newJs.SetAny(js.Interface()) //用json对象设置值

	var newJs1 gojson.Json
	newJs1.Set("copyJs").SetAny(js.Interface())

    // json object to json string
	jsStr1, err1 := newJs1.Encode()
	fmt.Println("------>>>", string(jsStr1), err1)


    // get key's different type value
	fmt.Println("将js对象转化json字符串///////////////////////")
	jsStr, err := js.Encode()
	fmt.Println("======>>>", string(jsStr), err)
	fmt.Println("转化为字符串类型：///////////////////////")
	fmt.Println(js.Get("ints").String())
	fmt.Println(js.Get("strs").String())
	fmt.Println(js.Get("floats").String())
	fmt.Println(js.Get("timestamp").String()) //time类型转成字符串
	fmt.Println(js.Get("boolen").String())
	fmt.Println(js.Get("array").GetIndex(0).String())

	fmt.Println("转化为整数类型：///////////////////////")
	fmt.Println(js.Get("ints").Int())
	fmt.Println(js.Get("strs").Int())
	fmt.Println(js.Get("floats").Int())
	fmt.Println(js.Get("timestamp").Int())
	fmt.Println(js.Get("boolen").Int())
	fmt.Println("转化为布尔类型：///////////////////////")
	fmt.Println(js.Get("ints").Bool())
	fmt.Println(js.Get("strs").Bool())
	fmt.Println(js.Get("floats").Bool())
	fmt.Println(js.Get("timestamp").Bool())
	fmt.Println(js.Get("boolen").Bool())

	fmt.Println("转化为浮点数类型：///////////////////////")
	fmt.Println(js.Get("ints").Float64())
	fmt.Println(js.Get("strs").Float64())
	fmt.Println(js.Get("floats").Float64())
	fmt.Println(js.Get("timestamp").Float64())
	fmt.Println(js.Get("boolen").Float64())

	fmt.Println("转化为时间类型：///////////////////////")

	fmt.Println(js.Get("timestamp").Time())
	fmt.Println(js.Get("timestampstr").Time()) //字符串转成time类型

}

```
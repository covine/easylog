# Easy Logger

[![godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/govine/easylog)[![license](https://img.shields.io/github/license/govine/easylog)](https://github.com/govine/easylog/blob/master/LICENSE)[![Build Status](https://travis-ci.com/govine/easylog.svg?branch=master)](https://travis-ci.com/govine/easylog)[![Coverage](http://gocover.io/_badge/github.com/govine/easylog)](http://gocover.io/github.com/govine/easylog)
## 特性

* 快速
* 支持模块层级
* 低内存分配
* 自定义处理函数
* 支持上下文
* 日志缓存

## 安装

```go
go get -u github.com/govine/easylog
```

## 开始

### 简单使用

```go
package easylog_test

import (
	"github.com/govine/easylog"
	"github.com/govine/easylog/handler"
)

func ExampleEasylog_Simple() {
	stdoutHandler := handler.NewStdoutHandler(nil) // 将输出定向到 stdout
	easylog.AddHandler(stdoutHandler)              // 添加到handlers
	easylog.SetLevel(easylog.DEBUG)                // 设置日志级别
	easylog.Debug().Msg("hello world")
	// Output: hello world
}
```
> Note: 未制定logger对象则默认使用root logger

### 自定义输出格式

**easylog** 允许在日志上下文中添加键值对，并自定义输出格式:

```go
package easylog_test

import (
	"encoding/json"

	"github.com/govine/easylog"
	"github.com/govine/easylog/handler"
)

func format(record *easylog.Record) string {
	type output struct {
		Fields map[string]interface{} `json:"fields"`
		Msg    string                 `json:"msg"`
	}
	o := output{
		Fields: record.FieldMap,
		Msg:    record.Message,
	}
	b, err := json.Marshal(o)
	if err == nil {
		return string(b[:])
	}
	return ""
}

func ExampleEasylog_SimpleFormat() {
	stdoutHandler := handler.NewStdoutHandler(format) // 将输出定向到 stdout
	easylog.AddHandler(stdoutHandler)                 // 添加到handlers
	easylog.SetLevel(easylog.DEBUG)                   // 设置日志级别
	easylog.Debug().Fields(map[string]interface{}{"name": "dog"}).Msg("hello")
	// Output:
	// {"fields":{"name":"dog"},"msg":"hello"}
}
```

### 输出到文件

```go
package easylog_test

import (
	"encoding/json"

	"github.com/govine/easylog"
	"github.com/govine/easylog/handler"
)

func ExampleEasylog_File() {
	fileHandler, err := handler.NewFileHandler("./log/file.log", nil)
	if err != nil {
		return
	}
	easylog.AddHandler(fileHandler)
	easylog.SetLevel(easylog.DEBUG)
	easylog.Debug().Fields(map[string]interface{}{"name": "dog"}).Msg("hello")
	// Output:
}
```
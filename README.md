# Easy Logger

[![godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/govine)[![license](https://img.shields.io/github/license/govine/easylog)](https://github.com/govine/blob/master/LICENSE)[![Build Status](https://travis-ci.com/govine/easylog.svg?branch=master)](https://travis-ci.com/govine/easylog)[![Coverage](https://gocover.io/_badge/github.com/govine)](https://gocover.io/github.com/govine)

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

### 输出到切割日志文件

```go
package easylog_test

import (
	"encoding/json"

	"github.com/govine/easylog"
	"github.com/govine/easylog/handler"
)

func ExampleEasylog_RotateFile() {
	l := easylog.GetLogger("rotateFile")
	l.SetLevel(easylog.DEBUG)
	l.SetPropagate(false)
	
	w, err := writer.NewRotateFileWriter(30, "./log/debug.log", 409600)
	if err != nil {
		return
	}
	rotateFileHandler := handler.NewRotateFileHandler(nil, w)
	rotateFileHandler.SetLevel(easylog.DEBUG)
	rotateFileHandler.AddFilter(&filter.LevelEqualFilter{Level: easylog.DEBUG})
	
	l.AddHandler(rotateFileHandler)
}
```

### 记录程序执行栈
```go
	l := easylog.GetLogger("test")
	l.EnableFrame(easylog.WARN)
```

### 给handler添加过滤器
```go
	rotateFileHandler.AddFilter(&filter.LevelEqualFilter{Level: easylog.DEBUG})
```

### 给handler设置级别
```go
	rotateFileHandler.SetLevel(easylog.DEBUG)
```

## 高级用法

### 划分模块层级
```
包结构
a 
│
└───b
│   │
│   └───c
│   
└───d
│   │
│   └───e
```

对应层级logger

```go
	a := easylog.GetLogger("a")
	ab := easylog.GetLogger("a.b")
	ad := easylog.GetLogger("a.d")
	ade := easylog.GetLogger("a.d.e")
	// ab, ad, ade 收集的日志都会流经 a 的handlers
	// ade 收集的日志会流经 ad 的 handlers
```

```go
// a 包
package c

import (
	"github.com/govine/easylog"
)

var logger *easylog.Logger

func init() {
	logger = easylog.GetLogger("a")
}

func Test() {
	logger.Debug().Msg("hello world")
}
```

```go
// c 包
package c

import (
	"github.com/govine/easylog"
)

var logger *easylog.Logger

func init() {
	logger = easylog.GetLogger("a.b.c")
}

func Test() {
	logger.Debug().Msg("hello world")
}
```

### 自定义handler

```go
/*
实现
type IHandler interface {
	Handle(*Record)
	Flush()
	Close()
}
接口
使用 easylog.NewEasyLogHandler() 包装
*/
package handler

import (
	"fmt"
	"os"

	"github.com/govine/easylog"
)

type StderrHandler struct {
	format easylog.Formatter
}

func (s *StderrHandler) Handle(record *easylog.Record) {
	var str string
	if s.format != nil {
		str = s.format(record)
	} else {
		str = fmt.Sprintf(record.Message, record.Args...)
	}

	_, _ = os.Stderr.Write([]byte(str + "\n"))
}

func (s *StderrHandler) Flush() {
}

func (s *StderrHandler) Close() {
}

func NewStderrHandler(format easylog.Formatter) easylog.IEasyLogHandler {
	return easylog.NewEasyLogHandler(&StderrHandler{
		format: format,
	})
}

```

### 请求内异步

> Note: 异步的前提是日志的顺序得到保证

```go
{
	// 设置接收的父级logger
	logger := easylog.NewCachedLogger(parentLogger)
	// 打开向父级传递
	logger.SetPropagate(true)
	// 设置等级
	logger.SetLevelByString("info")
	// 打开某些级别的执行上下文记录
	logger.EnableFrame(easylog.DEBUG)
	logger.EnableFrame(easylog.WARN)
	logger.EnableFrame(easylog.WARNING)
	logger.EnableFrame(easylog.FATAL)
	
	defer func() {
		// 在请求结束后异步按顺序处理
		// 会触发父级的处理
		go logger.Close()
	}()
}
```

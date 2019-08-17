package writer

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type RotateFileWriter struct {
	file       string
	interval   int64 // second
	bufSize    int
	mu         sync.RWMutex
	ctx        context.Context
	cancel     context.CancelFunc
	fileWriter *os.File
	bufWriter  *bufio.Writer
}

func NewRotateFileWriter(interval int64, file string, bufSize int) (*RotateFileWriter, error) {
	// 防止/0
	if interval < minInterval {
		return nil, errors.New("gap little than 30 seconds")
	}

	dir := filepath.Dir(file)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, errors.New(fmt.Sprintf("failed to create directory: %v", err))
	}

	fileWriter, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("open file %s fail: %v", file, err))
	}

	var bs int
	if bufSize >= maxBufSize {
		bs = maxBufSize
	} else if bufSize <= 0 {
		bs = defaultBufSize
	} else {
		bs = bufSize
	}

	bufWriter := bufio.NewWriterSize(fileWriter, bs)

	ctx, cancel := context.WithCancel(context.Background())

	r := &RotateFileWriter{
		file:       file,
		interval:   interval,
		bufSize:    bufSize,
		ctx:        ctx,
		cancel:     cancel,
		fileWriter: fileWriter,
		bufWriter:  bufWriter,
	}
	go r.rotate()

	return r, nil
}

func (r *RotateFileWriter) cut() {
	dateNow := time.Now().Unix() - r.interval
	dateNow = (int64(float64(dateNow)+0.5*float64(r.interval)) / r.interval) * r.interval
	dateFormat := time.Unix(dateNow, 0).Format("200601021504")
	rotateFile := r.file + "." + dateFormat + "00"

	_, err := os.Stat(rotateFile)
	if err != nil {
		if r.bufWriter != nil {
			_ = r.bufWriter.Flush()
			r.bufWriter = nil
		}
		if r.fileWriter != nil {
			_ = r.fileWriter.Close()
			r.fileWriter = nil
		}

		// 切换文件，检查日志文件 xx.log
		_, e := os.Stat(r.file)
		if e != nil {
			// 找不到 xx.log 文件，创建切割文件，创建日志文件
			r.fileWriter, e = os.OpenFile(r.file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
			if e != nil {
				log.Printf("create log file fail: %v", err)
				r.fileWriter = nil
				r.bufWriter = nil
			}
			r.bufWriter = bufio.NewWriterSize(r.fileWriter, r.bufSize)

			// 因为没法rename，就直接创建
			rf, e := os.OpenFile(rotateFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
			if e != nil {
				// 如果创建xx.20190522150400 不成功，不影响日志写入，使用原有文件 xx.log
				log.Printf("create rotate file fail: %v", err)
			}
			if rf != nil {
				_ = rf.Close()
			}
		} else {
			// 存在 xx.log 文件, 重命名 xx.log -> xx.20190522150400
			e := os.Rename(r.file, rotateFile)
			if e != nil {
				log.Printf("rename file %s to %s fail: %v ", r.file, rotateFile, e.Error())
			} else {
				// 重新创建 xx.log
				r.fileWriter, e = os.OpenFile(r.file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
				if e != nil {
					log.Printf("create log file fail: %v", err)
					r.fileWriter = nil
					r.bufWriter = nil
				}
				r.bufWriter = bufio.NewWriterSize(r.fileWriter, r.bufSize)
			}
		}
	} else {
		// 已经切割过，或是其它原因xx.20190522150400已经存在
	}
}

func (r *RotateFileWriter) calSleepTime() time.Duration {
	nowTime := time.Now().Unix()
	nextTime := (nowTime/r.interval + 1) * r.interval
	return time.Duration(nextTime - nowTime)
}

func (r *RotateFileWriter) rotate() {
	for {
		select {
		case <-r.ctx.Done():
			return
		default:
			time.Sleep(r.calSleepTime() * time.Second)
			r.mu.Lock()
			r.cut()
			r.mu.Unlock()
		}
	}
}

func (r *RotateFileWriter) Write(p []byte) (n int, err error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// 在文件切换的时候，如果发生故障
	if r.bufWriter == nil {
		return 0, errors.New("buffer writer is nil")
	}

	return r.bufWriter.Write(p)
}

func (r *RotateFileWriter) Flush() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.bufWriter == nil {
		return errors.New("buffer writer is nil")
	}

	return r.bufWriter.Flush()
}

func (r *RotateFileWriter) Close() {
	r.mu.Lock()
	defer r.mu.Unlock()

	// 关闭文件切割
	r.cancel()

	if r.bufWriter != nil {
		_ = r.bufWriter.Flush()
		r.bufWriter = nil
	}

	// 关闭文件
	if r.fileWriter != nil {
		_ = r.fileWriter.Close()
		r.fileWriter = nil
	}
}

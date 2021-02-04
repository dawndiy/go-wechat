package upload

import (
	"bytes"
	"io"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"time"
)

// Uploader 上传资源接口
type Uploader interface {
	// Name 资源名称
	Name() string
	// Reader 资源内容
	Reader() io.Reader
	// Err 资源读取错误
	Err() error
}

// bytesUploader 实现 Uploader 接口
type bytesUploader struct {
	filename string
	b        []byte
	err      error
}

// Name 资源名称
func (u *bytesUploader) Name() string { return u.filename }

// Reader 资源内容
func (u *bytesUploader) Reader() io.Reader { return bytes.NewReader(u.b) }

// Err 资源读取错误
func (u *bytesUploader) Err() error { return u.err }

// UploadBytes 上传 bytes
func UploadBytes(filename string, b []byte) Uploader {
	return &bytesUploader{filename, b, nil}
}

// UploadFile 上传文件
func UploadFile(name string) Uploader {
	_, filename := filepath.Split(name)
	b, err := ioutil.ReadFile(name)
	if err != nil {
		return &bytesUploader{err: err}
	}

	return &bytesUploader{
		filename: filename,
		b:        b,
		err:      nil,
	}
}

// Result 上传临时素材结果
type Result struct {
	// 媒体文件类型
	Type string `json:"type"`

	// 媒体文件上传后获取的唯一标识，3天内有效
	MediaID string `json:"media_id"`

	// 媒体文件上传时间戳
	CreatedAt int64 `json:"created_at"`
}

// GetCreatedAt 返回 time.Time 类型
func (r *Result) GetCreatedAt() time.Time {
	return time.Unix(r.CreatedAt, 0)
}

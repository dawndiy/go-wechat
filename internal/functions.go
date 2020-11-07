package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"strings"

	"github.com/dawndiy/go-wechat/httpclient"
)

// NewJSONRequest 新建一个 JSON 接口请求
func NewJSONRequest(ctx context.Context, method string, u *url.URL, body interface{}) (*http.Request, error) {
	var buf io.ReadWriter
	if body != nil {
		buf = &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		err := enc.Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return req, nil
}

// NewUploadRequest 文件上传请求
func NewUploadRequest(
	ctx context.Context, u *url.URL, fieldname, filename string, reader io.Reader) (*http.Request, error) {

	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(nil)
	mpWriter := multipart.NewWriter(buf)
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"; filelength=%d`,
			escapeQuotes(fieldname), escapeQuotes(filename), len(b)))
	h.Set("Content-Type", http.DetectContentType(b))

	w, err := mpWriter.CreatePart(h)
	if err != nil {
		return nil, err
	}
	if _, err := w.Write(b); err != nil {
		return nil, err
	}
	if err = mpWriter.Close(); err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, u.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", mpWriter.FormDataContentType())
	return req, nil
}

type ErrCheckFunc func(*http.Response) error

// DoJSONRequest 执行一个接口请求
func DoJSONRequest(
	h httpclient.RequestHandler,
	req *http.Request,
	v interface{},
	checkFunc ErrCheckFunc) (*http.Response, error) {

	resp, err := h.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	err = checkFunc(resp)
	if err != nil {
		return resp, err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			io.Copy(w, resp.Body)
		} else {
			decErr := json.NewDecoder(resp.Body).Decode(v)
			if decErr == io.EOF {
				decErr = nil // ignore EOF errors caused by empty response body
			}
			if decErr != nil {
				err = decErr
			}
		}
	}

	return resp, err
}

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}

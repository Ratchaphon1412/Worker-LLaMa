package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type API struct {
	BaseURL string
}

type Params struct {
	Headers map[string]string // ถ้าต้องการใช้เป็น query params หรือ header แยกเพิ่มเติมได้
	Query   map[string]string
}

// Body กำหนด contract ให้แต่ละชนิดสามารถคืน Content-Type และ io.Reader ได้
type Body interface {
	ContentType() string
	Reader() (io.Reader, error)
}

// FormBody ใช้ส่งข้อมูลแบบ application/x-www-form-urlencoded
type FormBody url.Values

func (f FormBody) ContentType() string {
	return "application/x-www-form-urlencoded"
}

func (f FormBody) Reader() (io.Reader, error) {
	v := url.Values(f)
	return strings.NewReader(v.Encode()), nil
}

// JSONBody ใช้ส่งข้อมูลแบบ application/json
type JSONBody map[string]interface{}

func (j JSONBody) ContentType() string {
	return "application/json"
}

func (j JSONBody) Reader() (io.Reader, error) {
	b, err := json.Marshal(j)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(b), nil
}

func (r *API) Get(params Params) (*http.Response, error) {

	q := url.Values{}
	for key, value := range params.Query {
		q.Add(key, value)
		fmt.Println("key:", key, "value:", value)
	}
	fmt.Println("endcode", q.Encode())
	req, err := http.NewRequest("GET", fmt.Sprintf("%s?%s", r.BaseURL, q.Encode()), nil) // สร้าง GET request[1]
	if err != nil {
		return nil, err
	}

	for key, value := range params.Headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r *API) Post(params Params, body Body) (*http.Response, error) {
	u, err := url.Parse(r.BaseURL)
	if err != nil {
		return nil, err
	}
	q := u.Query()
	for key, value := range params.Headers {
		q.Add(key, value)
	}
	u.RawQuery = q.Encode()

	reader, err := body.Reader()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", u.String(), reader) // สร้าง POST request พร้อม body[1]
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", body.ContentType())

	for key, value := range params.Headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

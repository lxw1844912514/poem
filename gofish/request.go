package gofish

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Request struct {
	Url     string
	Method  string
	Headers *http.Header
	Body    io.Reader
	Handle  Handle
	Client  http.Client
}

func (r *Request) Do() error {
	request, err := http.NewRequest(r.Method, r.Url, r.Body)
	if err != nil {
		return err
	}

	request.Header = *r.Headers
	resp, err := r.Client.Do(request)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status error code: %d", resp.StatusCode)
	}

	r.Handle.Worker(resp.Body,r.Url)

	defer resp.Body.Close()
	return nil
}

func  NewRequest(method, Url, userAgent string, handle Handle, body io.Reader) (*Request, error) {
	//解析URL是否正确
	_, err := url.Parse(Url)
	if err != nil {
		return nil, err
	}

	//添加header头
	header := http.Header{}
	if userAgent != "" {
		header.Add("User-Agent", userAgent)
	} else {
		header.Add("User-Agent", UserAgent)
	}

	//client
	client := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return nil
		},
	}

	return &Request{
		Url:     Url,
		Method:  method,
		Headers: &header,
		Handle:  handle,
		Body:    body,
		Client:  client,
	}, nil
}

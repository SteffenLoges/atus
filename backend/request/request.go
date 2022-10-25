package request

import (
	"atus/backend/config"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type Request struct {
	Raw *http.Request
}

func NewWithContext(ctx context.Context, method, url string, body io.Reader) (*Request, error) {

	r, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	return &Request{
		Raw: r,
	}, nil
}

func New(method, url string, body io.Reader) (*Request, error) {
	return NewWithContext(context.Background(), method, url, body)
}

func (r *Request) Do() (*http.Response, error) {

	for key, value := range genericHeaders {
		r.Raw.Header.Set(key, value)
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: config.Base.Requests.InsecureSkipVerify,
		},
	}

	client := &http.Client{
		Transport: tr,
	}

	resp, err := client.Do(r.Raw)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()

		buf, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("server returned error %d: %s", resp.StatusCode, buf)
	}

	return resp, nil

}

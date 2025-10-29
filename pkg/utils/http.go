package utils

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

type DoRequestValues struct {
	Method  string
	URL     string
	Body    io.Reader
	Headers map[string]string
}

func DoRequest(ctx context.Context, values DoRequestValues) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, values.Method, values.URL, values.Body)
	if err != nil {
		return nil, fmt.Errorf("new http request %+v | %w", values, err)
	}

	for k, v := range values.Headers {
		req.Header.Set(k, v)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do http request %+v | %w", values, err)
	}

	return resp, nil
}

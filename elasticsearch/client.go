package elasticsearch

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"time"

	esv8 "github.com/elastic/go-elasticsearch/v8"
	esapiv8 "github.com/elastic/go-elasticsearch/v8/esapi"
)

var (
	ErrEsNotFound = errors.New("elasticsearch: not found")
)

type EsClient struct {
	client *esv8.Client
}

func NewClient(hosts []string, username, password string) (client *EsClient, err error) {
	c, err := esv8.NewClient(esv8.Config{
		Addresses: hosts,
		Username:  username,
		Password:  password,
		Transport: &http.Transport{
			MaxIdleConnsPerHost:   http.DefaultMaxIdleConnsPerHost,
			ResponseHeaderTimeout: 5 * time.Second,
			DialContext:           (&net.Dialer{Timeout: 5 * time.Second}).DialContext,
			TLSClientConfig: &tls.Config{
				MinVersion: tls.VersionTLS11,
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	return &EsClient{client: c}, nil
}

func (es *EsClient) handleResponse(resp *esapiv8.Response) (map[string]interface{}, error) {
	var (
		r map[string]interface{}
	)
	if resp.StatusCode == 404 {
		return r, ErrEsNotFound
	}
	// Check response status
	if resp.IsError() {
		return nil, errors.New(resp.String())
	}
	// Deserialize the response into a map.
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, errors.New(fmt.Sprintf("Error parsing the response body: %s", err))
	}
	return r, nil
}

func doRequest(ctx context.Context, transport esapiv8.Transport, req esapiv8.Request, out interface{}) error {
	resp, err := req.Do(ctx, transport)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.IsError() {
		bytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return &Error{
			StatusCode: resp.StatusCode,
			Header:     resp.Header,
			body:       string(bytes),
		}
	}
	if out != nil {
		err = json.NewDecoder(resp.Body).Decode(out)
	}
	return err
}

// Error holds the details for a failed Elasticsearch request.
//
// Error is only returned for request is serviced, and not when
// a client or network failure occurs.
type Error struct {
	// StatusCode holds the HTTP response status code.
	StatusCode int

	// Header holds the HTTP response headers.
	Header http.Header

	body string
}

func (e *Error) Error() string {
	if e.body != "" {
		return e.body
	}
	return http.StatusText(e.StatusCode)
}

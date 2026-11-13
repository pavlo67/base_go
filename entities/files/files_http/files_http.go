package files_http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/pavlo67/base_go/entities/files"
	"github.com/pavlo67/base_go/lib/errors"
	"github.com/pavlo67/base_go/lib/logger"
)

var _ files.Operator = &filesHTTP{}

type filesHTTP struct {
	baseURL string
	client  *http.Client
}

var l logger.Operator

const onNew = "on files_http.New():"

func New(baseURL string, l_ logger.Operator) (files.Operator, error) {
	if l_ == nil {
		return nil, errors.New("", onNew+" l_ == nil")
	}
	l = l_

	baseURL = strings.TrimRight(strings.TrimSpace(baseURL), "/")
	if baseURL == "" {
		return nil, errors.New("", onNew+" baseURL is empty")
	}
	if _, err := url.ParseRequestURI(baseURL); err != nil {
		return nil, errors.Wrap(err, onNew)
	}

	return &filesHTTP{
		baseURL: baseURL,
		client:  http.DefaultClient,
	}, nil
}

// files.Operator -----------------------------------------------------------------------------

const onSave = "on files_http.Save():"

func (op *filesHTTP) Save(data files.Data) error {
	body, err := json.Marshal(data)
	if err != nil {
		return errors.Wrap(err, onSave)
	}

	return errors.Wrap(op.doJSON(http.MethodPost, "/save", nil, bytes.NewReader(body), nil), onSave)
}

const onRead = "on files_http.Read():"

func (op *filesHTTP) Read(path string) (*files.Item, error) {
	q := url.Values{}
	q.Set("path", path)

	var resp readResponse
	if err := op.doJSON(http.MethodGet, "/read", q, nil, &resp); err != nil {
		return nil, errors.Wrap(err, onRead)
	}
	return resp.Item, nil
}

const onRemove = "on files_http.Remove():"

func (op *filesHTTP) Remove(path string, forceRecursion bool) error {
	q := url.Values{}
	q.Set("path", path)
	q.Set("force", strconv.FormatBool(forceRecursion))

	return errors.Wrap(op.doJSON(http.MethodPost, "/remove", q, nil, nil), onRemove)
}

const onList = "on files_http.List():"

func (op *filesHTTP) List(path string, depth int) ([]files.Item, error) {
	q := url.Values{}
	q.Set("path", path)
	q.Set("depth", strconv.Itoa(depth))

	var resp listResponse
	if err := op.doJSON(http.MethodGet, "/list", q, nil, &resp); err != nil {
		return nil, errors.Wrap(err, onList)
	}
	return resp.Items, nil
}

// HTTP helpers --------------------------------------------------------------------------------

type readResponse struct {
	Item *files.Item `json:"item"`
}

type listResponse struct {
	Items []files.Item `json:"items"`
}

type errorResponse struct {
	Error string `json:"error,omitempty"`
}

func (op *filesHTTP) doJSON(method, route string, q url.Values, body io.Reader, out any) error {
	reqURL := op.baseURL + route
	if len(q) > 0 {
		reqURL += "?" + q.Encode()
	}

	req, err := http.NewRequest(method, reqURL, body)
	if err != nil {
		return err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")

	resp, err := op.client.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			l.Errorf("on files_http.doJSON(): on resp.Body.Close(): %v", err)
		}
	}()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		var errResp errorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err == nil && errResp.Error != "" {
			return fmt.Errorf("http status %d: %s", resp.StatusCode, errResp.Error)
		}
		return fmt.Errorf("http status %d", resp.StatusCode)
	}

	if out == nil {
		return nil
	}
	return json.NewDecoder(resp.Body).Decode(out)
}

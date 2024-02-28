package rpc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type APIClient struct {
	client  http.Client
	traceID string
}

func (c *APIClient) SetTraceID(traceID string) {
	c.traceID = traceID
}

func (c *APIClient) GetTraceID() (traceID string) {
	if c.traceID == "" {
		return XHeaderLogID
	}
	return c.traceID
}

func NewClientWithTimeout(timeout time.Duration) *APIClient {
	return &APIClient{
		client: http.Client{Timeout: timeout},
	}
}

func (c *APIClient) RawClient() *http.Client {
	return &c.client
}

func (c *APIClient) Call(ctx context.Context, reqUrl, method string, header http.Header, query url.Values, body []byte, apiRet APIRet) (err error) {
	reqURI, pErr := url.Parse(reqUrl)
	if pErr != nil {
		err = fmt.Errorf("parse request url error, %s", reqUrl)
		return
	}
	if len(query) > 0 {
		reqURI.RawQuery = query.Encode()
	}
	req, newErr := http.NewRequest(method, reqURI.String(), bytes.NewBuffer(body))
	if newErr != nil {
		err = fmt.Errorf("new request error, %s", newErr.Error())
		return
	}
	// add X-ReqId if set in context
	if ctx != nil && ctx.Value(c.GetTraceID()) != nil {
		req.Header.Add(c.GetTraceID(), ctx.Value(c.GetTraceID()).(string))
	}
	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	// copy the extra headers
	for key, values := range header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	// fire the request
	resp, callErr := c.client.Do(req)
	if callErr != nil {
		err = fmt.Errorf("call api failed, %s", callErr.Error())
		return
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	// parse api ret & check response code
	jsonDecoder := json.NewDecoder(resp.Body)
	if apiRet == nil {
		apiRet = &BaseAPIRet{}
	}
	// set X-Request-ID
	apiRet.SetRequestID(resp.Header.Get(c.GetTraceID()))
	// parse body
	decodeErr := jsonDecoder.Decode(apiRet)
	if resp.StatusCode/100 != 2 {
		// check for normal logic
		if decodeErr != nil {
			err = fmt.Errorf("status=%s", resp.Status)
		} else {
			err = fmt.Errorf("status=%s, error=%s, %s", resp.Status, apiRet.RetCode(), apiRet.RetMessage())
		}
	} else if !apiRet.IsOk() {
		err = fmt.Errorf("%s, %s", apiRet.RetCode(), apiRet.RetMessage())
	}
	return
}

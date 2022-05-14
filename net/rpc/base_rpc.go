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
	client http.Client
}

func NewClientWithTimeout(timeout time.Duration) *APIClient {
	return &APIClient{
		client: http.Client{Timeout: timeout},
	}
}

func (c APIClient) Call(ctx context.Context, reqUrl, method string, header http.Header, query url.Values, body []byte) (apiRet APIRet, err error) {
	reqURI, pErr := url.Parse(reqUrl)
	if pErr != nil {
		err = fmt.Errorf("parse request url failed, %s", pErr.Error())
		return
	}
	// set the query
	if len(query) > 0 {
		reqURI.RawQuery = query.Encode()
	}
	// create new request
	req, newErr := http.NewRequest(method, reqURI.String(), bytes.NewBuffer(body))
	if newErr != nil {
		err = fmt.Errorf("new request failed, %s", newErr.Error())
		return
	}
	// add X-ReqId if set in context
	if ctx != nil && ctx.Value(XHeaderLogId) != nil {
		req.Header.Add(XHeaderLogId, ctx.Value(XHeaderLogId).(string))
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
	// parse api ret
	jsonDecoder := json.NewDecoder(resp.Body)
	decodeErr := jsonDecoder.Decode(&apiRet)
	if decodeErr != nil {
		err = fmt.Errorf("parse body error, %s", decodeErr.Error())
		return
	}
	apiRet.LogId = resp.Header.Get(XHeaderLogId)
	// check status code
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("call api error, %s %s", apiRet.Code, apiRet.Message)
		return
	}
	// marshal data to bytes for later use
	apiRet.Data, err = json.Marshal(apiRet.Data)
	if err != nil {
		err = fmt.Errorf("encode data failed, %s", err.Error())
		return
	}
	return
}

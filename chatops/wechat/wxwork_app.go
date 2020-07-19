package wechat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

// WxWorkAppMessageAPI is the api to get the app access token
const WxWorkAppTokenAPI = "https://qyapi.weixin.qq.com/cgi-bin/gettoken"

// WxWorkAppMessageAPI is the api to send the app messages
const WxWorkAppMessageAPI = "https://qyapi.weixin.qq.com/cgi-bin/webhook/send"

// WxWorkAppTimeout is the wxwork app default timeout
const WxWorkAppTimeout = time.Second * 30
const WxWorkAppStatusOK = 0

type WxWorkAppTokenResp struct {
	ErrCode     int    `json:"errcode"`
	ErrMessage  string `json:"errmsg"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

type WxWorkAppMessageResp struct {
	ErrCode    int    `json:"errcode"`
	ErrMessage string `json:"errmsg"`
}

type WxWorkApp struct {
	agentID          string
	corpID           string // see doc https://work.weixin.qq.com/api/doc/90000/90135/91039
	corpSecret       string // see doc https://work.weixin.qq.com/api/doc/90000/90135/90665#secret
	client           *http.Client
	tokenRefreshLock sync.RWMutex // lock to refresh the access token which can expire in a period of time
	accessToken      string       // cached access token
	expiredAt        time.Time    // token expire time
}

func (r *WxWorkApp) IsAccessTokenExpired() bool {
	return time.Now().After(r.expiredAt)
}

// NewWxWorkApp create a new wxwork app
func NewWxWorkApp(corpID, corpSecret, agentID string) *WxWorkApp {
	return NewWxWorkAppWithTimeout(corpID, corpSecret, agentID, WxWorkAppTimeout)
}

// NewWxWorkAppWithTimeout create a new wxwork app with timeout
func NewWxWorkAppWithTimeout(corpID, corpSecret, agentID string, timeout time.Duration) *WxWorkApp {
	client := http.Client{}
	client.Timeout = timeout
	return &WxWorkApp{corpID: corpID, corpSecret: corpSecret, agentID: agentID, client: &client, tokenRefreshLock: sync.RWMutex{}}
}

// NewWxWorkAppWithClient create a new wxwork app with http.Client
func NewWxWorkAppWithClient(corpID, corpSecret, agentID string, client *http.Client) *WxWorkApp {
	return &WxWorkApp{corpID: corpID, corpSecret: corpSecret, agentID: agentID, client: client, tokenRefreshLock: sync.RWMutex{}}
}
func (r *WxWorkApp) SendTextMessage() (err error) {

}

func (r *WxWorkApp) refreshAccessToken() (err error) {
	reqURL := fmt.Sprintf("%s?corpid=%s&corpsecret=%s", WxWorkAppTokenAPI, r.corpID, r.corpSecret)
	req, newErr := http.NewRequest(http.MethodGet, reqURL, nil)
	if newErr != nil {
		err = fmt.Errorf("create request error, %s", newErr.Error())
		return
	}
	resp, getErr := r.client.Do(req)
	if getErr != nil {
		err = fmt.Errorf("get response error, %s", getErr.Error())
		io.Copy(ioutil.Discard, resp.Body)
		return
	}
	defer resp.Body.Close()
	// check http code
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("wxwork request error, %s", resp.Status)
		io.Copy(ioutil.Discard, resp.Body)
		return
	}
	// parse response body
	decoder := json.NewDecoder(resp.Body)
	var wxTokenResp WxWorkAppTokenResp
	if decodeErr := decoder.Decode(&wxTokenResp); decodeErr != nil {
		err = fmt.Errorf("parse response error, %s", decodeErr.Error())
		return
	}
	if wxTokenResp.ErrCode != WxWorkAppStatusOK {
		err = fmt.Errorf("call wxwork app api error, %d %s", wxTokenResp.ErrCode, wxTokenResp.ErrMessage)
		return
	}
	// set access token and expired at
	r.accessToken = wxTokenResp.AccessToken
	r.expiredAt = time.Now().Add(time.Second * time.Duration(wxTokenResp.ExpiresIn))
	return
}

func (r *WxWorkApp) sendMessage(messageObj interface{}) (err error) {
	// check the token expired or not
	if r.accessToken == "" || r.IsAccessTokenExpired() {
		r.tokenRefreshLock.Lock()
		if r.accessToken == "" || r.IsAccessTokenExpired() {
			err = r.refreshAccessToken()
		}
		r.tokenRefreshLock.Unlock()
		if err != nil {
			err = fmt.Errorf("refresh access token error, %s", err.Error())
			return
		}
	}
	reqURL := fmt.Sprintf("%s?access_token=%s", WxWorkAppMessageAPI, r.accessToken)
	reqBody, _ := json.Marshal(messageObj)

	req, newErr := http.NewRequest(http.MethodPost, reqURL, bytes.NewReader(reqBody))
	if newErr != nil {
		err = fmt.Errorf("create request error, %s", newErr.Error())
		return
	}
	req.Header.Add("Content-Type", "application/json")
	resp, getErr := r.client.Do(req)
	if getErr != nil {
		err = fmt.Errorf("get response error, %s", getErr.Error())
		io.Copy(ioutil.Discard, resp.Body)
		return
	}
	defer resp.Body.Close()
	// check http code
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("wxwork request error, %s", resp.Status)
		io.Copy(ioutil.Discard, resp.Body)
		return
	}
	// parse response body
	decoder := json.NewDecoder(resp.Body)
	var wxMessageResp WxWorkAppMessageResp
	if decodeErr := decoder.Decode(&wxMessageResp); decodeErr != nil {
		err = fmt.Errorf("parse response error, %s", decodeErr.Error())
		return
	}
	if wxMessageResp.ErrCode != WxWorkAppStatusOK {
		err = fmt.Errorf("call wxwork app api error, %d %s", wxMessageResp.ErrCode, wxMessageResp.ErrMessage)
		return
	}
	return
}

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
const WxWorkAppMessageAPI = "https://qyapi.weixin.qq.com/cgi-bin/appchat/send"

// WxWorkAppCreateGroupAPI is the api to create the wxwork group
const WxWorkAppCreateGroupAPI = "https://qyapi.weixin.qq.com/cgi-bin/appchat/create"

// WxWorkAppTimeout is the wxwork app default timeout
const WxWorkAppTimeout = time.Second * 30
const WxWorkAppStatusOK = 0

const (
	WxWorkAppMessageTypeText = "text"
)

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

type WxWorkAppCreateGroupResp struct {
	ErrCode    int    `json:"errcode"`
	ErrMessage string `json:"errmsg"`
	ChatID     string `json:"chatid"`
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

func (r *WxWorkApp) CreateGroupChat(name, chatID, ownerID string, userIDList []string) (newChatID string, err error) {
	createGroupReqObject := make(map[string]interface{})
	createGroupReqObject["name"] = name
	createGroupReqObject["chatid"] = chatID
	createGroupReqObject["owner"] = ownerID
	createGroupReqObject["userlist"] = userIDList
	var createGroupResp WxWorkAppCreateGroupResp
	err = r.fireRequest(WxWorkAppCreateGroupAPI, &createGroupReqObject, &createGroupResp)
	if err != nil {
		return
	}
	if createGroupResp.ErrCode != WxWorkAppStatusOK {
		err = fmt.Errorf("call wxwork app create group api error, %d %s", createGroupResp.ErrCode, createGroupResp.ErrMessage)
		return
	}
	newChatID = createGroupResp.ChatID
	return
}

func (r *WxWorkApp) SendGroupTextMessage(chatID, content string, safe bool) (err error) {
	messageObj := make(map[string]interface{})
	messageObj["chatid"] = chatID
	messageObj["msgtype"] = WxWorkAppMessageTypeText
	messageObj["text"] = map[string]string{
		"content": content,
	}
	if safe {
		messageObj["safe"] = 1
	} else {
		messageObj["safe"] = 0
	}
	return r.sendMessage(&messageObj)
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
	var messageResp WxWorkAppMessageResp
	err = r.fireRequest(WxWorkAppMessageAPI, messageObj, &messageResp)
	if err != nil {
		return
	}
	if messageResp.ErrCode != WxWorkAppStatusOK {
		err = fmt.Errorf("call wxwork app message api error, %d %s", messageResp.ErrCode, messageResp.ErrMessage)
		return
	}
	return
}

func (r *WxWorkApp) fireRequest(reqURL string, reqBodyObject interface{}, respObject interface{}) (err error) {
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
	reqURL = fmt.Sprintf("%s?access_token=%s", reqURL, r.accessToken)
	reqBody, _ := json.Marshal(reqBodyObject)

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
	if decodeErr := decoder.Decode(respObject); decodeErr != nil {
		err = fmt.Errorf("parse response error, %s", decodeErr.Error())
		return
	}
	return
}

package wechat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// WxWorkAppMessageAPI is the api to get the app access token
const WxWorkAppTokenAPI = "https://qyapi.weixin.qq.com/cgi-bin/gettoken"

// WxWorkAppMessageAPI is the api to send the app messages
const WxWorkAppMessageAPI = "https://qyapi.weixin.qq.com/cgi-bin/appchat/send"

// WxWorkAppCreateGroupAPI is the api to create the wxwork group
const WxWorkAppCreateGroupAPI = "https://qyapi.weixin.qq.com/cgi-bin/appchat/create"

// WxWorkAppUpdateGroupAPI is the api to update the wxwork group
const WxWorkAppUpdateGroupAPI = "https://qyapi.weixin.qq.com/cgi-bin/appchat/update"

// WxWorkAppGetGroupAPI is the api to get the wxwork group
const WxWorkAppGetGroupAPI = "https://qyapi.weixin.qq.com/cgi-bin/appchat/get"

// WxWorkAppTimeout is the wxwork app default timeout
const WxWorkAppTimeout = time.Second * 30
const WxWorkAppStatusOK = 0

const (
	WxWorkAppMessageTypeText     = "text"
	WxWorkAppMessageTypeImage    = "image"
	WxWorkAppMessageTypeVoice    = "voice"
	WxWorkAppMessageTypeVideo    = "video"
	WxWorkAppMessageTypeFile     = "file"
	WxWorkAppMessageTypeTextCard = "textcard"
	WxWorkAppMessageTypeNews     = "news"
	WxWorkAppMessageTypeMpNews   = "mpnews"
	WxWorkAppMessageTypeMarkdown = "markdown"
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

type WxWorkAppUpdateGroupResp struct {
	ErrCode    int    `json:"errcode"`
	ErrMessage string `json:"errmsg"`
}

type WxWorkAppGetGroupResp struct {
	ErrCode    int            `json:"errcode"`
	ErrMessage string         `json:"errmsg"`
	ChatInfo   WxWorkAppGroup `json:"chat_info"`
}

type WxWorkAppGroup struct {
	ChatID   string   `json:"chatid"`
	Name     string   `json:"name"`
	Owner    string   `json:"owner"`
	UserList []string `json:"userlist"`
}

type WxWorkAppNewsMessageArticle struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	PictureURL  string `json:"picurl"`
}

type WxWorkAppMpNewsMessageArticle struct {
	Title            string `json:"title"`
	ThumbMediaID     string `json:"thumb_media_id"`
	Author           string `json:"author"`
	ContentSourceURL string `json:"content_source_url"`
	Content          string `json:"content"`
	Digest           string `json:"digest"`
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
	err = r.fireRequest(http.MethodPost, WxWorkAppCreateGroupAPI, nil, &createGroupReqObject, &createGroupResp)
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

func (r *WxWorkApp) UpdateGroupChat(name, chatID, ownerID string, addUserList []string, delUserList []string) (err error) {
	updateGroupReqObject := make(map[string]interface{})
	updateGroupReqObject["name"] = name
	updateGroupReqObject["chatid"] = chatID
	updateGroupReqObject["owner"] = ownerID
	updateGroupReqObject["add_user_list"] = addUserList
	updateGroupReqObject["del_user_list"] = delUserList
	var updateGroupResp WxWorkAppUpdateGroupResp
	err = r.fireRequest(http.MethodPost, WxWorkAppUpdateGroupAPI, nil, &updateGroupReqObject, &updateGroupResp)
	if err != nil {
		return
	}
	if updateGroupResp.ErrCode != WxWorkAppStatusOK {
		err = fmt.Errorf("call wxwork app update group api error, %d %s", updateGroupResp.ErrCode, updateGroupResp.ErrMessage)
		return
	}
	return
}

func (r *WxWorkApp) GetGroupChat(chatID string) (group WxWorkAppGroup, err error) {
	var getGroupResp WxWorkAppGetGroupResp
	err = r.fireRequest(http.MethodGet, WxWorkAppGetGroupAPI, map[string]string{"chatid": chatID}, nil, &getGroupResp)
	if err != nil {
		return
	}
	if getGroupResp.ErrCode != WxWorkAppStatusOK {
		err = fmt.Errorf("call wxwork app get group api error, %d %s", getGroupResp.ErrCode, getGroupResp.ErrMessage)
		return
	}
	group = getGroupResp.ChatInfo
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

func (r *WxWorkApp) SendGroupMarkdownMessage(chatID, content string, safe bool) (err error) {
	messageObj := make(map[string]interface{})
	messageObj["chatid"] = chatID
	messageObj["msgtype"] = WxWorkAppMessageTypeMarkdown
	messageObj["markdown"] = map[string]string{
		"content": content,
	}
	if safe {
		messageObj["safe"] = 1
	} else {
		messageObj["safe"] = 0
	}
	return r.sendMessage(&messageObj)
}

func (r *WxWorkApp) SendGroupImageMessage(chatID, mediaID string, safe bool) (err error) {
	messageObj := make(map[string]interface{})
	messageObj["chatid"] = chatID
	messageObj["msgtype"] = WxWorkAppMessageTypeImage
	messageObj["image"] = map[string]string{
		"media_id": mediaID,
	}
	if safe {
		messageObj["safe"] = 1
	} else {
		messageObj["safe"] = 0
	}
	return r.sendMessage(&messageObj)
}

func (r *WxWorkApp) SendGroupVoiceMessage(chatID, mediaID string, safe bool) (err error) {
	messageObj := make(map[string]interface{})
	messageObj["chatid"] = chatID
	messageObj["msgtype"] = WxWorkAppMessageTypeVoice
	messageObj["voice"] = map[string]string{
		"media_id": mediaID,
	}
	if safe {
		messageObj["safe"] = 1
	} else {
		messageObj["safe"] = 0
	}
	return r.sendMessage(&messageObj)
}

func (r *WxWorkApp) SendGroupVideoMessage(chatID, mediaID, mediaTitle, mediaDescription string, safe bool) (err error) {
	messageObj := make(map[string]interface{})
	messageObj["chatid"] = chatID
	messageObj["msgtype"] = WxWorkAppMessageTypeVideo
	messageObj["video"] = map[string]string{
		"media_id":    mediaID,
		"title":       mediaTitle,
		"description": mediaDescription,
	}
	if safe {
		messageObj["safe"] = 1
	} else {
		messageObj["safe"] = 0
	}
	return r.sendMessage(&messageObj)
}

func (r *WxWorkApp) SendGroupFileMessage(chatID, mediaID string, safe bool) (err error) {
	messageObj := make(map[string]interface{})
	messageObj["chatid"] = chatID
	messageObj["msgtype"] = WxWorkAppMessageTypeFile
	messageObj["file"] = map[string]string{
		"media_id": mediaID,
	}
	if safe {
		messageObj["safe"] = 1
	} else {
		messageObj["safe"] = 0
	}
	return r.sendMessage(&messageObj)
}

func (r *WxWorkApp) SendGroupTextCardMessage(chatID, title, description, url, btnText string, safe bool) (err error) {
	messageObj := make(map[string]interface{})
	messageObj["chatid"] = chatID
	messageObj["msgtype"] = WxWorkAppMessageTypeTextCard
	messageObj["textcard"] = map[string]string{
		"title":       title,
		"description": description,
		"url":         url,
		"btntext":     btnText,
	}
	if safe {
		messageObj["safe"] = 1
	} else {
		messageObj["safe"] = 0
	}
	return r.sendMessage(&messageObj)
}

func (r *WxWorkApp) SendGroupNewsMessage(chatID, articles []WxWorkAppNewsMessageArticle, safe bool) (err error) {
	messageObj := make(map[string]interface{})
	messageObj["chatid"] = chatID
	messageObj["msgtype"] = WxWorkAppMessageTypeNews
	messageObj["news"] = map[string]interface{}{
		"articles": articles,
	}
	if safe {
		messageObj["safe"] = 1
	} else {
		messageObj["safe"] = 0
	}
	return r.sendMessage(&messageObj)
}

func (r *WxWorkApp) SendGroupMpNewsMessage(chatID, articles []WxWorkAppMpNewsMessageArticle, safe bool) (err error) {
	messageObj := make(map[string]interface{})
	messageObj["chatid"] = chatID
	messageObj["msgtype"] = WxWorkAppMessageTypeMpNews
	messageObj["mpnews"] = map[string]interface{}{
		"articles": articles,
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

// See doc https://work.weixin.qq.com/api/doc/90000/90135/90248
func (r *WxWorkApp) sendMessage(messageObj interface{}) (err error) {
	var messageResp WxWorkAppMessageResp
	err = r.fireRequest(http.MethodPost, WxWorkAppMessageAPI, nil, messageObj, &messageResp)
	if err != nil {
		return
	}
	if messageResp.ErrCode != WxWorkAppStatusOK {
		err = fmt.Errorf("call wxwork app message api error, %d %s", messageResp.ErrCode, messageResp.ErrMessage)
		return
	}
	return
}

func (r *WxWorkApp) fireRequest(reqMethod, reqURL string, reqParams map[string]string, reqBodyObject interface{}, respObject interface{}) (err error) {
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

	queryString := url.Values{}
	queryString.Add("access_token", r.accessToken)
	if reqParams != nil {
		for k, v := range reqParams {
			queryString.Add(k, v)
		}
	}

	reqURL = fmt.Sprintf("%s?%s", reqURL, queryString.Encode())
	var reqBodyReader io.Reader
	if reqBodyObject != nil {
		reqBody, _ := json.Marshal(reqBodyObject)
		reqBodyReader = bytes.NewReader(reqBody)
	}

	req, newErr := http.NewRequest(reqMethod, reqURL, reqBodyReader)
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
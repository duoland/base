package bytedance

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

// FeiShuAppTenantAccessTokenAPI is the api to get access token
const FeiShuAppTenantAccessTokenAPI = "https://open.feishu.cn/open-apis/auth/v3/tenant_access_token/internal/"

// FeiShuAppSendMessageAPI is the api to send message to user or group
const FeiShuAppSendMessageAPI = "https://open.feishu.cn/open-apis/message/v4/send/"

// FeiShuAppCreateGroupAPI is the api to create group
const FeiShuAppCreateGroupAPI = "https://open.feishu.cn/open-apis/chat/v4/create/"

// FeiShuAppTimeout is the default timeout to api call
const FeiShuAppTimeout = time.Second * 10

// FeiShuAppStatusOK is the ok status of api call
const FeiShuAppStatusOK = 0

// See doc https://open.feishu.cn/document/ukTMukTMukTM/ugjM14COyUjL4ITN
const FeishuCodeAccessTokenExpired = 99991663

type FeiShuAppGetTokenResp struct {
	Code              int    `json:"code"`
	Message           string `json:"msg"`
	TenantAccessToken string `json:"tenant_access_token"`
	Expire            int    `json:"expire"`
}

type FeiShuAppMessageSendTarget struct {
	OpenID string `json:"open_id,omitempty"`
	UserID string `json:"user_id,omitempty"`
	Email  string `json:"email,omitempty"`
	ChatID string `json:"chat_id,omitempty"`
}

const (
	FeiShuAppPostMessageText  = "text"
	FeiShuAppPostMessageHref  = "a"
	FeiShuAppPostMessageAt    = "at"
	FeiShuAppPostMessageImage = "img"
)

type FeishuAppPostMessageContent struct {
	Title   string                              `json:"title"`
	Content [][]FeishuAppPostMessageContentItem `json:"content"`
}

type FeishuAppPostMessageContentItem struct {
	Tag      string `json:"tag"`
	UnEscape bool   `json:"un_escape,omitempty"`
	Text     string `json:"text,omitempty"`
	Lines    int    `json:"lines,omitempty"`
	Href     string `json:"href,omitempty"`
	UserID   string `json:"user_id,omitempty"`
	ImageKey string `json:"image_key,omitempty"`
	Height   int    `json:"height,omitempty"`
	Width    int    `json:"width,omitempty"`
}

const (
	FeiShuAppMessageTypeText  = "text"
	FeiShuAppMessageTypeImage = "image"
	FeiShuAppMessageTypePost  = "post"
)

const (
	FeiShuAppI18nChinese  = "zh_cn"
	FeiShuAppI18nJapanese = "ja_jp"
	FeiShuAppI18nEnglish  = "en_us"
)

// See doc https://open.feishu.cn/document/ukTMukTMukTM/uMDMxEjLzATMx4yMwETM
type FeiShuAppMessageSendReq struct {
	OpenID      string      `json:"open_id,omitempty"`
	UserID      string      `json:"user_id,omitempty"`
	Email       string      `json:"email,omitempty"`
	ChatID      string      `json:"chat_id,omitempty"`
	AppID       string      `json:"robot_id,omitempty"`
	MessageType string      `json:"msg_type"`
	Content     interface{} `json:"content"`
}

type FeiShuAppMessageSendResp struct {
	Code    int                          `json:"code"`
	Message string                       `json:"msg"`
	Data    FeiShuAppMessageSendRespData `json:"data"`
}

type FeiShuAppMessageSendRespData struct {
	MessageID string `json:"message_id"`
}

type FeiShuAppCreateGroupOptions struct {
	OpenIDs        []string          `json:"open_ids"`
	I18nNames      map[string]string `json:"i18n_names"`
	OnlyOwnerAdd   bool              `json:"only_owner_add"`
	OnlyOwnerAtAll bool              `json:"only_owner_at_all"`
	OnlyOwnerEdit  bool              `json:"only_owner_edit"`
	ShareAllowed   bool              `json:"share_allowed"`
}

type FeiShuAppCreateGroupResp struct {
	Code    int                          `json:"code"`
	Message string                       `json:"msg"`
	Data    FeiShuAppCreateGroupRespData `json:"data"`
}

type FeiShuAppCreateGroupRespData struct {
	ChatID         string   `json:"chat_id"`
	InvalidOpenIDs []string `json:"invalid_open_ids"`
	InvalidUserIDs []string `json:"invalid_user_ids"`
}

// See doc https://open.feishu.cn/document/ukTMukTMukTM/uIjNz4iM2MjLyYzM

type FeiShuApp struct {
	appID            string
	appSecret        string
	client           *http.Client
	tokenRefreshLock sync.RWMutex // lock to refresh the access token which can expire in a period of time
	accessToken      string       // cached access token
	expiredAt        time.Time    // token expire time
}

func (r *FeiShuApp) IsAccessTokenExpired() bool {
	return time.Now().After(r.expiredAt)
}

// NewFeiShuApp create a new feishu app
func NewFeiShuApp(appID, appSecret string) *FeiShuApp {
	return NewFeiShuAppWithTimeout(appID, appSecret, FeiShuAppTimeout)
}

// NewFeiShuAppWithTimeout create a new feishu app with timeout
func NewFeiShuAppWithTimeout(appID, appSecret string, timeout time.Duration) *FeiShuApp {
	client := http.Client{}
	client.Timeout = timeout
	return &FeiShuApp{appID: appID, appSecret: appSecret, client: &client, tokenRefreshLock: sync.RWMutex{}}
}

// NewFeiShuAppWithClient create a new feishu app with http.Client
func NewFeiShuAppWithClient(appID, appSecret string, client *http.Client) *FeiShuApp {
	return &FeiShuApp{appID: appID, appSecret: appSecret, client: client, tokenRefreshLock: sync.RWMutex{}}
}

func (r *FeiShuApp) CreateGroupChat(name, description string, userIDList []string, options *FeiShuAppCreateGroupOptions) (newChatID string, err error) {
	createGroupReqObject := make(map[string]interface{})
	createGroupReqObject["name"] = name
	createGroupReqObject["description"] = description
	createGroupReqObject["user_ids"] = userIDList
	if options != nil {
		createGroupReqObject["open_ids"] = options.OpenIDs
		createGroupReqObject["i18n_names"] = options.I18nNames
		createGroupReqObject["only_owner_add"] = options.OnlyOwnerAdd
		createGroupReqObject["only_owner_at_all"] = options.OnlyOwnerAtAll
		createGroupReqObject["only_owner_edit"] = options.OnlyOwnerEdit
	}
	var createGroupResp FeiShuAppCreateGroupResp
	err = r.fireRequest(http.MethodPost, FeiShuAppCreateGroupAPI, &createGroupReqObject, &createGroupResp)
	if err != nil {
		return
	}
	if createGroupResp.Code != FeiShuAppStatusOK {
		if createGroupResp.Code == FeishuCodeAccessTokenExpired {
			// reset the access token
			r.accessToken = ""
		}
		err = fmt.Errorf("call feishu app create group api error, %d %s", createGroupResp.Code, createGroupResp.Message)
		return
	}
	newChatID = createGroupResp.Data.ChatID
	return
}

// See doc https://open.feishu.cn/document/ukTMukTMukTM/uUjNz4SN2MjL1YzM
// options only support root_id now if you want to replay a specified message
func (r *FeiShuApp) SendTextMessage(target *FeiShuAppMessageSendTarget, content string, options map[string]string) (messageResp FeiShuAppMessageSendResp, err error) {
	messageReq := FeiShuAppMessageSendReq{
		OpenID: target.OpenID,
		UserID: target.UserID,
		Email:  target.Email,
		ChatID: target.ChatID,
	}
	if options != nil {
		if v, exists := options["robot_id"]; exists {
			messageReq.AppID = v
		}
	}
	messageReq.MessageType = FeiShuAppMessageTypeText
	messageReq.Content = map[string]string{"text": content}
	return r.sendMessage(&messageReq)
}

func (r *FeiShuApp) SendImageMessage(target *FeiShuAppMessageSendTarget, imageKey string, options map[string]string) (messageResp FeiShuAppMessageSendResp, err error) {
	messageReq := FeiShuAppMessageSendReq{
		OpenID: target.OpenID,
		UserID: target.UserID,
		Email:  target.Email,
		ChatID: target.ChatID,
	}
	if options != nil {
		if v, exists := options["robot_id"]; exists {
			messageReq.AppID = v
		}
	}
	messageReq.MessageType = FeiShuAppMessageTypeImage
	messageReq.Content = map[string]string{"image_key": imageKey}
	return r.sendMessage(&messageReq)
}

func (r *FeiShuApp) SendPostMessage(target *FeiShuAppMessageSendTarget, title, i18nKey string, messageLines [][]FeishuAppPostMessageContentItem,
	options map[string]string) (messageResp FeiShuAppMessageSendResp, err error) {
	messageReq := FeiShuAppMessageSendReq{
		OpenID: target.OpenID,
		UserID: target.UserID,
		Email:  target.Email,
		ChatID: target.ChatID,
	}
	if options != nil {
		if v, exists := options["robot_id"]; exists {
			messageReq.AppID = v
		}
	}
	messageReq.MessageType = FeiShuAppMessageTypePost
	messageReq.Content = map[string]interface{}{
		"post": map[string]FeishuAppPostMessageContent{
			i18nKey: {
				Title:   title,
				Content: messageLines,
			},
		},
	}
	return r.sendMessage(&messageReq)
}

func (r *FeiShuApp) refreshAccessToken() (err error) {
	reqBody := map[string]string{
		"app_id":     r.appID,
		"app_secret": r.appSecret,
	}
	reqBodyBytes, _ := json.Marshal(&reqBody)
	req, newErr := http.NewRequest(http.MethodPost, FeiShuAppTenantAccessTokenAPI, bytes.NewReader(reqBodyBytes))
	if newErr != nil {
		err = fmt.Errorf("create request error, %s", newErr.Error())
		return
	}
	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	resp, getErr := r.client.Do(req)
	if getErr != nil {
		err = fmt.Errorf("get response error, %s", getErr.Error())
		return
	}
	defer resp.Body.Close()

	// check http code
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("feishu request error, %s", resp.Status)
		io.Copy(ioutil.Discard, resp.Body)
		return
	}
	// parse response body
	decoder := json.NewDecoder(resp.Body)
	var accessTokenResp FeiShuAppGetTokenResp
	if decodeErr := decoder.Decode(&accessTokenResp); decodeErr != nil {
		err = fmt.Errorf("parse response error, %s", decodeErr.Error())
		return
	}
	if accessTokenResp.Code != FeiShuAppStatusOK {
		err = fmt.Errorf("call feishu app api error, %d %s", accessTokenResp.Code, accessTokenResp.Message)
		return
	}
	// set access token and expired at
	r.accessToken = accessTokenResp.TenantAccessToken
	r.expiredAt = time.Now().Add(time.Second * time.Duration(accessTokenResp.Expire))
	return
}

func (r *FeiShuApp) fireRequest(reqMethod, reqURL string, reqBodyObject interface{}, respObject interface{}) (err error) {
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
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", r.accessToken))
	req.Header.Add("Content-Type", "application/json")
	resp, getErr := r.client.Do(req)
	if getErr != nil {
		err = fmt.Errorf("get response error, %s", getErr.Error())
		return
	}
	defer resp.Body.Close()
	// check http code
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("feishu request error, %s", resp.Status)
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

func (r *FeiShuApp) sendMessage(messageObj interface{}) (messageResp FeiShuAppMessageSendResp, err error) {
	err = r.fireRequest(http.MethodPost, FeiShuAppSendMessageAPI, messageObj, &messageResp)
	if err != nil {
		return
	}
	if messageResp.Code != FeiShuAppStatusOK {
		if messageResp.Code == FeishuCodeAccessTokenExpired {
			// reset the access token
			r.accessToken = ""
		}
		err = fmt.Errorf("call feishu app message api error, %d %s", messageResp.Code, messageResp.Message)
		return
	}
	return
}

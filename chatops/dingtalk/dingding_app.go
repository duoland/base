package dingtalk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

// DingDingAppTokenAPI is the api to get app access token
const DingDingAppTokenAPI = "https://oapi.dingtalk.com/gettoken"

// DingDingAppCreateGroupAPI is the api to create dingding group
const DingDingAppCreateGroupAPI = "https://oapi.dingtalk.com/chat/create"

// DingDingAppUpdateGroupAPI is the api to update dingding group
const DingDingAppUpdateGroupAPI = "https://oapi.dingtalk.com/chat/update"

// DingDingAppGetGroupAPI is the api to get dingding group
const DingDingAppGetGroupAPI = "https://oapi.dingtalk.com/chat/get"

// DingDingAppSendMessageAPI is the api to send message to end users
const DingDingAppSendMessageAPI = "https://oapi.dingtalk.com/topapi/message/corpconversation/asyncsend_v2"

// DingDingAppGetMessageSendProgressAPI is the api to get the message send progress
const DingDingAppGetMessageSendProgressAPI = "https://oapi.dingtalk.com/topapi/message/corpconversation/getsendprogress"

// DingDingAppGetMessageSendResultAPI is the api to get the message send result
const DingDingAppGetMessageSendResultAPI = "https://oapi.dingtalk.com/topapi/message/corpconversation/getsendresult"

// DingDingAppRecallMessageAPI is the api to recall the message
const DingDingAppRecallMessageAPI = "https://oapi.dingtalk.com/topapi/message/corpconversation/recall"

// DingDingAppSendGroupMessageAPI is the api to send message to group
const DingDingAppSendGroupMessageAPI = "https://oapi.dingtalk.com/chat/send"

// DingDingAppTimeout is the dingding app default timeout
const DingDingAppTimeout = time.Second * 10

// DingDingAppStatusOK is the ok status of api call
const DingDingAppStatusOK = 0

// See doc https://ding-doc.dingtalk.com/doc#/faquestions/rftpfg
const DingDingCodeAccessTokenExpired = 42001

const (
	DingDingOptionYes = 1
	DingDingOptionNo  = 0
)

const (
	DingDingAppMessageTypeText       = "text"
	DingDingAppMessageTypeImage      = "image"
	DingDingAppMessageTypeVoice      = "voice"
	DingDingAppMessageTypeFile       = "file"
	DingDingAppMessageTypeLink       = "link"
	DingDingAppMessageTypeOA         = "oa"
	DingDingAppMessageTypeMarkdown   = "markdown"
	DingDingAppMessageTypeActionCard = "action_card"
)

type DingDingAppTokenResp struct {
	ErrCode     int    `json:"errcode"`
	ErrMessage  string `json:"errmsg"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

type DingDingAppCreateGroupResp struct {
	ErrCode         int    `json:"errcode"`
	ErrMessage      string `json:"errmsg"`
	ChatID          string `json:"chatid"`
	ConversationTag int    `json:"conversationTag"`
}

type DingDingAppCreateGroupOptions struct {
	ShowHistoryType     int
	Searchable          int
	ValidationType      int
	MentionAllAuthority int
	ChatBannedType      int
	ManagementType      int
}

type DingDingAppUpdateGroupOptions struct {
	Name        string
	Owner       string
	Icon        string
	AddUserList []string
	DelUserList []string
	DingDingAppCreateGroupOptions
}

type DingDingAppUpdateGroupResp struct {
	ErrCode    int    `json:"errcode"`
	ErrMessage string `json:"errmsg"`
}

type DingDingAppGetGroupResp struct {
	ErrCode    int              `json:"errcode"`
	ErrMessage string           `json:"errmsg"`
	ChatInfo   DingDingAppGroup `json:"chat_info"`
}

type DingDingAppGroup struct {
	ChatID              string   `json:"chatid"`
	Name                string   `json:"name"`
	Owner               string   `json:"owner"`
	UserIDList          []string `json:"useridlist"`
	Icon                string   `json:"icon"`
	ShowHistoryType     int      `json:"showHistoryType"`
	Searchable          int      `json:"searchable"`
	ValidationType      int      `json:"validationType"`
	MentionAllAuthority int      `json:"mentionAllAuthority"`
	ChatBannedType      int      `json:"chatBannedType"`
	ManagementType      int      `json:"managementType"`
}

type DingDingAppMessageSendResp struct {
	ErrCode    int    `json:"errcode"`
	ErrMessage string `json:"errmsg"`
	TaskID     int    `json:"task_id"`
}

type DingDingAppMessageSendProgressResp struct {
	ErrCode    int                            `json:"errcode"`
	ErrMessage string                         `json:"errmsg"`
	Progress   DingDingAppMessageSendProgress `json:"progress"`
}

type DingDingAppMessageSendProgress struct {
	ProgressInPercent int `json:"progress_in_percent"`
	Status            int `json:"status"`
}

type DingDingAppMessageSendResultResp struct {
	ErrCode    int                          `json:"errcode"`
	ErrMessage string                       `json:"errmsg"`
	SendResult DingDingAppMessageSendResult `json:"send_result"`
}

type DingDingAppMessageSendResult struct {
	InvalidUserIDList   []string `json:"invalid_user_id_list"`
	ForbiddenUserIDList []string `json:"forbidden_user_id_list"`
	FailedUserIDList    []string `json:"failed_user_id_list"`
	ReadUserIDList      []string `json:"read_user_id_list"`
	UnReadUserIDList    []string `json:"unread_user_id_list"`
	InvalidDeptIDList   []string `json:"invalid_dept_id_list"`
}

type DingDingAppMessageRecallResp struct {
	ErrCode    int    `json:"errcode"`
	ErrMessage string `json:"errmsg"`
}

type DingDingAppLinkMessage struct {
	Text       string `json:"text"`
	Title      string `json:"title"`
	PicURL     string `json:"picUrl"`
	MessageURL string `json:"messageUrl"`
}

type DingDingAppActionCardMessage struct {
	Title             string                        `json:"title"`
	Markdown          string                        `json:"markdown"`
	SingleTitle       string                        `json:"single_title,omitempty"` // single jump action card fields
	SingleURL         string                        `json:"single_url,omitempty"`
	ButtonOrientation string                        `json:"btn_orientation"` // standalone jump action card fields
	Buttons           []DingDingAppActionCardButton `json:"btn_json_list,omitempty"`
}

type DingDingAppActionCardButton struct {
	Title     string `json:"title"`
	ActionURL string `json:"action_url"`
}

type DingDingAppGroupMessageSendResp struct {
	ErrCode    int    `json:"errcode"`
	ErrMessage string `json:"errmsg"`
	MessageID  string `json:"messageId"`
}

type DingDingApp struct {
	agentID          string
	appKey           string
	appSecret        string
	client           *http.Client
	tokenRefreshLock sync.RWMutex
	accessToken      string
	expiredAt        time.Time
}

func (r *DingDingApp) IsAccessTokenExpired() bool {
	return time.Now().After(r.expiredAt)
}

// NewDingDingApp create a new dingding app
func NewDingDingApp(appKey, appSecret, agentID string) *DingDingApp {
	return NewDingDingAppWithTimeout(appKey, appSecret, agentID, DingDingAppTimeout)
}

// NewDingDingAppWithTimeout create a new dingding app with timeout
func NewDingDingAppWithTimeout(appKey, appSecret, agentID string, timeout time.Duration) *DingDingApp {
	client := http.Client{}
	client.Timeout = timeout
	return &DingDingApp{appKey: appKey, appSecret: appSecret, agentID: agentID, client: &client, tokenRefreshLock: sync.RWMutex{}}
}

// NewDingDingAppWithClient create a new dingding app with http.Client
func NewDingDingAppWithClient(appKey, appSecret, agentID string, client *http.Client) *DingDingApp {
	return &DingDingApp{appKey: appKey, appSecret: appSecret, agentID: agentID, client: client, tokenRefreshLock: sync.RWMutex{}}
}

func (r *DingDingApp) SendTextMessage(userIDList []string, departmentIDList []string, toAllUser bool, content string) (
	resp DingDingAppMessageSendResp, err error) {
	messageObj := make(map[string]interface{})
	messageObj["agent_id"] = r.agentID
	if len(userIDList) > 0 {
		messageObj["userid_list"] = strings.Join(userIDList, ",")
	}
	if len(departmentIDList) > 0 {
		messageObj["dept_id_list"] = strings.Join(departmentIDList, ",")
	}
	messageObj["to_all_user"] = toAllUser
	messageObj["msg"] = map[string]interface{}{
		"msgtype": DingDingAppMessageTypeText,
		"text":    map[string]string{"content": content},
	}
	return r.sendMessage(&messageObj)
}

func (r *DingDingApp) SendMarkdownMessage(userIDList []string, departmentIDList []string, toAllUser bool, title, content string) (
	resp DingDingAppMessageSendResp, err error) {
	messageObj := make(map[string]interface{})
	messageObj["agent_id"] = r.agentID
	if len(userIDList) > 0 {
		messageObj["userid_list"] = strings.Join(userIDList, ",")
	}
	if len(departmentIDList) > 0 {
		messageObj["dept_id_list"] = strings.Join(departmentIDList, ",")
	}
	messageObj["to_all_user"] = toAllUser
	messageObj["msg"] = map[string]interface{}{
		"msgtype":  DingDingAppMessageTypeMarkdown,
		"markdown": map[string]string{"title": title, "text": content},
	}
	return r.sendMessage(&messageObj)
}

func (r *DingDingApp) SendImageMessage(userIDList []string, departmentIDList []string, toAllUser bool, mediaID string) (
	resp DingDingAppMessageSendResp, err error) {
	messageObj := make(map[string]interface{})
	messageObj["agent_id"] = r.agentID
	if len(userIDList) > 0 {
		messageObj["userid_list"] = strings.Join(userIDList, ",")
	}
	if len(departmentIDList) > 0 {
		messageObj["dept_id_list"] = strings.Join(departmentIDList, ",")
	}
	messageObj["to_all_user"] = toAllUser
	messageObj["msg"] = map[string]interface{}{
		"msgtype": DingDingAppMessageTypeImage,
		"image":   map[string]string{"media_id": mediaID},
	}
	return r.sendMessage(&messageObj)
}

func (r *DingDingApp) SendVoiceMessage(userIDList []string, departmentIDList []string, toAllUser bool, mediaID string, duration int) (
	resp DingDingAppMessageSendResp, err error) {
	messageObj := make(map[string]interface{})
	messageObj["agent_id"] = r.agentID
	if len(userIDList) > 0 {
		messageObj["userid_list"] = strings.Join(userIDList, ",")
	}
	if len(departmentIDList) > 0 {
		messageObj["dept_id_list"] = strings.Join(departmentIDList, ",")
	}
	messageObj["to_all_user"] = toAllUser
	messageObj["msg"] = map[string]interface{}{
		"msgtype": DingDingAppMessageTypeVoice,
		"voice":   map[string]string{"media_id": mediaID, "duration": strconv.Itoa(duration)},
	}
	return r.sendMessage(&messageObj)
}

func (r *DingDingApp) SendFileMessage(userIDList []string, departmentIDList []string, toAllUser bool, mediaID string) (
	resp DingDingAppMessageSendResp, err error) {
	messageObj := make(map[string]interface{})
	messageObj["agent_id"] = r.agentID
	if len(userIDList) > 0 {
		messageObj["userid_list"] = strings.Join(userIDList, ",")
	}
	if len(departmentIDList) > 0 {
		messageObj["dept_id_list"] = strings.Join(departmentIDList, ",")
	}
	messageObj["to_all_user"] = toAllUser
	messageObj["msg"] = map[string]interface{}{
		"msgtype": DingDingAppMessageTypeFile,
		"file":    map[string]string{"media_id": mediaID},
	}
	return r.sendMessage(&messageObj)
}

func (r *DingDingApp) SendLinkMessage(userIDList []string, departmentIDList []string, toAllUser bool, linkMessage *DingDingAppLinkMessage) (
	resp DingDingAppMessageSendResp, err error) {
	messageObj := make(map[string]interface{})
	messageObj["agent_id"] = r.agentID
	if len(userIDList) > 0 {
		messageObj["userid_list"] = strings.Join(userIDList, ",")
	}
	if len(departmentIDList) > 0 {
		messageObj["dept_id_list"] = strings.Join(departmentIDList, ",")
	}
	messageObj["to_all_user"] = toAllUser
	messageObj["msg"] = map[string]interface{}{
		"msgtype": DingDingAppMessageTypeLink,
		"link":    linkMessage,
	}
	return r.sendMessage(&messageObj)
}

func (r *DingDingApp) SendActionCardMessage(userIDList []string, departmentIDList []string, toAllUser bool, actionCardMessage *DingDingAppActionCardMessage) (
	resp DingDingAppMessageSendResp, err error) {
	messageObj := make(map[string]interface{})
	messageObj["agent_id"] = r.agentID
	if len(userIDList) > 0 {
		messageObj["userid_list"] = strings.Join(userIDList, ",")
	}
	if len(departmentIDList) > 0 {
		messageObj["dept_id_list"] = strings.Join(departmentIDList, ",")
	}
	messageObj["to_all_user"] = toAllUser
	messageObj["msg"] = map[string]interface{}{
		"msgtype":     DingDingAppMessageTypeActionCard,
		"action_card": actionCardMessage,
	}
	return r.sendMessage(&messageObj)
}

func (r *DingDingApp) sendMessage(messageObj interface{}) (messageResp DingDingAppMessageSendResp, err error) {
	err = r.fireRequest(http.MethodPost, DingDingAppSendMessageAPI, nil, messageObj, &messageResp)
	if err != nil {
		return
	}
	if messageResp.ErrCode != DingDingAppStatusOK {
		if messageResp.ErrCode == DingDingCodeAccessTokenExpired {
			// reset the access token
			r.accessToken = ""
		}
		err = fmt.Errorf("call dingding app message api error, %d %s", messageResp.ErrCode, messageResp.ErrMessage)
		return
	}
	return
}

func (r *DingDingApp) GetMessageSendProgress(taskID int) (sendProgressResp DingDingAppMessageSendProgressResp, err error) {
	agentID, _ := strconv.Atoi(r.agentID)
	reqBody := map[string]int{
		"agent_id": agentID,
		"task_id":  taskID,
	}
	err = r.fireRequest(http.MethodPost, DingDingAppGetMessageSendProgressAPI, nil, &reqBody, &sendProgressResp)
	if err != nil {
		return
	}
	if sendProgressResp.ErrCode != DingDingAppStatusOK {
		if sendProgressResp.ErrCode == DingDingCodeAccessTokenExpired {
			// reset the access token
			r.accessToken = ""
		}
		err = fmt.Errorf("call dingding app message api error, %d %s", sendProgressResp.ErrCode, sendProgressResp.ErrMessage)
		return
	}
	return
}

func (r *DingDingApp) GetMessageSendResult(taskID int) (sendResultResp DingDingAppMessageSendResultResp, err error) {
	agentID, _ := strconv.Atoi(r.agentID)
	reqBody := map[string]int{
		"agent_id": agentID,
		"task_id":  taskID,
	}
	err = r.fireRequest(http.MethodPost, DingDingAppGetMessageSendResultAPI, nil, &reqBody, &sendResultResp)
	if err != nil {
		return
	}
	if sendResultResp.ErrCode != DingDingAppStatusOK {
		if sendResultResp.ErrCode == DingDingCodeAccessTokenExpired {
			// reset the access token
			r.accessToken = ""
		}
		err = fmt.Errorf("call dingding app message api error, %d %s", sendResultResp.ErrCode, sendResultResp.ErrMessage)
		return
	}
	return
}

func (r *DingDingApp) RecallMessage(taskID int) (revokeResp DingDingAppMessageRecallResp, err error) {
	agentID, _ := strconv.Atoi(r.agentID)
	reqBody := map[string]int{
		"agent_id": agentID,
		"task_id":  taskID,
	}
	err = r.fireRequest(http.MethodPost, DingDingAppRecallMessageAPI, nil, &reqBody, &revokeResp)
	if err != nil {
		return
	}
	if revokeResp.ErrCode != DingDingAppStatusOK {
		if revokeResp.ErrCode == DingDingCodeAccessTokenExpired {
			// reset the access token
			r.accessToken = ""
		}
		err = fmt.Errorf("call dingding app message api error, %d %s", revokeResp.ErrCode, revokeResp.ErrMessage)
		return
	}
	return
}

func (r *DingDingApp) SendGroupTextMessage(chatID string, content string) (
	resp DingDingAppGroupMessageSendResp, err error) {
	messageObj := make(map[string]interface{})
	messageObj["chatid"] = chatID
	messageObj["msg"] = map[string]interface{}{
		"msgtype": DingDingAppMessageTypeText,
		"text":    map[string]string{"content": content},
	}
	return r.sendGroupMessage(&messageObj)
}

func (r *DingDingApp) SendGroupMarkdownMessage(chatID, title, content string) (
	resp DingDingAppGroupMessageSendResp, err error) {
	messageObj := make(map[string]interface{})
	messageObj["chatid"] = chatID
	messageObj["msg"] = map[string]interface{}{
		"msgtype":  DingDingAppMessageTypeMarkdown,
		"markdown": map[string]string{"title": title, "text": content},
	}
	return r.sendGroupMessage(&messageObj)
}

func (r *DingDingApp) SendGroupImageMessage(chatID, mediaID string) (
	resp DingDingAppGroupMessageSendResp, err error) {
	messageObj := make(map[string]interface{})
	messageObj["chatid"] = chatID
	messageObj["msg"] = map[string]interface{}{
		"msgtype": DingDingAppMessageTypeImage,
		"image":   map[string]string{"media_id": mediaID},
	}
	return r.sendGroupMessage(&messageObj)
}

func (r *DingDingApp) SendGroupVoiceMessage(chatID, mediaID string, duration int) (
	resp DingDingAppGroupMessageSendResp, err error) {
	messageObj := make(map[string]interface{})
	messageObj["chatid"] = chatID
	messageObj["msg"] = map[string]interface{}{
		"msgtype": DingDingAppMessageTypeVoice,
		"voice":   map[string]string{"media_id": mediaID, "duration": strconv.Itoa(duration)},
	}
	return r.sendGroupMessage(&messageObj)
}

func (r *DingDingApp) SendGroupFileMessage(chatID, mediaID string) (
	resp DingDingAppGroupMessageSendResp, err error) {
	messageObj := make(map[string]interface{})
	messageObj["chatid"] = chatID
	messageObj["msg"] = map[string]interface{}{
		"msgtype": DingDingAppMessageTypeFile,
		"file":    map[string]string{"media_id": mediaID},
	}
	return r.sendGroupMessage(&messageObj)
}

func (r *DingDingApp) SendGroupLinkMessage(chatID string, linkMessage *DingDingAppLinkMessage) (
	resp DingDingAppGroupMessageSendResp, err error) {
	messageObj := make(map[string]interface{})
	messageObj["chatid"] = chatID
	messageObj["msg"] = map[string]interface{}{
		"msgtype": DingDingAppMessageTypeLink,
		"link":    linkMessage,
	}
	return r.sendGroupMessage(&messageObj)
}

func (r *DingDingApp) SendGroupActionCardMessage(chatID string, actionCardMessage *DingDingAppActionCardMessage) (
	resp DingDingAppGroupMessageSendResp, err error) {
	messageObj := make(map[string]interface{})
	messageObj["chatid"] = chatID
	messageObj["msg"] = map[string]interface{}{
		"msgtype":     DingDingAppMessageTypeActionCard,
		"action_card": actionCardMessage,
	}
	return r.sendGroupMessage(&messageObj)
}

func (r *DingDingApp) sendGroupMessage(messageObj interface{}) (messageResp DingDingAppGroupMessageSendResp, err error) {
	err = r.fireRequest(http.MethodPost, DingDingAppSendGroupMessageAPI, nil, messageObj, &messageResp)
	if err != nil {
		return
	}
	if messageResp.ErrCode != DingDingAppStatusOK {
		if messageResp.ErrCode == DingDingCodeAccessTokenExpired {
			// reset the access token
			r.accessToken = ""
		}
		err = fmt.Errorf("call dingding app group message api error, %d %s", messageResp.ErrCode, messageResp.ErrMessage)
		return
	}
	return
}

// CreateGroupChat create a new group chat
func (r *DingDingApp) CreateGroupChat(name, ownerID string, userIDList []string, options *DingDingAppCreateGroupOptions) (newChatID string, err error) {
	createGroupReqObject := make(map[string]interface{})
	createGroupReqObject["name"] = name
	createGroupReqObject["owner"] = ownerID
	createGroupReqObject["useridlist"] = userIDList
	if options != nil {
		createGroupReqObject["showHistoryType"] = options.ShowHistoryType
		createGroupReqObject["searchable"] = options.Searchable
		createGroupReqObject["validationType"] = options.ValidationType
		createGroupReqObject["mentionAllAuthority"] = options.MentionAllAuthority
		createGroupReqObject["chatBannedType"] = options.ChatBannedType
		createGroupReqObject["managementType"] = options.ManagementType
	}
	var createGroupResp DingDingAppCreateGroupResp
	err = r.fireRequest(http.MethodPost, DingDingAppCreateGroupAPI, nil, &createGroupReqObject, &createGroupResp)
	if err != nil {
		return
	}
	if createGroupResp.ErrCode != DingDingAppStatusOK {
		if createGroupResp.ErrCode == DingDingCodeAccessTokenExpired {
			// reset the access token
			r.accessToken = ""
		}
		err = fmt.Errorf("call dingding app create group api error, %d %s", createGroupResp.ErrCode, createGroupResp.ErrMessage)
		return
	}
	newChatID = createGroupResp.ChatID
	return
}

func (r *DingDingApp) UpdateGroupChat(chatID string, options *DingDingAppUpdateGroupOptions) (err error) {
	updateGroupReqObject := make(map[string]interface{})
	updateGroupReqObject["chatid"] = chatID
	if options != nil {
		updateGroupReqObject["name"] = options.Name
		updateGroupReqObject["icon"] = options.Icon
		updateGroupReqObject["owner"] = options.Owner
		updateGroupReqObject["add_useridlist"] = options.AddUserList
		updateGroupReqObject["del_useridlist"] = options.DelUserList
		updateGroupReqObject["showHistoryType"] = options.ShowHistoryType
		updateGroupReqObject["searchable"] = options.Searchable
		updateGroupReqObject["validationType"] = options.ValidationType
		updateGroupReqObject["mentionAllAuthority"] = options.MentionAllAuthority
		updateGroupReqObject["chatBannedType"] = options.ChatBannedType
		updateGroupReqObject["managementType"] = options.ManagementType
	}
	var updateGroupResp DingDingAppUpdateGroupResp
	err = r.fireRequest(http.MethodPost, DingDingAppUpdateGroupAPI, nil, &updateGroupReqObject, &updateGroupResp)
	if err != nil {
		return
	}
	if updateGroupResp.ErrCode != DingDingAppStatusOK {
		if updateGroupResp.ErrCode == DingDingCodeAccessTokenExpired {
			// reset the access token
			r.accessToken = ""
		}
		err = fmt.Errorf("call dingding app update group api error, %d %s", updateGroupResp.ErrCode, updateGroupResp.ErrMessage)
		return
	}
	return
}

func (r *DingDingApp) GetGroupChat(chatID string) (group DingDingAppGroup, err error) {
	var getGroupResp DingDingAppGetGroupResp
	err = r.fireRequest(http.MethodGet, DingDingAppGetGroupAPI, map[string]string{"chatid": chatID}, nil, &getGroupResp)
	if err != nil {
		return
	}
	if getGroupResp.ErrCode != DingDingAppStatusOK {
		if getGroupResp.ErrCode == DingDingCodeAccessTokenExpired {
			// reset the access token
			r.accessToken = ""
		}
		err = fmt.Errorf("call dingding app get group api error, %d %s", getGroupResp.ErrCode, getGroupResp.ErrMessage)
		return
	}
	group = getGroupResp.ChatInfo
	return
}

func (r *DingDingApp) refreshAccessToken() (err error) {
	reqURL := fmt.Sprintf("%s?appkey=%s&appsecret=%s", DingDingAppTokenAPI, r.appKey, r.appSecret)
	req, newErr := http.NewRequest(http.MethodGet, reqURL, nil)
	if newErr != nil {
		err = fmt.Errorf("create request error, %s", newErr.Error())
		return
	}
	resp, getErr := r.client.Do(req)
	if getErr != nil {
		err = fmt.Errorf("get response error, %s", getErr.Error())
		return
	}
	defer resp.Body.Close()

	// check http code
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("dingding request error, %s", resp.Status)
		io.Copy(ioutil.Discard, resp.Body)
		return
	}
	// parse response body
	decoder := json.NewDecoder(resp.Body)
	var accessTokenResp DingDingAppTokenResp
	if decodeErr := decoder.Decode(&accessTokenResp); decodeErr != nil {
		err = fmt.Errorf("parse response error, %s", decodeErr.Error())
		return
	}
	if accessTokenResp.ErrCode != DingDingAppStatusOK {
		err = fmt.Errorf("call dingding app api error, %d %s", accessTokenResp.ErrCode, accessTokenResp.ErrMessage)
		return
	}
	// set access token and expired at
	r.accessToken = accessTokenResp.AccessToken
	r.expiredAt = time.Now().Add(time.Second * time.Duration(accessTokenResp.ExpiresIn))
	return
}

func (r *DingDingApp) fireRequest(reqMethod, reqURL string, reqParams map[string]string, reqBodyObject interface{}, respObject interface{}) (err error) {
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
		return
	}
	defer resp.Body.Close()
	// check http code
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("dingding request error, %s", resp.Status)
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

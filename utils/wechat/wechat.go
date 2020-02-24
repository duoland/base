package wechat

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// NotifyConfig holds the config of wechat notification
type NotifyConfig struct {
	APIAddress   string `yaml:"apiAddress"`
	AppID        string `yaml:"appID"`
	MessageToken string `yaml:"messageToken"`
}

// GroupCreateConfig holds the config of wechat group creation
type GroupCreateConfig struct {
	APIAddress   string `yaml:"apiAddress"`
	AppID        string `yaml:"appID"`
	MessageToken string `yaml:"messageToken"`
}

// Error is the backend API error
// Doc: http://wiki.lianjia.com/pages/viewpage.action?pageId=308971537
type Error struct {
	Code      int    `json:"code"`
	Message   string `json:"message,omitempty"`
	Path      string `json:"path,omitempty"`
	Timestamp int64  `json:"timestamp,omitempty"`
}

// UserCardNotify send the card message to a single user
func UserCardNotify(cfg *NotifyConfig, staffID, title, message, targetURL, buttonText string) (err error) {
	timestamp := time.Now().Unix()
	signature := md5HexString([]byte(fmt.Sprintf("%d%s%s", timestamp, cfg.AppID, cfg.MessageToken)))

	// pack request body
	params := url.Values{}
	params.Add("appId", cfg.AppID)
	params.Add("timestamp", fmt.Sprintf("%d", timestamp))
	params.Add("signature", signature)
	params.Add("touser", staffID)
	params.Add("title", title)
	params.Add("description", message)
	params.Add("url", targetURL)
	params.Add("btntxt", buttonText)

	return wechatRequest(cfg.APIAddress, params)
}

// UserMarkdownNotify send the markdown message to a single user
func UserMarkdownNotify(cfg *NotifyConfig, staffID, message string) (err error) {
	timestamp := time.Now().Unix()
	signature := md5HexString([]byte(fmt.Sprintf("%d%s%s", timestamp, cfg.AppID, cfg.MessageToken)))

	// pack request body
	params := url.Values{}
	params.Add("appId", cfg.AppID)
	params.Add("touser", staffID)
	params.Add("timestamp", fmt.Sprintf("%d", timestamp))
	params.Add("signature", signature)
	params.Add("content", message)

	return wechatRequest(cfg.APIAddress, params)
}

// UserTextNotify send the text message to a single user
func UserTextNotify(cfg *NotifyConfig, staffID, message string) (err error) {
	timestamp := time.Now().Unix()
	signature := md5HexString([]byte(fmt.Sprintf("%d%s%s", timestamp, cfg.AppID, cfg.MessageToken)))

	// pack request body
	params := url.Values{}
	params.Add("appId", cfg.AppID)
	params.Add("touser", staffID)
	params.Add("timestamp", fmt.Sprintf("%d", timestamp))
	params.Add("signature", signature)
	params.Add("content", message)

	return wechatRequest(cfg.APIAddress, params)
}

// GroupCardNotify send the card message to a wechat group
func GroupCardNotify(cfg *NotifyConfig, chatID, title, message, targetURL, buttonText string) (err error) {
	timestamp := time.Now().Unix()
	signature := md5HexString([]byte(fmt.Sprintf("%d%s%s%s", timestamp, cfg.AppID, chatID, cfg.MessageToken)))

	// pack request body
	params := url.Values{}
	params.Add("timestamp", fmt.Sprintf("%d", timestamp))
	params.Add("signature", signature)
	params.Add("chatId", chatID)
	params.Add("title", title)
	params.Add("description", message)
	params.Add("url", targetURL)
	params.Add("btntxt", buttonText)

	return wechatRequest(cfg.APIAddress, params)
}

// GroupMarkdownNotify send the markdown message to a wechat group
func GroupMarkdownNotify(cfg *NotifyConfig, chatID, message string) (err error) {
	timestamp := time.Now().Unix()
	signature := md5HexString([]byte(fmt.Sprintf("%d%s%s%s", timestamp, cfg.AppID, chatID, cfg.MessageToken)))

	// pack request body
	params := url.Values{}
	params.Add("chatId", chatID)
	params.Add("timestamp", fmt.Sprintf("%d", timestamp))
	params.Add("signature", signature)
	params.Add("content", message)

	return wechatRequest(cfg.APIAddress, params)
}

// GroupTextNotify send the text message to a wechat group
func GroupTextNotify(cfg *NotifyConfig, chatID, message string) (err error) {
	timestamp := time.Now().Unix()
	signature := md5HexString([]byte(fmt.Sprintf("%d%s%s%s", timestamp, cfg.AppID, chatID, cfg.MessageToken)))

	// pack request body
	params := url.Values{}
	params.Add("chatId", chatID)
	params.Add("timestamp", fmt.Sprintf("%d", timestamp))
	params.Add("signature", signature)
	params.Add("content", message)

	return wechatRequest(cfg.APIAddress, params)
}

// GroupCreate create a new wechat group
func GroupCreate(cfg *GroupCreateConfig, chatID, groupName, ownerStaffID, userStaffIDs string) (err error) {
	timestamp := time.Now().Unix()
	signature := md5HexString([]byte(fmt.Sprintf("%d%s%s%s", timestamp, cfg.AppID, chatID, cfg.MessageToken)))

	// pack request body
	params := url.Values{}
	params.Add("appId", cfg.AppID)
	params.Add("timestamp", fmt.Sprintf("%d", timestamp))
	params.Add("signature", signature)
	params.Add("name", groupName)
	params.Add("owner", ownerStaffID)
	params.Add("userlist", strings.Join([]string{ownerStaffID, userStaffIDs}, ","))
	params.Add("chatId", chatID)

	return wechatRequest(cfg.APIAddress, params)
}

func md5HexString(data []byte) string {
	md5Hasher := md5.New()
	md5Hasher.Write(data)
	return hex.EncodeToString(md5Hasher.Sum(nil))
}

func wechatRequest(reqURL string, reqParams url.Values) (err error) {
	var req *http.Request
	var resp *http.Response

	// prepare req
	req, err = http.NewRequest("POST", reqURL, bytes.NewReader([]byte(reqParams.Encode())))
	if err != nil {
		return
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")

	// timeout in 10 seconds
	client := http.DefaultClient
	client.Timeout = time.Second * 10

	// parse resp
	resp, err = client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var respData []byte
		respData, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return
		}

		// parse error
		var wechatError Error
		json.Unmarshal(respData, &wechatError)
		if wechatError.Code == -1 {
			err = errors.New(wechatError.Message)
		} else {
			err = errors.New(string(respData))
		}

		return
	}
	return
}

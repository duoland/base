package dingtalk

// See doc at https://ding-doc.dingtalk.com/doc#/serverapi2/qf2nxq
import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/duoland/base/hash"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// DingDingRobotMessageAPI is the api to send the robot messages
const DingDingRobotMessageAPI = "https://oapi.dingtalk.com/robot/send"

// DingDingRobotTimeout is the dingding robot default timeout
const DingDingRobotTimeout = time.Second * 10

// DingDingRobotStatusOK is the ok status of api call
const DingDingRobotStatusOK = 0

const (
	DingDingRobotMessageTypeText       = "text"
	DingDingRobotMessageTypeLink       = "link"
	DingDingRobotMessageTypeMarkdown   = "markdown"
	DingDingRobotMessageTypeActionCard = "actionCard"
	DingDingRobotMessageTypeFeedCard   = "feedCard"
)

const (
	DingDingActionCardMessageButtonOrientationVertical   = "0"
	DingDingActionCardMessageButtonOrientationHorizontal = "1"
)

type DingDingRobotMessageResp struct {
	ErrCode    int    `json:"errcode"`
	ErrMessage string `json:"errmsg"`
}

type DingDingRobot struct {
	client *http.Client
}

// NewDingDingRobot create a new dingding robot
func NewDingDingRobot() *DingDingRobot {
	return NewDingDingRobotWithTimeout(DingDingRobotTimeout)
}

// NewDingDingRobotWithTimeout create a new dingding robot with timeout
func NewDingDingRobotWithTimeout(timeout time.Duration) *DingDingRobot {
	client := http.Client{}
	client.Timeout = timeout
	return &DingDingRobot{client: &client}
}

// NewDingDingRobotWithClient create a new dingding robot with http.Client
func NewDingDingRobotWithClient(client *http.Client) *DingDingRobot {
	return &DingDingRobot{client: client}
}

type DingDingRobotMentionAt struct {
	AtMobiles []string `json:"atMobiles"`
	IsAtAll   bool     `json:"isAtAll"`
}

type DingDingSecuritySettings struct {
	AccessToken string
	SecureToken string
}

type DingDingRobotLinkMessage struct {
	Text       string `json:"text"`
	Title      string `json:"title"`
	PicURL     string `json:"picUrl"`
	MessageURL string `json:"messageUrl"`
}

type DingDingRobotMarkdownMessage struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

type DingDingRobotActionCardMessage struct {
	Title string `json:"title"`
	Text  string `json:"text"`
	// single jump action card fields
	SingleTitle string `json:"singleTitle,omitempty"`
	SingleURL   string `json:"singleURL,omitempty"`
	// standalone jump action card fields
	ButtonOrientation string                          `json:"btnOrientation"`
	Buttons           []DingDingRobotActionCardButton `json:"btns,omitempty"`
}

type DingDingRobotActionCardButton struct {
	Title     string `json:"title"`
	ActionURL string `json:"actionURL"`
}

type DingDingRobotFeedCardMessage struct {
	Title      string `json:"title"`
	MessageURL string `json:"messageURL"`
	PicURL     string `json:"picURL"`
}

func (r *DingDingRobot) SendTextMessage(securitySettings *DingDingSecuritySettings, content string) (err error) {
	return r.SendTextMessageWithMention(securitySettings, content, nil, false)
}

func (r *DingDingRobot) SendTextMessageWithMention(securitySettings *DingDingSecuritySettings, content string, mentionedMobileList []string, atAll bool) (err error) {
	messageObj := make(map[string]interface{})
	messageObj["msgtype"] = DingDingRobotMessageTypeText
	messageObj["text"] = map[string]string{"content": content}
	messageObj["at"] = DingDingRobotMentionAt{
		AtMobiles: mentionedMobileList,
		IsAtAll:   atAll,
	}
	return r.sendMessage(securitySettings, &messageObj)
}

func (r *DingDingRobot) SendMarkdownMessage(securitySettings *DingDingSecuritySettings, markdownMessage *DingDingRobotMarkdownMessage) (err error) {
	return r.SendMarkdownMessageWithMention(securitySettings, markdownMessage, nil, false)
}

func (r *DingDingRobot) SendMarkdownMessageWithMention(securitySettings *DingDingSecuritySettings, markdownMessage *DingDingRobotMarkdownMessage,
	mentionedMobileList []string, atAll bool) (err error) {
	messageObj := make(map[string]interface{})
	messageObj["msgtype"] = DingDingRobotMessageTypeMarkdown
	messageObj["markdown"] = markdownMessage
	messageObj["at"] = DingDingRobotMentionAt{
		AtMobiles: mentionedMobileList,
		IsAtAll:   atAll,
	}
	return r.sendMessage(securitySettings, &messageObj)
}

func (r *DingDingRobot) SendLinkMessage(securitySettings *DingDingSecuritySettings, linkMessage *DingDingRobotLinkMessage) (err error) {
	messageObj := make(map[string]interface{})
	messageObj["msgtype"] = DingDingRobotMessageTypeLink
	messageObj["link"] = linkMessage
	return r.sendMessage(securitySettings, messageObj)
}

func (r *DingDingRobot) SendActionCardMessage(securitySettings *DingDingSecuritySettings, actionCardMessage *DingDingRobotActionCardMessage) (err error) {
	messageObj := make(map[string]interface{})
	messageObj["msgtype"] = DingDingRobotMessageTypeActionCard
	messageObj["actionCard"] = actionCardMessage
	return r.sendMessage(securitySettings, messageObj)
}

func (r *DingDingRobot) SendFeedCardMessage(securitySettings *DingDingSecuritySettings, feedCardMessages []DingDingRobotFeedCardMessage) (err error) {
	messageObj := make(map[string]interface{})
	messageObj["msgtype"] = DingDingRobotMessageTypeFeedCard
	messageObj["feedCard"] = map[string]interface{}{
		"links": feedCardMessages,
	}
	return r.sendMessage(securitySettings, messageObj)
}

func (r *DingDingRobot) sendMessage(securitySettings *DingDingSecuritySettings, messageObj interface{}) (err error) {
	reqParams := url.Values{}
	reqParams.Add("access_token", securitySettings.AccessToken)
	if securitySettings.SecureToken != "" {
		tsNow := time.Now().UnixNano() / 1000000
		dataToSign := fmt.Sprintf("%d\n%s", tsNow, securitySettings.SecureToken)
		sign := base64.StdEncoding.EncodeToString(hash.HmacSha256([]byte(dataToSign), []byte(securitySettings.SecureToken)))
		reqParams.Add("timestamp", fmt.Sprintf("%d", tsNow))
		reqParams.Add("sign", sign)
	}
	reqURL := fmt.Sprintf("%s?%s", DingDingRobotMessageAPI, reqParams.Encode())
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
		return
	}
	defer resp.Body.Close()
	// check http code
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("dingtalk request error, %s", resp.Status)
		io.Copy(ioutil.Discard, resp.Body)
		return
	}
	// parse response body
	decoder := json.NewDecoder(resp.Body)
	var dingdingMessageResp DingDingRobotMessageResp
	if decodeErr := decoder.Decode(&dingdingMessageResp); decodeErr != nil {
		err = fmt.Errorf("parse response error, %s", decodeErr.Error())
		return
	}
	if dingdingMessageResp.ErrCode != DingDingRobotStatusOK {
		err = fmt.Errorf("call dingtalk robot api error, %d %s", dingdingMessageResp.ErrCode, dingdingMessageResp.ErrMessage)
		return
	}
	return
}

package bytedance

// See doc https://getfeishu.cn/hc/zh-cn/articles/360040566333

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

// FeiShuRobotMessageAPI is the api to send the shortcut messages
const FeiShuRobotMessageAPI = "https://www.feishu.cn/flow/api/trigger-webhook/"

// FeiShuRobotTimeout is the feishu shortcut default timeout
const FeiShuRobotTimeout = time.Second * 10
const FeiShuRobotStatusOK = 0

// FeiShuRobot is a robot to send feishu shortcut messages
type FeiShuRobot struct {
	client *http.Client
}

type FeiShuRobotMessageResp struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
}

// NewFeiShuRobot create a new feishu shortcut robot
func NewFeiShuRobot() *FeiShuRobot {
	return NewFeiShuRobotWithTimeout(FeiShuRobotTimeout)
}

// NewFeiShuRobotWithTimeout create a new feishu shortcut robot with timeout
func NewFeiShuRobotWithTimeout(timeout time.Duration) *FeiShuRobot {
	client := http.Client{}
	client.Timeout = timeout
	return &FeiShuRobot{client: &client}
}

// NewFeiShuRobotWithClient create a new feishu shortcut robot with http.Client
func NewFeiShuRobotWithClient(client *http.Client) *FeiShuRobot {
	return &FeiShuRobot{client: client}
}

// SendTextMessage send the text message
func (r *FeiShuRobot) SendTextMessage(key, title, content string) (err error) {
	messageObj := map[string]string{"title": title, "content": content}
	return r.sendMessage(key, &messageObj)
}

func (r *FeiShuRobot) sendMessage(key string, messageObj interface{}) (err error) {
	reqURL := fmt.Sprintf("%s%s", FeiShuRobotMessageAPI, key)
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
	io.Copy(os.Stdout, resp.Body)
	// check http code
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("feishu shortcut request error, %s", resp.Status)
		io.Copy(ioutil.Discard, resp.Body)
		return
	}
	// parse response body
	decoder := json.NewDecoder(resp.Body)
	var messageResp FeiShuRobotMessageResp
	if decodeErr := decoder.Decode(&messageResp); decodeErr != nil {
		err = fmt.Errorf("parse response error, %s", decodeErr.Error())
		return
	}
	if messageResp.Code != FeiShuRobotStatusOK {
		err = fmt.Errorf("call feishu shortcut api error, %d %s", messageResp.Code, messageResp.Message)
		return
	}
	return
}

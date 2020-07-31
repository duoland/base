package dingtalk

import (
	"fmt"
	"os"
	"testing"
)

var appKey = os.Getenv("DINGDING_APP_KEY")
var appSecret = os.Getenv("DINGDING_APP_SECRET")
var agentID = os.Getenv("DINGDING_APP_AGENT_ID")

var userIDList = []string{"manager2159"}
var departmentIDList = []string{"381323914"}
var toAllUser = false

func init() {
	fmt.Println("==> AgentID:", agentID)
	fmt.Println("==> AppKey:", appKey)
	fmt.Println("==> AppSecret:", appSecret)
	fmt.Println("")
}
func TestDingDingApp_refreshAccessToken(t *testing.T) {
	dingdingApp := NewDingDingApp(appKey, appSecret, agentID)
	err := dingdingApp.refreshAccessToken()
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func TestDingDingApp_SendTextMessage(t *testing.T) {
	dingdingApp := NewDingDingApp(appKey, appSecret, agentID)
	content := "hello, master, i am robot for your service"
	resp, err := dingdingApp.SendTextMessage(userIDList, nil, toAllUser, content)
	t.Logf("%v", resp)
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func TestDingDingApp_GetMessageSendProgress(t *testing.T) {
	taskID := 240072433293
	dingdingApp := NewDingDingApp(appKey, appSecret, agentID)
	resp, err := dingdingApp.GetMessageSendProgress(taskID)
	t.Logf("%v", resp)
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func TestDingDingApp_GetMessageSendResult(t *testing.T) {
	taskID := 240072433293
	dingdingApp := NewDingDingApp(appKey, appSecret, agentID)
	resp, err := dingdingApp.GetMessageSendResult(taskID)
	t.Logf("%v", resp)
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func TestDingDingApp_RecallMessage(t *testing.T) {
	taskID := 240072433293
	dingdingApp := NewDingDingApp(appKey, appSecret, agentID)
	resp, err := dingdingApp.RecallMessage(taskID)
	t.Logf("%v", resp)
	if err != nil {
		t.Fatalf(err.Error())
	}
}

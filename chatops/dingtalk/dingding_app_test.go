package dingtalk

import (
	"fmt"
	"os"
	"testing"
)

var chatID = "chat52999a8e1bdedfe94bb6e9841d581c9e"
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

func TestDingDingApp_SendMarkdownMessage(t *testing.T) {
	dingdingApp := NewDingDingApp(appKey, appSecret, agentID)
	content := "# hello, master, i am robot for your service"
	resp, err := dingdingApp.SendMarkdownMessage(userIDList, nil, toAllUser, "hello master", content)
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

// ChatID: chat52999a8e1bdedfe94bb6e9841d581c9e
func TestDingDingApp_CreateGroupChat(t *testing.T) {
	dingdingApp := NewDingDingApp(appKey, appSecret, agentID)
	name := "一个不简单的测试群"
	resp, err := dingdingApp.CreateGroupChat(name, userIDList[0], userIDList, nil)
	t.Logf("%v", resp)
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func TestDingDingApp_SendGroupTextMessage(t *testing.T) {
	dingdingApp := NewDingDingApp(appKey, appSecret, agentID)
	content := "hello, master, i am robot for your service"
	resp, err := dingdingApp.SendGroupTextMessage(chatID, content)
	t.Logf("%v", resp)
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func TestDingDingApp_SendGroupMarkdownMessage(t *testing.T) {
	dingdingApp := NewDingDingApp(appKey, appSecret, agentID)
	content := "# hello, master, i am robot for your service"
	resp, err := dingdingApp.SendGroupMarkdownMessage(chatID, "hello master", content)
	t.Logf("%v", resp)
	if err != nil {
		t.Fatalf(err.Error())
	}
}

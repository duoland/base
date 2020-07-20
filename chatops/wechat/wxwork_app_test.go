package wechat

import "testing"

var chatID = "test20200207"
var corpID = ""
var corpSecret = ""
var agentID = "1000002"

func TestWxWorkApp_CreateGroupChat(t *testing.T) {
	wxworkApp := NewWxWorkApp(corpID, corpSecret, agentID)
	newChatID, err := wxworkApp.CreateGroupChat("一个简单的测试群", chatID, "jinxinxin001", []string{"jinchengxi001", "jinxinxin001"})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("ChatID: %s\n", newChatID)
}

func TestWxWorkApp_SendGroupTextMessage(t *testing.T) {
	wxworkApp := NewWxWorkApp(corpID, corpSecret, agentID)
	err := wxworkApp.SendGroupTextMessage(chatID, "hello, master", false)
	if err != nil {
		t.Fatal(err)
	}
}

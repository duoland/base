package wechat

import "testing"

var chatID = "test20200207"
var corpID = "qqfc50e757c8e5ee4a"
var corpSecret = "rnrsvlPyjUstiIveyHmhuSJazEgMFx8bDq88s3yD6nk"
var agentID = "1000002"

func TestWxWorkApp_CreateGroupChat(t *testing.T) {
	wxworkApp := NewWxWorkApp(corpID, corpSecret, agentID)
	newChatID, err := wxworkApp.CreateGroupChat("一个简单的测试群", chatID, "jinxinxin001", []string{"jinchengxi001", "jinxinxin001"})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("ChatID: %s\n", newChatID)
}

func TestWxWorkApp_UpdateGroupChat(t *testing.T) {
	wxworkApp := NewWxWorkApp(corpID, corpSecret, agentID)
	name := "一个不简单的测试群"
	err := wxworkApp.UpdateGroupChat(name, chatID, "jinxinxin001", []string{"jinchengxi001"}, []string{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestWxWorkApp_GetGroupChat(t *testing.T) {
	wxworkApp := NewWxWorkApp(corpID, corpSecret, agentID)
	group, err := wxworkApp.GetGroupChat(chatID)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("GroupInfo: %v\n", group)
}

func TestWxWorkApp_SendGroupTextMessage(t *testing.T) {
	wxworkApp := NewWxWorkApp(corpID, corpSecret, agentID)
	err := wxworkApp.SendGroupTextMessage(chatID, "hello, master", false)
	if err != nil {
		t.Fatal(err)
	}
}

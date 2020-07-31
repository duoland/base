package bytedance

import (
	"os"
	"testing"
)

var userID = "da6e7g5d"
var chatID = "oc_84971ebabfe5bd9c8cb3a1cb76b6248a"
var appID = os.Getenv("FEISHU_APP_ID")
var appSecret = os.Getenv("FEISHU_APP_SECRET")
var groupName = "一个不简单的测试群"
var groupDescription = "我就是一个机器人创建的测试群"
var groupUserIDList = []string{"da6e7g5d"}

func TestFeiShuApp_CreateGroupChat(t *testing.T) {
	feishuApp := NewFeiShuApp(appID, appSecret)

	chatId, err := feishuApp.CreateGroupChat(groupName, groupDescription, groupUserIDList, nil)
	if err != nil {
		t.Fatal(err.Error())
		return
	}
	t.Logf("Chat ID: %s", chatId)
}

func TestFeiShuApp_SendTextMessage(t *testing.T) {
	feishuApp := NewFeiShuApp(appID, appSecret)
	target := FeiShuAppMessageSendTarget{
		//UserID: userID,
		ChatID: chatID,
	}
	content := "<h1>hello master<h1><a href='https://www.baidu.com'>baidu</a>"
	resp, err := feishuApp.SendTextMessage(&target, content, nil)
	t.Logf("%v\n", resp)
	if err != nil {
		t.Fatal(err.Error())
		return
	}
}

func TestFeiShuApp_SendPostMessage(t *testing.T) {
	feishuApp := NewFeiShuApp(appID, appSecret)
	title := "hello, i am a robot"
	target := FeiShuAppMessageSendTarget{
		//UserID: userID,
		ChatID: chatID,
	}
	postItems := [][]FeishuAppPostMessageContentItem{
		[]FeishuAppPostMessageContentItem{
			FeishuAppPostMessageContentItem{
				Tag:      FeiShuAppPostMessageText,
				UnEscape: false,
				Text:     "<h1>hello master<h1><a href='https://www.baidu.com'>baidu</a>",
			},
		},
		[]FeishuAppPostMessageContentItem{
			FeishuAppPostMessageContentItem{
				Tag:    FeiShuAppPostMessageAt,
				UserID: userID,
			},
		},
	}
	resp, err := feishuApp.SendPostMessage(&target, title, FeiShuAppI18nChinese, postItems, nil)
	t.Logf("%v\n", resp)
	if err != nil {
		t.Fatal(err.Error())
		return
	}
}

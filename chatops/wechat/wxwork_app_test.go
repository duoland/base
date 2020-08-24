package wechat

import (
	"encoding/base64"
	"os"
	"testing"
)

var chatID = "citest"
var corpID = os.Getenv("WXWORK_APP_CORP_ID")
var corpSecret = os.Getenv("WXWORK_APP_CORP_SECRET")
var agentID = os.Getenv("WXWORK_APP_AGENT_ID")
var userIDList = []string{"jinxinxin001", "jinchengxi001"}
var partyIDList = []string{}
var tagIDList = []string{}
var mediaID = "3oa0dnZ6N3dH0Y1DJFx1Mm7BJFv5UF5jE1Cni_R6uc6w"
var newsArticle = WxWorkAppNewsMessageArticle{
	Title:       "你好，我爱开源技术",
	Description: "这是一个技术网站",
	URL:         "https://oschina.net",
	PictureURL:  "https://wework.qpic.cn/wwpic/12732_8Z9RVL3rS7-S472_1595229725/0",
}
var mpNewsArticle = WxWorkAppMpNewsMessageArticle{
	Title:            "你好，我爱开源技术",
	ThumbMediaID:     mediaID,
	Author:           "小青蛙",
	ContentSourceURL: "https://oschina.net",
	Content:          "这是一个技术网站，不管你相信不相信",
	Digest:           "这是一个技术网站",
}

func TestWxWorkApp_SendTextMessage(t *testing.T) {
	wxworkApp := NewWxWorkApp(corpID, corpSecret, agentID)
	_, err := wxworkApp.SendTextMessage(userIDList, nil, nil, "hello, master",
		&WxWorkAppMessageSendOptions{EnableIDTrans: true})
	if err != nil {
		t.Fatal(err)
	}
}

func TestWxWorkApp_SendImageMessage(t *testing.T) {
	wxworkApp := NewWxWorkApp(corpID, corpSecret, agentID)
	_, err := wxworkApp.SendImageMessage(userIDList, nil, nil, mediaID,
		&WxWorkAppMessageSendOptions{EnableIDTrans: true})
	if err != nil {
		t.Fatal(err)
	}
}

func TestWxWorkApp_SendMarkdownMessage(t *testing.T) {
	wxworkApp := NewWxWorkApp(corpID, corpSecret, agentID)
	_, err := wxworkApp.SendMarkdownMessage(userIDList, nil, nil, `# hello
> big brother, i love you!
`,
		&WxWorkAppMessageSendOptions{EnableIDTrans: true})
	if err != nil {
		t.Fatal(err)
	}
}

func TestWxWorkApp_SendFileMessage(t *testing.T) {
	wxworkApp := NewWxWorkApp(corpID, corpSecret, agentID)
	_, err := wxworkApp.SendFileMessage(userIDList, nil, nil, mediaID,
		&WxWorkAppMessageSendOptions{EnableIDTrans: true})
	if err != nil {
		t.Fatal(err)
	}
}

func TestWxWorkApp_SendVoiceMessage(t *testing.T) {
	wxworkApp := NewWxWorkApp(corpID, corpSecret, agentID)
	_, err := wxworkApp.SendVoiceMessage(userIDList, nil, nil, mediaID,
		&WxWorkAppMessageSendOptions{EnableIDTrans: true})
	if err != nil {
		t.Fatal(err)
	}
}

func TestWxWorkApp_SendVideoMessage(t *testing.T) {
	wxworkApp := NewWxWorkApp(corpID, corpSecret, agentID)
	_, err := wxworkApp.SendVideoMessage(userIDList, nil, nil, mediaID, "人民", "伟大的人民",
		&WxWorkAppMessageSendOptions{EnableIDTrans: true})
	if err != nil {
		t.Fatal(err)
	}
}

func TestWxWorkApp_SendTextCardMessageC(t *testing.T) {
	wxworkApp := NewWxWorkApp(corpID, corpSecret, agentID)
	_, err := wxworkApp.SendTextCardMessage(userIDList, nil, nil, "人民", "伟大的人民",
		"https://oschina.net", "看看",
		&WxWorkAppMessageSendOptions{EnableIDTrans: true})
	if err != nil {
		t.Fatal(err)
	}
}

func TestWxWorkApp_SendNewsMessage(t *testing.T) {
	wxworkApp := NewWxWorkApp(corpID, corpSecret, agentID)
	_, err := wxworkApp.SendNewsMessage(userIDList, nil, nil, []WxWorkAppNewsMessageArticle{newsArticle},
		&WxWorkAppMessageSendOptions{EnableIDTrans: true})
	if err != nil {
		t.Fatal(err)
	}
}

func TestWxWorkApp_SendMpNewsMessage(t *testing.T) {
	wxworkApp := NewWxWorkApp(corpID, corpSecret, agentID)
	_, err := wxworkApp.SendMpNewsMessage(userIDList, nil, nil, []WxWorkAppMpNewsMessageArticle{mpNewsArticle},
		&WxWorkAppMessageSendOptions{EnableIDTrans: true})
	if err != nil {
		t.Fatal(err)
	}
}

func TestWxWorkApp_SendTaskCardMessage(t *testing.T) {
	wxworkApp := NewWxWorkApp(corpID, corpSecret, agentID)
	_, err := wxworkApp.SendTaskCardMessage(userIDList, nil, nil, "task123", "我要请假", "回家休息",
		"http://oschina.net", []WxWorkAppTaskCardMessageButton{
			{
				Key:   "keyOk",
				Name:  "批准",
				Color: "blue",
			},
			{
				Key:   "keyNo",
				Name:  "驳回",
				Color: "red",
			},
		},
		&WxWorkAppMessageSendOptions{EnableIDTrans: true})
	if err != nil {
		t.Fatal(err)
	}
}

//// Upload Media & Image ///
func TestWxWorkApp_UploadImage(t *testing.T) {
	wxworkApp := NewWxWorkApp(corpID, corpSecret, agentID)
	var imageData, _ = base64.StdEncoding.DecodeString(imageBase64Data)
	imageURL, err := wxworkApp.UploadImage(imageData, "golang.png")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Image URL: %s\n", imageURL)
}

func TestWxWorkApp_UploadMedia(t *testing.T) {
	wxworkApp := NewWxWorkApp(corpID, corpSecret, agentID)
	var imageData, _ = base64.StdEncoding.DecodeString(imageBase64Data)
	mediaID, _, err := wxworkApp.UploadMedia(imageData, "golang.png", WxWorkAppMediaTypeImage)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Media ID: %s\n", mediaID)
}

//// Group ////

func TestWxWorkApp_CreateGroupChat(t *testing.T) {
	wxworkApp := NewWxWorkApp(corpID, corpSecret, agentID)
	newChatID, err := wxworkApp.CreateGroupChat("一个简单的测试群", "jinxinxin001", []string{"jinchengxi001", "jinxinxin001"},
		&WxWorkAppCreateGroupOptions{ChatID: chatID})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("ChatID: %s\n", newChatID)
}

func TestWxWorkApp_UpdateGroupChat(t *testing.T) {
	wxworkApp := NewWxWorkApp(corpID, corpSecret, agentID)
	name := "一个不简单的测试群"
	err := wxworkApp.UpdateGroupChat(chatID, &WxWorkAppUpdateGroupOptions{Name: name})
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
	err := wxworkApp.SendGroupTextMessage(chatID, "hello, master", &WxWorkAppMessageSendOptions{Safe: true})
	if err != nil {
		t.Fatal(err)
	}
}

func TestWxWorkApp_SendGroupImageMessage(t *testing.T) {
	wxworkApp := NewWxWorkApp(corpID, corpSecret, agentID)
	err := wxworkApp.SendGroupImageMessage(chatID, mediaID, &WxWorkAppMessageSendOptions{Safe: true})
	if err != nil {
		t.Fatal(err)
	}
}

// must be Safe=false
func TestWxWorkApp_SendGroupMarkdownMessage(t *testing.T) {
	wxworkApp := NewWxWorkApp(corpID, corpSecret, agentID)
	err := wxworkApp.SendGroupMarkdownMessage(chatID, `# hello
> big brother, i love you!
`,
		&WxWorkAppMessageSendOptions{Safe: false})
	if err != nil {
		t.Fatal(err)
	}
}

func TestWxWorkApp_SendGroupFileMessage(t *testing.T) {
	wxworkApp := NewWxWorkApp(corpID, corpSecret, agentID)
	err := wxworkApp.SendGroupFileMessage(chatID, mediaID,
		&WxWorkAppMessageSendOptions{Safe: true})
	if err != nil {
		t.Fatal(err)
	}
}

func TestWxWorkApp_SendGroupVoiceMessage(t *testing.T) {
	wxworkApp := NewWxWorkApp(corpID, corpSecret, agentID)
	err := wxworkApp.SendGroupVoiceMessage(chatID, mediaID,
		&WxWorkAppMessageSendOptions{Safe: false})
	if err != nil {
		t.Fatal(err)
	}
}

func TestWxWorkApp_SendGroupVideoMessage(t *testing.T) {
	wxworkApp := NewWxWorkApp(corpID, corpSecret, agentID)
	err := wxworkApp.SendGroupVideoMessage(chatID, mediaID, "人民", "伟大的人民",
		&WxWorkAppMessageSendOptions{Safe: true})
	if err != nil {
		t.Fatal(err)
	}
}

func TestWxWorkApp_SendGroupTextCardMessageC(t *testing.T) {
	wxworkApp := NewWxWorkApp(corpID, corpSecret, agentID)
	err := wxworkApp.SendGroupTextCardMessage(chatID, "人民", "伟大的人民",
		"https://wework.qpic.cn/wwpic/12732_8Z9RVL3rS7-S472_1595229725/0", "看看",
		&WxWorkAppMessageSendOptions{Safe: true})
	if err != nil {
		t.Fatal(err)
	}
}

// must be Safe=false
func TestWxWorkApp_SendGroupNewsMessage(t *testing.T) {
	wxworkApp := NewWxWorkApp(corpID, corpSecret, agentID)
	err := wxworkApp.SendGroupNewsMessage(chatID, []WxWorkAppNewsMessageArticle{newsArticle},
		&WxWorkAppMessageSendOptions{Safe: false})
	if err != nil {
		t.Fatal(err)
	}
}

func TestWxWorkApp_SendGroupMpNewsMessage(t *testing.T) {
	wxworkApp := NewWxWorkApp(corpID, corpSecret, agentID)
	err := wxworkApp.SendGroupMpNewsMessage(chatID, []WxWorkAppMpNewsMessageArticle{mpNewsArticle},
		&WxWorkAppMessageSendOptions{Safe: true})
	if err != nil {
		t.Fatal(err)
	}
}

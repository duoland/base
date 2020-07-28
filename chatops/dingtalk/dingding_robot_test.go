package dingtalk

import (
	"os"
	"testing"
)

var accessToken = os.Getenv("DINGDING_ROBOT_ACCESS_TOKEN")
var secretToken = os.Getenv("DINGDING_ROBOT_SECRET_TOKEN")
var securitySettings = DingDingSecuritySettings{
	AccessToken: accessToken,
	SecureToken: secretToken,
}

var markdownMessage = DingDingRobotMarkdownMessage{
	Text: `Long long ago,  
**this is a mouse**
> he has a long ear!
`,
	Title: "hello, this is a markdown test message",
}

func TestDingDingRobot_SendTextMessage(t *testing.T) {
	ddRobot := NewDingDingRobot()
	if err := ddRobot.SendTextMessage(&securitySettings, "hello, master"); err != nil {
		t.Fatal(err)
	}
}
func TestDingDingRobot_SendTextMessageWithMention(t *testing.T) {
	ddRobot := NewDingDingRobot()
	if err := ddRobot.SendTextMessageWithMention(&securitySettings, "hello, master", []string{"17817213491"}, true); err != nil {
		t.Fatal(err)
	}
}

func TestDingDingRobot_SendMarkdownMessage(t *testing.T) {
	ddRobot := NewDingDingRobot()
	if err := ddRobot.SendMarkdownMessage(&securitySettings, &markdownMessage); err != nil {
		t.Fatal(err)
	}
}
func TestDingDingRobot_SendMarkdownMessageWithMention(t *testing.T) {
	ddRobot := NewDingDingRobot()
	if err := ddRobot.SendMarkdownMessageWithMention(&securitySettings, &markdownMessage, []string{"17817213491"}, false); err != nil {
		t.Fatal(err)
	}
}

func TestDingDingRobot_SendLinkMessage(t *testing.T) {
	ddRobot := NewDingDingRobot()
	linkMsg := DingDingRobotLinkMessage{
		Text:       "bla bla bla, a abstract bla bla bla",
		Title:      "i am a demo blog title",
		PicURL:     "http://res.mail.qq.com/node/ww/wwopenmng/images/independent/doc/test_pic_msg1.png",
		MessageURL: "http://www.oschina.net",
	}
	if err := ddRobot.SendLinkMessage(&securitySettings, &linkMsg); err != nil {
		t.Fatal(err)
	}
}

func TestDingDingRobot_SendActionCardMessage(t *testing.T) {
	ddRobot := NewDingDingRobot()
	actionCardMsg := DingDingRobotActionCardMessage{
		Title: "乔布斯 20 年前想打造一间苹果咖啡厅，而它正是 Apple Store 的前身",
		Text: `![screenshot](https://gw.alicdn.com/tfs/TB1ut3xxbsrBKNjSZFpXXcXhFXa-846-786.png) 
### 乔布斯 20 年前想打造的苹果咖啡厅
Apple Store 的设计正从原来满满的科技感走向生活化，而其生活化的走向其实可以追溯到 20 年前苹果一个建立咖啡馆的计划`,
		ButtonOrientation: DingDingActionCardMessageButtonOrientationVertical,
		SingleTitle:       "Read More>>",
		SingleURL:         "http://www.oschina.net",
	}
	if err := ddRobot.SendActionCardMessage(&securitySettings, &actionCardMsg); err != nil {
		t.Fatal(err)
	}
}

func TestDingDingRobot_SendActionCardMessage2(t *testing.T) {
	ddRobot := NewDingDingRobot()
	actionCardMsg := DingDingRobotActionCardMessage{
		Title: "乔布斯 20 年前想打造一间苹果咖啡厅，而它正是 Apple Store 的前身",
		Text: `![screenshot](https://gw.alicdn.com/tfs/TB1ut3xxbsrBKNjSZFpXXcXhFXa-846-786.png) 
### 乔布斯 20 年前想打造的苹果咖啡厅
Apple Store 的设计正从原来满满的科技感走向生活化，而其生活化的走向其实可以追溯到 20 年前苹果一个建立咖啡馆的计划`,
		ButtonOrientation: DingDingActionCardMessageButtonOrientationHorizontal,
		Buttons: []DingDingRobotActionCardButton{
			DingDingRobotActionCardButton{
				Title:     "内容不错",
				ActionURL: "http://oschina.net",
			},
			DingDingRobotActionCardButton{
				Title:     "不感兴趣",
				ActionURL: "http://www.baidu.com",
			},
		},
	}
	if err := ddRobot.SendActionCardMessage(&securitySettings, &actionCardMsg); err != nil {
		t.Fatal(err)
	}
}

func TestDingDingRobot_SendFeedCardMessage(t *testing.T) {
	ddRobot := NewDingDingRobot()
	feedCardMsgs := []DingDingRobotFeedCardMessage{
		{
			Title:      "时代的火车向前开",
			MessageURL: "http://www.oschina.net",
			PicURL:     "https://gw.alicdn.com/tfs/TB1ut3xxbsrBKNjSZFpXXcXhFXa-846-786.png",
		},
		{
			Title:      "时代的火车向前开2",
			MessageURL: "http://www.baidu.com",
			PicURL:     "https://gw.alicdn.com/tfs/TB1ut3xxbsrBKNjSZFpXXcXhFXa-846-786.png",
		},
	}
	if err := ddRobot.SendFeedCardMessage(&securitySettings, feedCardMsgs); err != nil {
		t.Fatal(err)
	}
}

package wechat

import (
	"encoding/base64"
	"os"
	"testing"
)

var key = os.Getenv("WXWORK_ROBOT_KEY")
var imageBase64Data = "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"

func TestWxWorkRobot_SendTextMessage(t *testing.T) {
	wxRobot := NewWxWorkRobot()
	if err := wxRobot.SendTextMessage(key, "hello, master"); err != nil {
		t.Fatal(err)
	}
}

func TestWxWorkRobot_SendMarkdownMessage(t *testing.T) {
	wxRobot := NewWxWorkRobot()
	content := `
# profile
## i am a robot
> i come, i see
`
	if err := wxRobot.SendMarkdownMessage(key, content); err != nil {
		t.Fatal(err)
	}
}

func TestWxWorkRobot_SendImageMessage(t *testing.T) {
	wxRobot := NewWxWorkRobot()
	imageData, _ := base64.StdEncoding.DecodeString(imageBase64Data)
	if err := wxRobot.SendImageMessage(key, imageData); err != nil {
		t.Fatal(err)
	}
}

func TestWxWorkRobot_SendNewsMessage(t *testing.T) {
	wxRobot := NewWxWorkRobot()
	articles := []WxWorkRobotNewsMessageArticle{
		WxWorkRobotNewsMessageArticle{
			Title:       "i am a demo blog title",
			Description: "i serve my master",
			URL:         "http://www.oschina.net",
			PictureURL:  "http://res.mail.qq.com/node/ww/wwopenmng/images/independent/doc/test_pic_msg1.png",
		},
	}
	if err := wxRobot.SendNewsMessage(key, articles); err != nil {
		t.Fatal(err)
	}
}

func TestWxWorkRobot_UploadFile(t *testing.T) {
	wxRobot := NewWxWorkRobot()
	fileBody := "hello, my master!"
	fileName := "hello.txt"
	if mediaID, createdAt, err := wxRobot.UploadFile(key, []byte(fileBody), fileName); err != nil {
		t.Fatal(err)
	} else {
		t.Logf("media id: %s, created at: %d", mediaID, createdAt)
	}
}

func TestWxWorkRobot_SendFileMessage(t *testing.T) {
	mediaID := "3pbvPgCn9jOzu7YmUx0o4BDxJErfgOZY-_DxbWv_m6kA"
	wxRobot := NewWxWorkRobot()
	if err := wxRobot.SendFileMessage(key, mediaID); err != nil {
		t.Fatal(err)
	}
}

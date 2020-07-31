package bytedance

import (
	"os"
	"testing"
)

var feishuShortcutKey = os.Getenv("FEISHU_SHORTCUT_KEY")

func TestNewFeiShuRobot_SendTextMessage(t *testing.T) {
	robot := NewFeiShuRobot()
	title := "this is a robot message"
	content := "great dreams comes from little steps"
	err := robot.SendTextMessage(feishuShortcutKey, title, content)
	if err != nil {
		t.Fatal(err)
		return
	}
}

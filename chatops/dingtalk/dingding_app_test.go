package dingtalk

import (
	"os"
	"testing"
)

var appKey = os.Getenv("DINGDING_APP_KEY")
var appSecret = os.Getenv("DINGDING_APP_SECRET")
var agentID = os.Getenv("DINGDING_APP_AGENT_ID")

func TestDingDingApp_refreshAccessToken(t *testing.T) {
	dingdingApp := NewDingDingApp(appKey, appSecret, agentID)
	err := dingdingApp.refreshAccessToken()
	if err != nil {
		t.Fatalf(err.Error())
	}
}

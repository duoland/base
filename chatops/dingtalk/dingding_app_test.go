package dingtalk

import (
	"testing"
)

const appKey = "dingufoij8zeywx2mtjq"
const appSecret = "4nN2SXn321Sqx3sXPgckMziO8QFay79OT-s1yV1Sae6673Q2zUsAcE7rKmEphsEa"
const agentID = "2302002"

func TestDingDingApp_refreshAccessToken(t *testing.T) {
	dingdingApp := NewDingDingApp(appKey, appSecret, agentID)
	err := dingdingApp.refreshAccessToken()
	if err != nil {
		t.Fatalf(err.Error())
	}
}

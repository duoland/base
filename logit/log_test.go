package logit

import (
	"testing"
)

func TestLogit(t *testing.T) {
	cfg := Config{
		PrintCaller: true,
		PrintJSON:   true,
		PrintStdout: true,
		BuiltinFields: map[string]string{
			"host_name": "parrot01",
		},
	}
	err := InitLogs(&cfg)
	if err != nil {
		t.Fatal(err)
		return
	}
	name := "jemy"
	Logger.Info("what is your name? ", name)
	Logger.Infof("what is your name? %s", name)
	Logger.Infow("what is your name?", "name", name)

	logId := []interface{}{"logid", "a34sb1312432"}
	Logger.Infow("hello world", logId...)

	Debug("hello world")
	Info("hello world")
	Warn("hello world")
	Error("hello world")
	Fatal("hello world") // exit with 1
	Panic("hello world") // print stack
}

package log

import (
	"testing"
)

func TestLog(t *testing.T) {
	InitLog("[ttt]", "", LevelError)
	Debug("1 this is Debug")
	Debugf("%d this is Debugf", 2)
	Info("3 this is Info")
	Infof("%d this is Infof", 4)
	Error("5 this is Error")
	Errorf("%d this is Errorf", 6)
	//	Fatal("7 this is Fatal")
	//	Fatalf("%d this is Fatalf", 8)
}

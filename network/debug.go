package network

import (
	"github.com/sirupsen/logrus"
)

var (
	DeLog *logrus.Logger = logrus.New()
	Log   *logrus.Logger = logrus.New()
)

const (
	INFOPREFIX = "[=== INFO]"
	WARNPREFIX = "[=== WARN]"
)

func init() {
	//DeLog = logrus.New()
	DeLog.SetLevel(logrus.DebugLevel)
	//Log = logrus.New()
	Log.SetLevel(logrus.ErrorLevel)
}

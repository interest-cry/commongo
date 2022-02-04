package network

import (
	"github.com/sirupsen/logrus"
)

var (
	DeLog *logrus.Logger
	Log   *logrus.Logger
)

const (
	INFOPREFIX = "[=== INFO]"
)

func init() {
	DeLog = logrus.New()
	DeLog.SetLevel(logrus.DebugLevel)
	Log = logrus.New()
	Log.SetLevel(logrus.ErrorLevel)
}

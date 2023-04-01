package logs

import "github.com/sirupsen/logrus"

func Init() {
	setFormatter(logrus.StandardLogger())
	logrus.SetLevel(logrus.DebugLevel)
}

func setFormatter(logger *logrus.Logger) {
	logger.SetFormatter(&logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "time",
			logrus.FieldKeyLevel: "severity",
			logrus.FieldKeyMsg:   "message",
		},
	})
}

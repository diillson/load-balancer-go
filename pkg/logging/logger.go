package logging

import (
	"github.com/sirupsen/logrus"
)

func InitLogger() {
	// Configurações do logrus para produção.
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.InfoLevel)
}

package logging

import (
	"github.com/sirupsen/logrus"
	"log"
)

type LogrusAdapter struct {
	logger *logrus.Logger
}

func (l *LogrusAdapter) Print(v ...interface{}) {
	l.logger.Info(v...)
}

func (l *LogrusAdapter) Printf(format string, v ...interface{}) {
	l.logger.Infof(format, v...)
}

func (l *LogrusAdapter) Println(v ...interface{}) {
	l.logger.Infoln(v...)
}

func (l *LogrusAdapter) Write(p []byte) (n int, err error) {
	l.logger.Info(string(p))
	return len(p), nil
}

func InitLogger() {
	// Configurações do logrus para produção.
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.InfoLevel)
}

// Esta função retorna um logger adaptado para ser usado com pacotes que esperam um log padrão.
func GetLogrusAdapter() *log.Logger {
	return log.New(&LogrusAdapter{logger: logrus.StandardLogger()}, "", 0)
}

package main

import (
	"github.com/diillson/load-balancer-go/api"
	"github.com/diillson/load-balancer-go/internal/loadbalancer"
	"github.com/diillson/load-balancer-go/pkg/config"
	"github.com/diillson/load-balancer-go/pkg/logging"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/url"
)

func main() {
	logging.InitLogger()

	if err := config.LoadConfig("."); err != nil {
		logrus.Warn("No configuration file found. Continuing without it:", err)
	}

	serversConfig := viper.GetStringSlice("servers")

	var servers []*loadbalancer.Server
	for _, s := range serversConfig {
		u, err := url.Parse(s)
		if err != nil {
			logrus.Fatal("Invalid server URL:", s)
		}
		servers = append(servers, &loadbalancer.Server{URL: u})
	}

	lb := &loadbalancer.SimpleLoadBalancer{Servers: servers}
	lb.Initialize()
	handler := &api.Handler{LB: lb}

	r := gin.Default() // Inicializa o router do Gin

	r.GET("/", handler.ProxyHandler)                       // Rota para o proxy
	r.GET("/list", handler.ListServersHandler)             // Rota para listar servidores
	r.POST("/addServer", handler.AddServerHandler)         // Rota para adicionar servidores
	r.DELETE("/removeServer", handler.RemoveServerHandler) // Rota para remover servidores (usando DELETE para ser mais RESTful)
	r.GET("/health", handler.HealthCheckHandler)           // Rota para checar a saúde do servidor

	logrus.Info("Load Balancer running on :3000")
	if err := r.Run(":3000"); err != nil {
		logrus.Fatal("Failed to start server:", err)
	}
}

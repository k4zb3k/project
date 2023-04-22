package main

import (
	"github.com/gin-gonic/gin"
	"github.com/k4zb3k/project/config"
	"github.com/k4zb3k/project/internal/db"
	"github.com/k4zb3k/project/internal/handler"
	"github.com/k4zb3k/project/internal/repository"
	"github.com/k4zb3k/project/internal/service"
	"github.com/k4zb3k/project/pkg/logger"
	"github.com/k4zb3k/project/pkg/redis"
	"github.com/k4zb3k/project/utils"
	"net"
)

func main() {
	router := gin.Default()

	utils.PutAdditionalSettings()
	logger.Init()

	cfg := config.GetConfig()
	logger.Info.Println(cfg)

	redisClient, err := redis.InitRedis(cfg.CacheConn)
	if err != nil {
		logger.Error.Println("failed to connect Redis: ", err)
		return
	}

	dbConn, err := db.GetDBConnection(cfg.DatabaseConn)
	if err != nil {
		logger.Error.Println("failed to connect DB: ", err)
		return
	}

	newRepository := repository.NewRepository(dbConn)

	newService := service.NewService(newRepository, redisClient)

	newHandler := handler.NewHandler(router, newService)
	newHandler.InitRoutes()

	addr := net.JoinHostPort(cfg.Listen.BindIP, cfg.Listen.Port)

	logger.Error.Fatalln(router.Run(addr))
}

package main

import (
	"apollo/server/internal/cronx"
	mw "apollo/server/internal/middleware"
	"apollo/server/internal/model"
	"apollo/server/internal/router"
	"apollo/server/pkg/app"
	"apollo/server/pkg/config"
	"apollo/server/pkg/data"
	"apollo/server/pkg/validate"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	conf := config.New("configs/config.yaml")

	db := data.Init(conf.MysqlConfig)
	model.InitMigrate(db)

	e := app.Engine()

	e.Logger.SetLevel(log.INFO)

	// 配置 lumberjack 日志切割
	logger := &lumberjack.Logger{
		Filename:   "logs/output.log", // 日志文件名
		MaxSize:    100,               // 每个日志文件的最大大小（MB）
		MaxBackups: 20,                // 保留的旧日志文件数量
		MaxAge:     30,                // 保留旧日志文件的最大天数
		Compress:   true,              // 是否压缩旧日志文件
	}

	// 设置日志输出到文件和控制台
	multiWriter := io.MultiWriter(os.Stdout)
	if config.IsProd() {
		multiWriter = io.MultiWriter(logger)
	}
	e.Logger.SetOutput(multiWriter)

	// 使用 echo 的日志中间件
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Output: multiWriter, // 将日志输出到文件和控制台
	}))

	// 添加CORS中间件以支持跨域请求
	e.Use(middleware.CORS())
	
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())

	e.Validator = validate.NewReqValidator()

	e.Use(middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
		c.Logger().Info("requestBody=", string(reqBody))
		c.Logger().Info("responseBody=", string(resBody))
	}))

	cronx.Init(e.Logger)

	e.Use(mw.HandleErrorMiddleware())

	router.Init(e)

	port := fmt.Sprintf(":%d", conf.WebHttpConfig.Port)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		if err := e.Start(port); err != nil && !errors.Is(err, http.ErrServerClosed) {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	e.Logger.Info("shutting down the server")
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal("Shutdown error", err)
	}
	e.Logger.Info("shutting down successfully")
}

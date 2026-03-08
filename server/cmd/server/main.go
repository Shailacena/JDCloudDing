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
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

type AutoConvertBinder struct {
}

func (b *AutoConvertBinder) Bind(i interface{}, c echo.Context) error {
	rawData := make(map[string]interface{})

	if c.Request().Body != nil {
		if err := json.NewDecoder(c.Request().Body).Decode(&rawData); err == nil {
			newb, _ := json.Marshal(rawData)
			c.Request().Body = io.NopCloser(bytes.NewReader(newb))
		}
	}

	queryParams := c.QueryParams()
	for key, values := range queryParams {
		if len(values) > 0 && rawData[key] == nil {
			rawData[key] = values[0]
		}
	}

	data := rawData

	if len(data) > 0 {
		val := reflect.ValueOf(i).Elem()
		rt := val.Type()
		for i := 0; i < val.NumField(); i++ {
			field := val.Field(i)
			fieldType := rt.Field(i)
			if !field.CanSet() {
				continue
			}

			key := fieldType.Tag.Get("json")
			if rawVal, ok := data[key]; ok && rawVal != nil {
				switch field.Kind() {
				case reflect.String:
					if field.Type() == reflect.TypeOf("") {
						switch v := rawVal.(type) {
						case string:
							field.SetString(v)
						case float64:
							field.SetString(strconv.FormatFloat(v, 'f', -1, 64))
						}
					}
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					switch v := rawVal.(type) {
					case float64:
						field.SetInt(int64(v))
					case string:
						if num, err := strconv.ParseInt(v, 10, 64); err == nil {
							field.SetInt(num)
						}
					}
				case reflect.Float32, reflect.Float64:
					switch v := rawVal.(type) {
					case float64:
						field.SetFloat(v)
					case string:
						if num, err := strconv.ParseFloat(v, 64); err == nil {
							field.SetFloat(num)
						}
					}
				case reflect.Slice:
					if arr, ok := rawVal.([]interface{}); ok {
						elemType := field.Type().Elem()
						slice := reflect.MakeSlice(field.Type(), 0, len(arr))
						for _, item := range arr {
							if itemMap, ok := item.(map[string]interface{}); ok {
								elem := reflect.New(elemType)
								for i := 0; i < elemType.NumField(); i++ {
									structField := elemType.Field(i)
									jsonKey := structField.Tag.Get("json")
									if rawVal, ok := itemMap[jsonKey]; ok && rawVal != nil {
										fieldVal := elem.Elem().Field(i)
										switch structField.Type.Kind() {
										case reflect.String:
											if structField.Type == reflect.TypeOf("") {
												switch v := rawVal.(type) {
												case string:
													fieldVal.SetString(v)
												case float64:
													fieldVal.SetString(strconv.FormatFloat(v, 'f', -1, 64))
												}
											}
										case reflect.Float32, reflect.Float64:
											switch v := rawVal.(type) {
											case float64:
												fieldVal.SetFloat(v)
											case string:
												if num, err := strconv.ParseFloat(v, 64); err == nil {
													fieldVal.SetFloat(num)
												}
											}
										}
									}
								}
								slice = reflect.Append(slice, elem.Elem())
							} else {
								var val reflect.Value
								switch elemType.Kind() {
								case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
									if f, ok := item.(float64); ok {
										val = reflect.New(elemType)
										val.Elem().SetInt(int64(f))
									}
								case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
									if f, ok := item.(float64); ok {
										val = reflect.New(elemType)
										val.Elem().SetUint(uint64(f))
									}
								case reflect.String:
									if s, ok := item.(string); ok {
										val = reflect.New(elemType)
										val.Elem().SetString(s)
									}
								}
								if val.IsValid() {
									slice = reflect.Append(slice, val.Elem())
								}
							}
						}
						field.Set(slice)
					}
				}
			}
		}
	}

	return nil
}

func main() {

	conf := config.New("configs/config.yaml")

	db := data.Init(conf.MysqlConfig)
	model.InitMigrate(db)

	e := app.Engine()

	e.Binder = &AutoConvertBinder{}

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

package main

import (
	"apollo/server/cmd/payment/api"
	"apollo/server/pkg/app"
	"apollo/server/pkg/config"
	"apollo/server/pkg/data"
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
	"strings"
	"time"

	mw "apollo/server/internal/middleware"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	llog "github.com/labstack/gommon/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

type AutoConvertBinder struct{
	rawData map[string]interface{}
}

func (b *AutoConvertBinder) Bind(i interface{}, c echo.Context) error {
	// 1. 先保存原始请求数据
	err0 := b.captureRawData(c);
	if  err0 != nil {
		return err0
	}
	// 2. 执行默认绑定
	err := (&echo.DefaultBinder{}).Bind(i, c);
	if  err != nil {
		if !strings.Contains(err.Error(), "code=400") {
			return err
		}
	}

	if err == nil {
		return nil
	}

	logStructInfo(c, i)

	// 3. 基于原始数据执行反射遍历结构体字段进行类型转换
	val := reflect.ValueOf(i).Elem()
	rt := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := rt.Field(i)
		if !field.CanSet() {
			continue
		}

		// 4. 仅处理未初始化的零值字段
		if !reflect.DeepEqual(field.Interface(), reflect.Zero(field.Type()).Interface()) {
			continue
		}

		// 5. 根据目标类型自动转换
		// 获取字段对应的JSON key
		key := fieldType.Tag.Get("json")
		switch field.Kind() {
		case reflect.String:
			if field.Type() == reflect.TypeOf("") {
				b.convertToString(field, b.rawData[key])
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			b.convertToInt(field, b.rawData[key])
		case reflect.Float32, reflect.Float64:
			b.convertToFloat(field, b.rawData[key])
		}
	}
	// log.Println(i)
	logStructInfo(c, i)
	return nil
}

func (b *AutoConvertBinder) convertToString(field reflect.Value, rawVal interface{}) {
	if rawVal == nil {
		return
	}
	
	switch val := rawVal.(type) {
	case int, int8, int16, int32, int64:
		field.SetString(b.toString(val))
	case float32, float64:
		field.SetString(b.toString(val))
	case bool:
		if val {
			field.SetString("true")
		} else {
			field.SetString("false")
		}
	default:
		// 尝试转换为字符串
		if str := b.toString(val); str != "" {
			field.SetString(str)
		}
	}
}

func (b *AutoConvertBinder) convertToInt(field reflect.Value, rawVal interface{}) {
	if rawVal == nil {
		return
	}
	
	switch val := rawVal.(type) {
	case string:
		if num, err := strconv.ParseInt(val, 10, 64); err == nil {
			field.SetInt(num)
		}
	case float64:
		// JSON 解析后的数字通常是 float64
		field.SetInt(int64(val))
	case float32:
		field.SetInt(int64(val))
	case int, int8, int16, int32, int64:
		field.SetInt(reflect.ValueOf(val).Int())
	case uint, uint8, uint16, uint32, uint64:
		field.SetInt(int64(reflect.ValueOf(val).Uint()))
	}
}

func (b *AutoConvertBinder) convertToFloat(field reflect.Value, rawVal interface{}) {
	if rawVal == nil {
		return
	}
	
	switch val := rawVal.(type) {
	case string:
		if num, err := strconv.ParseFloat(val, 64); err == nil {
			field.SetFloat(num)
		}
	case float64, float32:
		field.SetFloat(reflect.ValueOf(val).Float())
	case int, int8, int16, int32, int64:
		field.SetFloat(float64(reflect.ValueOf(val).Int()))
	case uint, uint8, uint16, uint32, uint64:
		field.SetFloat(float64(reflect.ValueOf(val).Uint()))
	}
}

func (b *AutoConvertBinder) captureRawData(c echo.Context) error {
	data := make(map[string]interface{})
	if c.Request().Body != nil {
		if err := json.NewDecoder(c.Request().Body).Decode(&data); err == nil {
			b.rawData = data
			newb, _ := json.Marshal(data)
			c.Request().Body = io.NopCloser(bytes.NewReader(newb)) // 重置body
		}
	}
	return nil
}

func (b *AutoConvertBinder) toString(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	case int:
		return strconv.Itoa(val)
	case int8:
		return strconv.FormatInt(int64(val), 10)
	case int16:
		return strconv.FormatInt(int64(val), 10)
	case int32:
		return strconv.FormatInt(int64(val), 10)
	case int64:
		return strconv.FormatInt(val, 10)
	case uint:
		return strconv.FormatUint(uint64(val), 10)
	case uint8:
		return strconv.FormatUint(uint64(val), 10)
	case uint16:
		return strconv.FormatUint(uint64(val), 10)
	case uint32:
		return strconv.FormatUint(uint64(val), 10)
	case uint64:
		return strconv.FormatUint(val, 10)
	case float32:
		return strconv.FormatFloat(float64(val), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(val)
	default:
		return ""
	}
}

func logStructInfo(c echo.Context, dest interface{}) {
	val := reflect.ValueOf(dest).Elem()
	typ := val.Type()

	c.Logger().Info("=== 结构体字段信息 ===")
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		c.Logger().Info(fmt.Sprintf("字段: %s 类型: %s 值: %v",
			fieldType.Name,
			field.Type().String(),
			field.Interface()))
	}
}

func main() {
	conf := config.New("configs/config.yaml")

	data.Init(conf.MysqlConfig)

	e := app.Engine()

	// 扩展Bind
	e.Binder = &AutoConvertBinder{}

	e.Logger.SetLevel(llog.INFO)

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
	
	// 添加CORS中间件
    // e.Use(middleware.CORS())

	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())

	e.Use(middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
		c.Logger().Info("requestBody=", string(reqBody))
		c.Logger().Info("responseBody=", string(resBody))
	}))

	e.Use(mw.HandleErrorMiddleware())

	paymentGroup := e.Group("/api")
	merchantOrderGroup := paymentGroup.Group("/order")
	merchantOrderGroup.POST("/create", api.CreateOrder)
	merchantOrderGroup.POST("/query", api.QueryOrder)
	merchantOrderGroup.POST("/query/balance", api.QueryBalance)

	notifyGroup := paymentGroup.Group("/notify")
	notifyGroup.POST("/success", api.NotifySuccess)

	notifyGroup.POST("/agiso", api.AgisoNotify)
	notifyGroup.POST("/anssy", api.AnssyNotify)
	notifyGroup.GET("/anssy/auth", api.AnssyAuthNotify)

	port := fmt.Sprintf(":%d", conf.PaymentHttpConfig.Port)
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

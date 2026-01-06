package common

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"sort"
	"strings"

	"github.com/labstack/echo/v4"
)

type AgisoSign struct {
	appSecret string
	jsonStr   string
	timestamp string
}

func NewAgisoSign(appSecret string, jsonStr string, timestamp string) *AgisoSign {
	return &AgisoSign{
		appSecret: appSecret,
		jsonStr:   jsonStr,
		timestamp: timestamp,
	}
}

func (s *AgisoSign) Check(c echo.Context, sign string) bool {
	newSignStr := fmt.Sprintf("%sjson%stimestamp%s%s", s.appSecret, s.jsonStr, s.timestamp, s.appSecret)

	hash := md5.Sum([]byte(newSignStr))
	newSign := fmt.Sprintf("%X", hash)

	if !strings.EqualFold(sign, newSign) {
		c.Logger().Infof("AgisoSign check: newSign=%s, sign=%s", newSign, sign)
	}

	return strings.EqualFold(sign, newSign)
}

func (s *AgisoSign) Generate(params any) string {
	fields, _ := findFiledOfParams(params)

	sort.Slice(fields, func(i, j int) bool {
		return fields[i].Name < fields[j].Name
	})

	fieldList := lo.Map(fields, func(f Filed, index int) string {
		return fmt.Sprintf("%s%s", f.Name, cast.ToString(f.Value))
	})

	fieldStr := s.appSecret + strings.Join(fieldList, "") + s.appSecret

	firstHash := md5.Sum([]byte(fieldStr))
	hashString := hex.EncodeToString(firstHash[:])
	hashString = strings.ToLower(hashString)

	return hashString
}

package common

import (
	"apollo/server/cmd/payment/pkg/types"
	"apollo/server/pkg/contextx"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

type Sign struct {
	privateKey string
}

func NewSign(privateKey string) *Sign {
	return &Sign{
		privateKey: privateKey,
	}
}

func (s *Sign) Check(c echo.Context, params any) bool {
	fields, sign := findFiledOfParams(params)
	newSign := makeHash(s.privateKey, fields)

	if !strings.EqualFold(newSign, sign) {
		c.Logger().Infof("checkSign: newSign=%s, sign=%s", newSign, sign)
		newSign2 := makeHash2(s.privateKey, fields)
		if !strings.EqualFold(newSign2, sign) {
			c.Logger().Infof("checkSign: newSign2=%s, sign=%s", newSign2, sign)
		} else {
			return true
		}
	} else {
		return true
	}

	return false
}

func (s *Sign) Generate(c contextx.Context, params any) string {
	fields, _ := findFiledOfParams(params)

	return makeHash(s.privateKey, fields)
}

type Filed struct {
	Name  string
	Value any
}

func makeHash(privateKey string, fields []Filed) string {
	sort.Slice(fields, func(i, j int) bool {
		return fields[i].Name < fields[j].Name
	})

	fieldList := lo.Map(fields, func(f Filed, index int) string {
		return fmt.Sprintf("%s=%s", f.Name, cast.ToString(f.Value))
	})

	// 加入SignKey
	fieldList = append(fieldList, fmt.Sprintf("key=%s", privateKey))
	// 组装
	fieldStr := strings.Join(fieldList, "&")

	// md5
	firstHash := md5.Sum([]byte(fieldStr))
	hashString := hex.EncodeToString(firstHash[:])
	hashString = strings.ToUpper(hashString)

	return hashString
}

func makeHash2(privateKey string, fields []Filed) string {
	sort.Slice(fields, func(i, j int) bool {
		return fields[i].Name < fields[j].Name
	})

	fieldList := lo.Map(fields, func(f Filed, index int) string {
		if f.Name == "amount" {
			return fmt.Sprintf("%s=%s", f.Name, fmt.Sprintf("%.2f", f.Value))
		}
		return fmt.Sprintf("%s=%s", f.Name, cast.ToString(f.Value))
	})

	// 加入SignKey
	fieldList = append(fieldList, fmt.Sprintf("key=%s", privateKey))
	// 组装
	fieldStr := strings.Join(fieldList, "&")

	// md5
	firstHash := md5.Sum([]byte(fieldStr))
	hashString := hex.EncodeToString(firstHash[:])
	hashString = strings.ToUpper(hashString)

	return hashString
}

func findFiledOfParams(params any) ([]Filed, string) {
	v := reflect.ValueOf(params)
	t := v.Type()
	fields := make([]Filed, 0, v.NumField())

	var sign string
	for i := 0; i < v.NumField(); i++ {
		name := t.Field(i).Tag.Get("json")
		if name == types.FiledSign {
			sign = v.Field(i).String()
			continue
		}

		keys := strings.Split(name, ",")
		if len(keys) > 0 {
			name = keys[0]
		}

		fields = append(fields, Filed{
			Name:  name,
			Value: v.Field(i).Interface(),
		})
	}

	return fields, sign
}

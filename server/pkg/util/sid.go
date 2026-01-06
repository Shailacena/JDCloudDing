package util

import (
	"github.com/sony/sonyflake"
	"github.com/spf13/cast"
)

var (
	SFlake *Sid
)

type Sid struct {
	sf *sonyflake.Sonyflake
}

func init() {
	SFlake = NewSid()
}

func NewSid() *Sid {
	sf := sonyflake.NewSonyflake(sonyflake.Settings{})
	if sf == nil {
		panic("sonyflake not created")
	}

	return &Sid{sf}
}

func (s Sid) GenString() (string, error) {
	id, err := s.sf.NextID()
	if err != nil {
		return "", err
	}

	return cast.ToString(id), nil
}

func (s Sid) GenUint64() (uint64, error) {
	return s.sf.NextID()
}

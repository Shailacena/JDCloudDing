package timex

import (
	"github.com/golang-module/carbon"
	"time"
)

func GetPRCNowTime() carbon.Carbon {
	return carbon.Now(carbon.PRC)
}

func CreatePRCTimeFromTs(ts int64) carbon.Carbon {
	return carbon.CreateFromTimestamp(ts, carbon.PRC)
}

func GetPRCStartOfToday() carbon.Carbon {
	today := GetPRCNowTime().Format(time.DateOnly)
	return carbon.ParseByLayout(today, time.DateOnly)
}

package util

import (
	// "encoding/json"
	"regexp"
	"strings"
	"time"
)

func FomatNowDate() string {
	now := time.Now()
	return now.Format("2006-01-02")
}

func Stations(stationTetx []byte) map[string]string {
	comp, _ := regexp.Compile("([\u4e00-\u9fa5]+)\\|([A-Z]+)")

	datas := comp.FindAll(stationTetx, -1)

	cityMap2Code := make(map[string]string, len(datas))
	for _, v := range datas {
		temp := strings.Split(string(v), "|")
		cityMap2Code[temp[0]] = temp[1]
	}
	return cityMap2Code
}

func trainList() {

}

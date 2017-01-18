package util

import (
	"encoding/json"
	"io/ioutil"
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

var tl = make(map[string]string)

func TrainList() error {
	bytes, err := ioutil.ReadFile("trainlist.json")

	if err != nil {
		return err
	}

	var v interface{}
	if err := json.Unmarshal(bytes, &v); err != nil {
		return err
	}

	if m, ok := v.(map[string]interface{}); ok {
		for _, endDate := range m {
			for _, trainType := range endDate.(map[string]interface{}) {
				for _, trains := range trainType.([]interface{}) {
					obj := trains.(map[string]interface{})
					stc := strings.Split(obj["station_train_code"].(string), "(")
					tl[stc[0]] = obj["train_no"].(string)
				}
			}
		}
	}

	marshl, err := json.Marshal(tl)

	if err != nil {
		return err
	}

	ioutil.WriteFile("compress.data", marshl, 0644)

	return nil
}

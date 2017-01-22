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

type Station struct {
	Train_no string
	From     string
	To       string
}

var All = make(map[string]map[string]*Station)

// Update
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
		for k, endDate := range m {
			if endDate == nil {
				continue
			}

			var oneday = make(map[string]*Station)

			for _, trainType := range endDate.(map[string]interface{}) {
				for _, trains := range trainType.([]interface{}) {
					obj := trains.(map[string]interface{})
					stc := strings.Split(obj["station_train_code"].(string), "(")
					cityTocity := strings.TrimSuffix(stc[1], ")")

					combain := strings.Split(cityTocity, "-")

					there := combain[0]
					home := combain[1]

					oneday[stc[0]] = &Station{
						Train_no: obj["train_no"].(string),
						From:     there,
						To:       home,
					}
				}
			}

			All[k] = oneday
		}
	}

	marshl, err := json.Marshal(All)

	if err != nil {
		return err
	}

	ioutil.WriteFile("compress.data", marshl, 0644)

	return nil
}

func compare(t1, t2 string) bool {
	time1, _ := time.Parse("2006-01-02", t1)
	time2, _ := time.Parse("2006-01-02", t2)

	if time1.After(time2) {
		return true
	}
	return false
}

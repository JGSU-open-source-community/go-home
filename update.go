package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"strings"
)

const (
	train_list = "train_list"
)

// Update
func TrainList(cmd *Command, args []string) int {
	if err := cleanOldData(station_lat_lgt); err != nil {
		log.Fatal(err)
		return 2
	}
	latitudeAndLongitude()
	// when you update meta data, old data will be clean.
	if err := cleanOldData(train_list); err != nil {
		log.Fatal(err)
		return 2
	}

	v, err := trainlistJson()

	if err != nil {
		log.Fatal(err)
		return 2
	}

	var buffer bytes.Buffer

	if m, ok := v.(map[string]interface{}); ok {
		for day, endDate := range m {
			if endDate == nil {
				continue
			}

			for _, trainType := range endDate.(map[string]interface{}) {
				for k, trains := range trainType.([]interface{}) {
					obj := trains.(map[string]interface{})
					stc := strings.Split(obj["station_train_code"].(string), "(")
					cityTocity := strings.TrimSuffix(stc[1], ")")

					combain := strings.Split(cityTocity, "-")

					code := stc[0]
					trainNo := obj["train_no"].(string)
					there := combain[0]
					home := combain[1]

					batchSQL := `'` + code + `'` + "," +
						`'` + trainNo + `'` + "," +
						`'` + there + `'` + "," +
						`'` + home + `'` + "," +
						`'` + day + `'`

					if k == (len(trainType.([]interface{})) - 1) {
						buffer.WriteString(`(` + batchSQL + `)` + ",")
						combine := strings.TrimRight(buffer.String(), ",")
						sql := `insert into ` + train_list + `
				         	  (code, train_no, there, home, depart_date) values ` + combine

						if err := insert(sql); err != nil {
							log.Fatal(err)
							return 2
						}

						buffer.Reset()
					} else {
						buffer.WriteString(`(` + batchSQL + `)` + ",")
					}
				}
			}
		}
	}
	return 0
}

func trainlistJson() (interface{}, error) {
	client := newClient()

	url := "https://kyfw.12306.cn/otn/resources/js/query/train_list.js?scriptVersion=1.5462"
	resp, err := client.Get(url)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	strs := strings.Split(string(body), "=")

	list := strs[1]

	var v interface{}
	if err := json.Unmarshal([]byte(list), &v); err != nil {
		return nil, err
	}
	return v, nil
}

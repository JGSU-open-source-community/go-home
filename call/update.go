package call

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"strings"
)

type Station struct {
	Train_no string
	From     string
	To       string
}

var All = make(map[string]map[string]*Station)

// Update
func TrainList(cmd *Command, args []string) int {
	client := newClient()

	url := "https://kyfw.12306.cn/otn/resources/js/query/train_list.js?scriptVersion=1.5462"
	resp, err := client.Get(url)

	if err != nil {
		log.Fatal(err)
		return 2
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
		return 2
	}

	strs := strings.Split(string(body), "=")

	list := strs[1]

	var v interface{}
	if err := json.Unmarshal([]byte(list), &v); err != nil {
		log.Fatal(err)
		return 2
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
		log.Fatal(err)
		return 2
	}

	ioutil.WriteFile("compress.data", marshl, 0644)
	return 0
}

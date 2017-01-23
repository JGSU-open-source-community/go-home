package call

import (
	"encoding/json"
	"github.com/olekukonko/tablewriter"
	"io/ioutil"
	"os"
	"testing"
)

// func TestCall(t *testing.T) {
// 	datas := Call("G4775")

// 	t.Log(string(datas))
// }

const (
	start = "\x1b[91m(始)\x1b[0m"
	pass  = "\x1b[93m(过)\x1b[0m"
	end   = "\x1b[92m(终)\x1b[0m"
)

func TestShowLeftTricket(t *testing.T) {
	table := tablewriter.NewWriter(os.Stdout)

	f, err := ioutil.ReadFile("./testlefttricketdata.json")

	if err != nil {
		t.Fatal(err)
	}

	var v interface{}

	if err := json.Unmarshal(f, &v); err != nil {
		t.Fatal(err)
	}

	table.SetHeader([]string{"车次", "出发站", "到达站", "出发时间", "到达时间", "历时", "商务座", "特等座", "一等座", "二等座", "高级软卧", "软卧", "硬卧", "软座", "硬座", "无座", "其他"})
	if m, ok := v.(map[string]interface{}); ok {
		if m["httpstatus"].(float64) == 200 {
			if data, ok := m["data"].([]interface{}); ok {
				for _, queryLeftNewDTO := range data {
					if ql, ok := queryLeftNewDTO.(map[string]interface{}); ok {
						raw := ql["queryLeftNewDTO"]
						detail := raw.(map[string]interface{})

						// 始发站
						start_station_name := detail["start_station_name"]
						// 终点站
						end_station_name := detail["end_station_name"]

						// 车次
						station_train_code := detail["station_train_code"].(string)
						// 出发站
						from_station_name := detail["from_station_name"].(string)

						// 到达站
						to_station_name := detail["to_station_name"].(string)
						// 出发时间

						if start_station_name == from_station_name {
							from_station_name = start + from_station_name
						} else {
							from_station_name = pass + from_station_name
						}

						if end_station_name == to_station_name {
							to_station_name = end + to_station_name
						} else {
							to_station_name = pass + to_station_name
						}

						satrt_time := detail["start_time"].(string)
						// 到达时间
						arrive_time := detail["arrive_time"].(string)
						// 历时
						lishi := detail["lishi"].(string)
						// 商务座
						swz_nun := detail["swz_num"].(string)
						// 特等座
						tz_num := detail["tz_num"].(string)
						// 一等座
						zy_num := detail["zy_num"].(string)
						// 二等座
						ze_num := detail["ze_num"].(string)
						// 高级软卧
						gr_num := detail["gr_num"].(string)
						// 软卧
						rw_num := detail["rw_num"].(string)
						// 硬卧
						yw_num := detail["yw_num"].(string)
						// 软座
						rz_num := detail["rz_num"].(string)
						// 硬座
						yz_num := detail["yz_num"].(string)
						// 无座
						wz_num := detail["wz_num"].(string)
						// 其他
						qt_num := detail["qt_num"].(string)

						row := []string{
							station_train_code,
							from_station_name,
							to_station_name,
							satrt_time,
							arrive_time,
							lishi,
							swz_nun,
							tz_num,
							zy_num,
							ze_num,
							gr_num,
							rw_num,
							yw_num,
							rz_num,
							yz_num,
							wz_num,
							qt_num,
						}
						table.Append(row)
					}
				}
			}
		} else {
			t.Fatal("invalid left tricket message!")
			return
		}
	}
	table.Render()
}

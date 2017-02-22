package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/liyu4/tablewriter"
)

const (
	start   = "\x1b[93m(始)\x1b[0m"
	pass    = "\x1b[91m(过)\x1b[0m"
	end     = "\x1b[92m(终)\x1b[0m"
	newpath = "/src/github.com/JingDa-open-source-community"
)

type Command struct {
	UsageLine string
	Run       func(cmd *Command, args []string) int

	// Flag is a set of flags specific to this command.
	Flag flag.FlagSet
}

func newClient() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return &http.Client{Transport: tr}
}

// get city map to code
func stationName() []byte {

	client := newClient()

	url := "https://kyfw.12306.cn/otn/resources/js/framework/station_name.js?station_version=1.8994"
	resp, err := client.Get(url)

	if err != nil {
		log.Fatal(err)
		os.Exit(2)
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	return body
}

func (c *Command) Name() string {
	name := c.UsageLine
	return name
}

var (
	cmdSchedule = &Command{
		UsageLine: "train",
	}

	cmdLeftTicket = &Command{
		UsageLine: "left",
	}

	cmdUpdate = &Command{
		UsageLine: "update",
	}

	Commands = []*Command{
		cmdLeftTicket,
		cmdSchedule,
		cmdUpdate,
	}

	Train  string
	Date1  string
	Date2  string
	From   string
	To     string
	Update string
)

func init() {
	cmdSchedule.Run = ShowSchedule
	cmdSchedule.Flag.StringVar(&Train, "train", "", "the train number is your will seat")
	cmdSchedule.Flag.StringVar(&Date1, "date1", "", "special depart date when you query train schedule")

	cmdLeftTicket.Run = ShowLeftTicket
	cmdLeftTicket.Flag.StringVar(&Date2, "date2", "", "special depart date when you query left ticket")
	cmdLeftTicket.Flag.StringVar(&To, "to", "", "arrive station")
	cmdLeftTicket.Flag.StringVar(&From, "from", "", "start station")

	cmdUpdate.Run = TrainList
	cmdUpdate.Flag.StringVar(&Update, "update", "", "update basic data")
}

var cityMapToCode = stations(stationName())

func schedule(train, date string) (datas []byte) {

	train = strings.ToUpper(train)

	client := newClient()

	var v interface{}

	execFileRelativePath, _ := exec.LookPath(os.Args[0])

	var newexecFileRelativePath string

	if runtime.GOOS == "windows" {
		execFileRelativePath = strings.TrimSuffix(execFileRelativePath, ".exe")
		newexecFileRelativePath = strings.Replace(execFileRelativePath, "bin", newpath, 1)
		newexecFileRelativePath = newexecFileRelativePath + "\\compress.data"
	} else {

		if !strings.Contains(execFileRelativePath, "./") {
			newexecFileRelativePath = strings.Replace(execFileRelativePath, "bin", newpath, 1)
			newexecFileRelativePath = newexecFileRelativePath + "/compress.data"
		} else {
			newexecFileRelativePath = "compress.data"
		}

	}

	f, err := ioutil.ReadFile(newexecFileRelativePath)

	if err != nil {
		log.Fatal(err)
		os.Exit(2)
	}

	if err := json.Unmarshal(f, &v); err != nil {
		log.Fatal(err)
		os.Exit(2)
	}

	first := v.(map[string]interface{})
	compress := first[date].(map[string]interface{})

	if combain, ok := compress[train].(map[string]interface{}); !ok {
		fmt.Println("请输入正确的车次信息！")
		return nil
	} else {
		no := combain["Train_no"].(string)
		from := cityMapToCode[combain["From"].(string)]
		to := cityMapToCode[combain["To"].(string)]

		url := fmt.Sprintf("https://kyfw.12306.cn/otn/czxx/queryByTrainNo?train_no=%s&from_station_telecode=%s&to_station_telecode=%s&depart_date=%s", no, from, to, date)

		resp, err := client.Get(url)

		if err != nil {
			log.Fatal(err)
			os.Exit(2)
		}

		defer resp.Body.Close()

		datas, err = ioutil.ReadAll(resp.Body)

		if err != nil {
			log.Fatal(err)
			os.Exit(2)
		}
	}
	return datas
}

func ShowSchedule(cmd *Command, args []string) int {
	data := schedule(args[0], args[1])

	if data == nil {
		return 2
	}

	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		log.Fatal(err)
		return 2
	}

	table := tablewriter.NewColorWriter(os.Stdout)
	table.SetHeader([]string{"站序", "站名", "到站时间", "出站时间", "停留时间"})

	if m, ok := v.(map[string]interface{}); ok {
		if m["httpstatus"].(float64) == 200 {
			if subdata, ok := m["data"].(map[string]interface{}); ok {
				if elements, ok := subdata["data"].([]interface{}); ok {
					for _, vv := range elements {
						element := vv.(map[string]interface{})
						station_no := element["station_no"].(string)
						station_name := element["station_name"].(string)
						arrive_time := element["arrive_time"].(string)
						start_time := element["start_time"].(string)
						stopover_time := element["stopover_time"].(string)
						row := []string{station_no, station_name, arrive_time, start_time, stopover_time}
						table.Append(row)
					}
				}
			}
		} else {
			log.Fatal("invalid train schedule message!")
			return 2
		}
	}
	table.Render()
	return 1
}

// query left ticket in 12306
// form start city
// to arrive city
func leftTicket(from, to, date string) []byte {

	fromCode := cityMapToCode[from]
	toCode := cityMapToCode[to]
	url := fmt.Sprintf("https://kyfw.12306.cn/otn/leftTicket/query?leftTicketDTO.train_date=%s&leftTicketDTO.from_station=%s&leftTicketDTO.to_station=%s&purpose_codes=ADULT", date, fromCode, toCode)

	client := newClient()

	resp, err := client.Get(url)

	if err != nil {
		log.Fatal(err)
		os.Exit(2)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
		os.Exit(2)
	}

	return body
}

func ShowLeftTicket(cmd *Command, args []string) int {
	leftTicketDatas := leftTicket(args[0], args[1], args[2])

	var v interface{}

	if err := json.Unmarshal(leftTicketDatas, &v); err != nil {
		log.Fatal(err)
		return 2
	}

	// new table redict to terminal
	table := tablewriter.NewColorWriter(os.Stdout)
	table.SetHeader([]string{"车次", "出发站", "到达站", "出发时间", "到达时间", "历时", "商务座", "特等座", "一等座", "二等座", "高级软卧", "软卧", "硬卧", "软座", "硬座", "无座", "其他"})
	if m, ok := v.(map[string]interface{}); ok {
		if m["status"].(bool) != true {
			fmt.Println("12306接口访问异常")
			return 2
		}
		if m["httpstatus"].(float64) == 200 {
			if data, ok := m["data"].([]interface{}); ok {
				for _, queryLeftNewDTO := range data {
					if ql, ok := queryLeftNewDTO.(map[string]interface{}); ok {
						raw := ql["queryLeftNewDTO"]
						detail := raw.(map[string]interface{})

						// 始发站
						start_station_name := detail["start_station_name"].(string)
						// 终点站
						end_station_name := detail["end_station_name"].(string)

						// 车次
						station_train_code := detail["station_train_code"].(string)
						// 出发站
						from_station_name := detail["from_station_name"].(string)

						// 到达站
						to_station_name := detail["to_station_name"].(string)

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

						// 出发时间
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
			log.Fatal("invalid left tricket message!")
			return 2
		}
	}
	table.Render()
	return 1
}

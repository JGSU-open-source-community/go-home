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
	"sort"
	"strconv"
	"strings"

	"github.com/liyu4/tablewriter"
)

const (
	start   = "\x1b[93m(始)\x1b[0m"
	pass    = "\x1b[91m(过)\x1b[0m"
	end     = "\x1b[92m(终)\x1b[0m"
	newpath = "/src/github.com/JingDa-open-source-community"
)

var top = 10

type Command struct {
	UsageLine string
	Run       func(cmd *Command, args []string) int

	// Flag is a set of flags specific to this command.
	Flag flag.FlagSet
}

type info struct {
	City     string
	Distance float64
}

type newList []*info

func (I newList) Len() int {
	return len(I)
}
func (I newList) Less(i, j int) bool {
	return I[i].Distance < I[j].Distance
}
func (I newList) Swap(i, j int) {
	I[i], I[j] = I[j], I[i]
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

	cmdLeftTicket.Run = ShowTransferPlan
	cmdLeftTicket.Flag.StringVar(&Date2, "date2", "", "special depart date when you query left ticket")
	cmdLeftTicket.Flag.StringVar(&To, "to", "", "arrive station")
	cmdLeftTicket.Flag.StringVar(&From, "from", "", "start station")

	cmdUpdate.Run = TrainList
	cmdUpdate.Flag.StringVar(&Update, "update", "", "update basic data")
}

var cityMapToCode = stations(stationName())

func shortestcity(from, to, date string) *newList {
	sql := `select station, latitude, longitude from  station_lat_lgt where station in (select distinct home from train_list where there like '%` + from + `%' and depart_date='` + date + `')`

	maps, err := query(sql)

	if err != nil {
		log.Fatal(err)
		os.Exit(2)
	}

	l := new(newList)
	for _, v := range *maps {
		station := v["station"].(string)
		latitude, _ := strconv.ParseFloat(v["latitude"].(string), 64)
		longitude, _ := strconv.ParseFloat(v["longitude"].(string), 64)
		dis := earthDistance(29.0200802067, 115.8154807999, latitude, longitude)

		info := &info{
			City:     station,
			Distance: dis,
		}
		*l = append(*l, info)
	}

	sort.Sort(newList(*l))
	return l
}

func ShowTransferPlan(cmd *Command, args []string) int {

	ret, err := isThrougt(args[0], args[1], args[2])

	if ret {
		// through from start station to your hometown
		return ShowLeftTicket(cmd, args)
	} else {
		if err != nil {
			log.Fatal(err)
			return 2
		}
	}

	// Can't go through
	l := shortestcity(args[0], args[1], args[2])

	ch := make(chan string, 1)
	for i := 0; i < top; i++ {
		// revsert query
		if *l == nil {
			l = shortestcity(args[1], args[0], args[2])
			for i := 0; i < top; i++ {
				if *l != nil {
					topCity := (*l)[i].City
					ret, err := isThrougt(topCity, args[0], args[2])
					if err != nil {
						log.Fatal(err)
					}

					if ret {
						ch <- topCity
						break
					}
				}
			}
			break
		} else {
			topCity := (*l)[i].City
			ret, err := isThrougt(topCity, args[1], args[2])
			if err != nil {
				log.Fatal(err)
			}

			if ret {
				ch <- topCity
				break
			}
		}
	}

	city := <-ch
	args1 := []string{args[0], city, args[2]}
	fmt.Printf("===================到达中转站-%s=======================\n", city)
	ShowLeftTicket(cmd, args1)
	args2 := []string{city, args[1], args[2]}
	fmt.Printf("===================中转站出发-%s=======================\n", city)
	ShowLeftTicket(cmd, args2)
	return 1
}

func isThrougt(from, to, date string) (bool, error) {
	data := leftTicket(from, to, date)
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return false, err
	}

	if m, ok := v.(map[string]interface{}); ok {
		// in interface nil maybe not equal nil
		if len(m["data"].(map[string]interface{})) == 0 {
			return false, nil
		}
	}
	return true, nil
}

func schedule(train, date string) (data []byte) {

	train = strings.ToUpper(train)

	sql := `select train_no, there, home from ` + train_list + ` where depart_date='` + date + `' and code='` + train + `'`
	maps, err := query(sql)
	if err != nil {
		log.Fatal(err)
		os.Exit(2)
	}
	client := newClient()

	if *maps == nil {
		fmt.Println("请输入正确的车次信息或者出发日期!")
		return nil
	}

	for _, v := range *maps {
		no := v["train_no"].(string)
		there := cityMapToCode[v["there"].(string)]
		home := cityMapToCode[v["home"].(string)]
		url := fmt.Sprintf("https://kyfw.12306.cn/otn/czxx/queryByTrainNo?train_no=%s&from_station_telecode=%s&to_station_telecode=%s&depart_date=%s", no, there, home, date)
		fmt.Println(url)
		resp, err := client.Get(url)

		if err != nil {
			log.Fatal(err)
			os.Exit(2)
		}

		defer resp.Body.Close()

		data, err = ioutil.ReadAll(resp.Body)

		if err != nil {
			log.Fatal(err)
			os.Exit(2)
		}
	}
	return data
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

func b4(ct []interface{}, cv map[string]interface{}) []interface{} {
	var cs []interface{}

	for cr := 0; cr < len(ct); cr = cr + 1 {
		//cw := make(map[string]interface{})

		//var cq = ct[cr].split("|")
		cq := strings.Split(ct[cr].(string), "|")

		//cw.secretHBStr = cq[36]
		//cw.secretStr = cq[0]
		//cw.buttonTextInfo = cq[1]

		cu := make(map[string]interface{})
		cu["train_no"] = cq[2]
		cu["station_train_code"] = cq[3]
		cu["start_station_telecode"] = cq[4]
		cu["end_station_telecode"] = cq[5]

		if cv[cq[4]] == nil {
			cu["start_station_name"] = " "
		} else {
			cu["start_station_name"] = cv[cq[4]]
		}

		if cv[cq[5]] == nil {
			cu["end_station_name"] = " "
		} else {
			cu["end_station_name"] = cv[cq[5]]
		}

		cu["from_station_telecode"] = cq[6]
		cu["to_station_telecode"] = cq[7]
		cu["start_time"] = cq[8]
		cu["arrive_time"] = cq[9]
		cu["lishi"] = cq[10]
		cu["canWebBuy"] = cq[11]
		cu["yp_info"] = cq[12]
		cu["start_train_date"] = cq[13]
		cu["train_seat_feature"] = cq[14]
		cu["location_code"] = cq[15]
		cu["from_station_no"] = cq[16]
		cu["to_station_no"] = cq[17]
		cu["is_support_card"] = cq[18]
		cu["controlled_train_flag"] = cq[19]
		//cu["gg_num"] = len(cq[20])!=0 ? cq[20] : "--"
		//cu["gr_num"] = len(cq[21])!=0 ? cq[21] : "--"
		//cu["qt_num"] = len(cq[22])!=0 ? cq[22] : "--"
		//cu["rw_num"] = len(cq[23])!=0 ? cq[23] : "--"
		//cu["rz_num"] = len(cq[24])!=0 ? cq[24] : "--"
		//cu["tz_num"] = len(cq[25])!=0 ? cq[25] : "--"
		//cu["wz_num"] = len(cq[26])!=0 ? cq[26] : "--"
		//cu["yb_num"] = len(cq[27])!=0 ? cq[27] : "--"
		//cu["yw_num"] = len(cq[28])!=0 ? cq[28] : "--"
		//cu["yz_num"] = len(cq[29])!=0 ? cq[29] : "--"
		//cu["ze_num"] = len(cq[30])!=0 ? cq[30] : "--"
		//cu["zy_num"] = len(cq[31])!=0 ? cq[31] : "--"
		//cu["swz_num"] = len(cq[32])!=0 ? cq[32] : "--"
		//cu["srrb_num"] = len(cq[33])!=0 ? cq[33] : "--"
		cu["gg_num"] = cq[20]
		cu["gr_num"] = cq[21]
		cu["qt_num"] = cq[22]
		cu["rw_num"] = cq[23]
		cu["rz_num"] = cq[24]
		cu["tz_num"] = cq[25]
		cu["wz_num"] = cq[26]
		cu["yb_num"] = cq[27]
		cu["yw_num"] = cq[28]
		cu["yz_num"] = cq[29]
		cu["ze_num"] = cq[30]
		cu["zy_num"] = cq[31]
		cu["swz_num"] = cq[32]
		cu["srrb_num"] = cq[33]
		cu["yp_ex"] = cq[34]
		cu["seat_types"] = cq[35]
		cu["exchange_train_flag"] = cq[36]

		cu["from_station_name"] = cv[cq[6]]
		cu["to_station_name"] = cv[cq[7]]
		//cu["from_station_name"] = cq[6]
		//cu["to_station_name"] = cq[7]
		//start_station_name

		//cw.queryLeftNewDTO = cu
		//cs.push(cw)
		cs = append(cs, cu)
	}
	return cs
}

// query left ticket in 12306
// form start city
// to arrive city
func leftTicket(from, to, date string) []byte {

	fromCode := cityMapToCode[from]
	toCode := cityMapToCode[to]
	url := fmt.Sprintf("https://kyfw.12306.cn/otn/leftTicket/queryX?leftTicketDTO.train_date=%s&leftTicketDTO.from_station=%s&leftTicketDTO.to_station=%s&purpose_codes=ADULT", date, fromCode, toCode)
	client := newClient()
	fmt.Println(url)

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
	leftTicketData := leftTicket(args[0], args[1], args[2])

	var v interface{}

	if err := json.Unmarshal(leftTicketData, &v); err != nil {
		log.Fatal(err)
		return 2
	}
	return readerTable(v)
}

func readerTable(v interface{}) int {
	// new table redict to terminal
	table := tablewriter.NewColorWriter(os.Stdout)
	table.SetHeader([]string{"车次", "出发站", "到达站", "出发时间", "到达时间", "历时", "商务座", "特等座", "一等座", "二等座", "高级软卧", "软卧", "硬卧", "软座", "硬座", "无座", "其他"})
	if m, ok := v.(map[string]interface{}); ok {
		if m["status"].(bool) != true {
			fmt.Println("12306接口访问异常")
			return 2
		}
		if m["httpstatus"].(float64) == 200 {
			if data, ok := m["data"].(map[string]interface{}); ok {
				if resdata, ok := data["result"].([]interface{}); ok {
					mapcity, _ := data["map"].(map[string]interface{})
					data2 := b4(resdata, mapcity)
					for _, queryLeftNewDTO := range data2 {
						if detail, ok := queryLeftNewDTO.(map[string]interface{}); ok {
							//raw := ql["queryLeftNewDTO"]
							//detail := raw.(map[string]interface{})

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

			}
		} else {
			log.Fatal("invalid left tricket message!")
			return 2
		}
	}
	table.Render()
	return 1
}

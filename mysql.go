package main

import (
	"encoding/csv"
	"log"
	"os"
	"strings"

	"github.com/astaxie/beego/orm"
)

// func writeRawDataToMysql() {
// 	staions := allChinaRailwayStations(stationName())
// 	data := []byte(strings.Join(staions, "站\n"))
// 	ioutil.WriteFile("stations.txt", data, 0644)
// }

const (
	station_lat_lgt = "station_lat_lgt"
)

func latitudeAndLongitude() {
	// Load a csv file
	f, err := os.Open("lt.csv")
	if err != nil {
		log.Fatal(err)
		os.Exit(2)
	}

	defer f.Close()

	// Create a new reader
	r := csv.NewReader(f)

	// A successful call returns err == nil this is different with r.Read()
	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
		os.Exit(2)
	}

	for index, singleRecord := range records {
		if index == 0 {
			continue
		}

		sql := `insert into  ` + station_lat_lgt + ` (station, latitude, longitude)` +
			`values('` + strings.TrimSuffix(singleRecord[8], "站") + `', '` + singleRecord[2] + `', '` + singleRecord[3] + `')
			`
		if err := insert(sql); err != nil {
			log.Println(err)
		}
	}
}

func insert(sql string) error {
	o := orm.NewOrm()
	if _, err := o.Raw(sql).Exec(); err != nil {
		return err
	}
	return nil
}

func query(sql string) (*[]orm.Params, error) {
	o := orm.NewOrm()

	var maps []orm.Params
	if _, err := o.Raw(sql).Values(&maps); err != nil {
		return nil, err
	}
	return &maps, nil
}

func cleanOldData(tablename string) error {
	o := orm.NewOrm()
	sql := `truncate table ` + tablename + ``
	if _, err := o.Raw(sql).Exec(); err != nil {
		return err
	}
	return nil
}

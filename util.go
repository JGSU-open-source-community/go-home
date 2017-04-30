package main

import (
	"math"
	"regexp"
	"strings"
	"time"
)

func FomatNowDate() string {
	now := time.Now()
	return now.Format("2006-01-02")
}

func stations(stationTetx []byte) map[string]string {
	comp, _ := regexp.Compile("([\u4e00-\u9fa5]+)\\|([A-Z]+)")

	data := comp.FindAll(stationTetx, -1)

	cityMap2Code := make(map[string]string, len(data))
	for _, v := range datas {
		temp := strings.Split(string(v), "|")
		cityMap2Code[temp[0]] = temp[1]
	}
	return cityMap2Code
}

func allChinaRailwayStations(stationTetx []byte) []string {
	comp, _ := regexp.Compile("([\u4e00-\u9fa5]+)")

	data := comp.FindAll(stationTetx, -1)

	crs := make([]string, 0, len(data))
	for _, v := range data {
		crs = append(crs, string(v))
	}
	return crs
}

func compare(t1, t2 string) bool {
	time1, _ := time.Parse("2006-01-02", t1)
	time2, _ := time.Parse("2006-01-02", t2)

	if time1.After(time2) {
		return true
	}
	return false
}

// 根据经纬度计算两地距离
func earthDistance(lat1, lng1, lat2, lng2 float64) float64 {
	var radius float64 = 6371000 // 6378137
	rad := math.Pi / 180.0

	lat1 = lat1 * rad
	lng1 = lng1 * rad
	lat2 = lat2 * rad
	lng2 = lng2 * rad

	theta := lng2 - lng1
	dist := math.Acos(math.Sin(lat1)*math.Sin(lat2) + math.Cos(lat1)*math.Cos(lat2)*math.Cos(theta))

	return dist * radius
}

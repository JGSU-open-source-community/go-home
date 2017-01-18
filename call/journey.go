package call

import (
	"crypto/tls"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	// "go-home/util"
)

func Call() (datas []byte) {

	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}

	client := &http.Client{Transport: tr}

	// date := util.FomatNowDate()

	resp, err := client.Get(`https://kyfw.12306.cn/otn/resources/js/query/train_list.js?scriptVersion=1.8994`)

	if err != nil {
		log.Fatal(err)
		os.Exit(2)
	}

	datas, err = ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
		os.Exit(2)
	}
	return datas
}

func CallAllStation() {

}

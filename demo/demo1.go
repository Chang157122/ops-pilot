package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/xuri/excelize/v2"
	"io"
	"net/http"
	url2 "net/url"
	"strconv"
	"time"
)

type MetricsData struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Metric struct {
			} `json:"metric"`
			Value []interface{} `json:"value"`
		} `json:"result"`
	} `json:"data"`
}

var (
	PromNodeMemQL = "sum((1 - (node_memory_MemAvailable_bytes / (node_memory_MemTotal_bytes)))* 100) /count(node_memory_MemAvailable_bytes)"
	PromNodeCPUQL = "sum(100 - (avg(irate(node_cpu_seconds_total{mode=\"idle\"}[5m])) by (instance) * 100)) / count(100 - (avg(irate(node_cpu_seconds_total{mode=\"idle\"}[5m])) by(instance) * 100))"
	PromUrl       string
	excelFileName = "site.xlsx"
	excelSheet    = "巡检"
)

func initParams() {
	flag.StringVar(&PromUrl, "prom", "http://10.0.20.180:9090", "prom接口地址")
	flag.Parse()
}

func main() {
	// 初始化参数
	initParams()

	nowDateTime := time.Now().Format(time.DateTime)
	cpu := getCPU()
	mem := getMEM()

	// 读取excel表格
	func() {

		f, err := excelize.OpenFile(excelFileName)
		if err != nil {
			panic(err)
		}
		f.SetCellValue(excelSheet, "A2", nowDateTime)
		f.SetCellValue(excelSheet, "E2", mem+"%")
		f.SetCellValue(excelSheet, "E3", cpu+"%")
		if err = f.Save(); err != nil {
			panic(err)
		}
	}()

}

func getCPU() string {
	parm := url2.QueryEscape(PromNodeCPUQL)
	url := PromUrl + fmt.Sprintf("/api/v1/query?query=%s", parm)
	resp := okHttp(url)
	f, err := strconv.ParseFloat(resp.Data.Result[0].Value[1].(string), 64)
	if err != nil {
		panic(err)
	}
	data := fmt.Sprintf("%.2f", f)
	return data
}

func getMEM() string {
	parm := url2.QueryEscape(PromNodeMemQL)
	url := PromUrl + fmt.Sprintf("/api/v1/query?query=%s", parm)
	resp := okHttp(url)
	f, err := strconv.ParseFloat(resp.Data.Result[0].Value[1].(string), 64)
	if err != nil {
		panic(err)
	}
	data := fmt.Sprintf("%.2f", f)
	return data
}

func okHttp(url string) MetricsData {
	var data MetricsData
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	response, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(response, &data); err != nil {
		panic(err)
	}
	return data
}

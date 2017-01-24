package main

import (
	"fmt"
	"net/http"
	"flag"
	"net/url"
	"strings"
	"io/ioutil"
	"encoding/json"
)

var opType, from, to, date, page *string
var state, t *int

type Result struct {
	Success bool
	Content Content
}

type Content struct {
	BusNumberList []BusNumber
}

type BusNumber struct {
	BeginStationName string
	LeaveTime string
	RemainSeat string
}

func init()  {
	opType = flag.String("type", "ticket", "please input your opeartion type!")
	from = flag.String("from", "杭州", "please input from city")
	to = flag.String("to", "绍兴", "please input to city")
	date = flag.String("date", "2017-01-25", "please input when to go")
	page = flag.String("page", "1", "please input search page")
	state = flag.Int("state", 1,"please choose which begin station")
	t = flag.Int("time", 0,"please choose which begin station")
	flag.Parse()
}

func main()  {
	v := url.Values{}
	v.Set("type", *opType)
	v.Set("from", *from)
	v.Set("to", *to)
	v.Set("date", *date)
	v.Set("page", *page)
	if *state == 1 {
		v.Set("startStations", "010106")
	}
	switch *t {
		case 1 : v.Set("leaveTimes", "00:00-06:00")
		case 2 : v.Set("leaveTimes", "06:00-12:00")
		case 3 : v.Set("leaveTimes", "12:00-18:00")
		case 4 : v.Set("leaveTimes", "18:00-24:00")
	}
	body := strings.NewReader(v.Encode())

	client := http.DefaultClient
	req, err := http.NewRequest("POST", "http://www.bababus.com/ticket/ticketList.htm", body)
	if err != nil {
		fmt.Printf("无法访问该网址")
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Host", "www.bababus.com")
	req.Header.Set("Origin", "http://www.bababus.com")
	req.Header.Set("Referer", "http://www.bababus.com/ticket/searchbus.htm")

	resp, er := client.Do(req)
	if er != nil {
		fmt.Printf("不能执行此请求")
		return
	}

	ans, _ := ioutil.ReadAll(resp.Body)
	result := &Result{}
	json.Unmarshal(ans, &result)

	if result.Success == true {
		busList := result.Content.BusNumberList
		if len(busList) == 0 {
			fmt.Printf("no any record!")
			return
		}
		for i := 0 ; i < len(busList); i ++ {
			bus := busList[i]
			fmt.Printf(bus.BeginStationName + "(" + bus.LeaveTime + "): " + bus.RemainSeat + "张\n")
		}
	} else {
		fmt.Printf("请检查你是否在短时间内执行次数过多, 需要登录网站输入图形校验码(待改进 todo)")
	}
}



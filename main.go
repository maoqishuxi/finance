package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

func fetchUrl(code string, page int) ([]map[string][]byte, error) {
	url := "https://fundf10.eastmoney.com/F10DataApi.aspx?type=lsjz&code=%s&page=%d&per=20"
	response, err := http.Get(fmt.Sprintf(url, code, page))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	express := `<td[a-z\s=\']*?>([0-9\.%-]+)</td>`
	data := regexp.MustCompile(express).FindAllSubmatch(body, -1)

	result := make([]map[string][]byte, 0)

	for i := 0; i < len(data); i += 4 {
		result = append(result, map[string][]byte{
			"date":                 data[i][1],
			"unit_net_value":       data[i+1][1],
			"cumulative_net_value": data[i+2][1],
			"daily_growth_rate":    data[i+3][1],
		})
	}

	return result, nil
}

func getData() {
	code := "000942"
	page := 1
	for i := 0; i < 5; i++ {
		result, err := fetchUrl(code, page)
		if err == nil {
			for _, v := range result {
				fmt.Println(string(v["date"]))
				fmt.Println(string(v["unit_net_value"]))
				fmt.Println(string(v["cumulative_net_value"]))
				fmt.Println(string(v["daily_growth_rate"]))
			}
			break
		}

	}
}

func main() {
	fmt.Println("welcome to you")
	getData()
}

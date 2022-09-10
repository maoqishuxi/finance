package main

import (
	"context"
	"database/sql"
	"finance/writeToDatabase"
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
			"time":   data[i][1],
			"valueV": data[i+1][1],
			"valueC": data[i+2][1],
			"gain":   data[i+3][1],
		})
	}

	return result, nil
}

func getData(code string, page int) []map[string][]byte {
	var result []map[string][]byte
	for i := 0; i < 5; i++ {
		ret, err := fetchUrl(code, page)
		if err == nil {
			result = ret
			break
		}
		log.Println("request over timeout")
	}

	return result
}

func main() {
	//code := "000942"
	//page := 1

	fmt.Println("welcome to you")

	//data := getData(code, page)
	db, err := sql.Open("sqlite3", "./data/file.db")
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	var ctx context.Context
	var data = map[string][]byte{
		"id":     []byte("0"),
		"time":   []byte("2022-09-10"),
		"valueV": []byte("2"),
		"valueC": []byte("2"),
		"gain":   []byte("1.0"),
	}

	//writeToDatabase.CreateTable(db, "file")
	writeToDatabase.InsertData(db, ctx, data)
}

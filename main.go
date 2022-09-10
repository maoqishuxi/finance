package main

import (
	"database/sql"
	"finance/writeToDatabase"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

func fetchUrl(code string, page int) ([]writeToDatabase.Item, error) {
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

	result := make([]writeToDatabase.Item, 0)

	for i := 0; i < len(data); i += 4 {
		result = append(result, writeToDatabase.Item{
			Id:     i / 4,
			Time:   data[i][1],
			ValueV: data[i+1][1],
			ValueC: data[i+2][1],
			Gain:   data[i+3][1],
		})
	}

	return result, nil
}

func getData(code string, page int) []writeToDatabase.Item {
	var result []writeToDatabase.Item
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
	db, err := sql.Open("sqlite3", "./data/file.db")
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	//data := getData(code, page)
	//writeToDatabase.CreateTable(db, "file")
	//for _, item := range data {
	//	writeToDatabase.InsertData(db, item)
	//}

	//result := writeToDatabase.QueryData(db, 10)
	//for _, item := range result {
	//	fmt.Println(item.Id)
	//	fmt.Println(string(item.ValueV))
	//	fmt.Println(string(item.Gain))
	//}
	//id := writeToDatabase.QueryID(db)
	//fmt.Println(id)
}

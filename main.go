package main

import (
	"database/sql"
	"finance/writeToDatabase"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"time"
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
		log.Println("getData request over timeout")
	}

	return result
}

func UpdateData(db *sql.DB, code string, page int, tableName string) {
	for {
		now := time.Now()
		next := now.Add(time.Hour * 24)
		hour := 14
		minute := 30
		second := 00
		nSecond := 00

		next = time.Date(next.Year(), next.Month(), next.Day(), hour, minute, second, nSecond, next.Location())
		t := time.NewTimer(next.Sub(now))
		<-t.C
		for range t.C {
			data := getData(code, page)
			lastID := writeToDatabase.QueryID(db, tableName)
			id := writeToDatabase.InsertData(db, tableName, data[0], lastID+1)
			log.Printf("插入%v数据成功\n", id)
		}
	}
}

func InitData(db *sql.DB, tableName string, code string) {
	writeToDatabase.CreateTable(db, tableName)
	lastID := writeToDatabase.QueryID(db, tableName)

	page := 3
	var id int64 = 0
	for i := 1; i <= page; i++ {
		data := getData(code, i)
		for _, item := range data {
			id = writeToDatabase.InsertData(db, tableName, item, lastID+id)
			if id == -1 {
				break
			}
			fmt.Println(id)
		}
	}

}

func tradeData(gain float64) {

}

func Average60(db *sql.DB, tableName string) (float64, float64) {
	result := writeToDatabase.QueryData(db, tableName, 60)
	var sumValue, sumGain float64
	for _, item := range result {
		value, err := strconv.ParseFloat(string(item.ValueV), 64)
		gain, err := strconv.ParseFloat(string(item.Gain)[:len(string(item.Gain))-1], 64)
		if err != nil {
			log.Println(err)
			return -1, -1
		}

		sumValue += value
		sumGain += math.Abs(gain)
	}

	return sumValue / 60, sumGain * 3 / 100 / 60
}

func main() {
	code := "159938"
	db, err := sql.Open("sqlite3", "./data/file.db")
	if err != nil {
		log.Println(err)
	}
	defer db.Close()
	tableName := "A" + code

	avgValue, avgGain := Average60(db, tableName)
	fmt.Println(avgValue, avgGain)

}

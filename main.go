package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type Ticker struct {
	Ask       string `json:"ask"`
	Bid       string `json:"bid"`
	High      string `json:"high"`
	LastPrice string `json:"last_price"`
	Low       string `json:"low"`
	Market    string `json:"market"`
	Timestamp string `json:"timestamp"`
	Volume    string `json:"volume"`
}

type ApiResp struct {
	Data   []Ticker `json:"data"`
	Status string   `json:"status"`
}

func main() {
	fileName := "newfile.csv"
	CreateNewFile(fileName)

	client := &http.Client{}
	client.Timeout = time.Second * 15

	var thisResp ApiResp

	var coins = []string{"ETHCLP", "ETHARS", "ETHEUR", "ETHBRL", "ETHMXN", "XLMCLP", "XLMARS", "XLMEUR", "XLMBRL", "XLMMXN",
		"BTCCLP", "BTCARS", "BTCEUR", "BTCBRL", "BTCMXN", "EOSCLP", "EOSARS", "EOSEUR", "EOSBRL", "EOSMXN"}

	var respChan = make(chan ApiResp, 0)

	for i := 0; i < len(coins); i++ {
		go func(reqUrl string, respChan chan ApiResp) {
			req, _ := http.NewRequest("GET", reqUrl, nil)
			resp, _ := client.Do(req)
			defer resp.Body.Close()
			bodyArrBytes, readerr := ioutil.ReadAll(resp.Body)
			if readerr != nil {
				log.Printf("Error on decode bytes. %s", readerr)
			}
			fmt.Println("STATUSCODE:", resp.StatusCode)
			unmarshalerr := json.Unmarshal(bodyArrBytes, &thisResp)
			if unmarshalerr != nil {
				log.Printf("Error on unmarshal. %s", unmarshalerr)
			}
			respChan <- thisResp
		}(fmt.Sprintf("https://api.cryptomkt.com/v1/ticker?market=%s", coins[i]), respChan)
	}

	for i := 0; i < len(coins); i++ {
		response := <-respChan
		WriteLine(fileName, fmt.Sprintf("%v", response))
	}

}

func CreateNewFile(newFileName string) {
	newFile, err := os.Create(newFileName)
	if err != nil {
		log.Fatal(err)
	}
	newFile.Close()
}

func OpenFile(inputFileName string) *os.File {
	file, err := os.Open(inputFileName)
	if err != nil {
		log.Fatal(err)
	}
	return file
}

func WriteLine(inputFileName, line string) {
	newFile, err := os.OpenFile(inputFileName, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	defer newFile.Close()
	newFile.WriteString(line)
	if err != nil {
		panic(err)
	}
}

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type vtotalDomainReport struct {
}

type vtotalError struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func main() {

	url := "https://www.virustotal.com/api/v3/domains/ya.ru"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("accept", "application/json")
	req.Header.Add("x-apikey", "")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
	}

	if res.StatusCode >= 400 {
		var vtotalErrorResp vtotalError
		if err := json.NewDecoder(res.Body).Decode(&vtotalErrorResp); err != nil {
			log.Fatalln(err)
		}
		defer res.Body.Close()
		fmt.Printf("Error %d %s -> %s", res.StatusCode, vtotalErrorResp.Error.Code, vtotalErrorResp.Error.Message)
	} else {
		var vtotalDomainReportResp vtotalDomainReport
		if err := json.NewDecoder(res.Body).Decode(&vtotalDomainReportResp); err != nil {
			log.Fatalln(err)
		}
		defer res.Body.Close()
		fmt.Printf("%s", vtotalDomainReportResp.Data.Type)
	}

}

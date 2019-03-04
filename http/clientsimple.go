package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type GetRequest struct {
	url string // URL al cual mando un GET
	body string // Este es el string mandado
}


func main() {
	print(1,"\n")
	ipRequest := GetRequest {
		url:"https://ifconfig.co",
		body: "",
	}
	for {
		response, err := ipRequest.getResponse()
		if err != nil {
			fmt.Println("Big error:%s",fmt.Sprintf("%s",err))
		} else {
			fmt.Printf("%s", response)
		}
		time.Sleep(time.Second * 3)
	}

}


func (rq *GetRequest) getResponse() ([]byte,error) {
	client := &http.Client{}
	request, err := http.NewRequest("GET", "https://ifconfig.co", nil)
	if err != nil {
		return nil,err
	}
	response, err := client.Do(request)
	if err != nil {
		return nil,err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return body,err
	}
	return body,nil
}
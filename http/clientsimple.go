package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"os"
	"time"
)

type GetRequest struct {
	url string // URL al cual mando un GET
	body string // Este es el string mandado
	typeRequested string
	timeout time.Duration
	debug bool
}

// NewGetRequest para crear un get request con tiempo razonable
//func NewGetRequest(url string) GetRequest {
//	rq := GetRequest{
//		url: url,
//		timeout: time.Second*5,
//	}
//	return rq
//}


func main() {
	ipRequest := GetRequest {
		url:"https://ifconfig.co",
		body: "",
		typeRequested:"",//"application/json", //Puede ser en blanco : ""
		timeout: time.Millisecond * 3000,
		debug: true,
	}

	for {
		response, err := ipRequest.getResponse()
		if err != nil {
			fmt.Println("Big error:%s",fmt.Sprintf("%s",err))
		} else {
			fmt.Printf("%s", response)
		}
		time.Sleep(time.Millisecond * 3000)
	}
}


func (rq *GetRequest) getResponse() ([]byte,error) {
	debug := os.Getenv("DEBUG") // better way of debugging. Setting env var to DEBUG = 1
	if rq.debug || debug=="1" {
		fmt.Printf("[DEBUG] Starting getResponse()")
	}
	client := &http.Client{
		Timeout: rq.timeout,
	}

	request, err := http.NewRequest("GET", "https://ifconfig.co", nil)
	if rq.typeRequested!="" {
		request.Header.Add("Accept",rq.typeRequested)
	}
	if err != nil {
		return nil,err
	}

	if debug == "1" || rq.debug {
		debugRequest, err := httputil.DumpRequestOut(request, true)
		if err != nil {
			return nil, err
		}
		fmt.Printf("%s", debugRequest)
	}
	response, err := client.Do(request)
	if err != nil {
		return nil,err
	}

	defer response.Body.Close()

	if debug == "1" || rq.debug {
		debugResponse, err := httputil.DumpResponse(response, true)
		if err != nil {
			return nil,err
		}
		fmt.Printf("%s", debugResponse)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return body,err
	}
	return body,nil
}
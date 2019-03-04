package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// IPget obtiene direccion IP externa.
func IPget() ([]byte,error) {
	response, err := http.Get("https://ifconfig.co/")
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

type PostRequest struct {
	Request string
	url string
	contentType string
}
func main() {
	A:
	myIP,err := IPget()
	if err!=nil {
		fmt.Println("Error getting IP adress")
		//log.Fatal(err)
	} else {
		fmt.Printf("%s",myIP)
	}
	query := PostRequest{
		Request: fmt.Sprintf("%s",myIP),//`{"some":"json"}`,
		url: "https://httpbin.org/post",
		contentType: "application/json",
	}

	response,err := postThis(query)
	if err!=nil {
		fmt.Println("Error posting. Retrying all steps...")
	} else{
		fmt.Printf("%s",response)
	}

	time.Sleep(time.Second*2)
	goto A


}

// postTo POSTer simple
func postThis(myPost PostRequest) ([]byte, error) {
	postData := strings.NewReader(myPost.Request)
	response, err := http.Post(myPost.url, myPost.contentType, postData)
	if err != nil {
		return nil,err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil,err
	}
	return body,nil
}
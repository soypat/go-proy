package main

import (
	"encoding/binary"
	"io/ioutil"
	"net/http"
)

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



// itob returns an 8-byte big endian representation of v.
func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
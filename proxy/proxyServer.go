package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
)

func main() {
	server := http.Server{
		Addr:    ":8989",
		Handler: &myHandler{},
	}
	server.ListenAndServe()
}

type myHandler struct {
}

func (*myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// output request
	request, err := httputil.DumpRequest(r, false)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(request))

	// create client to get site
	client := http.Client{}
	r.RequestURI = "" // pay attention : it is very important
	response, err := client.Do(r)
	if err != nil {
		log.Panic(err, "not get response")
	}
	defer response.Body.Close()
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Panic(err, "not get Body")
	}
	fmt.Println("Get Body")

	// save Header
	for i, j := range response.Header {
		for _, k := range j {
			w.Header().Add(i, k)
		}
	}
	w.WriteHeader(response.StatusCode)
	w.Write(data)
}

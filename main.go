package main

import (
	"io/ioutil"
	"net/http"
)

type data struct {
    Body []byte
    Error error
}

func futureDate(url string) <- chan data {
    c := make(chan data, 1)

    go func(){
        resp, err := http.Get(url)
        if err != nil {
            c <- data{
            	Body:  nil,
            	Error: err,
            }
            return
        }

        body, err := ioutil.ReadAll(resp.Body)
        resp.Body.Close()

        if err != nil {
            c <- data{
            	Body:  nil,
            	Error: err,
            }
            return
        }

        c <- 
    }()

    return nil
}

func main() {

}

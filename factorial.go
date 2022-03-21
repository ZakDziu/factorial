package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type answJson struct {
	A int `json:"a!"`
	B int `json:"b!"`
}

type Input struct {
	A int `json:"a"`
	B int `json:"b"`
}

func createJson(a, b int, arr []int) ([]byte, error) {
	var answ answJson

	if arr[0] == a {
		answ.A = arr[1]
		answ.B = arr[3]
	}
	answ.A = arr[3]
	answ.B = arr[1]
	if answ.A <= 0 || answ.B <= 0 {
		return nil, errors.New(`{"error": "Incorrect message, very big number"}`)
	}
	newJson, _ := json.Marshal(answ)
	return newJson, nil

}

func calculateFactorial(n int, c chan int) {
	factorial := 1

	for i := 1; i <= n; i++ {
		factorial *= i
	}
	c <- n
	c <- factorial

}

func calculate(a, b int) ([]byte, error) {
	if a <= 0 || b <= 0 {
		return nil, errors.New(`{"error": "Incorrect message"}`)
	}
	c := make(chan int)
	arr := []int{}

	go calculateFactorial(a, c)
	go calculateFactorial(b, c)

	for {
		select {
		case val := <-c:
			arr = append(arr, val)
		}
		if len(arr) == 4 {
			answ, err := createJson(a, b, arr)
			if err != nil {
				return nil, err
			}
			return answ, nil
		}
	}
}

func Calculate(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req Input

	body, error := ioutil.ReadAll(r.Body)

	if error != nil {
		fmt.Fprintln(w, error.Error())
	}

	error = json.Unmarshal(body, &req)

	if error != nil {
		fmt.Fprintln(w, error.Error())
	}

	response, err := calculate(req.A, req.B)

	if err != nil {
		http.Error(w, err.Error(), 400)
	}

	fmt.Fprintln(w, string(response))

}

func main() {
	router := httprouter.New()
	router.GET("/calculate", Calculate)
	fmt.Println("Server listening!")
	log.Fatal(http.ListenAndServe(":8989", router))
}

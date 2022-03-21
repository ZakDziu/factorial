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

type Output struct {
	A int `json:"a!"`
	B int `json:"b!"`
}

type Input struct {
	A int `json:"a"`
	B int `json:"b"`
}

func createJson(a int, arr []int) ([]byte, error) {
	var answ Output

	if arr[0] == a {
		answ.A = arr[1]
		answ.B = arr[3]
	}
	answ.A = arr[3]
	answ.B = arr[1]

	if arr[1] <= 0 || arr[3] <= 0 {
		return nil, errors.New(`{"error": "Incorrect input, very big number"}`)
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
	defer close(c)
	for {
		select {
		case val := <-c:
			arr = append(arr, val)
		}
		if len(arr) == 4 {
			answ, err := createJson(a, arr)
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
		http.Error(w, error.Error(), 400)
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

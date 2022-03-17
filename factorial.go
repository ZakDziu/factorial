package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

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

type Message struct {
	Err string `json:"error"`
}

func createJson(a, b int, arr []int) []byte {
	var answ answJson

	if a > b {
		if arr[0] > arr[1] {
			answ.A = arr[0]
			answ.B = arr[1]
		} else {
			answ.A = arr[1]
			answ.B = arr[0]
		}
	} else {
		if arr[0] < arr[1] {
			answ.A = arr[0]
			answ.B = arr[1]
		} else {
			answ.A = arr[1]
			answ.B = arr[0]
		}
	}
	if arr[0] > 0 && arr[1] > 0 {
		newJson, _ := json.Marshal(answ)
		return newJson
	} else {
		e := Message{Err: "Incorrect input, big number"}
		je, _ := json.Marshal(e)
		return je
	}

}

func calculateFactorial(n int, c chan int) {
	factorial := 1

	for i := 1; i <= n; i++ {
		factorial *= i
	}

	c <- factorial

}

func calculate(a, b int) []byte {
	c := make(chan int, 2)
	until := time.After(1 * time.Second)
	arr := []int{}

	go calculateFactorial(a, c)
	go calculateFactorial(b, c)

	for {
		select {
		case val := <-c:
			arr = append(arr, val)
		case <-until:
			return createJson(a, b, arr)
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

	if req.A <= 0 || req.B <= 0 {
		e := Message{Err: "Incorrect input"}
		je, error := json.Marshal(e)
		if error != nil {
			fmt.Fprintln(w, error.Error())
		}
		w.WriteHeader(400)
		fmt.Fprintln(w, string(je))
	} else {
		if string(calculate(req.A, req.B)) == `{"error":"Incorrect input, big number"}` {
			w.WriteHeader(400)
			fmt.Fprintln(w, string(calculate(req.A, req.B)))
		} else {
			fmt.Fprintln(w, string(calculate(req.A, req.B)))
		}
	}
}

func main() {
	router := httprouter.New()
	router.GET("/calculate", Calculate)
	fmt.Println("Server listening!")
	log.Fatal(http.ListenAndServe(":8989", router))
}

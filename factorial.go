package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

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
type Number struct {
	answ int
	err  error
}

var wg sync.WaitGroup

func (num *Number) Factorial(a, b int) *Number {
	start := time.Now()
	for {
		num.answ = a
		for i := a; i <= b; i++ {
			num.answ *= i
			if time.Since(start) > 3*time.Second {
				break
			}
		}
		if time.Since(start) < 3*time.Second {
			break
		}
		if time.Since(start) > 3*time.Second {
			num.err = err()
			break
		}
	}
	defer wg.Done()
	return num
}

func err() error {
	return errors.New(`{"error": "Incorrect message"}`)
}

func calculate(a, b int) ([]byte, error, int) {
	var answ Output
	var factA = new(Number)
	var factB = new(Number)
	if a <= 0 || b <= 0 {
		return nil, err(), 400
	}
	if a > 20 || b > 20 {
		return nil, errors.New(`{"error": "Incorrect number, very big"}`), 400
	}

	go factA.Factorial(1, a)
	go factB.Factorial(1, b)

	wg.Add(2)
	wg.Wait()

	if factA.err != nil || factB.err != nil {
		return nil, errors.New(`{"error":"Service hang up"}`), 500
	}
	answ.A = factA.answ
	answ.B = factB.answ
	newJson, _ := json.Marshal(answ)

	return newJson, nil, 200
}

func Calculate(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req Input
	err := err()
	body, error := ioutil.ReadAll(r.Body)
	if error != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	error = json.Unmarshal(body, &req)
	if error != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	response, err, code := calculate(req.A, req.B)
	if err != nil {
		http.Error(w, err.Error(), code)
		return
	}

	fmt.Fprintln(w, string(response))

}

func main() {
	router := httprouter.New()
	router.GET("/calculate", Calculate)
	fmt.Println("Server listening!")
	log.Fatal(http.ListenAndServe(":8989", router))
}

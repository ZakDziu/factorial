package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"github.com/julienschmidt/httprouter"
)

const (
	maxIntegerFactorial = 20
	minIntegerFactorial = 0
)

var (
	errIncorrectMessage   = errors.New(`{"error": "Incorrect message"}`)
	errServiceUnavailable = errors.New(`{"error": "Service error"}`)
)

type Answer struct {
	A    int `json:"a!"`
	B    int `json:"b!"`
	errA error
	errB error
}

type Request struct {
	A int `json:"a"`
	B int `json:"b"`
}

func calculateFactorial(to int) (int, error) {
	factorial := 1
	if to > maxIntegerFactorial || to <= minIntegerFactorial {
		return 0, errIncorrectMessage
	}

	for i := 1; i <= to; i++ {
		factorial *= i
	}

	return factorial, nil
}

func calculate(a, b int) (string, error, int) {
	answ := &Answer{}
	wg := &sync.WaitGroup{}

	wg.Add(2)

	go func() {
		defer wg.Done()
		answ.A, answ.errA = calculateFactorial(a)
	}()

	go func() {
		defer wg.Done()
		answ.B, answ.errB = calculateFactorial(b)
	}()

	wg.Wait()

	if answ.errA != nil || answ.errB != nil {
		return "", errIncorrectMessage, http.StatusBadRequest
	}

	newJson, err := json.Marshal(answ)
	if err != nil {
		return "", errServiceUnavailable, http.StatusInternalServerError
	}

	return string(newJson), nil, http.StatusAccepted
}

func Calculate(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req Request

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, errServiceUnavailable.Error(), http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, &req)
	if err != nil {
		http.Error(w, errIncorrectMessage.Error(), http.StatusBadRequest)
		return
	}

	response, err, statusCode := calculate(req.A, req.B)
	if err != nil {
		http.Error(w, err.Error(), statusCode)
		return
	}

	fmt.Fprintln(w, response)

}

func main() {
	router := httprouter.New()
	router.GET("/calculate", Calculate)
	fmt.Println("Server listening!")
	log.Fatal(http.ListenAndServe(":8989", router))
}

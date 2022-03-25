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

var wg sync.WaitGroup

const maxIntegerFactorial = 20
const minIntegerFactorial = 0

var errIncorrectMessage = errors.New(`{"error": "Incorrect message"}`)
var errServiceUnavailable = errors.New(`{"error": "Service error"}`)

type output struct {
	A int `json:"a!"`
	B int `json:"b!"`
}

type input struct {
	A int `json:"a"`
	B int `json:"b"`
}

type number struct {
	answ int
	err  error
}

func (num *number) factorial(to int) *number {
	num.answ, num.err = calculateFactorial(to)
	defer wg.Done()
	return num
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
	var answ output
	var factA = new(number)
	var factB = new(number)

	wg.Add(2)
	go factA.factorial(a)
	go factB.factorial(b)
	wg.Wait()

	if factA.err != nil || factB.err != nil {
		return "", errIncorrectMessage, http.StatusBadRequest
	}

	answ.A = factA.answ
	answ.B = factB.answ

	newJson, err := json.Marshal(answ)
	if err != nil {
		return "", errServiceUnavailable, http.StatusInternalServerError
	}

	return string(newJson), nil, http.StatusAccepted
}

func Calculate(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req input

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

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/sync/errgroup"
)

const (
	maxIntegerFactorial = 20
	minIntegerFactorial = 0
)

var (
	ErrIncorrectMessage   = errors.New(`{"error": "Incorrect message"}`)
	ErrServiceUnavailable = errors.New(`{"error": "Service error"}`)
)

type Answer struct {
	A int `json:"a!"`
	B int `json:"b!"`
}

type Request struct {
	A int `json:"a"`
	B int `json:"b"`
}

func CalculateFactorial(to int) (int, error) {
	factorial := 1
	if to > maxIntegerFactorial || to <= minIntegerFactorial {
		return 0, ErrIncorrectMessage
	}

	for i := 1; i <= to; i++ {
		factorial *= i
	}

	return factorial, nil
}

func calculateF(a, b int) (string, int, error) {
	answ := &Answer{}
	group := new(errgroup.Group)
	var err error
	group.Go(func() error {
		answ.A, err = CalculateFactorial(a)
		return err
	})

	group.Go(func() error {
		answ.B, err = CalculateFactorial(b)
		return err
	})

	if err := group.Wait(); err != nil {
		return "", http.StatusBadRequest, ErrIncorrectMessage
	}

	newJson, err := json.Marshal(answ)
	if err != nil {
		return "", http.StatusInternalServerError, ErrServiceUnavailable
	}

	return string(newJson), http.StatusAccepted, nil
}

func calculate(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req Request

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, ErrServiceUnavailable.Error(), http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, &req)
	if err != nil {
		http.Error(w, ErrIncorrectMessage.Error(), http.StatusBadRequest)
		return
	}

	response, statusCode, err := calculateF(req.A, req.B)
	if err != nil {
		http.Error(w, err.Error(), statusCode)
		return
	}

	fmt.Fprintln(w, response)

}

func router() *httprouter.Router {
	router := httprouter.New()
	router.GET("/calculate", calculate)
	return router
}

func main() {
	fmt.Println("Server listening!")
	log.Fatal(http.ListenAndServe(":8989", router()))
}

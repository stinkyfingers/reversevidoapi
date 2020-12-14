package main

import (
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/stinkyfingers/lambdify"
)

func main() {
	lambdaFunction := lambdify.Lambdify(mux())
	lambda.Start(lambdaFunction)
}

func mux() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/handle", handler)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		log.Print("health called")
		_, err := w.Write([]byte("OK"))
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	})
	return mux
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(r.URL.String()))
}

package main

import (
	"errors"
	"fmt"

	"github.com/alexbathome/errefl"
)

type MyHttpError struct {
	errefl.Err `errtpl:"*{endpoint}* is not responding on: {protocol}"`

	StatusCode int
	Endpoint   string `errefl:"endpoint"`
	Protocol   string `errefl:"protocol"`
}

type MyApplicationError struct {
	errefl.Err `errtpl:"unable to run process: {process}"`

	Process string `errefl:"process"`
}

type wrapper interface {
	Wrap(error)
}

var _ wrapper = (*MyHttpError)(nil)
var _ wrapper = (*MyApplicationError)(nil)

func main() {
	err := errefl.New2[MyHttpError](500, "https://example.com", "HTTPS")
	errWrapped := errefl.NewWrapped2[MyApplicationError](err, "myprocess")

	var myHttpError MyHttpError
	var myApplicationError MyApplicationError
	errors.As(errWrapped, &myHttpError)
	errors.As(errWrapped, &myApplicationError)

	fmt.Println(myApplicationError.Error(), myApplicationError.Process)
	fmt.Println(myHttpError.Error(), myHttpError.Endpoint, myHttpError.Protocol, myHttpError.StatusCode)
}

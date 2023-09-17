package main

import (
	"fmt"

	"github.com/alexbathome/errefl"
)

type MyHttpError struct {
	errefl.Err `errtpl:"{endpoint} is not responding on {protocol}"`

	StatusCode int
	Endpoint   string `errefl:"endpoint"`
	Protocol   string `errefl:"protocol"`
}

type MyApplicationError struct {
	errefl.Err `errtpl:"unable to run process {process}"`

	Process string `errefl:"process"`
}

func main() {
	err := errefl.New[MyHttpError](500, "https://example.com", "HTTPS")

	select {
	case v, ok := errefl.Catch[MyHttpError](err), <- _:
		fmt.Printf("%s; %s, %s", v.Error(), v.Protocol, v.Endpoint)
	}
}

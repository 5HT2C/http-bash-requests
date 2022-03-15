package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"strings"
)

var (
	port   = flag.String("port", "6016", "port to use for serving")
	bin    = flag.String("bin", "/bin/bash", "binary to run")
	binArg = flag.String("binArg", "-c", "binary args to use")
)

func handle(_ http.ResponseWriter, req *http.Request) {
	var buf strings.Builder
	if _, err := io.Copy(&buf, req.Body); err != nil {
		log.Printf("err with io.Copy: %s\n", err)
	}

	if out, err := exec.Command(*bin, *binArg, buf.String()).Output(); err != nil {
		log.Printf("err with exec.Command: %s\n", err)

	} else {
		fmt.Printf("%s", out)
	}
}

func main() {
	flag.Parse()
	http.HandleFunc("/", handle)

	if err := http.ListenAndServe("localhost:"+*port, nil); err != nil {
		log.Printf("err with request: %s\n", err)
	}
}

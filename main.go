package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"strings"
)

var (
	argPort   = flag.String("port", "6016", "port to use for serving")
	argBin    = flag.String("bin", "/bin/bash", "binary to run")
	argBinArg = flag.String("binArg", "-c", "binary args to use")
)

func handle(w http.ResponseWriter, req *http.Request) {
	var buf strings.Builder
	if _, err := io.Copy(&buf, req.Body); err != nil {
		log.Printf("err with io.Copy: %s\n", err)
		return
	}

	body := buf.String()
	bin := *argBin
	binArg := *argBinArg

	binPath := req.Header.Get("X-Bin-Path")
	if len(binPath) > 0 {
		bin = binPath
		binArg = req.Header.Get("X-Bin-Arg")
	}

	args := []string{body}

	if req.Header.Get("X-Body-Split") == "true" {
		r := csv.NewReader(strings.NewReader(body))
		r.Comma = ' ' // space
		fields, err := r.Read()

		if err != nil { // will usually fail if body is empty
			log.Printf("err with r.Read: %v\n", err)
			return
		}

		args = fields
	}

	// If binArg is not empty, prepend to slice
	if len(binArg) > 0 {
		args = append([]string{binArg}, args...)
	}

	if out, err := exec.Command(bin, args...).CombinedOutput(); err != nil {
		_, _ = fmt.Fprintln(w, err.Error())
		log.Printf("err with exec.Command: %s\n", err)
	} else {
		_, _ = fmt.Fprint(w, string(out))
		fmt.Printf("%s", out)
	}
}

func main() {
	flag.Parse()
	http.HandleFunc("/", handle)

	if err := http.ListenAndServe("localhost:"+*argPort, nil); err != nil {
		log.Printf("err with request: %s\n", err)
	}
}

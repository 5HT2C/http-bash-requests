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
	argPort   = flag.String("port", "6016", "port to use for serving")
	argBin    = flag.String("bin", "/bin/bash", "binary to run")
	argBinArg = flag.String("binArg", "-c", "binary args to use")
)

func handle(w http.ResponseWriter, req *http.Request) {
	var buf strings.Builder
	if _, err := io.Copy(&buf, req.Body); err != nil {
		log.Printf("err with io.Copy: %s\n", err)
	}

	bin := *argBin
	binArg := *argBinArg

	binPath := req.Header.Get("X-Bin-Path")
	if len(binPath) > 0 {
		bin = binPath
		binArg = req.Header.Get("X-Bin-Arg")
	}

	args := make([]string, 0)
	if binArg == "" {
		args = []string{buf.String()}
	} else {
		args = []string{binArg, buf.String()}
	}

	if out, err := exec.Command(bin, args...).CombinedOutput(); err != nil {
		_, _ = fmt.Fprintf(w, "%s\n", err)
		log.Printf("err with exec.Command: %s\n", err)
	} else {
		_, _ = fmt.Fprintf(w, string(out))
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

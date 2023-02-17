package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/chyroc/detect-ip/internal"
)

func main() {
	url := os.Args[1]
	cli := &http.Client{Timeout: time.Second * 5}
	fmt.Println(internal.DetectIP(cli, url, false, os.Getenv("DEBUG") != ""))
	fmt.Println(internal.DetectIP(cli, url, true, os.Getenv("DEBUG") != ""))
}

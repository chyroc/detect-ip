package main

import (
	"fmt"
	"github.com/chyroc/detect-ip/internal"
	"net/http"
	"os"
	"time"
)

func main() {
	url := os.Args[1]
	cli := &http.Client{Timeout: time.Second * 5}
	fmt.Println(internal.DetectIP(cli, url, false, false))
	fmt.Println(internal.DetectIP(cli, url, true, false))
}

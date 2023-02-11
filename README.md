# detect-ip

detect-ip is a tool to detect your public ip(v4/v6).

## Install

- *For Cli*

```shell
go install github.com/chyroc/detect-ip/cli/detect-ip@latest
```

- *For Go Package*

```shell
go get github.com/chyroc/detect-ip/detect_ip
```

## Usage

- *For Cli*

```shell
detect-ip

detect-ip -timeout 3000 # 3000 = 3s
```

- *For Go Package*

```go
package main

import (
    "fmt"

    "github.com/chyroc/detect-ip/detect_ip"
)

func main() {
    ip := detect_ip.PublicIPV4() // PublicIPV6
	// ip := detect_ip.PublicIPV4(detect_ip.WithTimeout(time.Second*3))
    fmt.Println(ip)
}
```


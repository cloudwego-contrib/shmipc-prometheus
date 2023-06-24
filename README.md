## shmipc-prometheus is a Prometheus monitoring for [shmipc-go](https://github.com/cloudwego/shmipc-go) (*This is a community driven project*)

## How to use with shmipc-go server?

**[example/shmipc_server/main.go](example/shmipc_server/main.go)**

```go
package main

import (
	shmipcprometheus "github.com/cloudwego-contrib/shmipc-prometheus"
	"github.com/cloudwego/shmipc-go"
)

func main() {
	...
	conf := shmipc.DefaultSessionManagerConfig()
	conf.Monitor = shmipcprometheus.NewPrometheusMonitor(":9094", "/metrics")
	smgr, _ := shmipc.NewSessionManager(conf)
	...
}
```

## How to use with shmipc-go client?

**[example/shmipc_client/main.go](example/shmipc_client/main.go)**

```go
package main

import (
	shmipcprometheus "github.com/cloudwego-contrib/shmipc-prometheus"
	"github.com/cloudwego/shmipc-go"
)

func main() {
	...
	conf := shmipc.DefaultConfig()
	conf.Monitor = shmipcprometheus.NewPrometheusMonitor(":9095", "/metrics")
	server, err := shmipc.Server(conn, conf)
	...
}
```

## example

**[example for shmipc-prometheus](example)**


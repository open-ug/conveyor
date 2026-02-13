---
sidebar_position: 6
---

# Embeding into your System

Conveyor CI can be embedded into you system via two ways depending on its requirements and scale

## 1. Intergrating via API

Option one is to intergrate it into you system as an external microservice that can be called via an API. See the [Using the API Guide](./using-api.mdx) to understand the API reference. Client SDKs also often absrtact this API into utility functions making it easy to build programs that communicate with Conveyor CI as a microservice.

## 2. Embedding in a Go Application

Conveyor CI is built in Go and is also published as a Go package that you can embed into you go Application. This ensures that you can compile both Conveyor CI server and your application into one program. Conveyor CI also allows you to extend the API HTTP Server with more custom routes. This is possible because the HTTP Server uses gofiber v2 under the hood.

Below is an example code snippet of embeding Conveyor CI in a Go application and starting the server.

```go
package main

import (
 "github.com/open-ug/conveyor/pkg/server"
 "github.com/open-ug/conveyor/pkg/types"
)

func main() {
 config := types.ServerConfig{}

 config.API.Port = 8000
 config.API.AuthEnabled = false
 config.API.Data = "/data"

 config.NATS.Port = 4222

 apiContext, err := server.Setup(&config)
 if err != nil {
  panic(err)
 }

  // Extending with custom /health route
 apiContext.App.Get("/health", healthHandler)

 apiContext.Start()
}
```

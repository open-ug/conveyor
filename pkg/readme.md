# Conveyor CI SDK

This is the official Golang SDK for Conveyor CI. Conveyor CI is an open-source engine and framework for building CI/CD platforms.

This Go package provides an API Client for interacting with Conveyor API Server and a driver runtime to assist you in driver development.

## Installation

```sh
go get -u github.com/open-ug/conveyor
```

## Usage

A simple codebase using this library, you can find full codebase at [https://github.com/open-ug/simple-runner](https://github.com/open-ug/simple-runner)

```go
package main

import (
 "context"
 "fmt"

 c "github.com/open-ug/conveyor/pkg/client"
 runtime "github.com/open-ug/conveyor/pkg/driver-runtime"
 cmd "github.com/open-ug/runner/cmd"
 utils "github.com/open-ug/runner/cmd/utils"
)

func main() {
 client := c.NewClient()

 // Register the Pipeline resource definition with the client
 _, err := client.CreateOrUpdateResourceDefinition(context.Background(), utils.PipelineResourceDefinition)

 if err != nil {
  panic(err)
 }

 // Create a new driver instance
 driver := &runtime.Driver{
  Name: "command-runner",
  Resources: []string{
   utils.PipelineResourceDefinition.Name},
  Reconcile: cmd.Reconcile,
 }

 // Create a new driver manager with the driver
 driverManager, err := client.NewDriverManager(driver, []string{"*"})
 if err != nil {
  fmt.Println("Error creating driver manager: ", err)
  return
 }

 // Start the driver manager
 err = driverManager.Run()
 if err != nil {
  fmt.Println("Error running driver manager: ", err)
 }
}
```

## License

Apache License 2.0, see [LICENSE](./LICENSE).

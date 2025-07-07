---
sidebar_position: 3
---

# Building the Driver

:::danger

Documentation still in early stages

:::


Assumming you have already followed the Installation Guide, this tutorial expects you to already have a Go project setup and and the Conveyor CI tools installed. We shall use that to create the driver.

A driver is a component that runs tasks in your CI/CD platform basing off the resource objects created. Visit the Driver Documentation vor more information about drivers.

In this tutorial we shall create a simple driver that compresses a a directory of code into a zip file and uploads it to a defined server. We shall use the resource defined in the Defining a Resource section of this tutorial. We shall name our driver `zipper`

Within your `main.go` file paste this code below. We shall explain what it does

```go
package main

import (
 "encoding/json"
 "fmt"
 "log"

 runtime "github.com/open-ug/conveyor/pkg/driver-runtime"
 dLogger "github.com/open-ug/conveyor/pkg/driver-runtime/log"
 types "github.com/open-ug/conveyor/pkg/types"
)

// Listen for messages from the runtime
func Reconcile(payload string, event string, driverName string, logger *dLogger.DriverLogger) error {

 return nil
}

func Listen() {
 driver := &runtime.Driver{
  Reconcile: Reconcile,
  Name:      "zipper",
 }

 driverManager, err := runtime.NewDriverManager(driver, []string{"*"})
 if err != nil {
  fmt.Println("Error creating driver manager: ", err)
  return
 }

 err = driverManager.Run()
 if err != nil {
  fmt.Println("Error running driver manager: ", err)
 }

}
```

---
sidebar_position: 1
---

# Driver Development

As a Platform developer using Conveyor CI to build your platform. Drivers that main component you will need to develop as they are the main components that carry out the work your platform is meant to perform.

This guide will take you throught the processes of developing a driver. We shall mainly use the Go SDK for implementation reference but the concepts remain similar throught.

Before continuing with this guide, we highly recommend that you read the [Drivers Concept](/docs/concepts/drivers) page to understand the [Components of the driver](/docs/concepts/drivers#components-of-a-driver) and the [Driver lifecycle](/docs/concepts/drivers#driver-lifecycle).

## Driver Overview

Drivers are software programs that execute CI/CD processes. When developing one, It requires certain properties that you need to define and these include:

- A unique name
- A list of resources the drivers listens for
- And a `Reconcile` function that contains the execution logic.

Conveyor CI also provides you with cerain helper utitlities you can use when developing a driver and are often packages in whats called the **Driver runtime**. These include:

- The Driver Manager
- The Driver Logger
- The Run ID
- The API Client
- The Driver Result

These components are explained in the [Drivers Concept](/docs/concepts/drivers) page.

## The `Reconcile` Fuction

The `Reconcile` function is a function in a driver that contains the custom execution logic that the driver runs upon a resource state change. This is run each time a new resource is created or updated. It runs under a loop that is continously listening for such resource mutation events and is triggered when a change occurs.

This function often has a set of paramentes that it exposes to your logic and these include:

- The resource payload: This is the resource JSON payload that was created. It also often contains extra metadata that what you defined.
- The event name: This is the name of the event that occured and can either be `create` or `update` reperesenting if the resource what just created or was merely updated.
- The Run ID: this is a unique UUID string that identifies the resource event uniquely i.e each time a resource is created or updated, a UUID string is generated and attached that event, so each time a driver reconciles, it should expect a different run id.
- The Driver Logger: This is used to collect logs that you might find useful to store in long time storage fromm your driver execution.

The function also returns a return type called a Driver result. This is an object that contains information if the execution was a success or a failure, a message explaining the what happened and any extra abitrary data that about the execution that might be useful to store.

In Go lang, the reconcile function would look like this:

```go
import (

 "fmt"
 "time"

 logger "github.com/open-ug/conveyor/pkg/driver-runtime/log"
 "github.com/open-ug/conveyor/pkg/types"
)

// Listen for messages from the runtime
func Reconcile(payload string, event string, runID string, logger *logger.DriverLogger) types.DriverResult {

 /// simulate doing some work by looping 5 times
 for i := 0; i < 5; i++ {
  // using the Driver logger to save logs
  logger.Log(map[string]string{"step": fmt.Sprintf("%d", i+1)}, fmt.Sprintf("Sample Driver Working... %d/5", i+1))
 }

  // The function returns a Driver Result type
 return types.DriverResult{
  Success: true,
  Message: "Sample Driver Reconciled Successfully",
  Data:    nil,
 }
}
```

## Driver Example

Lets looks at an example of a driver code base. We shall first show the code and explain it step by step.

```go showLineNumbers
package main

import (
 "fmt"
 "time"

 runtime "github.com/open-ug/conveyor/pkg/driver-runtime"
 logger "github.com/open-ug/conveyor/pkg/driver-runtime/log"
 "github.com/open-ug/conveyor/pkg/types"
)

// Listen for messages from the runtime
// highlight-next-line
func Reconcile(payload string, event string, runID string, logger *logger.DriverLogger) types.DriverResult {

 /// simulate doing some work by looping 5 times
 for i := 0; i < 5; i++ {
  // highlight-next-line
  logger.Log(map[string]string{"step": fmt.Sprintf("%d", i+1)}, fmt.Sprintf("Sample Driver Working... %d/5", i+1))
  time.Sleep(2 * time.Second)
 }

 // highlight-next-line
 return types.DriverResult{
  Success: true,
  Message: "Sample Driver Reconciled Successfully",
  Data:    nil,
 }
}

func main() {
 // highlight-start
 driver := &runtime.Driver{
  Reconcile: Reconcile,
  Name:      "sample-driver",
  Resources: []string{"resource-1"},
 }
 // highlight-end

 // highlight-next-line
 driverManager, err := runtime.NewDriverManager(driver, []string{"*"})

 if err != nil {
  fmt.Println("Error creating driver manager: ", err)
  return
 }

 // highlight-next-line
 err = driverManager.Run()
 if err != nil {
  fmt.Println("Error running driver manager: ", err)
 }

}
```

Lets jump to the entry point of the program in the `main()` function. W begin by defining our driver using the `Driver` struct on lines 29 - 33 in which we define the driver name, resources it listends too and the Reconcile function. We then move on to intiate the driver manager on line 35 and passing in the driver we created. We then move on to line 42 where we call  the `Run()` method on the driver manager, this starts the rvent listening loop waiting for events that will be passed to the Reconcile function.

Moving on to the `Reconcile` function. We defined it on line 13 and it takes in the parameters as mentioned before and specifies the return type of Driver Result. We can also observer te usage of the Drirver logger on line 17 which send logs on to the Conveyor CI API Server for long term storage.

## Streaming Driver logs

The logs that are collected from the Driver by the Driver logger can be streamed or collected from the API server in realtime or not.

To do this, you simply have to open a Websocket connection to the `/logs/streams/{DROVER_NAME}/{RUN_ID}` route i.e `ws://localhost:8080/logs/streams/sample-driver/fbab31f6-a278-4a8f-96be-ac49b007ca65`.

This connection will then return JSON data containing log lines and there associated labels you might have defined.

Besides opening a websocket connection for realtime logs, You can also fetch these logs using HTTP by quering the the API on the `/logs` route. You can also specify query paramenters like `driver`, `runid`, `pipeline` forexample `GET /logs?driver=...&runid=...`
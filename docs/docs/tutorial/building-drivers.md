---
sidebar_position: 4
---

# Building the Drivers

Moving on, lets build the drivers.

As mentioned before the drivers contain the the logic for executing the CI/CD processes. We had outlines that wee need 4 drivers both in the architecture and the pipeline definition and these include:

- `container-start`: A driver to create and start the build container.
- `git-cloner`: A driver to clone the repository in the container.
- `builder`: A driver to trigger the build process.
- `container-stop`: A driver to stop and delete the container.

## Driver Basics

Drivers in Conveyor CI have a similar high level structure. They are composed of mainly three key properties.

- A unique name that defines and identifes the driver.
- A list of resources that the driver listens to.
- A `Reconcile` function that is run each time there is a new resources or mutation of a resource created.

Within the `Reconcile` function, custom logic can be added where the driver decodes the resource and carries out the appropriate computations to execute the CI/CD process. This function is always listening in realtime for events triggered when a resaource changes. It returns a return type called a `DriverResult` that contains information about the results of the runction, whether it was successfull of not.

## Driver Code Example

To build a driver, you can use one of the growing list of SDKs provided by Conveyor CI, however in this tutorial, we shall use the Go SDK mainly because its the most stable and mature.

We shall demosnstrate code for only the `container-start` drover but you can find all code in the [https://github.com/open-ug/simple-runner](https://github.com/open-ug/simple-runner) repository.

Lets look at the code.

```go showLineNumbers
/// https://github.com/open-ug/simple-runner/blob/main/cmd/container-start/app.go
package containerstart

import (
 //"encoding/json"

 "context"
 "fmt"

 "github.com/docker/docker/client"
 runtime "github.com/open-ug/conveyor/pkg/driver-runtime"
 logger "github.com/open-ug/conveyor/pkg/driver-runtime/log"
 "github.com/open-ug/conveyor/pkg/types"
 "github.com/open-ug/runner/cmd/utils"
)

// Listen for messages from the runtime
// highlight-next-line
func Reconcile(payload string, event string, runID string, logger *logger.DriverLogger) types.DriverResult {

 // Initialize Docker client
 ctx := context.Background()
 cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
 if err != nil {
  return types.DriverResult{
   Success: false,
   Message: fmt.Sprintf("failed to create Docker client: %v", err),
  }
 }

 // 1. Create and start the container
 containerID, err := utils.CreateAndStartContainer(ctx, cli, "ghcr.io/cirruslabs/flutter:3.38.5", nil)
 if err != nil {
  return types.DriverResult{
   Success: false,
   Message: fmt.Sprintf("failed to create/start container: %v", err),
  }
 }

 return types.DriverResult{
  Success: true,
  Message: "Sample Driver Reconciled Successfully",
  Data: map[string]interface{}{
   "containerID": containerID,
  },
 }
}

func Listen() {
 // highlight-start
 driver := &runtime.Driver{
  Reconcile: Reconcile,
  Name:      "container-start",
  Resources: []string{utils.FlutterBuilderResourceDefinition.Name},
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

If we look at the above code sample from the `container-start` driver. we can identify the key aspects of building a driver. At line 49 - 53, we are defining the driver and its key properties which are the name, resources and the Reconcile function(whose implementation starts at line 18). We then pass it into the driver manager at line 55 and trigger the run function which listens for events and calles the reconcile function.

The above code is only for the `container-start` driver but in the [https://github.com/open-ug/simple-runner](https://github.com/open-ug/simple-runner) repository. All the code for other drivers can be found there.

The code at the above mentioned repository is built as a command line program with each driver written as a subcommand. YOu can clone the repository and run `go run main.go --help` to see the available subcommands.

```sh
$ go run main.go --help

Flutter CI is a continuous integration tool for Flutter applications.

Usage:
  flutter-ci [flags]
  flutter-ci [command]

Available Commands:
  builder         Start the Builder Driver
  completion      Generate the autocompletion script for the specified shell
  container-start Start the Container Start Driver
  container-stop  Start the Container Stop Driver
  git-cloner      Start the Git Cloner Driver
  help            Help about any command

Flags:
  -h, --help      help for flutter-ci
  -t, --toggle    Help message for toggle
  -v, --version   version for flutter-ci

Use "flutter-ci [command] --help" for more information about a command.
```

You can then start the individual driver instances.

```sh
$ go run main.go container-start
$ go run main.go git-cloner 
$ go run main.go builder
$ go run main.go container-stop
```

> In conveyor ci, one driver can have multiple instances forexample you can have more than one processes of the `container-start` driver running at once.

Once you have the drivers running, we can move on to the next section of triggering a workflow or CI/CD process.
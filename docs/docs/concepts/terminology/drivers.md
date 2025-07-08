---
sidebar_position: 2
---

# Drivers

:::danger

Documentation still in early stages

:::

A Driver is a component that runs pipeline jobs depending on the configuration defined in the [Resource](resources). When you create a Resource object, the driver is always listening for new objects and once one is created, it calls a function that reconciles the object state to desired pipline actions.

A Driver requests for what objects and Events to listen for and when a new resource object is created. it then runs a ppipeline job. This process is called **Driver Run**. Each Driver Run stores a unique `run-id` that is a UUID unique to each time a resource onject undergoes a CRUD process.

## Driver Manager

All drivers are composed of a components called a Driver manager that manages the Driver process. The manager is incharge for connecting to the Event stream and listening for new objects and events that occur in the Conveyor CI System and passing them to the Reconcile Function.

## Driver Logger

Drivers are also composed of a Driver logger that is incharge of collecting and storing pipeline Job logs. These logs can then be retrieved or streamed via a websocket client on the COnveyor CI API Server.

## Driver Runtime

The Driver runtime is a set of software libraries that are used to build Conveyor CI Drivers. It is abstracts internal workings of Drivers providing you with a simplified way to focus driver logic. The driver runtime is officially ditributed for the Following laguage Client libraries.

| Language | Status      | Source Code                                                              | Documentation                                                      |
| -------- | ----------- | ------------------------------------------------------------------------ | ------------------------------------------------------------------ |
| Go       | Beta        | [github.com/open-ug/conveyor](https://github.com/open-ug/conveyor)       | [Docs](https://pkg.go.dev/github.com/open-ug/conveyor@v0.1.14/pkg) |
| Python   | Coming Soon | [github.com/open-ug/conveyor.py](https://github.com/open-ug/conveyor.py) | Todo                                                               |
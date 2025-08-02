---
sidebar_position: 3
---

# Drivers

Drivers in Conveyor CI are pluggable software components that transform [Resource](./resources) state into desired CI/CD executions. Drivers watch for new changes in Resource state and carry out appropriate executions depending on the state defined.

Drivers run under an execution environment called the Driver Runtime. This environment provides a pre-defined lifecycle that Drivers follow throughout execution.

## Driver Lifecycle

Throughout a driver's execution it follows this cycle:

- **Initiation**: The driver first registers and defines some metadata when it's being created. This metadata includes its name, events and the resources it wants to listen to.
- **Listen for Events**: The next step in the cycle is to listen for events from Conveyor CI. These events occur when a new Resource is created or a change in Resource state is detected.
- **Reconcile Resource State**: Once an event is received, the Driver runs a Reconcile function that includes logic to transform Resource state into desired CI/CD processes.

## Components of a Driver

A fully functional driver consists of multiple important components that a developer needs to understand.

- **Driver**: The Driver is the core software component that executes CI/CD processes. It contains the Reconcile function that contains the custom code/logic that includes the drivers intended functionality. Each reconciliation done by a driver is called a *Driver Run* and can be identified via a unique UUID string called a *Run ID*.
- **Driver Manager**: Orchestrates the Driver’s lifecycle and communication with the Conveyor CI Event Stream. It is responsible for connecting to the Conveyor CI Event Stream. Receiving and forwarding Resource events to the Reconcile Function and Managing Driver processes and ensuring stable execution.
- **Driver Logger**: This component is in charge of handling logging for drivers, it collects, stores and streams necessary logs from driver executions. These logs can then be collected by or streamed via a Client.
- **Driver Runtime**: Provides a set of libraries and tools to simplify Driver development. It also acts as an execution environment that is responsible for abstracting low-level event handling and job execution allowing developers to focus solely on pipeline logic.

The Driver runtime has been implemented in multiple programming languages in order to ensure development of drivers is possible in multiple programming languages. Current and future implementations include the following.

| Language   | Status   | Source Code                                                           | Documentation                                                        |
| ---------- | -------- | --------------------------------------------------------------------- | -------------------------------------------------------------------- |
| Go         | Stable   | [github.com/open-ug/conveyor](https://github.com/open-ug/conveyor)    | [Go Doc](https://pkg.go.dev/github.com/open-ug/conveyor@v0.1.14/pkg) |
| Python     | Beta     | [github.com/open-ug/conveyor](https://github.com/open-ug/conveyor)    | Todo                                                                 |
| Rust       | Proposed | [github.com/open-ug/conveyor](https://github.com/open-ug/conveyor.py) | Todo                                                                 |
| JavaScript | Proposed | [github.com/open-ug/conveyor](https://github.com/open-ug/conveyor.py) | Todo                                                                 |
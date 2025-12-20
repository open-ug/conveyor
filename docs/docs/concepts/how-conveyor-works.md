---
sidebar_position: 1
---

# How Conveyor CI Works

Conveyor CI is designed to be a highly modular, Realtime and Event Driven CI/CD system. It orchestrates CI/CD processes by storing system state as Resources and processing those states through Drivers.

At its core, Conveyor CI provides a resource-driven architecture for its users to interact with. This model allows Conveyor CI to react instantly to changes in your pipelines while remaining scalable and modular.

To understand how these components work together, let’s explore a high-level view of the Resource-driven architecture in Conveyor CI.

## Conveyor CI’s Resource-Driven Architecture

Conveyor CI’s Resource Driven Architecture is centred around the concept that Resources determine every event that occurs in the system. This means that users define their CI/CD processes in terms of Resources, and all events and executions are driven by changes in those Resources.

Internally, Conveyor CI reacts to Resource mutations rather than relying on fixed pipelines or static workflows. If no Resource is created, updated, or deleted, no CI/CD processes are triggered.

At a high level, it works like this:

- [Resources](resources) store the state of your CI/CD processes.
- [Resource Definitions](resource-definitions) define the schema and validation rules for those Resources.
- [Drivers](drivers) watch for changes to Resources and transform that state into actual CI/CD executions.

This flow provides developers with the ability to build Event Driven CI/CD platforms in which actions only occur when Resources change. It also enables Stateful observability since each change can be traced back to a Resource change.

Although this is the Architecture that Conveyor CI provides, internally it uses an entirely different architecture and requires and contains multiple software components to pull this off.

## Conveyor CI System Architecture

Internally, conveyor CI contains multiple software components that are responsible for different functionalities like Data Storage, Event Streaming etc. 

These components include:

- **Conveyor API Server**: This acts as the entry point into the system. It carries out certain core functionalities like validating Resources against Resource definitions, Publishing events to event stream in case of resource change.
- **etcd**: This is a Key Value store in charge of storing all system state. It acts as a data store for Conveyor CI.
- **Nats**: Nats is a message broker that is in charge of event routing and processing. It enables decoupled and event driven execution.
- **Badger DB**: This is in charge of log storage. It collects and stores important logs from Drives and allows developers to be able to observe CI/CD execution logs.

### Execution Workflow

Conveyor CI follows a standard execution workflow that complements its architecture. The process occurs as follows:

1. Resource Submission

     - A user sends a Resource object to the API Server.
  
2. Validation and Storage
     - The API Server validates the Resource against its Resource Definition.
     - If valid, the Resource is stored in etcd.

3. Event Publication
   - The API Server publishes a Resource change event to NATS.

4. Event Routing

   - NATS routes the event to the appropriate Driver subscribed to that Resource type.

5. Driver Execution

   - The Driver executes the CI/CD process defined by the Resource.
   - Execution logs are sent back to Conveyor CI for storage.

6. Log Streaming (Optional)

   - If a user wants to stream logs, they open a WebSocket to the API Server.

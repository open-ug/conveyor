---
slug: feature-proposal--enable-api-server-embedding
title: Feature Proposal - Enable API Server Embedding in Go Applications
authors: jim-junior
tags: [conveyor-ci, proposal]
---

Proposal by Infralane Cloud team to enable API Server Embedding in Go Applications

<!-- truncate -->

## Summary

Refactor the Conveyor CI API server into a modular, embeddable Go package that can be imported and started programmatically within other Go applications.

Instead of requiring a standalone Conveyor CI server process, developers should be able to embed and control the API server directly from their own codebases.

---

## Rationale

Conveyor CI is designed to be minimal, lightweight, and suitable for air-gapped and constrained environments. Requiring operators to always deploy and manage a standalone API server introduces unnecessary operational complexity.

By enabling embedding:

* Conveyor CI becomes a headless, composable workflow engine.
* Developers gain tighter integration within existing systems.
* Operational overhead is reduced in embedded or appliance-style deployments.

This change aligns with the long-term vision of making Conveyor CI modular and cloud-native.

---

## Current Architectural Constraints

The current API server design is tightly coupled with internal stateful components, including:

* Built-in log storage
* Embedded etcd
* Embedded NATS
* Embedded BadgerDB

Because these components are bundled and managed internally:

* The API server behaves as a stateful system.
* Horizontal scaling requires state replication.
* Distributed deployments introduce additional complexity.

This limits Conveyor CI to primarily centralized deployment patterns.

---

## Proposed Direction

### 1. Modularization

Refactor the API server into a reusable Go module (e.g., `server` package) with a clearly defined initialization and lifecycle API.

Example: (this is not necessarily how it should/wound be implemented rather just an example for inspiration)

```go
srv := conveyor.NewServer(conveyor.Options{
    LogStore:    customLogStore,
    MetadataDB:  externalDB,
    MessageBus:  externalNATS,
})

if err := srv.Start(ctx); err != nil {
    log.Fatal(err)
}
```

---

### 2. Introduce Abstraction Interfaces

Define interfaces for all stateful dependencies, including:

* Log storage backend
* Metadata storage
* Messaging backend
* Coordination layer

This allows:

* Plugging into external NATS clusters
* Connecting to external etcd clusters
* Using custom or enterprise log storage systems
* Replacing embedded storage with managed services

---

### 3. Stateless Mode

Decouple state management from the core API server logic, enabling:

* Multiple API server instances
* Horizontal scaling
* High-availability configurations
* Distributed system topologies

---

## Expected Benefits

* Enables embedding inside Go applications
* Reduces operational complexity
* Allows pluggable storage and messaging backends
* Supports horizontally scalable deployments
* Prepares Conveyor CI for enterprise environments
* Improves alignment with cloud-native architectural principles

---

## Scope

### Phase 1

* Extract API server into embeddable Go module
* Define and introduce required abstraction interfaces
* Maintain backward compatibility with standalone binary mode

### Phase 2 (Future Evaluation)

* Full stateless mode support
* Externalized storage and coordination by default
* Distributed scaling enhancements

---

## Compatibility

The existing standalone Conveyor CI binary will remain supported. The CLI entrypoint will internally instantiate the embeddable server module to ensure backward compatibility.

---

## Timeline

Roadmap update planned following CNCF Sandbox evaluation (expected February 25, 2026).

---

## Proposed By

[Infralane Cloud](https://www.infralane.cloud/)

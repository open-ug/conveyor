---
sidebar_position: 1
---

# What is Conveyor CI

<div align="center">

<img src="/logos/icon.svg" 
alt="logo"
height="200"
width="200"
 />

</div>

### The Headless Control Plane for Building CI/CD Platforms

Conveyor CI is not a CI tool. It is a control plane for building CI systems. If Jenkins and GitHub Actions are “all-in-one” platforms, Conveyor CI is the orchestration engine underneath them minus the UI, minus opinionated runners, minus infrastructure lock-in. You don’t use Conveyor to run `go build`. You use Conveyor to build a system that runs `go build` for thousands of users.

## Who Is This For?

Conveyor CI is built for _platform engineers_ and SaaS builders, not end users. It is ideal for:

* Internal Developer Platforms (IDPs) embedding CI/CD capabilities
* SaaS products offering “Deploy”, “Build”, or “Automation” features
* Edge & air-gapped environments requiring custom execution logic
* Organizations with strict security or infrastructure control needs

If you need complete control over where and how workloads execute, Conveyor is for you.

## Architecture: Control Plane vs Execution

Conveyor CI enforces a strict architectural boundary.

| Component                     | Responsibility                                                                       | Provided By |
| ----------------------------- | ------------------------------------------------------------------------------------ | ----------- |
| **Control Plane**             | Job lifecycle management, state persistence, queuing, retries, log aggregation, APIs | Conveyor CI |
| **SDK**                       | Typed interface for interacting with the engine                                      | Conveyor CI |
| **Drivers (Execution Logic)** | The code that actually runs workloads                                                | **You**     |

Conveyor does **not** ship Drivers, you define execution. This separation gives you, Infrastructure independence, Security model control, Hardware flexibility, Multi-environment support.

## Core Capabilities

Conveyor CI provides the orchestration primitives required to build serious CI systems:

* Job Lifecycle Management: creation, scheduling, retries, cancellation.
* State Persistence: durable workflow state tracking
* Event-Driven Execution: reactive architecture for scaling systems.
* Real-Time Log Streaming: aggregated and streamed logs.
* Declarative Workflows: define pipelines as structured configurations.
* API-First Design: fully accessible via HTTP.
* Embeddable Runtime: small Go binary, minimal resource footprint.

No UI. No opinionated runners. No vendor lock-in.

Just orchestration.

## Why This Design?

Most CI/CD platforms bundle UI, Execution environment, Orchestration logic, Storage, Scheduling, Opinionated infrastructure assumptions. Conveyor extracts the **orchestration layer** and leaves execution to you. This enables:

* Running on edge devices
* Air-gapped military networks
* Custom hardware
* Multi-cloud abstraction
* Fully white-labeled CI platforms

It’s the missing middle layer between Infrastructure and Developer-facing CI platforms.

## Example Use Case

You are building an Internal Developer Platform.

You want:

* Teams to push code
* Pipelines to trigger
* Jobs to execute inside Kubernetes
* Logs to stream back to your web dashboard

With Conveyor:

* Conveyor handles orchestration, state, and events.
* Your Kubernetes driver handles execution.
* Your UI consumes Conveyor’s API.

Clean separation. Clear ownership.

To install Conveyor CI, checkout the [Installation Page](/docs/installation)

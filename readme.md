
<h1 align="center" style="border-bottom: none; height: 200px;">
    <a style="height: 200px; max-width: 200px;" href="https://conveyor.open.ug" target="_blank">
        <img alt="Conveyor CI Logo" src="https://conveyor.open.ug/logos/logo.svg" style="height: 200px; max-width: 200px;">
    </a>
</h1>

<div align="center">

![Docker Pulls](https://img.shields.io/docker/pulls/openug/conveyor.svg?maxAge=604800)
[![Go Report Card](https://goreportcard.com/badge/github.com/open-ug/conveyor)](https://goreportcard.com/report/github.com/open-ug/conveyor)
[![License](https://img.shields.io/github/license/open-ug/conveyor.svg)](https://github.com/open-ug/conveyor/blob/main/LICENSE)
[![GitHub release](https://img.shields.io/github/v/release/open-ug/conveyor)](https://github.com/open-ug/conveyor/releases)
[![Maintainability](https://qlty.sh/badges/229750f3-9423-4ea6-8528-8e0f8cf854b5/maintainability.svg)](https://qlty.sh/gh/open-ug/projects/conveyor)
[![OpenSSF Best Practices](https://www.bestpractices.dev/projects/10999/badge)](https://www.bestpractices.dev/projects/10999)

</div>

---

# Conveyor CI

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

Conveyor does **not** ship drivers, you define execution. This separation gives you, Infrastructure independence, Security model control, Hardware flexibility, Multi-environment support.

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

### Why This Design?

Most CI/CD platforms bundle UI, Execution environment, Orchestration logic, Storage, Scheduling, Opinionated infrastructure assumptions. Conveyor extracts the **orchestration layer** and leaves execution to you. This enables:

* Running on edge devices
* Air-gapped military networks
* Custom hardware
* Multi-cloud abstraction
* Fully white-labeled CI platforms

It’s the missing middle layer between Infrastructure and Developer-facing CI platforms.

## Installation

Conveyor CI is Linux-first and distributed as a lightweight Go binary.

```bash
curl -fsSL conveyor.open.ug/install | sh
```

More installation options and configuration guides are available in the [official documentation](https://conveyor.open.ug/docs/installation).


### Example Use Case

You are building an Internal Developer Platform. You want:

* Teams to push code
* Pipelines to trigger
* Jobs to execute inside Kubernetes
* Logs to stream back to your web dashboard

With Conveyor:

* Conveyor handles orchestration, state, and events.
* Your Kubernetes driver handles execution.
* Your UI consumes Conveyor’s API.

Clean separation. Clear ownership.

### Contributing

Please 🌟 star the project if you like it.

Contributions are welcome! Please read [contributing guide](https://conveyor.open.ug/docs/contributing/how-to-contribute) and follow the governance model to submit PRs, issues, or feature requests.

### License

Apache License 2.0, see [LICENSE](./LICENSE).
Copyright © 2024 - Present, Conveyor CI Authors.

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fopen-ug%2Fconveyor.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fopen-ug%2Fconveyor?ref=badge_large)

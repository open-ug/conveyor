
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
[![GitHub issues](https://img.shields.io/github/issues/open-ug/conveyor)](https://github.com/open-ug/conveyor/issues)
[![GitHub pull requests](https://img.shields.io/github/issues-pr/open-ug/conveyor)](https://github.com/open-ug/conveyor/pulls?q=is%3Aopen+is%3Apr)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fopen-ug%2Fconveyor.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fopen-ug%2Fconveyor?ref=badge_shield)
[![Maintainability](https://qlty.sh/badges/229750f3-9423-4ea6-8528-8e0f8cf854b5/maintainability.svg)](https://qlty.sh/gh/open-ug/projects/conveyor)
[![Code Coverage](https://qlty.sh/badges/229750f3-9423-4ea6-8528-8e0f8cf854b5/test_coverage.svg)](https://qlty.sh/gh/open-ug/projects/conveyor)
[![OpenSSF Best Practices](https://www.bestpractices.dev/projects/10999/badge)](https://www.bestpractices.dev/projects/10999)

</div>

---

# Conveyor CI

**Conveyor CI** is an open-source, **lightweight, embeddable engine for building distributed CI/CD systems**.  

It provides a modular toolkit of **APIs, SDKs, drivers, and a Go-based runtime** so you can build your own CI/CD platform without reinventing execution, event handling, scaling, and observability.

---

## What Conveyor CI *is* and *is not*

**What it *is***  

- An embeddable CI/CD engine you can integrate into your platform or run standalone.  
- Lightweight, **Go-based single-binary runtime** with built-in observability.  
- Event-driven (powered by NATS JetStream + etcd) for dynamic, responsive execution.  
- Ideal for **self-hosted, air-gapped, or custom platform environments**.  

**What it *is not***  

- Not a hosted SaaS like GitHub Actions — you deploy and manage it.  
- Not Kubernetes-dependent — it runs without clusters or CRDs.  
- Not a plugin-heavy monolith like Jenkins — it’s an engine, not a marketplace.  .
- Not a visual pipeline-builder or fully-fledged pipeline solution out of the box — it provides the **engine** to build pipelines, not a complete pre-built CI/CD platform

---

## Key Features

- **Embeddable Engine:** Run pipelines anywhere — VMs, containers, or edge devices.  
- **API-first & SDK-ready:** Programmatically control pipelines, triggers, and drivers.  
- **Built-in Observability:** Metrics, logging, and tracing integrated out-of-the-box.  
- **Declarative Pipelines:** Define workflows using flexible YAML or programmatic APIs.  
- **Event-driven Architecture:** Real-time, responsive execution across distributed systems.  
- **Horizontal Scaling:** Scale drivers and workloads effortlessly with minimal setup.  
- **Custom Drivers:** Execute tasks via Docker, SSH, systemd, or custom runtimes.  

---

## Installation

Conveyor CI is distributed as an OCI container and available on [Docker Hub](https://hub.docker.com/r/openug/conveyor). and also available as a binary. You can install it by running this command.

```sh
curl -fsSL conveyor.open.ug/install | sh
```

---

## Use Cases

Conveyor CI is ideal for teams and organizations that need a flexible, embeddable CI/CD engine to power custom platforms or workflows:

* **Internal Developer Platforms (IDPs):** Integrate CI/CD into your own developer platform with full control over pipelines and drivers.
* **Custom Mobile App CI/CD:** Build a specialized CI/CD platform tailored for mobile app development, including custom workflows for build, test, and release.
* **Enterprise Automation:** Automate internal workflows that go beyond standard software pipelines, e.g., deployment orchestration or data pipelines.
* **Self-hosted / Air-gapped environments:** Run pipelines fully within your infrastructure without relying on external CI/CD SaaS providers.
* **Edge / IoT Deployments:** Execute pipelines on distributed nodes or remote environments where lightweight, embeddable execution is required.
* **Experimentation and Innovation:** Quickly prototype new CI/CD workflows or custom pipeline behaviors using Conveyor’s API and driver model.

---

## Why Conveyor CI?

Compared to other CI/CD tools:

| Feature / Tool           | Conveyor CI                                | GitHub Actions                   | Jenkins                  | Argo CD                          |
| ------------------------ | ------------------------------------------ | -------------------------------- | ------------------------ | -------------------------------- |
| Embeddable               | ✔                                          | ❌                                | ⚠️ Possible               | ❌                                |
| Kubernetes Required      | ❌                                          | ❌                                | ❌                        | ✔                                |
| API-first                | ✔                                          | ⚠️ Partial                        | ⚠️ Partial                | ✔                                |
| Self-hosted / Air-gapped | ✔                                          | ❌                                | ✔                        | ⚠️ Limited                        |
| Observability Built-in   | ✔                                          | ⚠️ Partial                        | ⚠️ Plugin-dependent       | ✔                                |
| Best Use Case            | Build custom CI/CD platform, Edge, on-prem | Repo workflows, SaaS convenience | Highly custom automation | GitOps deployments on Kubernetes |

> Conveyor CI is **engine-first**, built for platform teams and technical operators who want full control without vendor lock-in or heavy dependencies.

---

## More Information

* [Official Documentation](https://conveyor.open.ug) — Architecture, SDK, and driver development.
* [GitHub Repository](https://github.com/open-ug/conveyor) — Open-source code and releases.
* [CONTRIBUTING.md](./CONTRIBUTING.md) — How to contribute.
* [Project Governance](https://conveyor.open.ug/docs/contributing/governance)

---

## Contributing

Contributions are welcome! Please read [CONTRIBUTING.md](./CONTRIBUTING.md) and follow the governance model to submit PRs, issues, or feature requests.

---

## License

Apache License 2.0, see [LICENSE](./LICENSE).
Copyright © 2024 - Present, Conveyor CI Contributors.

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fopen-ug%2Fconveyor.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fopen-ug%2Fconveyor?ref=badge_large)

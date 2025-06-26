<h1 align="center" style="border-bottom: none; height: 200px;">
    <a style="height: 200px; max-width: 200px;" href="https://conveyor.open.ug" target="_blank"><img alt="Prometheus" src="https://conveyor.open.ug/img/logo.png"
    style="height: 200px; max-width: 200px;"></a>
</h1>

<div align="center">


![Docker Pulls](https://img.shields.io/docker/pulls/openug/conveyor.svg?maxAge=604800)
[![Go Report Card](https://goreportcard.com/badge/github.com/open-ug/conveyor)](https://goreportcard.com/report/github.com/open-ug/conveyor)

</div>

---

Conveyor CI is a(the first) Software Framework for building CI/CD Platforms. It provides a set of tools(programs) and SDKs that abstract the complexities of building CI/CD Platforms by implementing common features that CI/CD platforms require.

Key features that Conveyor CI provides include:

- **Built-In Observability**: Conveyor provides seamless integration of **Metrics**, **Tracing**, and **Logging**, so you can monitor, trace, and debug workflows at every stage without external setups
- **Declarative API for CI/CD Workflows**: Conveyor CI and its API adopt a decelarative and this helps yor design pipelines that are highly customizable and extensible
- **Realtime Event System**: Conveyor’s powerful event-driven architecture lets drivers publish and respond to events in realtime, enabling fast, responsive execution of tasks across your system.
- **Effortless Horizontal Scaling**: Conveyor’s runtime enables automatic horizontal scaling of drivers across distributed systems, optimized for cloud-native deployments all with zero extra code.
- **Realtime Log Management**: Conveyor CI Drivers include a built-in Logger that supports realtime log streaming and storage.

## Installation

Conveyor CI is distributed as an OCI container and can be found on Docker Hub, however considering the fact that it has dependency containers which are etcd, loki, and nats. A standard deployment configuration was created for docker compose and Helm charts are comming soon.

To Install it on docker compose you can head over to the Releases page and download `compose.yml` and `loki.yml` or on a linux system you can download them using curl.

```sh
curl -s https://api.github.com/repos/open-ug/conveyor/releases/latest | grep browser_download_url | grep compose.yml | cut -d '"' -f 4 | xargs curl -L -o compose.yml

curl -s https://api.github.com/repos/open-ug/conveyor/releases/latest | grep browser_download_url | grep loki.yml | cut -d '"' -f 4 | xargs curl -L -o loki.yml
```

Next start the containers using docker compose.

```sh
docker compose up

# OR

docker compose up -d
```

The Conveyor API Server will be reachable on [http://localhost:8080](http://localhost:8080)

## More information

For more information you can check out the [official documentation](https://conveyor.open.ug).

## Contributing

Refer to [CONTRIBUTING.md](./CONTRIBUTING.md)

## License

Apache License 2.0, see [LICENSE](./LICENSE).

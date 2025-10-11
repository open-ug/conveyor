---
sidebar_position: 1
---

# Installation & Setup

The Conveyor CI engine is composed of a set of software components. This is because of its highly modular design. This can sometimes mean that installing can get tedious since each system requires its own specific configuration. However, the official team provides a simplified set of installation options each configured to run in a specific environment. These are the available options.

## Install using Docker Compose

Ensure you have Docker installed on your environment.

First you will have to download the configuration files. These include a compose file and the default Loki configuration. Head to the [Releases page](https://github.com/open-ug/conveyor/releases) and download the `compose.yml` and `loki.yml` file in the release assets. You can also run this command on a UNIX system to download the latest release files.

```bash
curl -s https://api.github.com/repos/open-ug/conveyor/releases/latest | grep browser_download_url | grep compose.yml | cut -d '"' -f 4 | xargs curl -L -o compose.yml

curl -s https://api.github.com/repos/open-ug/conveyor/releases/latest | grep browser_download_url | grep loki.yml | cut -d '"' -f 4 | xargs curl -L -o loki.yml
```

Once you have all the files, you can then start the system using Docker Compose

```sh
docker compose up -d
```

And *Voila!!*, Conveyor is up and running

## Install using Helm (Kubernetes)

Helm charts are available for deploying Conveyor CI to Kubernetes. The Helm chart can be found in the `/helm` directory of the repository.

Ensure you have Helm installed on your system and access to a Kubernetes cluster.

```bash
# Add the Conveyor Helm repository (once charts are published)
# Or install directly from the repository
helm install conveyor ./helm
```

For more details on Helm installation and configuration options, refer to the Helm chart documentation in the `/helm` directory.

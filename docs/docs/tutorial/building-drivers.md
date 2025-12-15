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


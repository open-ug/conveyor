---
sidebar_position: 2
---

# The Architecture

Lets begin by exploring what the high level architecture of the system will look like.

From the project overview, we can break down the tasks that the system has to do in steps:

- Setup an environment to build the application
- Fetch source code from a Git repo
- Build the Application

Conveyor CI provides intutive primitives that we can use to accomplish this and these are, Resources, Drivers and Pipelines.

- **Resource**: The resource is a primitive that defines the state of the artifact you are working on. In our scenario it can contain the Git repository and other confighurations like environment variables etc.
- **Driver**: These are programs whose logic carries out the actual execution of the of the CI/CD processes. For example setting up a container environment in which the flutter build will occur and then running the build process.
- **Pipelines**: A pipeline defines the order in which drivers will be triggered, for example: from the steps we mentioned above, we can have a driver for each step and then define a pipeline that outlines the order that the drivers will follow to execute.

## The Resource

So our resource should include the following

- The Github Repository
- The environment variables

## The Drivers and Pipeline

The drivers we shall use include and will be ordered the following pepeline:

- A driver to create and start the build container
- A driver to clone the repository in the container
- A driver to trigger the build process
- A driver to stop and delete the container

## Workflow

Hence the resultant workflow should look something like this:

1. A user creates a resource object that includes the Git repository and required environment variables
2. They then send this resource to the Conveyor CI API server
3. The Conveyor CI API server saves the resource and uses the predefined pipeline information to orchestrate the resource to all the drivers
4. When the drivers recieve the resource event, they then do their work respectively.

Below is a visual representation of this workflow

![](./conveyor-ci-tutorial-system-arch.png)

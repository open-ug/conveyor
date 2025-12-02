---
sidebar_position: 3
---

# Triggering a Workflow

Once you have the Driver installed and running. You can trigger a CI/CD process.

CI/CD processes are triggered using [Resources](/docs/concepts/resources). These are objects that define how the CI/CD process will be run.

## Creating the Resource

Drivers are responsible for defining the schema and format of the resources they manage and you have to follow this schema to create your resource. The [open-ug/simple-runner](https://github.com/open-ug/simple-runner) defines a resource that requires an image and a set of commands to be run in the image. Lets define an example resource

```json
{
  "name": "ubuntu-pipeline-6",
  "resource": "pipeline",
  "spec": {
    "image": "jimjuniorb/hello-node:1.2.3",
    "steps": [
      {
        "name": "print-working-directory",
        "command": "pwd"
      },
      {
        "name": "list-files",
        "command": "ls -l"
      },
      {
        "name": "show-os-info",
        "command": "cat /etc/os-release"
      }
    ]
  }
}
```

Once the resource is created, it's sent to the API Server using a POST request to the `/resources/`. The API Server will return a response containing `runid` field that contains a UUID. The response you get is in this format

```json
{
  "message": "Resource created successfully",
  "name": "ubuntu-pipeline-6",
  "runid": "fbdf5edf-66d4-46ce-859a-8ee44d6c9463"
}
```

## Viewing the CI/CD process

Once you have posted the resource, you can view the progress of the process as it happens in real time by streaming its logs via a websocket.

To stream the progress, use a websocket client to connect to the `ws://localhost:8080/logs/streams/<driver-name>/<runid>`. In this case the driver name is `command-runner`. An example would be

```sh
ws://localhost:8080/logs/streams/command-runner/fbdf5edf-66d4-46ce-859a-8ee44d6c9463
```

The progress will be in the format of a timestamp and a log message.

---

Having completed this tutorial, you have a highlevel understanding of Conveyor CI usage workflow. You can now continue to the [Concepts Documentation](/docs/category/concepts) to gain a Deep understanding of how Conveyor CI works.

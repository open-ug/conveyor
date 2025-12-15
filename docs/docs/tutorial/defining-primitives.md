---
sidebar_position: 3
---

# Defining Primitives

Next we need to define the Resource Schema and Pipeline strucuture


## Defining the Resource Schema

Conveyor CI allows flexibility in the structure of your resources, this is because CI/CD artifacts differ. To define the structure of a resource we use a primitive called [Resource Definitions](/docs/concepts/resource-definitions.md). This is a JSON object that specifies the structure of the Resource.

Our Resource needs two main properties a Git Repo and environment variables so its Resource definition will look like this.

```json
{
  "name": "flutter-builder",
  "version": "0.0.1",
  "schema": {
    "type": "object",
    "properties": {
      "repository": {
        "type": "string"
      },
      "env": {
        "type": "array",
        "items": {
          "type": "object",
          "properties": {
            "name": {
              "type": "string"
            },
            "value": {
              "type": "string"
            }
          },
          "required": ["name", "value"]
        }
      }
    },
    "required": ["repository", "env"]
  }
}
```

This will create a Resource type called `flutter-builder` and an example of a resource object using this definition would look like

```json
{
  "name": "build-app",
  "resource": "flutter-builder",
  "spec": {
    "repository": "https://github.com/org/repo",
    "env": [
      {
        "name": "EXAMPLE_ENV",
        "value": "xxxxxxxxxxxxxxx"
      }
    ]
  }
}
```

Once you have created the resource definition, you have to send a `POST` request to the `/resource-definitions/` on the Conveyor API server to save it.

## Defining the Pipeline

As mentioned before, we need to define the pipeline which specifies the order of execution our drivers will follow.

Pipelines are defined using JSON object which includes properties like `resource` and `steps` that the define the resource type that this pipeline applies to and the order in which the drives should execute respectively.

In our use case, our pipeline would look like this:

```json

{
  "name": "flutter-pipeline",
  "resource": "flutter-builder",
  "description": "A pipeline to build a flutter codebase into an APK",
  "steps": [
      {
        "id": "1",
        "name": "Start Build Environment",
        "driver": "container-start"
      },
      {
        "id": "2",
        "name": "Clone Git Repository",
        "driver": "git-cloner"
      },
      {
        "id": "3",
        "name": "Trigger Build Process",
        "driver": "builder"
      },
      {
        "id": "4",
        "name": "Stop and Delete container",
        "driver": "container-stop"
      }
  ]
}
```

In the above pipeline definition, each object in the `steps` array specifies a driver that will be executed with `id` standing for the order, `name` is a human readable name of the driver and `driver` is a unique driver identifier that the driver defines itself as.

Once you have created the pipeline definition, its then sent to the `/pipelines/` endpoint on the Conveyor CI API server using a `POST` request.

Lets move on to building the drivers.
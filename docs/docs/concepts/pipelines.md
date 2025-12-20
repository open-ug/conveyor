---
sidebar_position: 5
---

# Pipelines

Pipelines in Conveyor CI are objects that define the order of execution that drivers have to follow whenever a resource state changes.

They contain only one main propetry, `steps` that defines the hererical order that drivers will follow when executing.

## How Pipelines works

When a new resource is created and sent to the Conveyor CI API server, all drivers that are configured to listen to it immediately reconcile it, In most CI/CD scenarios, this is not the ideal functionality or workflow required, they mostly require some form of ordered steps to follow forexample you might need to clone the source code, run tests, build and complie the program, generate a docker image, and then push it to registery. Incase you have separate drivers that carry out all these actions, it wont be ideal to have this all those drivers run at once, rather have them run in a predefined order since some actions depend on others having been done completely and succesfully.

And thats why pipelines where introduced in Conveyor CI, they mainly just define the order of execution that drivers will follow.

## Creating and Using a Pipeline

Lets look at how to create and use a pipeline.

From the CI/CD scenario mentioned above where we are trying to build a docker image and buish it to a registry. We can have the following drivers that carry out those functions:

- `git-cloner`: A driver that clones the source code.
- `test-runner`: A driver that runs test.
- `builder`: A driver that builds and compiles the program.
- `image-generator`: A driver that builds the docker image.
- `image-pusher`: A driver that pushes the docker image to a registry.

Once you have the list of drivers, you can define the pipeline. ALl you need is a:

1. Unique name of the pipeline
2. The resource type that this pipeline apply to
3. A desription of the pipeline
4. The hererical order that the drivers will follow

With this information, you can define the pipeline using JSON like this:

```json
{
  "name": "docker-pipeline",
  "resource": "docker-resource",
  "description": "A pipeline to build and publish a docker image",
  "steps": [
      {
        "id": "1",
        "name": "Clone git repository",
        "driver": "git-cloner"
      },
      {
        "id": "2",
        "name": "Run tests",
        "driver": "test-runner"
      },
      {
        "id": "3",
        "name": "Build and Compile",
        "driver": "builder"
      },
      {
        "id": "4",
        "name": "Generate docker image",
        "driver": "image-generator"
      },
      {
        "id": "5",
        "name": "Publish to registry",
        "driver": "image-pusher"
      }
  ]
}
```

You can then save this pipeline to the Conveyor API server and when you create a resource all you have to do is to specify the `pipeline` field in the reosurce object and set it to the pipeline name. as follows

```json showLineNumbers
{
  "name": "build-app-1",
  "resource": "docker-resource",
  // highlight-next-line
  "pipeline": "docker-pipeline",
  ...
}
```

Once you send the resource to the API server. Conveyor CI will follow the order you specified in the pipeline to send events to the drivers to trigger the `Recocile` function.

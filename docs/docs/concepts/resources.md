---
sidebar_position: 1
---

# Resources

:::danger

Documentation still in early stages

:::


A Resource is the core object in a Conveyor CI Platform that defines the pipeline configuration. It tells the CI Driver what to do, when to do it and how to do it. A comparison would be the Workflow file in the GitHub Actions.

## Resource Definitions

By default, Conveyor CI does not come with any Resource in its API. You have to install them into Conveyor CI. To install a Resource, you must define a Resource Definetion. A Resource defenition is method of defining the Schema of a Resource. It determines how the Resouce will be defined, what properties it will have and also the validation schema for the resource.

Resource definitions are created by following the [Open API Specification](https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.3.md). They can be defined either in Yaml or JSON syntax and saved into the Conveyor CI API. Resource Definitions have only two mandatory requirementsa and these are:

An example of a Resource definition in yaml.

```yml
name: workflow
version: 0.0.1
schema:
  pipeline:
    type: object
    properties:
      name:
        type: string
      stages:
        type: array
        items:
          type: string
      distributed:
        type: boolean
      runners:
        type: array
        items:
          type: string
    required:
      - name
      - stages
      - distributed
      - runners
```

This above Resouce Definition defines a Resource named `application` and it can be used like this below.

```yml
name: example-resource
resource: workflow
pipeline:
  name: build-and-deploy
  stages:
    - test
    - build
    - deploy
  distributed: true
  runners:
    - cloud-native
```

A Resouce Definition have three mandatory fields.

- `name`: This defines the name of the Resource.
- `version`: This defines the version of the Resouce Definition
- `schema`: This contains the Open API Schema used to validate the resource

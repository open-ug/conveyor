---
sidebar_position: 2
---

# Resources

Resources in Conveyor CI are objects that store the state of CI/CD processes. They define information that is used by Drivers to know what to do and how to do it.

The data stored by resources does not follow any standard format. Its format must be predefined using [Resource Definitions](resource-definitions). These are other internal objects that define the syntax of Resource objects.

## Parts of a Resource

Resource objects contain specific fields that are required when creating them. These parts include:

- **Name**: The name field is used to specify the name of the Resource. It is meant to be unique across different Resources. It is an alpha-numerical string field that does not include spaces or any special characters except “-”
- **Resource**: The resource field is used to specify what Resource Definition schema to use. It corresponds and has to be equal to the name of an existing Resource Definition.
- **Spec**: The spec field contains the resource data. The data must follow the convection defined in the Resource Definition.

## Resources workflow

When a Resource is created, Conveyor CI stores it in the data store and then sends an event to the drivers that are associated with that Resource. The Drivers then read the Resource spec and use it to carry out executions depending on the Resource data.

## Creating a Resource

To create a Resource you have to follow a few steps:

- First you should create and register a Resource Definition for your resource or ensure one already exists.
- Then write the schema of your resource and send it to the API Server via a POST request to the `/resources/` route.
- This will validate the resource and save it to the data store.

An example of a resource is here

```json
{
  "name": "example-resource",
  "resource": "workflow",
  "spec": {
    "pipeline": {
      "name": "build-and-deploy",
      "stages": [
        "test",
        "build",
        "deploy"
      ],
      "distributed": true,
      "runners": [
        "cloud-native"
      ]
    }
  }
}
```

## Deleting resources

`DELETE /resources/{type}/{name}` returns `202` for a resource with a pipeline. The resource remains available while cleanup runs. Clients poll `GET /resources/{type}/{name}`: `404` confirms completion, while `metadata.deletion_failed: true` indicates a cleanup failure that can be retried. The delete operation reaches every pipeline stage; drivers must implement cleanup for this event. Versioned resource snapshots are removed only after the final successful stage. Updates are refused while a resource is being deleted, and stale queued pipeline stages are discarded.

Direct resources without a pipeline retain their `204` deletion behavior. A missing resource returns `404`. `/capabilities` advertises `resource_deletion_pipeline: 1` for clients that require the completion protocol. Deletion retries restart failed cleanup; a stalled operation can restart after two minutes.

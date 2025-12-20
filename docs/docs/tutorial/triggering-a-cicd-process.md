---
sidebar_position: 5
---

# Triggering a Workflow

Lets trigger a sample workflow from the system we have built.

Before we trigger a workflow, you need to ensure you have done the folowing tasks that where explained and demosntrated in the previous sections:

- Installed and started the Conveyor CI API server
- Defined the resource definitiaon and posted it to the API Server
- Defined the Pipeline object and Posted it to the API server
- Cloned the driver repository and started all 4 driver intances

Once you have the above tasks done. You can find a simple flutter application code base that is on a remote Git repository and use it to create a resource object. similar to this:

```json
{
  "name": "build-app",
  "resource": "flutter-builder",
  "pipeline": "flutter-pipeline",
  "spec": {
    "repository": "https://github.com/AmirBayat0/Sneakers-shop-app-Flutter",
    "env": []
  }
}
```

You can then send this JSON to the Conveyor CI API Server using a POST request to `/resources` route. This will save the resource in the database and send an event to the drivers in order depending on the order you defined in the Pipeline.

The logs that are collected from the Driver by the Driver logger from this process can be streamed or collected from the API server in realtime.

To do this, you simply have to open a Websocket connection to the `/logs/streams/{DROVER_NAME}/{RUN_ID}` route i.e `ws://localhost:8080/logs/streams/sample-driver/fbab31f6-a278-4a8f-96be-ac49b007ca65`.

---

Having completed this tutorial, you have a highlevel understanding of Conveyor CI usage workflow. You can now continue to the [Concepts Documentation](/docs/category/concepts) to gain a Deep understanding of how Conveyor CI works.

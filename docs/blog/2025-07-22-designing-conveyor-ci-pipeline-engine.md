---
slug: designing-the-conveyor-ci-pipeline-engine
title: Designing The Conveyor CI Pipeline Engine
authors: jim-junior
tags: [conveyor-ci, update, pipeline]
---

In CI/CD Systems, there is a concept of pipelines. Pipelines define the steps of how code changes throughout the entire CI/CD process. Multiple systems like Github Actions, GitLab CI/CD, Jenkins etc all have this functionality of pipelines natively embedded within them, something that [Conveyor CI](https://conveyor.open.ug/) currently lacks. You could engineer a walkaround to implement your own pipeline-like functionality in Conveyor but natively, this functionality does not exist and this is a serious downside existing within Conveyor.

<!-- truncate -->

This got me thinking, how could I implement it. As you might guess from the name, Conveyor CI was inspired by how a [conveyor system](https://en.wikipedia.org/wiki/Conveyor_system) in an industry works. Take an example of a car manufacturing factory, the skeleton of the car moves along the conveyor system and at each step there is a dedicated robot in charge of attaching a certain component to the car to carry out a certain function. By the time the skeleton reached the final stage, its a complete car. From the start I wanted to adopt this kind of concept as the core paradigm of Conveyor CI whereby a resource(e.g. source code) moves on a belt and at the end of the belt, all necessary actions are carried out.

With this in mind, I split the Conveyor CI into components in order to easily understand it adopting from components of a conveyor system ie. Package, conveyor belt, and Peripheral equipment. In that I came up with the following:

- **Resource**: The resource is an internal object in Conveyor that represents what is being acted upon throughout the CI/CD process. Think of it as the Package on a conveyor system. It can represent anything from source code, and application, a program etc.
- **Drivers**: Drivers are software components that carry out certain actions depending on the state defined in the Resource. Think of them as the peripheral devices that act on a package as it moves along a conveyor system. In this case the Packages are the Resources.

Notice I have not mentioned a component that corresponds to the conveyor belt. That is because currently there was none, and to create one was where pipelines come in.

So currently in Conveyor CI we have Resources and Drivers. They are generally mature enough as individual components and although you could create an entire CI/CD tool with these only, there are still some issues. Mainly is that there is no order of execution of drivers upon a Resource. This means that once a Resource event occurs, all drivers execute their corresponding actions whenever they receive the event and they do this in no orderly fashion. This means it is not possible to predict what action will occur at what point in the CI/CD process and you also can’t create dependency actions that depend on others being pre-executed. An example is if you have a workflow whereby you compile your source code then upload the output program to a distribution server, You might have two drivers, one for compiling and another for uploading. In the current implementation, these drivers will execute at once, yet the uploading driver should depend on the constraint that the compiling driver is done executing.

To fix this we have to come up with a way to define a workflow that drivers must follow when carrying out their executions, something that defines the steps and order that these drivers will follow throughout execution. Something to act as the conveyor system. This is where Pipelines Come in. They will define the order of execution followed by drivers. The pipeline can be an object containing the configuration defining the order followed by drivers.

Pipelines also introduce more possibilities like shared context among driver executions, meaning, a pipeline can define some metadata that is shared and used across all the drivers that are executing upon a resource in that pipeline.

## Implementing Pipelines

Now that I had come up with a high level design of what pipelines are expected to work, I had to move on to designing an implementation that would easily fit in into the existing Conveyor CI implementation. Inorder to create a implementation that is good, i set a few constraints that i had to follow:

1. Pipelines should not introduce breaking changes that might interfere with the Developer experience and development paradigm for developers developing drivers. It should be more of an incremental change, building upon the already existing development paradigm without breaking already existing codebases.
2. Minimize the risk of over engineering the system while trying to follow constraint one. This is because usually one trade off of maintaining a good developer experience is that you might end up over engineering the system.

With these design constraints in mind, I embarked on designing the Pipeline engine. First I started with the usage workflow from the user viewpoint i.e. the process that will be followed by a developer when creating a pipeline and using it.

### Use case Workflows

Considering that developer experience is one of the most important concepts for the success of a tool, I had to come up with a really easy to understand workflow and also not break the already existing usage workflows in Conveyor CI. I came up with this.

- First a pipeline is designed by the developer choosing their appropriate drivers and resources and arranging them in a desired order.
- The pipeline is then sent to the Conveyor CI API Server to be registered and saved.
- Once its saved then a Resource that is using that pipeline is created and sent to the API server
- Then the Pipeline Engine appropriately routes the resource to the drivers.
- Driver do there work

### How the Engine Works

With this in mind, I came up with a system design that can accomplish these functions but also act as an incrementation on the existing Conveyor base.

First a new state object called a *Pipeline* was introduced into the system. This object would store and represent information about a pipeline throughout the stages of execution. It would store the order of execution followed by drivers and additional context information/metadata that is required by the driver in said pipeline.

I also had to introduce a *Driver Result* event object. This object, acting as an event, was to be used by drivers to communicate the result of the `Reconcile` function(The function that is run when a driver receives a resource to act upon). This would be used to inform the Pipeline Engine if a desired driver operation was successful or not.

In order to maintain realtime seamless execution, a new execution process running concurrently with the API Server had to be introduced to Conveyor CI's core runtime program. This would be in charge of listening for new resources attached to pipelines, routing those resources to drivers in a pre-defined order as defined in the *Pipeline* objects, and watching for pipeline realworld state changes and reconciling it to the state stored in the Database. This process is what I named the **Pipeline Engine**

#### Execution Workflow

Having introduced the required components, I came up with this workflow.

- The user sends a POST request to the Conveyor API server to register their Pipeline, this is then saved to the ETCD Database.
- Then the user creates a pipeline Resource(a resource that depends on that pipeline), also via the API Server.
- When the API server receives and saves this Resource to the DB, it sends an Event to the Pipeline Engine.
- The Pipeline Engine receives the Event and using the event metadata it collects the Pipeline information from the DB and depending on the driver execution order defined by the pipeline it begins to send events to the drivers.
- Once the Driver has finished to carry out its functions, it returns a *Driver Result* and the Driver Manager sends that result as an event back to the Pipeline Engine.
- When the Result event is received, the Pipeline engine can then move on to sending the event to the next driver. If no result event has been received yet from the Driver, the Engine wont move on to the next Driver.
- Lastly, if the result event on one driver indicates that the Driver execution failed or an error occurred. The pipeline engine won't send events to the remaining drivers and will rather register the pipeline as stopped or done with a status of failed. Else it will finish execution of all drivers and still register the pipeline as done/complete.

##### A Simple Diagram representing this process

![Pipeline Ecexution flow](https://jim-junior.github.io/img/pipeline-engine.png)

#### Technical Details

I have managed to describe the high level working of how the Pipeline engine would work. But beyond that, integrating such functionality into the Conveyor Stack requires clarification on how to implement the different components technically.

Starting off with the Pipeline engine itself, we said it's a process running concurrently with the API Server in the same program. Considering that the Conveyor CI core software program is written in the [Go programming language](https://go.dev/) and Go as a language has inbuilt concurrency via the [Go routines](https://www.geeksforgeeks.org/go-language/goroutines-concurrency-in-golang/), this means that the Pipeline Engine would better be a go routine running alongside the API server on the same process. One upside of this mechanism of embedding the Pipeline engine into the same process with the API Server is that we will require little updates to the metrics and monitoring codebase in order to add metrics for monitoring the Pipeline Engine.

Inorder to achieve the realtime nature of execution among these separate execution runtimes, we utilize events via [NATS Jetstream](https://docs.nats.io/nats-concepts/jetstream), the message broker that is already being used by Conveyor. This means that the Pipeline Engine performs both the roles of an Event publisher and subscriber. When new Pipeline resources are created, an event is sent to the Pipeline engine via JetStream. In this scenario, the Pipeline Engine is running as a [Jetstream Consumer](https://docs.nats.io/nats-concepts/jetstream/consumers), listening for events via the `pipeline` [stream](https://docs.nats.io/nats-concepts/jetstream/streams). Upon receiving an event, it will collect the pipeline data from ETCD and collect all the driver names. Using the driver names it will now become an event publisher and publish to the `messages` stream(the default stream in which drivers listen for events) and use the subjects with the following semantics `{DRIVER_NAME}.resources.{RESOURCE_NAME}` with the placeholders `{DRIVER_NAME}` and `{RESOURCE_NAME}` referring to the name of the driver to publish to and the name of the resource being published respectively. All we have to do now is update the driver managers of the drivers, to include the new subject in the filtered subjects field of there Jetstream consumers. THis will require updating all the SDKs though. With this, we have a complete working Pipeline engine.

Finally to keep track of Driver executions that belong to one pipeline, we utilize an already existing concept in Conveyor CI which is the `Run ID`. This is a *UUID* string that identifies individual driver runs/executions upon a resource. In a pipeline, we maintain the same run ID throughout executions in a single pipeline execution.

## Conclusion

The whole point of Pipelines was to introduce ordered execution of drivers on a resource but this feature can act as a fundamental building block for more robust and important functionality to be integrated into Conveyor CI. With this I have to acknowledge that Conveyor CI is still a really immature project, it really still has a long way to go and still lacks important features e.g. resource versioning etc. But one step at a time. Just like Rome wasn't built in a day, You can build a Cloud native software framework in a short period of time.

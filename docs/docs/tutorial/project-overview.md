---
sidebar_position: 1
---

# Overview

Let’s learn by example.

In this tutorial we shall learn how to use Conveyor CI to build a CI/CD platform, we’ll walk you through the creation of a basic Mobile DevOps CI/CD platorm.

## The product

Asssuming you are working at a startup that wants to build a CI/CD pltform that builds and packages flutter mobile application in the cloud. A platform like Expo EAS but focusing on flutter applications. The product works as follows;

- A user develops a mobile application in flutter.
- They provide the app code in a Git repository
- Our platform fetches the app code and builds an `apk` in the cloud
- The user can then download the outpput `apk` file for distribution.
---
sidebar_position: 1
---

# Overview

Let’s learn by example.

In this tutorial we shall learn how to use Conveyor CI to build a CI/CD platform. We’ll walk you through the creation of a basic Mobile DevOps CI/CD platform.

## The product

Assuming you are working at a startup that wants to build a CI/CD platform that builds and packages [Flutter](https://flutter.dev/) mobile application in the cloud. A platform like [Expo EAS](https://expo.dev/services) but focusing on [Flutter](https://flutter.dev/) applications. The product works as follows;

- A user develops a mobile application in flutter.
- They provide the app code in a Git repository
- Our platform fetches the app code and builds an [Android](https://www.android.com/) `apk` in the cloud
- The user can then download the output `apk` file for distribution.

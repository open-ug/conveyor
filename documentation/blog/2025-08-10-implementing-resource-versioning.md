---
slug: implementing-resource-versioning-in-conveyor-ci
title: Implementing Resource Versioning in Conveyor CI
authors: jim-junior
tags: [conveyor-ci]
---

Resource versioning is a key feature in Ci/CD Platforms. It provides developers a way to track how a pipeline has changed over time, ensuring reproducibility, stability, and traceability of builds. Conveyor CI on the other hand, doesn't have this inbuilt into it. This has therefore been a great downside of it, leading developers who might use it miss out on the above mentioned features.

<!-- truncate -->

## How it is designed in other systems

The design or implementation can vary depending on multiple cases like the purpose, use case, or even the internal design of the system. But at the core, you need to be able to differentiate multiple variations of the same process as it is re-triggered. This means each execution of a CI/CD process has to have a unique Identifier that differentiates it from other executions under that same resource. The most common identifiers or ways in which this differentiation is done include:

- **Semantic Versioning (SemVer)**: This is a famous software versioning scheme that follows the `MAJOR.MINOR.PATCH` eg `1.2.4`. It is easy to understand and commonly used in CI/CD systems like GitHub Actions, Gitlab CI and Jenkins.
- **Commit SHA Pinning**: This is a method commonly used in Git based CI/CD systems, whereby they pin the Git Commit to a version and one major advantage of this is it ensures absolute reproducibility.
- **Custom Revision Generations**: In this case a unique string is attached to a resource version. It doesn't have to have any semantic meaning but has to be unique.

## Intended Design in Conveyor CI

Conveyor CI being a minimal State Based system that uses a resource driven architecture, it won't base on any external system to implement versioning rather on system state. This means we are free to implement our own versioning scheme. The goal is to use a scheme that is intuitive on the human side and semantically.

With this in mind, we can opt for an incremental integer based versioning scheme. Meaning we just keep incrementing a positive integer starting from zero with each new resource version created.

## Technical Implementation

The proposed design can seem like a simple and easy concept to implement, but just like all technical problems, choosing the right implementation might not be as straight forward as it seems. We had to explore different ways to implement this.

A straightforward way might be that each time a resource is updated, a new resource is created with a resource version increased by one and its stored in the database, then a pointer is created so that when a user requests for the current/latest resource version we provide the recent update added.

But this introduces an expected bottle neck. We are assured to blot the database as resources get more revisions. For example, assuming we have 1000 resources and each resource has 100 revisions. That's automatically one hundred thousand records in the database. And say our system becomes big and he has around one million resources each with 100 revisions. The database becomes huge with 100 million records.

So we have to then take into consideration our data store and investigate if it's able to handle this form of incrementation. Conveyor CI uses [etcd](https://etcd.io/), a key-value store, it is reliable and highly performant. As we investigated further into the architecture of `etcd`, we realized that internally `etcd` uses *Multi-Version Concurrency Control (MVCC)* which allows reads at specific revisions of a record or key. This means for every update to an `etcd` key you make, `etcd` stores the previous version of that record and you can query or read it later on. So technically, rather than having to create a new record each time we update a resource, we can rely on `etcd`’s internal MVCC to store revisions. Another upside to this is that `etcd` also uses the incremental revision numbering format so we can also rely on that.

However, the designers of `etcd` also released that storing revisions, introduces database blot as the records keep increasing in number and as they keep being updated. So they designed a concept of compacting whereby the database will periodically delete old revisions and keep only the latest version of a record. But this now becomes a danger to the data integrity of Conveyor CI, meaning we wont be able to access old revisions since they will be permanently deleted. Luckily, this feature of compacting automatically can be turned off and compacting can be done manually by a system administrator. So we can utilize this manual compacting to our benefit and design our own custom compact strategy.

### Compact Strategy

To prevent the `etcd` from becoming extremely huge and also capitalize on the compact feature in `etcd`. We can introduce a custom compact strategy for systems that have managed to generate an extremely large database and performance is degrading. Our strategy works as follows:

- Enable manual `etcd` compacting
- Introduce an External Audit Storage.
- Upon compacting, we store a snapshot of the revisions in the External Audit storage system. Then compact `etcd`
- When a user tries to query extremely old revisions, we can then refer and collect them from the external audit system

This ensures that ou `etcd` datastore remains efficient without Key Explosion and avoids the performance degradation and storage inefficiency of millions of custom historical keys. The external Audit storage system in some cases can also work as a backup system.

The above mentioned implementation is currently the best we have and has actually been proved to be production ready by large systems like kubernetes, and it is what will be patched to Conveyor CI, unless proven otherwise or a better implementation is proposed.


---
sidebar_position: 1
---

# Contributing

Conveyor CI is an open-source project licenced under the [Apache License 2.0](https://github.com/open-ug/conveyor/blob/main/LICENSE). We welcome anyone who would be interested in contributing to urunc. As a first step, please take a look at the following document. The current document provides a high level overview of Conveyor CI's code structure, along with a few guidelines regarding contributions to the project.


## Code Organisation

Conveyor CI is written in Go and it exposes an Application Programming Interface whose client libraries can be written in any language. We structure the codebase and other files as follows:

- `/`: The root directory contains the non-code files, such as the licence, readme and conveyor's entry point `main.go`.
- `/documentation`: This directory contains a [Docusaurous Site](https://docusaurus.io/) for all the documentation related to Conveyor CI.
- `/cmd/cli`: Contains the entrypoint for the command line and definitions for all CLI commands.
- `/cmd/api`: Contains the entrypoint for the API Server.
- `/internal/`: Contains majority of the Conveyor CI codebase.
  - `/internal/config/`: This directory contains code responsible for loading and handling Conveyor CI configuration.
  - `/internal/handlers`: Contains code .
  - `/internal/routes`: Contains code defining API Server routes.
  - `/internal/models`: Contains code defining `etcd` data models.
  - `/internal/streaming`: Contains handling websocket streaming.
  - `/internal/metrics`: Contains code for collecting metrics for the API Server.
  - `/internal/utils`: This directory contains utility functions.
- `/pkg`: This contains the Publicaly available Conveyor CI SDK library
  - `/pkg/client`: Contains the Go API Client
  - `/pkg/driver-runtime/`: Contains Driver runtime and its utility finctions
- `/sdk`: This directory is meant to contain any Conveyor SDK implementations in other languages
- `/helm`: This contains the Conveyor Helm Chart for deploying to Kubernetes.

## How to contribute

There are plenty of ways to contribute to an open source project, even without changing or touching the code. Therefore, anyone who is interested in this project is very welcome to contribute in one of the following ways:

- Using Coveyor CI. Try it out yourself and let us know your experience. Did everything work well? Were the instructions clear?
- Improve or suggest changes to the documentation of the project. Documentation is very important for every project, hence any ideas on how to improve the documentation to make it more clear are more than welcome.
- Request new features. Any proposals for improving or adding new features are very welcome.
- Find a bug and report it. Bugs are everywhere and some are hidden very well. As a result, we would really appreciate it if someone found a bug and reported it to the maintainers.
- Make changes to the code. Improve the code, add new functionalities and make urunc even more useful.

## Opening an issue

We use Github issues to track bugs and requests for new features. Anyone is welcome to open a new issue, which is either related to a bug or a request for a new feature.

### Reporting Bugs

If you find a bug, you can help fix it by submiting and issue to the appropriate repository. within your issue Include:

- A clear and descriptive title
- Steps to reproduce the issue
- What you expected to happen
- What actually happened
- If the repository has an issue tempate. you should follow it.

Before submitting, check if the issue already exists in the repository issue list.

### Suggesting Enhancements

We also welcome feature suggestions and ideas for improvement. When submitting a feature request:

- Explain why the feature is useful
- Provide example scenarios where it would help
- Suggest a possible implementation, if you have one
- Try to keep requests focused and concise.


### Submitting Code Changes

Once you have identified and issue or an enhancement you would like to work on. you can folloe this workflow to submit your code changes.

1. **Fork** the repository
2. **Clone** your fork:

   ```bash
   git clone https://github.com/your-username/project-name.git
   cd project-name
   ```

3. **Create a new branch**:

   ```bash
   git checkout -b feature/your-feature-name 
   
   # OR 

   git checkout -b fix/your-fixture-name 
   ```

4. **Make your changes**
5. **Write or update tests**, if applicable
6. **Commit** your changes with clear messages:

    ```bash
    git commit -m "feat: add new feature X"
    # OR
    git commit -m "fix: solved issue Y"
    ```

7. **Push** your changes:

    ```bash
    git push origin feature/your-feature-name
    ```

8. **Open a Pull Request** on the GitHub repository and describe what you’ve done

Pull requests should be:

- Focused on a single change
- Thoroughly tested
- Aligned with project’s code style and guidelines

### Improving Tests or Documentation

Improving test coverage and documentation is highly valuable. You can:

- Add test cases for untested components
- Update outdated documentation
- Fix typos or formatting issues

No contribution is too small!

### Golang code styde

We follow gofmt's rules on formatting GO code. Therefore, we ask all contributors to do the same. Go provides the gofmt tool, which can be used for formatting your code.

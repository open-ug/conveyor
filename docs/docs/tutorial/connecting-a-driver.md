---
sidebar_position: 2
---

# Connecting a Driver

[Drivers](/docs/concepts/drivers) in Conveyor CI are the components that actually execute the CI/CD process. So once you have the Conveyor CI engine up and running, you have to connect a driver to it.

If you are building a Driver with the official SDKs, The Driver will connect to the Conveyor CI engine automatically. In this tutorial we use an example driver that runs commands inside a Docker Container. The drivers code base can be found in the [open-ug/simple-runner](https://github.com/open-ug/simple-runner) repository on Github.

## Installing the Driver

Assuming you already have Conveyor CI running. Head over to the Github releases of the [open-ug/simple-runner](https://github.com/open-ug/simple-runner) and download the `runner-linux-amd64`. You could also just download it with `curl`.

```sh
curl -s https://api.github.com/repos/open-ug/simple-runner/releases/latest | grep browser_download_url | grep runner-linux-amd64 | cut -d '"' -f 4 | xargs curl -L -o runner-linux-amd64
```

## Running the Driver

Once you have it downloaded, adjust its permission to enable execution. And then run the binary.

```sh
# Set Permissions
chmod +x runner-linux-amd64

# Run the Binary
./runner-linux-amd64

```

You should get output like this

```log
Driver Manager is running for driver:  command-runner
```

You can move to the next page about triggering a workflow or if you want to dive deep into driver development head over to the Driver Development Page

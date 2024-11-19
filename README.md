# Edge Leap

EdgeLeap is a development tool that streamlines IoT Edge module deployment. It integrates seamlessly with Azure IoT Edge and Azure IoT Hub to enable rapid local development and testing directly on edge devices, bypassing the traditional CI/CD pipeline. This allows developers to iterate quickly and validate their IoT Edge solutions in real device environments without the overhead of production deployment processes.

## Installing Edge Leap client

The `edge-leap` client can be installed in a few different ways:

1. **Compiling from source**

The `edge-leap` client can be compiled from source using the provided `Makefile`. By default, compiled binaries are placed in the `bin` directory, but this can be overridden by setting the `GO_BINDIR` environment variable for the `make` command.

- To compile the client for the host platform and architecture, run:

    ```shell
    make build
    ```

- To compile for all architectures for a given platform, run:

    ```shell
    make build-[platform] # where [platform] is one of [macos | linux | windows]
    ```


2. **Using the pre-compiled binaries**

Another option to install the `edge-leap` client is to use the pre-compiled binaries available in the [releases page](https://github.com/unbrikd/edge-leap/releases). For each release, the client is compiled for: macOS, Linux, and Windows.

3. **Using the Docker image**

If you prefer to use Docker, the `edge-leap` docker image can be pulled from the container registry:

```shell
docker run ghcr.io/unbrikd/elcli:latest <COMMAND>
```

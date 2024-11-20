# Edge Leap

EdgeLeap is a development tool that streamlines IoT Edge module deployment. It integrates seamlessly with Azure IoT Edge and Azure IoT Hub to enable rapid local development and testing directly on edge devices, bypassing the traditional CI/CD pipeline. This allows developers to iterate quickly and validate their IoT Edge solutions in real device environments without the overhead of production deployment processes.

## Installing Edge Leap client

### Compiling from source

The `edge-leap` client can be compiled from source using the provided `Makefile` rulesm and the output binary will be placed in the `bin` directory:

```shell
# compile for the current platform
make build

# compile for a specific platform: macos, linux or windows
make build-<platform>
```


### Using the pre-compiled binaries

Another option to install the `edge-leap` client is to use the pre-compiled binaries available in the [releases page](https://github.com/unbrikd/edge-leap/releases). For each release, the client is compiled for: macOS, Linux, and Windows.

### Using the Docker image

If you prefer to use Docker, the `edge-leap` docker image can be pulled from the container registry:

```shell
docker run ghcr.io/unbrikd/elcli:latest <COMMAND>
```

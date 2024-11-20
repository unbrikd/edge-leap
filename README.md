# Edge Leap

EdgeLeap is a development tool that streamlines IoT Edge module deployment. It integrates seamlessly with Azure IoT Edge and Azure IoT Hub to enable rapid local development and testing directly on edge devices, bypassing the traditional CI/CD pipeline. This allows developers to iterate quickly and validate their IoT Edge solutions in real device environments without the overhead of production deployment processes.

## Installing Edge Leap client

### Compiling from source

The `edge-leap` client can be compiled from source using the provided `Makefile` rules and the output binary will be placed in the `bin` directory:

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

## Usage

The `edge-leap` client is a CLI tool that can be used to manage the development stages of an IoT Edge module. The tool is designed to operate in two modes: `draft` and `release`.
Each mode has its own set of subcommands and flags, that can be found by using the `--help` flag for each subcommand.

### Draft mode

The `draft` mode is used to manage development sessions. It allows developers to provide a configuration of the development environment and automatically handle the required operations, in order to build and deploy the module to the target device.

- `elcli draft new`: Initialize a new development session by creating a configuration file in the current directory.

  First you need to initialize a new development session: `elcli draft new`. This will place the `edge-leap.yaml` configuration file in the current directory, which is expected to be filled with the required information for the development session.

  > _The configuration file schema details can be found [here](./docs/configuration-schema-v1.md)._

- `elcli draft deploy`: Deploy the module to the target device.

  Once the configuration file is set, you can start the development session and as soon as you want to deploy the module to the target device, you can run `elcli draft deploy`. If no flags are provided the configuration file information will be used, otherwise the flags will override the configuration file values.

  Deploying the module to the target device involves the following steps:
  - pushing the module manifest to the IoT Hub as layered deployment
  - defining a unique deployment ID and using it as target condition
  - updating the device's device twin with to match the target condition

### Release mode

The `release` mode can be used to orchestrate the module release under the CI/CD pipeline. It allows developers to provide a configuration of the release environment and automatically handle the required operations, in order to deploy the module manifest to the target IoT Hub.

A GitHub action is provided to automate the release process. The action can be found [here](https://github.com/unbrikd/actions/tree/master/elcli). The action requires the `AZURE_TOKEN` to be set as an environment variable for the workflow.


## Contributing

In order to contribute to the project, please read the [CONTRIBUTING.md](./CONTRIBUTING.md) file.

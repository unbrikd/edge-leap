# edge-leap

EdgeLeap is a development tool that streamlines IoT Edge module deployment. It integrates seamlessly with Azure IoT Edge and Azure IoT Hub to enable rapid local development and testing directly on edge devices, bypassing the traditional CI/CD pipeline. This allows developers to iterate quickly and validate their IoT Edge solutions in real device environments without the overhead of production deployment processes.

## Installing EdgeLeap

### Compiling from source

The `edge-leap` client can be compiled from source using the provided `Makefile`. By default, compiled binaries are placed in the `bin` directory, but this can be overridden by setting the `GO_BINDIR` environment variable. 

- To compile the client for the host platform and architecture, run:

    ```shell
    make build
    ```

- To compile for all architectures for a given platform, run:

    ```shell
    make build-[platform] # where [platform] is one of [macos | linux | windows]
    ```


### Using the pre-compiled binaries

Download the pre-compiled binaries from the releases page and run the executable. It is advisable to add the client to your `$PATH` so that it can be run from any directory.

- For MacOS and Linux

    ```shell
    export EL_VERSION=0.1.0
    export OS=[darwin | linux]
    export ARCH=[amd64 | arm64]

    curl -L https://github.com/unbrikd/edge-leap/releases/download/v${EL_VERSION}/elcli-${EL_VERSION}.${OS}-${ARCH}
    ```

### Using the Docker image

The `edge-leap` client is also available as a Docker image. 

- Build the image using the `Makefile` with the following command:

    ```shell
    make docker-image # builds the image for the host platform and architecture
    ```

- Use the pre-built image from the container registry:

    ```shell
    docker run ghcr.io/unbrikd/elcli:latest <command> # where <command> is the edge-leap client command to run
    ```
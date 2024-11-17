package docker

import (
	"context"
	"io"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/pkg/errors"
)

const defaultDockerfile = "Dockerfile"

// ClientFactory is a factory for a Docker client.
// It returns a pointer to an instance of a Docker client and a nil error, if
// it succeeds. Otherwise, it returns a nil client and the error that caused the
// failure.
// The client must be closed by the caller of this factory.
var ClientFactory = func() (client.APIClient, error) {
	return client.NewClientWithOpts(client.FromEnv)
}

type BuildOpts struct {
	// ContextFolder is the path to the folder that contains the context that will
	// be used to build the image.
	ContextFolder string

	// Dockerfile is the name (without the path) of the dockerfile that will be
	// used to build the image.
	Dockerfile string

	// ImageTag is the tag that will be assigned to the image.
	ImageTag string
}

// BuildImage builds a docker image based on the given configuration options.
// Returns nil if the image is successfully built. Otherwise, it returns the
// error that caused the failure.
func BuildImage(ctx context.Context, cli client.APIClient, opts BuildOpts) error {
	if cli == nil {
		return errors.New("docker client cannot be nil")
	}

	buildOpts := types.ImageBuildOptions{
		Tags:       []string{opts.ImageTag},
		Dockerfile: dockerfilePath(opts),
	}

	buildContext, err := archive.Tar(opts.ContextFolder, archive.Uncompressed)
	if err != nil {
		return errors.Wrap(err, "create context archive")
	}
	defer buildContext.Close()

	resp, err := cli.ImageBuild(ctx, buildContext, buildOpts)
	if err != nil {
		return errors.Wrap(err, "build docker image")
	}
	defer resp.Body.Close()

	if _, err = io.Copy(os.Stdout, resp.Body); err != nil {
		return errors.Wrap(err, "copy docker image")
	}

	return nil
}

func dockerfilePath(opts BuildOpts) string {
	dockerfile := strings.TrimSpace(opts.Dockerfile)
	if len(dockerfile) == 0 {
		dockerfile = defaultDockerfile
	}

	return dockerfile
}

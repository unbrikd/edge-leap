package docker_test

import (
	"bytes"
	"context"
	"io"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/golang/mock/gomock"
	"github.com/unbrikd/edge-leap/internal/docker"
)

func TestBuildImageWithInvalidClient(t *testing.T) {
	ctx := context.TODO()
	opts := docker.BuildOpts{
		ContextFolder: "./testdata",
		ImageTag:      "test:0.0.1",
	}
	if err := docker.BuildImage(ctx, nil, opts); err == nil {
		t.Error("expect errors; got nil")
	}
}

func TestBuildImage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cli := docker.NewMockAPIClient(ctrl)
	cli.EXPECT().
		ImageBuild(gomock.Any(), gomock.Any(), gomock.Any()).
		Times(1).
		Return(
			types.ImageBuildResponse{
				Body: io.NopCloser(bytes.NewBufferString("Build successful")),
			},
			nil,
		)

	ctx := context.TODO()
	opts := docker.BuildOpts{
		ContextFolder: "./testdata",
		ImageTag:      "test:0.0.1",
	}
	if err := docker.BuildImage(ctx, cli, opts); err != nil {
		t.Errorf("didn't expect errors; got %v", err)
	}
}

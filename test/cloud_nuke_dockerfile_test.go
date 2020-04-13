package main

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/docker"
	"github.com/stretchr/testify/assert"
)

func TestDockerfile(t *testing.T) {
	t.Parallel()

	tag := "jfharden/fargate-cloud-nuke:dockerfile"

	buildOptions := &docker.BuildOptions{
		Tags: []string{tag},
	}

	docker.Build(t, "../", buildOptions)

	runOptions := &docker.RunOptions{}
	output := docker.Run(t, tag, runOptions)
	assert.Equal(t, "cloud-nuke version v0.1.17", output)

	runOptions.Command = []string{"whoami"}
	output = docker.Run(t, tag, runOptions)
	assert.Equal(t, "cloud-nuke", output)

	runOptions.Command = []string{"pwd"}
	output = docker.Run(t, tag, runOptions)
	assert.Equal(t, "/cloud-nuke", output)
}

func TestDockerfileOverrideVersion(t *testing.T) {
	t.Parallel()

	tag := "jfharden/fargate-cloud-nuke:dockerfile-override-version"

	buildOptions := &docker.BuildOptions{
		Tags: []string{tag},

		BuildArgs: []string{
			"CLOUD_NUKE_VERSION=v0.1.16",
		},
	}

	docker.Build(t, "../", buildOptions)

	output := docker.Run(t, tag, &docker.RunOptions{})
	assert.Equal(t, "cloud-nuke version v0.1.16", output)
}

func TestDockerfileOverrideBinary(t *testing.T) {
	t.Parallel()

	tag := "jfharden/fargate-cloud-nuke:dockerfile-override-binary"

	buildOptions := &docker.BuildOptions{
		Tags: []string{tag},

		BuildArgs: []string{
			"CLOUD_NUKE_BINARY=cloud-nuke_linux_386",
		},
	}

	docker.Build(t, "../", buildOptions)

	runOptions := &docker.RunOptions{
		Command: []string{
			"ls", "/cloud-nuke/cloud-nuke_linux_386",
		},
	}

	output := docker.Run(t, tag, runOptions)
	assert.Equal(t, "/cloud-nuke/cloud-nuke_linux_386", output)

	runOptions.Command = []string{
		"readlink", "-f", "/cloud-nuke/cloud-nuke",
	}

	output = docker.Run(t, tag, runOptions)
	assert.Equal(t, "/cloud-nuke/cloud-nuke_linux_386", output)
}

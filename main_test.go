package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGitHubWorkFlows(t *testing.T) {
	assert.Equalf(t, true, GitHubWorkFlows(), "check CI")
}

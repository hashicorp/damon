// Copyright IBM Corp. 2021, 2023
// SPDX-License-Identifier: MPL-2.0

// Package buildconfig holds tests that guard the build configuration:
// the Go toolchain versions and the GitHub Actions setup have to stay in sync
// with what the module requires, otherwise CI and the container build break
// long after the change that caused it.
package buildconfig

import (
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const repoRoot = "../.."

// goVersion is a Go release such as 1.26 or 1.27.1.
type goVersion []int

func parseGoVersion(t *testing.T, raw string) goVersion {
	t.Helper()

	var version goVersion
	for part := range strings.SplitSeq(strings.TrimSpace(raw), ".") {
		number, err := strconv.Atoi(part)
		require.NoErrorf(t, err, "cannot parse Go version %q", raw)

		version = append(version, number)
	}

	require.NotEmptyf(t, version, "cannot parse Go version %q", raw)

	return version
}

func (v goVersion) atLeast(other goVersion) bool {
	for i, number := range other {
		if i >= len(v) {
			return false
		}

		if v[i] != number {
			return v[i] > number
		}
	}

	return true
}

func readFile(t *testing.T, path string) string {
	t.Helper()

	content, err := os.ReadFile(repoRoot + "/" + path)
	require.NoError(t, err)

	return string(content)
}

func findSubmatch(t *testing.T, content, expression string) string {
	t.Helper()

	match := regexp.MustCompile(expression).FindStringSubmatch(content)
	require.NotNilf(t, match, "no match for %q", expression)

	return match[1]
}

func TestGoVersionsAreConsistent(t *testing.T) {
	r := require.New(t)

	required := parseGoVersion(t, findSubmatch(t, readFile(t, "go.mod"), `(?m)^go (\S+)$`))
	toolchain := parseGoVersion(t, readFile(t, ".go-version"))
	image := parseGoVersion(t, findSubmatch(t, readFile(t, "Dockerfile"), `(?m)^FROM golang:(\S+) AS devbuild$`))

	// CI installs the toolchain from .go-version, it has to satisfy go.mod.
	r.Truef(toolchain.atLeast(required), ".go-version %v is older than the go directive %v", toolchain, required)
	r.Truef(image.atLeast(required), "the golang image %v is older than the go directive %v", image, required)
}

func TestWorkflowsUseSupportedInfrastructure(t *testing.T) {
	r := require.New(t)

	workflows, err := os.ReadDir(repoRoot + "/.github/workflows")
	r.NoError(err)

	// GitHub removed these runner images, a workflow pinning them never starts.
	retiredRunners := regexp.MustCompile(`runs-on: (ubuntu-(1[68]|20)\.04|macos-(10\.15|11|12))`)

	// GitHub rejects runs that use the v3 artifact actions.
	artifactAction := regexp.MustCompile(`actions/(?:upload|download)-artifact@\S+ #\s*v(\d+)`)

	for _, workflow := range workflows {
		if !strings.HasSuffix(workflow.Name(), ".yml") {
			continue
		}

		content := readFile(t, ".github/workflows/"+workflow.Name())

		r.Emptyf(retiredRunners.FindString(content), "%s pins a retired runner image", workflow.Name())

		for _, match := range artifactAction.FindAllStringSubmatch(content, -1) {
			major, err := strconv.Atoi(match[1])
			r.NoError(err)
			r.GreaterOrEqualf(major, 4, "%s uses a deprecated artifact action: %s", workflow.Name(), match[0])
		}
	}
}

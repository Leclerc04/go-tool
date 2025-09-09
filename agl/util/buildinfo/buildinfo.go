package buildinfo

import (
	"fmt"
	"os"
	"time"

	"github.com/leclerc04/go-tool/agl/util/errs"
)

// GitCommit is the commit hash when this binary is built.
var GitCommit string

// GitCommitTag is the tag assigned to the commit when this bianry is built.
var GitCommitTag string

// GitCommitTime is the time when the commit was made.
var GitCommitTime string

// BuildTime is the time when this binary is built.
var BuildTime string

// BuilderInfo contains the builder information.
var BuilderInfo string

// GitRepo is the URL to clonethe Git repository.
var GitRepo string

// Summary returns a string that contains the build information.
func Summary() string {
	formatVar := func(v string) string {
		if v != "" {
			return v
		}
		return "<not set>"
	}

	return fmt.Sprint(
		"Builder: ", formatVar(BuilderInfo), "\n",
		"Build time: ", formatVar(BuildTime), "\n",
		"Repo: ", formatVar(GitRepo), "\n",
		"Commit: ", formatVar(GitCommit), "\n",
		"CommitTag: ", formatVar(GitCommitTag), "\n",
		"CommitTime: ", formatVar(GitCommitTime), "\n")
}

// Release returns a string that identifies the release.
func Release() string {
	if BuildTime == "" {
		return GitCommit
	}
	if GitCommit == "" {
		return BuildTime
	}
	return BuildTime + ":" + GitCommit
}

func LDFlags() string {
	pkgName := "github.com/leclerc04/go-tool/agl/util/buildinfo"
	hostname, err := os.Hostname()
	errs.Ignore(err)
	commitSHA := os.Getenv("CI_COMMIT_SHA")
	commitTag := os.Getenv("CI_COMMIT_TAG")
	commitTime := os.Getenv("CI_COMMIT_TIMESTAMP")
	gitRepo := os.Getenv("CI_REPOSITORY_URL")
	buildTime := time.Now().UTC().Format(time.RFC3339)
	return "-X '" + pkgName + ".BuilderInfo=host:" + hostname +
		"' -X '" + pkgName + ".GitCommit=" + commitSHA +
		"' -X '" + pkgName + ".GitCommitTag=" + commitTag +
		"' -X '" + pkgName + ".GitCommitTime=" + commitTime +
		"' -X '" + pkgName + ".GitRepo=" + gitRepo +
		"' -X '" + pkgName + ".BuildTime=" + buildTime +
		"'"
}

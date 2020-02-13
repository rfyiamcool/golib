package buildinfo

import (
	"fmt"
	"runtime"
	"strings"
)

var (
	GitCommitLog   = "unknown"
	GitStatus      = "unknown"
	BuildTime      = "unknown"
	BuildGoVersion = "unknown"
)

func StringifySingleLine() string {
	return fmt.Sprintf("GitCommitLog=%s. GitStatus=%s. BuildTime=%s. GoVersion=%s. runtime=%s/%s.",
		GitCommitLog, GitStatus, BuildTime, BuildGoVersion, runtime.GOOS, runtime.GOARCH)
}

func StringifyMultiLine() string {
	return fmt.Sprintf("GitCommitLog=%s\nGitStatus=%s\nBuildTime=%s\nGoVersion=%s\nruntime=%s/%s\n",
		GitCommitLog, GitStatus, BuildTime, BuildGoVersion, runtime.GOOS, runtime.GOARCH)
}

func beauty() {
	if GitStatus == "" {
		GitStatus = "cleanly"
	} else {
		GitStatus = strings.Replace(strings.Replace(GitStatus, "\r\n", " |", -1), "\n", " |", -1)
	}
}

func init() {
	beauty()
}

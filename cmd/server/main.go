package main

import (
    "github.com/configcenter/cmd/server/cmd"
)

var (
    Version      string
    GoVersion    string
    GitBranch    string
    GitCommit    string
    GitLatestTag string
    BuildTime    string
)

func main() {
    cmd.Version = Version
    cmd.GoVersion = GoVersion
    cmd.GitBranch = GitBranch
    cmd.GitCommit = GitCommit
    cmd.GitLatestTag = GitLatestTag
    cmd.BuildTime = BuildTime

    cmd.Execute()
}

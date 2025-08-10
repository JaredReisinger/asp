//go:build tools
// +build tools

package tools

import (
	_ "github.com/automation-co/husky"
	// See the comment in Taskfile.yml's "prepare" task.  We can't use task to
	// acquire/setup task itself. :sadpanda:
	// _ "github.com/go-task/task/v3/cmd/task"
	_ "github.com/conventionalcommit/commitlint"
	// lintingzhen/commitizen-goseems to be failing...
	// _ "github.com/lintingzhen/commitizen-go"
	_ "github.com/JosephNaberhaus/go-mitizen"
	// We don't need goreleaser as a tool, since it's only used in CI.. we're
	// using a GitHub Action to get it.
	_ "github.com/securego/gosec/v2/cmd/gosec"
)

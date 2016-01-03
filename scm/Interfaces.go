package scm

import "github.com/driusan/bug/bugs"

type SCMHandler interface {
	Commit(dir bugs.Directory, commitMsg string) error
	Purge(bugs.Directory) error
	GetSCMType() string
}

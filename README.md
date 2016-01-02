# Bug

Bug is an implementation of a distributed issue tracker using
git to manage issues on the filesystem following [poor man's issue tracker](https://github.com/driusan/PoormanIssueTracker) conventions.

The goal is to use the filesystem in a human readable way, similar to
how an organized person without any bug tracking software might, 
by keeping track of bugs in an `issues/` directory, one (descriptive)
subdirectory per issue. bug provides a tool to maintain the nearest 
`issues/` directory to your current working directory and provides hooks 
to commit (or remove) the issues from source control.

This differs from other distributed bug tracking tools, (which usually 
store a database in a hidden directory) in that you can still easily 
view, edit, or understand bugs even without access to the bug tool. bug
only acts as a way to streamline the process of maintaining them. Another 
benefit is that you can also have multiple `issues/` directories at 
different places in your directory tree to, for instance, keep separate 
bug repositories for different submodules or packages contained in a 
single git repository.

Because issues are stored as human readable plaintext files, they branch
and merge along with the rest of your code, and you can resolve conflicts 
using your standard tools.

# Sample Usage

If an environment variable named PMIT is set, that directory will be
used to create and maintain issues, otherwise the bug command will
walk up the tree until it finds somewhere with a subdirectory named
"issues" to track issues in.

Some sample usage (assuming you're already in a directory tracked by
git):

```
$ mkdir issues
$ bug help
Usage: bug command [options]

Use "bug help [command]" for more information about any command below

Valid commands

Issue editing commands:
	create	 File a new bug
	list	 List existing bugs
	edit	 Edit an existing bug
	tag	 Tag a bug with a category
	close	 Delete an existing bug
	rm	 Alias of close
	status	 View or edit a bug's status
	priority View or edit a bug's priority

Source control commands:
	commit	 Commit any new, changed or deleted bug to git
	purge	 Remove all issues not tracked by git

Other commands:
	env	 Show settings that bug will use if invoked from this directory
	dir	 Prints the issues directory to stdout (useful subcommand in the shell)
	pwd	 Alias of dir
	help	 Show this screen

$ bug create I don't know what I'm doing
# (Your standard editor will open here for you to enter a description, save it when you're done)

$ bug list
Issue 1: I don't know what I'm doing

$ bug list 1
Title: I don't know what I'm doing

Description:
The description that I entered

$ bug purge
Removing issues/I-don't-know-what-I'm-doing

$ bug create Need better formating for README
# (Your editor opens again)

$ bug list
Issue 1: Need better formating for README

$ bug commit
$ git push
```

# Installation
If you have go installed, install the latest development version with:

`go get github.com/driusan/bug`

Make sure `$GOPATH/bin` or `$GOBIN` are in your path (or copy
the "bug" binary somewhere that is.)

Otherwise, you can download a 64 bit release for OS X or Linux on the 
[releases](https://github.com/driusan/bug/releases/) page.


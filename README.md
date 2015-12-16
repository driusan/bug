# Issue Tracker 

This repo contains an implementation of a tool to create
bug's in a (poor man's issue tracker)[https://github.com/driusan/PoormanIssueTracker].

After compiling the main.go file, copy it somewhere in your PATH
and use it to create issues. If an environment variable named PMIT
is set, it will create issues in that directory, otherwise it will
walk up the path from the current working directory until it finds
somewhere with an "issues" subdirectory and use that as a location
for any issues that are created.

Some sample usage (assuming you compile it into a binary called "bug":

```bash
$ bug
Usage: bug command [options]

Valid commands
	create	File a new bug
	list	List existing bugs
	env	Show settings that bug will use if invoked from this directory
	help	Show this screen
$ bug create I don't know what I'm doing
# (An editor will open here for you to enter a description, save it when you're done)

$ bug list
Issue 1: I don't know what I'm doing

$ bug list 1

Title: I don't know what I'm doing

Description:
The description that I entered
```

Not that the issue numbers are not constant and only shown to simplify
command line usage. They will change as you create/add or delete new issues.

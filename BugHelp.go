package main

import (
	"fmt"
	"os"
)

func (a BugApplication) Help(args ...string) {
	var cmd string
	if args == nil {
		cmd = "help"

	}
	if len(args) <= 1 {
		cmd = "help"
	} else {
		cmd = args[1]
	}
	switch cmd {
	case "create":
		fmt.Printf("Usage: " + os.Args[0] + " create [-n] Issue Title\n\n")
		fmt.Printf(
			`This will create an issue with the title Issue Title.  An editor 
will be opened automatically for you to enter a more detailed 
description. If your EDITOR environment variable is set, it 
will be used, otherwise the default editor is vim.

If the first argument to create is "-n", then %s will not open 
any editor and create an empty Description
`, os.Args[0])
	case "list":
		fmt.Printf("Usage: " + os.Args[0] + " list [issue numbers]\n")
		fmt.Printf("       " + os.Args[0] + " list [tags]\n\n")
		fmt.Printf(
			`This will list the issues found in the current environment

With no arguments, titles will be printed to the screen along
with the issue number that can be used to reference this issue
on the command line.

If 1 or more issue numbers are provided, the whole issue including
description will be printed to stdout.

If, instead of issue numbers, you provide list with 1 or more tags, 
it will print any issues which have that tag (in short form).

Note that issue numbers are not intended to be stable, but only
to provide a quick way to reference issues on the command line.
They will change as you create, edit, and close other issues.

The subcommand "view" is an alias for "list".
`)

	case "edit":
		fmt.Printf("Usage: " + os.Args[0] + " edit IssueNumber\n\n")
		fmt.Printf(
			`This will launch your standard editor to edit the description 
of the bug numbered IssueNumber, where IssueNumber is a reference
to same index provided with a "bug list" command.
`)
	case "status":
		fmt.Printf("Usage: " + os.Args[0] + " status IssueNumber [NewStatus]\n\n")
		fmt.Printf(
			`This will edit or display the status of the bug numbered IssueNumber.
            
If NewStatus is provided, it will update the first line of the Status file
for the issue (creating the file as necessary). If not provided, it will 
display the first line of the Status file to STDOUT.

Note that you can manually edit the Status file in the issues/ directory
to provide further explanation (for instance, why that status is set.)
This command will preserve the explanation when updating a status.
`)
	case "priority":
		fmt.Printf("Usage: " + os.Args[0] + " priority IssueNumber [NewPriority]\n\n")
		fmt.Printf(
			`This will edit or display the priority of the bug numbered IssueNumber.
By convention, priorities should be an integer number (higher is more 
urgent), but that is not enforced by this command and NewPriority can
be any free-form text if you prefer.
            
If NewPriority is provided, it will update the first line of the Priority
file for the issue (creating the file as necessary). If not provided, it 
will display the first line of the Priority file to STDOUT.

Note that you can manually edit the Priority file in the issues/ directory
to provide further explanation (for instance, why that priority is set.)
This command will preserve the explanation when updating a priority.
`)
	case "milestone":
		fmt.Printf("Usage: " + os.Args[0] + " milestone IssueNumber [NewMilestone]\n\n")
		fmt.Printf(
			`This will edit or display the milestone of the bug numbered IssueNumber.

There are no restrictions on how milestones must be named, but
semantic versioning is a good convention to adopt. Failing that,
it's a good idea to use milestones that collate properly when
sorted as strings so that they appear properly in "%s roadmap".

If NewMilestone is provided, it will update the first line of the
Milestone file for the issue (creating the file as necessary). 
If not provided, it will display the first line of the Milestone 
file to STDOUT.

Note that you can manually edit the Milestone file in the issues/
directory to provide further explanation (for instance, why that 
milestone is set.)

This command will preserve the explanation when updating a priority.
`, os.Args[0])
	case "retitle":
		fallthrough
	case "mv":
		fallthrough
	case "relabel":
		fmt.Printf("Usage: " + os.Args[0] + " relabel IssueNumber New Title\n\n")
		fmt.Printf(
			`This will change the title of IssueNumber to "New Title". Use this
to rename an issue.

"%s mv", "%s retitle", and "%s rename" are all aliases for "%s relabel".
`, os.Args[0], os.Args[0], os.Args[0], os.Args[0])
	case "rm":
		fallthrough
	case "close":
		fmt.Printf("Usage: " + os.Args[0] + " close IssueNumber\n")
		fmt.Printf("       " + os.Args[0] + " rm IssueNumber\n\n")
		fmt.Printf(
			`This will delete the issue numbered IssueNumber. IssueNumbers
correspond to the number in the "bug list" command.

Note that closing a bug will cause all existing bugs to be
renumbered and IssueNumbers are not intended to be stable.

Also note that this does not remove the issue from git, but only 
from the file system. You'll need to execute "bug commit" to
remove the bug from source control.

"%s rm" is an alias for this "%s close"
`, os.Args[0], os.Args[0])
	case "purge":
		fmt.Printf("Usage: " + os.Args[0] + " purge\n\n")
		fmt.Printf(
			`This will delete any bugs that are not currently tracked by
git.
`)
	case "commit":
		fmt.Printf("Usage: " + os.Args[0] + " commit\n\n")
		fmt.Printf(`This will commit any new, modified, or removed issues to
git.

Your working tree and staging area should be otherwise
unaffected by using this command.
`)
	case "env":
		fmt.Printf("Usage: " + os.Args[0] + " env\n\n")
		fmt.Printf(`This will print the environment used by the bug command to stdout.

Use this command if you want to see what directory bug create is
using to store bugs, or what editor will be invoked by bug create/edit.
`)

	case "dir":
		fallthrough
	case "pwd":
		fmt.Printf("Usage: " + os.Args[0] + " dir\n\n")
		fmt.Printf(
			`This will print the undecorated bug directory to stdout, 
so you can use it as a subcommand for arguments to any 
arbitrary shell commands. For example "cd $(bug dir)"

"%s dir" is an alias for "%s pwd"
`, os.Args[0], os.Args[0])
	case "tag":
		fmt.Printf("Usage: " + os.Args[0] + " tag IssueNumber [tags]\n\n")
		fmt.Printf(`This will tag the given IssueNumber with the tags
given as parameters. At least one tag is required.

Tags can be any string which would make a valid file name.
`)
	case "roadmap":
		fmt.Printf("Usage: " + os.Args[0] + " roadmap [--simple]\n\n")
		fmt.Printf(
			`This will print a markdown formatted list of all open
issues, grouped by milestone.

If the argument --simple is provided, only the titles will be displayed.
Otherwise, the Status and Priority will be included as well as the title
`)
	case "about":
		fallthrough
	case "version":
		fmt.Printf("Usage: " + os.Args[0] + " version\n\n")
		fmt.Printf(
			`This will print information about the version of %s being
invoked.

"%s about" is an alias for "version".
`, os.Args[0], os.Args[0])

	case "help":
		fallthrough
	default:
		fmt.Printf("Usage: " + os.Args[0] + " command [options]\n\n")
		fmt.Printf("Use \"bug help [command]\" for more information about any command below\n\n")
		fmt.Printf("Valid commands\n")
		fmt.Printf("\nIssue editing commands:\n")
		fmt.Printf("\tcreate\t  File a new bug\n")
		fmt.Printf("\tlist\t  List existing bugs\n")
		fmt.Printf("\tedit\t  Edit an existing bug\n")
		fmt.Printf("\ttag\t  Tag a bug with a category\n")
		fmt.Printf("\trelabel\t  Rename the title of a bug\n")
		fmt.Printf("\tclose\t  Delete an existing bug\n")
		fmt.Printf("\tstatus\t  View or edit a bug's status\n")
		fmt.Printf("\tpriority  View or edit a bug's priority\n")
		fmt.Printf("\tmilestone View or edit a bug's milestone\n")

		fmt.Printf("\nSource control commands:\n")
		fmt.Printf("\tcommit\t Commit any new, changed or deleted bug to git\n")
		fmt.Printf("\tpurge\t Remove all issues not tracked by git\n")

		fmt.Printf("\nOther commands:\n")
		fmt.Printf("\tenv\t Show settings that bug will use if invoked from this directory\n")
		fmt.Printf("\tpwd\t Prints the issues directory to stdout (useful subcommand in the shell)\n")
		fmt.Printf("\troadmap\t Print list of open issues sorted by milestone\n")
		fmt.Printf("\tversion\t Print the version of this software\n")
		fmt.Printf("\thelp\t Show this screen\n")
	}
}

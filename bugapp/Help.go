package bugapp

import (
	"fmt"
	"os"
)

func Help(args ...string) {
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
		fmt.Printf("Usage: " + os.Args[0] + " create [-n] [options] Issue Title\n\n")
		fmt.Printf(
			`This will create an issue with the title Issue Title.  An editor 
will be opened automatically for you to enter a more detailed 
description. If your EDITOR environment variable is set, it 
will be used, otherwise the default editor is vim.

If the first argument to create is "-n", then %s will not open 
any editor and create an empty Description.

Options take a value and set a field on the bug at the same
time as creating it. Valid options are:
    --status     Sets the bug status to the next parameter
    --tag        Tags the bug with a tag on creation
    --priority   Sets the priority to the next parameter
    --milestone  Sets the milestone to the next parameter
    --identifier Sets the identifier to the next parameter
    --generate-id Automatically generate a stable bug identifier
`, os.Args[0])
	case "list":
		fmt.Printf("Usage: " + os.Args[0] + " list [BugIDs]\n")
		fmt.Printf("       " + os.Args[0] + " list [tags]\n\n")
		fmt.Printf(
			`This will list the issues found in the current environment

With no arguments, titles will be printed to the screen along
with the issue number that can be used to reference this issue
on the command line.

If 1 or more BugIDs are provided, the whole issue including
description will be printed to STDOUT.  See "bug help identifiers"
for a description of what makes a BugID.

If, instead of BugIDs, you provide list with 1 or more tags, 
it will print any issues which have that tag (in short form).

Note that BugIDs may change as you create, edit, and close other
unless yo have defined a stable identifier for the issue. Again,
see "bug help identifiers."

The subcommand "view" is an alias for "list".
`)

	case "edit":
		fmt.Printf("Usage: " + os.Args[0] + " edit [Filename] BugID\n\n")
		fmt.Printf(
			`This will launch your standard editor to edit the description 
of the bug identified by BugID.  See "bug help identifiers" for a 
description of what makes a BugID.

If the Filename option is provided, bug will instead launch an editor
to edit that file name within the bug directory. Files that have
special meaning to bug (Status, Milestone, Priority, Identifier) are
treated in a case insensitive manner, otherwise the filename is passed
directly to your editor.
`)
	case "status":
		fmt.Printf("Usage: " + os.Args[0] + " status BugID [NewStatus]\n\n")
		fmt.Printf(
			`This will edit or display the status of the bug identified by BugID.
See "bug help identifiers" for a description of what constitutes a BugID.
            
If NewStatus is provided, it will update the first line of the Status file
for the issue (creating the file as necessary). If not provided, it will 
display the first line of the Status file to STDOUT.

Note that you can edit the status in your standard editor with the
command "%s edit status BugID". If you provide a longer than 1 line
status with "bug edit status", "bug status" will preserve everything
after the first line when editing a status. You can use this to provide
further context on a status (for instance, why that status is setup.)
`, os.Args[0])
	case "priority":
		fmt.Printf("Usage: " + os.Args[0] + " priority BugID [NewPriority]\n\n")
		fmt.Printf(
			`This will edit or display the priority of BugID. See "bug help identifiers"
for a description of what constitutes a BugID.

By convention, priorities should be an integer number (higher is more 
urgent), but that is not enforced by this command and NewPriority can
be any free-form text if you prefer.
            
If NewPriority is provided, it will update the first line of the Priority
file for the issue (creating the file as necessary). If not provided, it 
will display the first line of the Priority file to STDOUT.

Note that you can manually edit the Priority file in the issues/ directory
by running "%s edit priority BugID", to provide further explanation (for 
instance, why that priority is set.) This command will preserve the 
explanation when updating a priority.
`, os.Args[0])
	case "milestone":
		fmt.Printf("Usage: " + os.Args[0] + " milestone BugID [NewMilestone]\n\n")
		fmt.Printf(
			`This will edit or display the milestone of the identified by BugID.
See "%s help identifiers" for a description of what constitutes a BugID.

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
milestone is set) with the command "bug edit milestone BugID"

This command will preserve the explanation when updating a priority.
`, os.Args[0], os.Args[0])
	case "retitle", "mv", "rename", "relabel":
		fmt.Printf("Usage: " + os.Args[0] + " relabel BugID New Title\n\n")
		fmt.Printf(
			`This will change the title of BugID to "New Title". Use this
to rename an issue.

"%s mv", "%s retitle", and "%s rename" are all aliases for "%s relabel".
`, os.Args[0], os.Args[0], os.Args[0], os.Args[0])
	case "rm", "close":
		fmt.Printf("Usage: " + os.Args[0] + " close BugID\n")
		fmt.Printf("       " + os.Args[0] + " rm BugID\n\n")
		fmt.Printf(
			`This will delete the issue identifier by BugID. See
"%s help identifiers" for details on what constitutes a BugID.

Note that closing a bug may cause existing BugIDs to change if
they do not have a stable identifier set (see "%s help identifiers",
again.)

Also note that this does not remove the issue from git, but only 
from the file system. You'll need to execute "bug commit" to
remove the bug from source control.

"%s rm" is an alias for this "%s close"
`, os.Args[0], os.Args[0], os.Args[0], os.Args[0])
	case "find":
	    fmt.Printf("Usage: %s find tag value1 [value2 ...]\n", os.Args[0])
	    fmt.Printf("Usage: %s find status value1 [value2 ...]\n", os.Args[0])
	    fmt.Printf("Usage: %s find priority value1 [value2 ...]\n", os.Args[0])
	    fmt.Printf("Usage: %s find milestone value1 [value2 ...]\n\n", os.Args[0])
		fmt.Printf(
            `This will search all bugs for multiple tags, statuses, priorities, or milestone.
The matching bugs will be printed.
`)
    case "purge":
		fmt.Printf("Usage: " + os.Args[0] + " purge\n\n")
		fmt.Printf(
			`This will delete any bugs that are not currently tracked by
git.
`)
	case "commit":
		fmt.Printf("Usage: " + os.Args[0] + " commit [--no-autoclose]\n\n")
		fmt.Printf(`This will commit any new, modified, or removed issues to
git or hg.

Your working tree and staging area should be otherwise
unaffected by using this command.

If the --no-autoclose option is passed to commit, bug will
not include a "Closes #x" line for each issue imported from
"bug-import --github." Otherwise, the commit message will
include the list of issues that were closed so that GitHub
will autoclose them when the changes are pushed upstream.
`)
	case "env":
		fmt.Printf("Usage: " + os.Args[0] + " env\n\n")
		fmt.Printf(`This will print the environment used by the bug command to stdout.

Use this command if you want to see what directory bug create is
using to store bugs, or what editor will be invoked by bug create/edit.
`)

	case "dir", "pwd":
		fmt.Printf("Usage: " + os.Args[0] + " dir\n\n")
		fmt.Printf(
			`This will print the undecorated bug directory to stdout, 
so you can use it as a subcommand for arguments to any 
arbitrary shell commands. For example "cd $(bug dir)"

"%s dir" is an alias for "%s pwd"
`, os.Args[0], os.Args[0])
	case "tag":
		fmt.Printf("Usage: " + os.Args[0] + " tag [--rm] BugID [tags]\n\n")
		fmt.Printf(`This will tag the given BugID with the tags
given as parameters. At least one tag is required.

Tags can be any string which would make a valid file name.

If the --rm option is provided before the BugID, all tags provided will
be removed instead of added.
`)
	case "roadmap":
		fmt.Printf("Usage: " + os.Args[0] + " roadmap [options]\n\n")
		fmt.Printf(
			`This will print a markdown formatted list of all open
issues, grouped by milestone.

Valid options are:
    --simple      Don't show anything other than the title in the output
    --no-status   Don't show the status of an issue
    --no-priority Don't show the priority of an issue
    --no-identifier Don't include the bug identifier of an issue
    --tags        Include the tags attached to a bug in it's output

    --filter tag           Only show bugs matching tag
    --filter tag1,tag2,etc Only show issues matching at least one of
                           the supplied tags

`)
	case "id", "identifier":
		fmt.Printf("Usage: " + os.Args[0] + " identifier BugID [--generate] [value]\n\n")
		fmt.Printf(
			`This will either set of retrieve the identifier for the bug
currently identified by BugID.

If value is provided as an argument, the bug identifier will be set
to the value passed in. You should take care to ensure that any
identifier used has at least 1 non-numeric character, to ensure there
are no conflicts with automatically generated issue numbers used for
a bug that has no explicit identifier set.

If the --generate option is passed instead of a static value, a
short identifier will be generated derived from the issue's current
title (however, the identifier will remain unchanged if the bug's title
is changed.)

If only a BugID is provided, the current identifier will be printed.

"%s id" is an alias for "%s identifier"
`, os.Args[0], os.Args[0])
	case "about", "version":
		fmt.Printf("Usage: " + os.Args[0] + " version\n\n")
		fmt.Printf(
			`This will print information about the version of %s being
invoked.

"%s about" is an alias for "version".
`, os.Args[0], os.Args[0])
	case "identifiers":
		fmt.Printf(
			`Bugs can be referenced in 2 ways on the commandline, either by
an index of where the bug directory is located inside the issues
directory, or by an identifier. "BugID" can be either of these,
and %s will try and intelligently guess which your command is
referencing.

By default, no identifiers are set for an issue. This means that
the issue number provided in "%s list" is an index into the directory,
and is unstable as bugs are created, modified, and closed. However,
the benefit is that they are easy to reference and remember, at least
in the short term.

If you have longer lasting issues that need a stable identifier,
they can be created by "%s identifier BugID NewIdentifier" to
set the identifier of BugID to NewIdentifier. From that point
forward, you can use NewIdentifier to reference the bug instead
of the directory index.

There are no rules for what constitutes a valid identifier, but
you should try and ensure that they have at least 1 non-numeric
character so that they don't conflict with directory indexes.

If you just want an identifier but don't care what it is, you
can use "%s identifier BugID --generate" to generate a new
identifier for BugID.

If there are no exact matches for the BugID provided, %s commands will
also try and look up the bug by a substring match on all the valid 
identifiers in the system before giving up.
`, os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0])

	case "help":
		fallthrough
	default:
		fmt.Printf("Usage: " + os.Args[0] + " command [options]\n\n")
		fmt.Printf("Use \"bug help [command]\" for more information about any command below\n\n")
		fmt.Printf("Valid commands\n")
		fmt.Printf("\nIssue editing commands:\n")
		fmt.Printf("\tcreate\t   File a new bug\n")
		fmt.Printf("\tlist\t   List existing bugs\n")
		fmt.Printf("\tedit\t   Edit an existing bug\n")
		fmt.Printf("\ttag\t   Tag a bug with a category\n")
		fmt.Printf("\tidentifier Set a stable identifier for the bug\n")
		fmt.Printf("\trelabel\t   Rename the title of a bug\n")
		fmt.Printf("\tclose\t   Delete an existing bug\n")
		fmt.Printf("\tstatus\t   View or edit a bug's status\n")
		fmt.Printf("\tpriority   View or edit a bug's priority\n")
		fmt.Printf("\tmilestone  View or edit a bug's milestone\n")
		fmt.Printf("\tfind\t   Search bugs for a tag, status, priority, or milestone\n")

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

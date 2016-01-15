package bugapp

import "github.com/driusan/bug/bugs"

func Status(args ArgumentList) {
	fieldHandler("status", args, bugs.Bug.SetStatus, bugs.Bug.Status)
}

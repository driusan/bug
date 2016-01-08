package bugapp

import (
	"github.com/driusan/bug/bugs"
)

func Milestone(args ArgumentList) {
	fieldHandler("milestone", args, bugs.Bug.SetMilestone, bugs.Bug.Milestone)
}

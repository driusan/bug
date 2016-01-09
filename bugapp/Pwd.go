package bugapp

import (
	"fmt"
	"github.com/driusan/bug/bugs"
)

func Pwd() {
	fmt.Printf("%s", bugs.GetIssuesDir())
}

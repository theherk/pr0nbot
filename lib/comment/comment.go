// Package comment post comments to reddit.
package comment

import (
	"fmt"

	"github.com/turnage/graw/reddit"
)

const tmpl = "__This image contains prawn__.\n\nWould somebody please think of " +
	"the children? Posting images containing prawn should met with scorn and " +
	"tagged appropriately. shame, shame, shame... _bong_. %s, you must confess."

// Do will reply to a reddit comment.
func Do(b reddit.Bot, p reddit.Post) error {
	return b.Reply(p.Name, fmt.Sprintf(tmpl, p.Author))
}

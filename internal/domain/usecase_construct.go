package domain

import (
	"fmt"
)

func (c *Construct) Title() string {
	return fmt.Sprintf("%s-%s", c.LeftPole, c.RightPole)
}

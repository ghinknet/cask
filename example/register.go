package example

import "go.gh.ink/cask/driver"

func Trigger() {
	driver.Register("example", Driver{})
}

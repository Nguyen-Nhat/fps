package constant

import "regexp"

var (
	FileNameRegex    = regexp.MustCompile(`(?:.+/)?(?P<fileName>(?P<Name>[\w\p{P}\p{S}]+)\.(?P<Extension>\w+))$`)
	PhoneNumberRegex = regexp.MustCompile(`(?m)^(?:(?:\+?84)|0)\d{9}$`)
)

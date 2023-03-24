package constant

import "regexp"

var (
	FileNameRegex = regexp.MustCompile(`(?:.+/)?(?P<fileName>(?P<Name>[\w\p{P}\p{S}\p{L}\p{M} ]+)\.(?P<Extension>\w+))$`)
)

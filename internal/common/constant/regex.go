package constant

import "regexp"

var (
	FileNameRegex = regexp.MustCompile(`.*/(?P<fileName>[\w\p{P}\p{S}]+\.\w+)$`)
)

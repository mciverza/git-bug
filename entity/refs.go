package entity

import "strings"

func RefsToIds(refs []string) []Id {
	ids := make([]Id, len(refs))

	for i, ref := range refs {
		ids[i] = RefToId(ref)
	}

	return ids
}

func RefToId(ref string) Id {
	split := strings.Split(ref, "/")
	return Id(split[len(split)-1])
}

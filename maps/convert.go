package maps

func FromStringSlice(s []string) map[string]struct{} {
	ret := make(map[string]struct{})
	empty := struct{}{}
	for _, val := range s {
		ret[val] = empty
	}
	return ret
}

package support

func EqualsStr(a []string, b []string) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil && b != nil || a != nil && b == nil {
		return false
	}
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

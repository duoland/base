package maps

// CopyStringMap makes a shallow copy of a string map.
func CopyStringMap(m map[string]string) map[string]string {
	if m == nil {
		return nil
	}

	newCopy := make(map[string]string, len(m))
	for k, v := range m {
		newCopy[k] = v
	}
	return newCopy
}

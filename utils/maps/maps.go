package maps

// MergeMaps merge two maps values which is used in helm values merge.
func MergeMaps(a, b map[string]any) map[string]any {
	out := make(map[string]any, len(a))
	for k, v := range a {
		out[k] = v
	}
	for k, v := range b {
		if v, ok := v.(map[string]any); ok {
			if bv, ok := out[k]; ok {
				if bv, ok := bv.(map[string]any); ok {
					out[k] = MergeMaps(bv, v)
					continue
				}
			}
		}
		out[k] = v
	}
	return out
}

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

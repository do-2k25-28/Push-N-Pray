package utils

func MergeMap[K comparable, V any, T ~map[K]V](map1 T, map2 T) T {
	merged := make(T)
	for k, v := range map1 {
		merged[k] = v
	}
	for k, v := range map2 {
		merged[k] = v
	}
	return merged
}

func MergeMaps[K1 comparable, K2 comparable, V any, T ~map[K1]map[K2]V](maps []T) T {
	merged := make(T)
	for _, m := range maps {
		for k1, submap := range m {
			if merged[k1] == nil {
				merged[k1] = make(map[K2]V)
			}
			for k2, v := range submap {
				merged[k1][k2] = v
			}
		}
	}
	return merged
}

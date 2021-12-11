package utils

func MergeMapsRight(baseMap map[string]string, otherMaps ...map[string]string) map[string]string {
	main := map[string]string{}

	for _, singleMap := range otherMaps {
		for k, v := range singleMap {
			main[k] = v
		}
	}

	for k, v := range baseMap {
		main[k] = v
	}

	return main
}

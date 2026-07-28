package query

func init() {
	filtered := statePaths[:0]
	for _, schema := range statePaths {
		if isDirectPrivateZonePath(schema.Pattern) {
			continue
		}
		filtered = append(filtered, schema)
	}
	statePaths = filtered
}

func isDirectPrivateZonePath(path []string) bool {
	return len(path) == 3 && path[0] == "players" && path[1] == "*" && path[2] == "private_zones"
}

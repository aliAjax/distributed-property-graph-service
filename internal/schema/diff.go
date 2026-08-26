package schema

type Change struct {
	Kind     string `json:"kind"`
	Path     string `json:"path"`
	Breaking bool   `json:"breaking"`
}

func DiffSchemas(before, after Schema) []Change {
	changes := []Change{}
	old := map[string]Property{}
	for _, v := range before.Vertices {
		for _, p := range v.Properties {
			old[v.Name+"."+p.Name] = p
		}
	}
	for _, v := range after.Vertices {
		for _, p := range v.Properties {
			key := v.Name + "." + p.Name
			previous, ok := old[key]
			if !ok {
				changes = append(changes, Change{"property_added", key, false})
				continue
			}
			if previous.Type != p.Type {
				changes = append(changes, Change{"property_type_changed", key, true})
			}
			delete(old, key)
		}
	}
	for key := range old {
		changes = append(changes, Change{"property_removed", key, true})
	}
	return changes
}
func HasBreakingChanges(changes []Change) bool {
	for _, change := range changes {
		if change.Breaking {
			return true
		}
	}
	return false
}

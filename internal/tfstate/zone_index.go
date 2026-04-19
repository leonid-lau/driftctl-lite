package tfstate

// ZoneIndex maps zone values to slices of resources.
type ZoneIndex struct {
	index map[string][]Resource
}

// BuildZoneIndex constructs a ZoneIndex from the given resources.
// Zone is read from Attributes["zone"] or Attributes["availability_zone"].
func BuildZoneIndex(resources []Resource) *ZoneIndex {
	idx := &ZoneIndex{index: make(map[string][]Resource)}
	for _, r := range resources {
		zone := zoneKey(r)
		if zone == "" {
			continue
		}
		key := toLower(zone)
		idx.index[key] = append(idx.index[key], r)
	}
	return idx
}

// Lookup returns all resources matching the given zone (case-insensitive).
func (z *ZoneIndex) Lookup(zone string) []Resource {
	if zone == "" {
		return nil
	}
	return z.index[toLower(zone)]
}

// Zones returns all indexed zone keys.
func (z *ZoneIndex) Zones() []string {
	keys := make([]string, 0, len(z.index))
	for k := range z.index {
		keys = append(keys, k)
	}
	return keys
}

func zoneKey(r Resource) string {
	if v, ok := r.Attributes["zone"]; ok {
		return v
	}
	if v, ok := r.Attributes["availability_zone"]; ok {
		return v
	}
	return ""
}

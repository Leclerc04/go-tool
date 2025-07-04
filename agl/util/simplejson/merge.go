package simplejson

// MergeFrom merge the given JSON recursively.
func (j *JSON) MergeFrom(from *JSON) {
	j.data = mergeFromInternal(j.data, from.data)
}

func mergeFromInternal(to interface{}, from interface{}) interface{} {
	if to == nil {
		return from
	}
	switch fromV := from.(type) {
	case map[string]interface{}:
		switch toV := to.(type) {
		case map[string]interface{}:
			for k, v := range fromV {
				oldV, ok := toV[k]
				if ok {
					toV[k] = mergeFromInternal(oldV, v)
					continue
				}
				toV[k] = v
			}
			return to
		default:
			return from
		}
	default:
		return from
	}
}

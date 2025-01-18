package main

// // // // // // // // // // // // // // // // // //

func merge(def, data any) any {
	switch defTyped := def.(type) {
	case map[string]any:
		res := make(map[string]any)
		dataMap, ok := data.(map[string]any)
		for k, v := range defTyped {
			if ok {
				if dv, exists := dataMap[k]; exists {
					res[k] = merge(v, dv)
				} else {
					res[k] = merge(v, nil)
				}
			} else {
				res[k] = merge(v, nil)
			}
		}
		return res

	default:
		if data != nil {
			return data
		}
		return ""
	}
}

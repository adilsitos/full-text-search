package main

// inverted index: map every word to an index.

// for example:

// map {
// 	"donut": [1]
// 	"glass": [2, 3]
// }

// every word will be mapped to several indexes
// eg. if we want to find all strings that contain the word cat
// we will map the word cat on the map, and it will return all the idx strings that contains it
type index map[string][]int

func (idx index) add(docs []document) {
	for _, doc := range docs {
		for _, token := range analyze(doc.Text) {
			ids := idx[token]
			if ids != nil && ids[len(ids)-1] == doc.ID {
				// don't add the same id twice
				continue
			}

			idx[token] = append(ids, doc.ID)
		}
	}
}

func (idx index) search(text string) []int {
	var r []int

	for _, token := range analyze(text) {
		if ids, ok := idx[token]; ok {
			if r == nil {
				r = ids
			} else {
				r = intersection(r, ids)
			}
		} else {
			// token does not exist
			return nil
		}
	}

	return r
}

func intersection(a, b []int) []int {
	maxLen := len(a)
	if len(b) > maxLen {
		maxLen = len(b)
	}
	r := make([]int, 0, maxLen)
	var i, j int
	for i < len(a) && j < len(b) {
		if a[i] < b[j] {
			i++
		} else if a[i] > b[j] {
			j++
		} else {
			r = append(r, a[i])
			i++
			j++
		}
	}

	return r
}

package domain

func ReverseSliceTerms(s []Term) []Term {
	res := make([]Term, 0)
	for i := len(s) - 1; i >= 0; i-- {
		res = append(res, s[i])
	}

	return res
}

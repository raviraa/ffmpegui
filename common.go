package main

func filterList(ss []string, f func(string) bool) (res []string) {
	for _, s := range ss {
		if f(s) {
			res = append(res, s)
		}
	}
	return res
}

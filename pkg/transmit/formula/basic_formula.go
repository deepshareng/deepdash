package formula

import "strconv"

type BasicFormula struct{}

func (f *BasicFormula) Requisite(param string, gran string) []string {
	return []string{"match/install_with_params", "match/open_with_params"}
}

func (f *BasicFormula) Calculate(param string, data map[string][]int, gran string, limit int) map[string][]string {
	res := make(map[string][]string)
	res["match/install_with_params"] = make([]string, limit)
	res["match/open_with_params"] = make([]string, limit)
	for i := 0; i < limit; i++ {
		res["match/install_with_params"][i] = strconv.Itoa(data["match/install_with_params"][i])
		res["match/open_with_params"][i] = strconv.Itoa(data["match/open_with_params"][i])
	}
	return res
}

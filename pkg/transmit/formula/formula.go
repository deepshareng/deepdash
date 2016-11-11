package formula

var TransmitFormulaCollection = map[string]TransmitFormula{
	"basic-event":                 &BasicFormula{},
	"3-day-retention":             &RetentionFormula{},
	"7-day-retention":             &RetentionFormula{},
	"sharelink-tem-overall-value": &SharelinkTemFormula{},
	"url-tem-overall-value":       &GenURLTemFormula{},
}

type TransmitFormula interface {
	// return a neccessary params to calculate pointed value
	Requisite(param string, gran string) []string

	// On the basis of channel params, caculate pointed value
	Calculate(param string, data map[string][]int, gran string, limit int) map[string][]string
}

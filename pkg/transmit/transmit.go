package transmit

import (
	"strconv"

	"github.com/MISingularity/deepdash/pkg/transmit/formula"
	"github.com/MISingularity/deepshare2/deepstats/aggregate"
)

func RequisiteAttrs(params []string, gran string) []string {
	if gran == "" {
		gran = "d"
	}
	requisite := make(map[string]bool)
	for _, v := range params {
		if _, ok := formula.TransmitFormulaCollection[v]; ok {
			attrs := formula.TransmitFormulaCollection[v].Requisite(v, gran)
			for _, a := range attrs {
				requisite[a] = true
			}
		} else {
			requisite[v] = true
		}
	}

	attrs := make([]string, len(requisite))
	size := 0
	for k, _ := range requisite {
		attrs[size] = k
		size++
	}
	return attrs
}

func min(a int, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func Calculate(params, convertParams []string, res []*aggregate.AggregateResult, gran string, limit int) []map[string]string {
	col := make(map[string][]int)
	for _, v := range convertParams {
		col[v] = make([]int, limit)
	}

	// reverse aggregate result
	// arrange newest day ahead.
	for _, ch := range res {
		col[ch.Event] = make([]int, limit)
		size := min(len(ch.Counts), limit)
		for i := 0; i < size; i++ {
			col[ch.Event][size-i-1] = ch.Counts[i].Count
		}
	}

	// calculate according event
	interm := make(map[string][]string)
	for _, v := range params {
		if _, ok := formula.TransmitFormulaCollection[v]; ok {
			res := formula.TransmitFormulaCollection[v].Calculate(v, col, gran, limit)
			for k, _ := range res {
				interm[k] = res[k]
			}
		} else {
			interm[v] = make([]string, len(col[v]))
			for i := 0; i < len(col[v]); i++ {
				interm[v][i] = strconv.Itoa(col[v][i])
			}
		}
	}
	if len(params) == 0 {
		for v, _ := range col {
			interm[v] = make([]string, len(col[v]))
			for i := 0; i < len(col[v]); i++ {
				interm[v][i] = strconv.Itoa(col[v][i])
			}
		}
	}

	// change key order, a["event"]["day"] -> a["day"]["event"]
	convert := make([]map[string]string, limit)
	for i := 0; i < limit; i++ {
		convert[i] = make(map[string]string)
	}
	for k, _ := range interm {
		size := min(len(interm[k]), limit)
		for i := 0; i < size; i++ {
			convert[i][k] = interm[k][i]
		}
	}
	return convert
}

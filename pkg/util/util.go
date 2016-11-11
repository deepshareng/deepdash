package util

import (
	"strings"
)

var convertTable = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func SequenceToShortHash(seq int, width int) string {
	// if seq overflow with, will be ignored
	base := len(convertTable)

	rst := make([]string, width)

	for i := width - 1; i >= 0; i-- {
		rst[i] = string(convertTable[seq%base])
		seq = seq / base
	}

	return strings.Join(rst, "")
}

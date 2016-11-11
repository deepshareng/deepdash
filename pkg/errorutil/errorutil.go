package errorutil

import (
	"strings"

	"github.com/MISingularity/deepdash/api"
	"github.com/MISingularity/deepdash/pkg/httputil"
	"github.com/MISingularity/deepdash/pkg/log"
)

func ProcessMongoError(f *httputil.Ginframework, err error) bool {
	if err != nil && !strings.Contains(err.Error(), "not found") {
		log.Error(err)
		if strings.Contains(err.Error(), "EOF") || strings.Contains(err.Error(), "Closed explicitly") {
			log.Fatal("Mongo Pipeline Broken!", err)
			f.WriteHTTPError(api.ErrMongoBrokenPipe, "")
		} else {
			f.WriteHTTPError(api.ErrMongoOperationFail, "")
		}
		return true
	}
	return false
}

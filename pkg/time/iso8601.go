package time

import (
	"time"

	"muck0120/youtube2csv/pkg/errors"

	"github.com/sosodev/duration"
)

func ParseISO8601Duration(iso8601Duration string) (time.Duration, error) {
	d, err := duration.Parse(iso8601Duration)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	return d.ToTimeDuration(), nil
}

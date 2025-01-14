package time

import (
	"time"

	"github.com/sosodev/duration"

	"github.com/muck0120/youtube2csv/internal/pkg/errors"
)

func ParseISO8601Duration(iso8601Duration string) (time.Duration, error) {
	d, err := duration.Parse(iso8601Duration)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	return d.ToTimeDuration(), nil
}

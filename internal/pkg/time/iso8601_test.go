package time_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	pkgtime "github.com/muck0120/youtube2csv/internal/pkg/time"
)

func TestParseISO8601Duration(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		arg       string
		want      time.Duration
		wantError bool
	}{
		{
			name:      "有効な ISO8601 duration 'PT2H30M45S' をパースできる",
			arg:       "PT2H30M45S",
			want:      (2 * time.Hour) + (30 * time.Minute) + (45 * time.Second),
			wantError: false,
		},
		{
			name:      "無効な ISO8601 duration 'InvalidDuration' をパースできない",
			arg:       "InvalidDuration",
			want:      0,
			wantError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got, err := pkgtime.ParseISO8601Duration(test.arg)
			if test.wantError {
				assert.Error(t, err)

				return
			}

			assert.Equal(t, test.want, got)
		})
	}
}

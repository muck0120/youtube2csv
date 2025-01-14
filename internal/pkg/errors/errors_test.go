package errors_test

import (
	"errors"
	"fmt"
	"log/slog"
	"testing"

	goerrors "github.com/go-errors/errors"
	"github.com/stretchr/testify/assert"

	pkgerrors "github.com/muck0120/youtube2csv/internal/pkg/errors"
)

func TestLogStackTrace(t *testing.T) {
	t.Parallel()

	err := errors.New("something went wrong")
	goerr := goerrors.Wrap(err, 0)

	tests := []struct {
		name string
		args error
		want slog.Attr
	}{
		{
			name: "エラーが nil の場合、空の stacktrace を返す",
			args: nil,
			want: slog.Any("stacktrace", []any{}),
		},
		{
			name: "エラーが go-errors/errors.Error でない場合、エラーの詳細を返す",
			args: err,
			want: slog.String("details", fmt.Sprintf("%+v", err)),
		},
		{
			name: "エラーが go-errors/errors.Error の場合、stacktrace を返す",
			args: goerr,
			want: slog.Any("stacktrace", goerr.StackFrames()),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := pkgerrors.LogStackTrace(test.args)
			assert.Equal(t, test.want, got)
		})
	}
}

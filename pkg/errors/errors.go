package errors

import (
	"errors"
	"fmt"
	"log/slog"

	goerrors "github.com/go-errors/errors"
)

func WithStack(err error) error {
	return goerrors.Wrap(err, 0)
}

func LogStackTrace(err error) slog.Attr {
	if err == nil {
		return slog.Any("stacktrace", []any{})
	}

	var goerror *goerrors.Error
	if !errors.As(err, &goerror) {
		return slog.String("details", fmt.Sprintf("%+v", err))
	}

	return slog.Any("stacktrace", goerror.StackFrames())
}

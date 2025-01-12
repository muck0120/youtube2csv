package youtube_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"muck0120/youtube2csv/internal/domain/youtube"
)

func TestVideo_GetNumber(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		model *youtube.Video
		want  int
	}{
		{
			name:  "タイトルに # が含まれている場合、その数字を取得できる",
			model: &youtube.Video{Title: "Title #123"},
			want:  123,
		},
		{
			name:  "タイトルに # が含まれていない場合、0 を返す",
			model: &youtube.Video{Title: "Title"},
			want:  0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, test.want, test.model.GetNumber())
		})
	}
}

func TestVideo_GetURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		model *youtube.Video
		want  string
	}{
		{
			name:  "Valid ID",
			model: &youtube.Video{ID: "abc123"},
			want:  "https://youtu.be/abc123",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, test.want, test.model.GetURL())
		})
	}
}

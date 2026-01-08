package utils

import "testing"

func TestExtractCommentFromLine(t *testing.T) {
	tests := []struct {
		name string
		line string
		want string
	}{
		{
			name: "no comment",
			line: "this is a line",
			want: "",
		},
		{
			name: "simple comment",
			line: "this is a line # comment",
			want: "comment",
		},
		{
			name: "comment with leading spaces",
			line: "this is a line #    comment",
			want: "comment",
		},
		{
			name: "comment at start",
			line: "# comment",
			want: "comment",
		},
		{
			name: "comment with hash inside",
			line: "this is a line # comment # inside",
			want: "comment # inside",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExtractCommentFromLine(tt.line); got != tt.want {
				t.Errorf("ExtractCommentFromLine() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRemoveCommentFromLine(t *testing.T) {
	tests := []struct {
		name string
		line string
		want string
	}{
		{
			name: "no comment",
			line: "this is a line",
			want: "this is a line",
		},
		{
			name: "simple comment",
			line: "this is a line # comment",
			want: "this is a line",
		},
		{
			name: "comment with spaces",
			line: "this is a line   # comment",
			want: "this is a line",
		},
		{
			name: "comment at start",
			line: "# comment",
			want: "",
		},
		{
			name: "empty line",
			line: "",
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RemoveCommentFromLine(tt.line); got != tt.want {
				t.Errorf("RemoveCommentFromLine() = %v, want %v", got, tt.want)
			}
		})
	}
}

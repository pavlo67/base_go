package filelib

import "testing"

func TestJoin(t *testing.T) {
	tests := []struct {
		name     string
		basePath string
		path     string
		want     string
	}{
		{
			name:     "",
			basePath: "/a/b/c",
			path:     "d",
			want:     "/a/b/c/d",
		},
		{
			name:     "",
			basePath: "/a/b/c",
			path:     "../d",
			want:     "/a/b/d",
		},
		{
			name:     "",
			basePath: "/a/b/c",
			path:     "/d",
			want:     "/d",
		},
		{
			name:     "",
			basePath: "/a/b/c",
			path:     "d/",
			want:     "/a/b/c/d",
		},
		{
			name:     "",
			basePath: "/a/b/c",
			path:     "..",
			want:     "/a/b",
		},
		{
			name:     "",
			basePath: "/a/b/c",
			path:     "",
			want:     "/a/b/c",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Join(tt.basePath, tt.path); got != tt.want {
				t.Errorf("Join() = %v, want %v", got, tt.want)
			}
		})
	}
}

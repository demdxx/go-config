package goconfig

import (
	"reflect"
	"testing"
)

func TestCliParser(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want map[string]string
		err  bool
	}{
		{
			name: "long-flag",
			args: []string{"--http-listen=addr:test"},
			want: map[string]string{"http-listen": "addr:test"},
		},
		{
			name: "short-flag",
			args: []string{"-v"},
			want: map[string]string{"v": "true"},
		},
		{
			name: "short-flag-value",
			args: []string{"-v", "1"},
			want: map[string]string{"v": "1"},
		},
		{
			name: "short-flag-value2",
			args: []string{"-v", "1", "-v", "2"},
			want: map[string]string{"v": "2"},
		},
		{
			name: "short-flag-value3",
			args: []string{"-v", "1", "-v", "2", "-v", "3"},
			want: map[string]string{"v": "3"},
		},
		{
			name: "mixed-flags",
			args: []string{"-v", "1", "--http-listen=addr:test", "-v", "2", "--http-listen=addr:test2"},
			want: map[string]string{"v": "2", "http-listen": "addr:test2"},
		},
		{
			name: "mixed-flags2",
			args: []string{"-v", "--http-listen=addr:test", "--grpc-listen", "addr:test2", "--log-level", "debug"},
			want: map[string]string{"v": "true", "http-listen": "addr:test", "grpc-listen": "addr:test2", "log-level": "debug"},
		},
		{
			name: "invalid-flag",
			args: []string{"invalid"},
			want: map[string]string{},
			err:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseCommandFlags(tt.args)
			if err != nil && !tt.err {
				t.Errorf("parseCommandFlags() error = %v", err)
			} else if err == nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseCommandFlags() = %v, want %v", got, tt.want)
			}
		})
	}
}

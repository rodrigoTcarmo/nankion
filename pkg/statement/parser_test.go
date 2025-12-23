package statement

import (
	"testing"
	"time"
)

func TestParseOFXDate(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    time.Time
		wantErr bool
	}{
		{
			name:    "valid date YYYYMMDD",
			input:   "20251213",
			want:    time.Date(2025, 12, 13, 0, 0, 0, 0, time.UTC),
			wantErr: false,
		},
		{
			name:    "valid date with time YYYYMMDDHHMMSS",
			input:   "20251213120000",
			want:    time.Date(2025, 12, 13, 12, 0, 0, 0, time.UTC),
			wantErr: false,
		},
		{
			name:    "valid date with leading spaces",
			input:   "  20240101",
			want:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			wantErr: false,
		},
		{
			name:    "valid date with trailing spaces",
			input:   "20231231  ",
			want:    time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC),
			wantErr: false,
		},
		{
			name:    "date from sample OFX - salary",
			input:   "20251219",
			want:    time.Date(2025, 12, 19, 0, 0, 0, 0, time.UTC),
			wantErr: false,
		},
		{
			name:    "date from sample OFX - pix",
			input:   "20251122",
			want:    time.Date(2025, 11, 22, 0, 0, 0, 0, time.UTC),
			wantErr: false,
		},
		{
			name:    "empty string",
			input:   "",
			want:    time.Time{},
			wantErr: true,
		},
		{
			name:    "too short",
			input:   "2025121",
			want:    time.Time{},
			wantErr: true,
		},
		{
			name:    "invalid characters",
			input:   "abcdefgh",
			want:    time.Time{},
			wantErr: true,
		},
		{
			name:    "invalid format - Brazilian style",
			input:   "13/12/2025",
			want:    time.Time{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseOFXDate(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseOFXDate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !got.Equal(tt.want) {
				t.Errorf("ParseOFXDate() = %v, want %v", got, tt.want)
			}
		})
	}
}

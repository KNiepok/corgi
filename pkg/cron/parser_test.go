package cron

import (
	"testing"
)

func TestParseRequestedIntervalToCron(t *testing.T) {
	tests := []struct {
		name     string
		interval string
		want     string
		wantErr  bool
	}{
		{
			name:     "daily @ 17:20 passes",
			interval: "daily @ 17:20",
			want:     "20 17 * * *",
			wantErr:  false,
		},
		{
			name:     "daily @ 9 passes",
			interval: "daily @ 9",
			want:     "0 9 * * *",
			wantErr:  false,
		},
		{
			name:     "daily @ 25:12 fails",
			interval: "daily @ 25:12",
			wantErr:  true,
		},
		{
			name:     "daily @ 12:89 fails",
			interval: "daily @ 12:89",
			wantErr:  true,
		},
		{
			name:     "daily1730 fails",
			interval: "daily1730",
			wantErr:  true,
		},
		{
			name:     "daily @ 17:xo fails",
			interval: "daily @ 17:xo",
			wantErr:  true,
		},
		{
			name:     "weekly @ FRI 17:20 passes",
			interval: "weekly @ FRI 17:20",
			want:     "20 17 * * FRI",
			wantErr:  false,
		},
		{
			name:     "weekly @ SAT 9 passes",
			interval: "weekly @ SAT 9",
			want:     "0 9 * * SAT",
			wantErr:  false,
		},
		{
			name:     "weekly @ 25:12 fails",
			interval: "weekly @ 25:12",
			wantErr:  true,
		},
		{
			name:     "weekly @ ABC 17:00 fails",
			interval: "weekly @ ABC 17:00",
			wantErr:  true,
		},
		{
			name:     "weekly @ FRI 17:xx fails",
			interval: "weekly @ FRI 17:xx",
			wantErr:  true,
		},
		{
			name:     "ugabuga fails",
			interval: "ugabuga",
			wantErr:  true,
		},
		{
			name:     "empty string fails",
			interval: "",
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseIntervalToCron(tt.interval)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseIntervalToCron() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseIntervalToCron() got = %v, want %v", got, tt.want)
			}
		})
	}
}

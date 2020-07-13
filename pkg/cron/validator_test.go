package cron

import (
	"context"
	"testing"
)

func TestValidator_Validate(t *testing.T) {
	tests := []struct {
		name     string
		interval string
		wantErr  bool
	}{
		{
			name:     "daily @ 17:20 passes",
			interval: "daily @ 17:20",
			wantErr:  false,
		},
		{
			name:     "daily @ 9 passes",
			interval: "daily @ 9",
			wantErr:  false,
		},
		{
			name:     "daily @ 25:12 returns error",
			interval: "daily @ 25:12",
			wantErr:  true,
		},
		{
			name:     "daily @ 12:89 returns error",
			interval: "daily @ 12:89",
			wantErr:  true,
		},
		{
			name:     "daily1730 returns error",
			interval: "daily1730",
			wantErr:  true,
		},
		{
			name:     "daily @ 17:xo returns error",
			interval: "daily @ 17:xo",
			wantErr:  true,
		},
		{
			name:     "weekly @ FRI 17:20 passes",
			interval: "weekly @ FRI 17:20",
			wantErr:  false,
		},
		{
			name:     "weekly @ SAT 9 passes",
			interval: "weekly @ SAT 9",
			wantErr:  false,
		},
		{
			name:     "weekly @ 25:12 returns error",
			interval: "weekly @ 25:12",
			wantErr:  true,
		},
		{
			name:     "weekly @ ABC 17:00 returns error",
			interval: "weekly @ ABC 17:00",
			wantErr:  true,
		},
		{
			name:     "weekly @ FRI 17:xx returns error",
			interval: "weekly @ FRI 17:xx",
			wantErr:  true,
		},
		{
			name:     "ugabuga returns error",
			interval: "ugabuga",
			wantErr:  true,
		},
		{
			name:     "empty string returns error",
			interval: "",
			wantErr:  true,
		},
		{
			name:     "daily @ 17:20 passes",
			interval: "daily @ 17:20",
			wantErr:  false,
		},
		{
			name:     "daily @ 9 passes",
			interval: "daily @ 9",
			wantErr:  false,
		},
		{
			name:     "daily @ 25:12 fails",
			interval: "daily @ 25:12",
			wantErr:  true,
		},
		{
			name:     "daily @ 12:89 returns error",
			interval: "daily @ 12:89",
			wantErr:  true,
		},
		{
			name:     "daily1730 returns error",
			interval: "daily1730",
			wantErr:  true,
		},
		{
			name:     "daily @ 17:xo returns error",
			interval: "daily @ 17:xo",
			wantErr:  true,
		},
		{
			name:     "weekly @ FRI 17:20 passes",
			interval: "weekly @ FRI 17:20",
			wantErr:  false,
		},
		{
			name:     "weekly @ SAT 9 passes",
			interval: "weekly @ SAT 9",
			wantErr:  false,
		},
		{
			name:     "weekly @ 25:12 returns error",
			interval: "weekly @ 25:12",
			wantErr:  true,
		},
		{
			name:     "weekly @ ABC 17:00 returns error",
			interval: "weekly @ ABC 17:00",
			wantErr:  true,
		},
		{
			name:     "weekly @ FRI 17:xx returns error",
			interval: "weekly @ FRI 17:xx",
			wantErr:  true,
		},
		{
			name:     "ugabuga returns error",
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
			v := NewValidator()
			if err := v.Validate(context.Background(), tt.interval); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

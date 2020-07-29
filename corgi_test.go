package corgi

import (
	"reflect"
	"testing"
	"time"
)

func TestNewSubscriptionDetails(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name    string
		args    args
		want    SubscriptionDetails
		wantErr bool
	}{
		{
			name: "daily @ 17:20 can be parsed",
			args: args{
				input: "daily @ 17:20",
			},
			want: SubscriptionDetails{
				Mode:   SubModeDaily,
				Hour:   17,
				Minute: 20,
			},
		},
		{
			name: "daily @ 9 can be parsed",
			args: args{
				input: "daily @ 9",
			},
			want: SubscriptionDetails{
				Mode:   SubModeDaily,
				Hour:   9,
				Minute: 0,
			},
		},
		{
			name: "daily @ 25:12 cannot be parsed",
			args: args{
				input: "daily @ 25:12",
			},
			wantErr: true,
		},
		{
			name: "daily @ 12:89 cannot be parsed",
			args: args{
				input: "daily @ 12:89",
			},
			wantErr: true,
		},
		{
			name: "daily1730 cannot be parsed",
			args: args{
				input: "daily1730",
			},
			wantErr: true,
		},
		{
			name: "weekly @ FRI 17:20 can be parsed",
			args: args{
				input: "weekly @ FRI 17:20",
			},
			want: SubscriptionDetails{
				Mode:      SubModeWeekly,
				DayOfWeek: time.Friday,
				Hour:      17,
				Minute:    20,
			},
		},
		{
			name: "weekly @ SAT 9 can be parsed",
			args: args{
				input: "weekly @ SAT 9",
			},
			want: SubscriptionDetails{
				Mode:      SubModeWeekly,
				DayOfWeek: time.Saturday,
				Hour:      9,
				Minute:    0,
			},
		},
		{
			name: "weekly @ 25:12 can not be parsed",
			args: args{
				input: "weekly @ 25:12",
			},
			wantErr: true,
		},
		{
			name: "weekly @ ABC 17:00 can not be parsed",
			args: args{
				input: "weekly @ ABC 17:00",
			},
			wantErr: true,
		},
		{
			name: "weekly @ FRI 17:xx can not be parsed",
			args: args{
				input: "weekly @ FRI 17:xx",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewSubscriptionDetails(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSubscriptionDetails() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSubscriptionDetails() got = %v, want %v", got, tt.want)
			}
		})
	}
}

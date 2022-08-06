package env

import "testing"

func TestFromString(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want Environment
	}{
		{"development", args{"development"}, EnvDevelopment},
		{"staging", args{"staging"}, EnvStaging},
		{"production", args{"production"}, EnvProduction},
		{"unknown", args{"unknown"}, EnvDevelopment},
		{"empty", args{""}, EnvDevelopment},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FromString(tt.args.str); got != tt.want {
				t.Errorf("FromString() = %v, want %v", got, tt.want)
			}
		})
	}
}

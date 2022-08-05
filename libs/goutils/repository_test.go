package goutils

import "testing"

func TestValidateRepositoryName(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"empty", args{""}, true},
		{"valid", args{"angular/angular-ja"}, false},
		{"invalid", args{"angular-ja"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateRepositoryName(tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("ValidateRepositoryName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

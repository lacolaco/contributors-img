package model_test

import (
	"reflect"
	"testing"

	"contrib.rocks/libs/goutils/model"
)

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
			if err := model.ValidateRepositoryName(tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("ValidateRepositoryName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepositoryString_Object(t *testing.T) {
	tests := []struct {
		name string
		r    model.RepositoryString
		want *model.Repository
	}{
		{
			name: "valid",
			r:    "angular/angular-ja",
			want: &model.Repository{Owner: "angular", RepoName: "angular-ja"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.Object(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RepositoryString.Object() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepository_String(t *testing.T) {
	tests := []struct {
		name string
		r    *model.Repository
		want string
	}{
		{
			name: "valid",
			r:    &model.Repository{Owner: "angular", RepoName: "angular-ja"},
			want: "angular/angular-ja",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.String(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Repository.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

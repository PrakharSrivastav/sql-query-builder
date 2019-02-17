/*
Package querybuilder helps to generate sql queries in different drivers.
This package can be best used with the scenarios where the structure of the domains models is unknown beforehand.
*/
package sql

import (
	"testing"

	"github.com/PrakharSrivastav/sql-query-builder/sql/core"
)

func TestNewQueryBuilder(t *testing.T) {
	t.Parallel()
	type args struct {
		driver int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "ok-driver",
			args:    args{driver: core.PGSQL},
			wantErr: false,
		},
		{
			name:    "bad-driver",
			args:    args{driver: 7},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewQueryBuilder(tt.args.driver)
			if err != nil && !tt.wantErr {
				t.Errorf("NewQueryBuilder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestNewSingletonQueryBuilder(t *testing.T) {
	t.Parallel()
	type args struct {
		driver int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "ok-driver",
			args:    args{driver: core.PGSQL},
			wantErr: false, // set a driver first, this will persist for next calls
		},
		{
			name:    "bad-driver",
			args:    args{driver: 100},
			wantErr: false, // once a driver is set, even wrong drivers with not error out
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewSingletonQueryBuilder(tt.args.driver)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSingletonQueryBuilder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

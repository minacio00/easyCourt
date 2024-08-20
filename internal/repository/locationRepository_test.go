package repository

import (
	"testing"

	"github.com/minacio00/easyCourt/internal/model"
	"gorm.io/gorm"
)

func Test_locationRepository_CreateLocation(t *testing.T) {
	type fields struct {
		db *gorm.DB
	}
	type args struct {
		location *model.Location
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &locationRepository{
				db: tt.fields.db,
			}
			if err := l.CreateLocation(tt.args.location); (err != nil) != tt.wantErr {
				t.Errorf("locationRepository.CreateLocation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

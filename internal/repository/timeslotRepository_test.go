package repository

import (
	"log"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/minacio00/easyCourt/internal/db"
	"github.com/minacio00/easyCourt/internal/model"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

func loadConfig() {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	projectRoot := filepath.Join(basepath, "../..")

	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(projectRoot)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

func Test_TimeslotRepository_GetTimeslotById(t *testing.T) {
	loadConfig()
	db.Init()
	testDB, _ := db.DB.DB()
	defer testDB.Close()

	cleanDB := db.DB.Begin()
	defer cleanDB.Rollback()

	court := &model.Court{}
	db.DB.First(court)

	startTime, _ := time.Parse(time.TimeOnly, "09:00:00")
	endTime, _ := time.Parse(time.TimeOnly, "10:00:00")
	testTimeslot := &model.Timeslot{
		StartTime: startTime,
		EndTime:   endTime,
		CourtID:   nil,
		Day:       model.SegundaFeira,
		IsActive:  false,
	}
	cleanDB.Create(testTimeslot)
	log.Printf(`
    %v\n,
    %v\n,
    %v\n,
    %v\n,
    %v\n,
    %v\n,
    `, testTimeslot.StartTime, testTimeslot.EndTime, testTimeslot.CourtID, testTimeslot.Day, testTimeslot.IsActive, testTimeslot.ID)

	type fields struct {
		db *gorm.DB
	}
	type args struct {
		id int
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *int
		wantErr bool
	}{
		{
			name:    "existing timeslot",
			fields:  fields{db: cleanDB},
			args:    args{id: testTimeslot.ID},
			want:    &testTimeslot.ID,
			wantErr: false,
		},
		{
			name:    "non-existent timeslot",
			fields:  fields{db: cleanDB},
			args:    args{id: 999},
			want:    nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewTimeslotRepository(tt.fields.db)

			got, err := repo.GetTimeslotByID(tt.args.id)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetTimeslotById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want == nil {
				if got != nil {
					t.Errorf("Wanted nil got %v", got)
				}
				return
			}

			if !tt.wantErr && got.ID != *tt.want {
				t.Errorf("GetTimeslotById() = %v, want %v", got.ID, tt.want)
			}
		})
	}
}

package handlers

import (
	"log"
	"os"
	"testing"

	"github.com/minacio00/easyCourt/database"
	"github.com/minacio00/easyCourt/types"
)

func TestMain(m *testing.M) {
	os.Setenv("APP_ENV", "test")
	log.Println("Do stuff BEFORE the tests!")
	database.Connectdb()
	os.Exit(m.Run())
}

type args struct {
	tenant *types.Tenant
}

func (a *args) defaultValues() {
	a.tenant = &types.Tenant{
		Name:        "amado222",
		Email:       "email@email.com",
		Password:    "senha",
		TrialPeriod: false,
	}
}

func Test_checkEmail(t *testing.T) {

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "novo email", args: args{}, wantErr: false},
		{name: "email iqual", args: args{}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.defaultValues()
			if err := checkEmail(tt.args.tenant); (err != nil) != tt.wantErr {
				t.Errorf("checkEmail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
		database.Db.Save(tt.args.tenant)
	}
}

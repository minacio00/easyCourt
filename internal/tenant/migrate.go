package tenant

import "github.com/minacio00/easyCourt/database"

func Migrate() {
	database.Db.AutoMigrate(&Tenant{})
}

package jobs

import "github.com/aikintech/scim-api/pkg/database"

func SeedRolesAndPermissions() {
	database.RunDatabaseSeeder()
}

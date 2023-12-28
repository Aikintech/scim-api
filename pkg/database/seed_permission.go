package database

import (
	"errors"
	"fmt"

	"github.com/aikintech/scim-api/pkg/models"
	mapSet "github.com/deckarep/golang-set/v2"
	"gorm.io/gorm"
)

type Perm struct {
	name        string
	displayName string
	module      string
}

func SeedPermissions() {
	permissions := mapSet.NewSet[Perm]()

	// Posts
	permissions.Add(Perm{name: "create-post", displayName: "Create post", module: "Post"})
	permissions.Add(Perm{name: "read-posts", displayName: "View all posts", module: "Post"})
	permissions.Add(Perm{name: "read-post", displayName: "View post", module: "Post"})
	permissions.Add(Perm{name: "update-post", displayName: "Update post", module: "Post"})
	permissions.Add(Perm{name: "delete-post", displayName: "Delete post", module: "Post"})

	// Prayer requests
	permissions.Add(Perm{name: "read-prayers", displayName: "View all prayer requests", module: "Prayer Requests"})
	permissions.Add(Perm{name: "read-prayer", displayName: "View prayer request", module: "Prayer Requests"})
	permissions.Add(Perm{name: "update-prayer", displayName: "Update prayer request", module: "Prayer Requests"})

	// Events
	permissions.Add(Perm{name: "create-event", displayName: "Create event", module: "Events"})
	permissions.Add(Perm{name: "read-events", displayName: "View all events", module: "Events"})
	permissions.Add(Perm{name: "read-event", displayName: "View event", module: "Events"})
	permissions.Add(Perm{name: "update-event", displayName: "Update event", module: "Events"})
	permissions.Add(Perm{name: "delete-event", displayName: "Delete event", module: "Events"})

	// Donations
	permissions.Add(Perm{name: "read-donations", displayName: "View all donations", module: "Donations"})
	permissions.Add(Perm{name: "read-donation", displayName: "View donation", module: "Donations"})

	// Insert
	for _, p := range permissions.ToSlice() {
		if err := DB.Where("name = ?", p.name).
			Assign(models.Permission{Name: p.name, DisplayName: p.displayName, Description: "", Module: p.module}).
			FirstOrCreate(&models.Permission{}).Error; err != nil {

			if !errors.Is(err, gorm.ErrRecordNotFound) {
				fmt.Printf("Permissions seeder: An error occurred while seeding %s. Error: %s\n", p.name, err.Error())
			}
		}
	}
}

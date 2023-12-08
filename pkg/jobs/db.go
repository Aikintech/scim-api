package jobs

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/aikintech/scim-api/pkg/constants"
)

// BackupDatabase is a job that backs up the database
func BackupDatabase() {
	// Perform backup using pg_dump
	now := time.Now().Format("2006-01-02T15:04:05")
	filename := fmt.Sprintf("%s_%s___database.sql", now, constants.DB_DATABASE)

	fmt.Println("Backing up database...")
	fmt.Println("Filename:", filename)

	cmd := exec.Command("pg_dump", fmt.Sprintf("-U%s", constants.DB_USERNAME), fmt.Sprintf("-h%s", constants.DB_HOST), fmt.Sprintf("-p%s", constants.DB_PORT), constants.DB_DATABASE)
	output, err := cmd.Output()

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Write the backup to a file
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	_, err = file.Write(output)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	fmt.Println("Backup complete. File:", filename)
}

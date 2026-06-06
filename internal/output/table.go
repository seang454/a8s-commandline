package output

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/yourname/a8s/internal/models"
)

func PrintUsersTable(users []models.User) {
	if len(users) == 0 {
		fmt.Println("No users found.")
		return
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Name", "Email", "Admin", "Active", "Created At"})
	table.SetBorder(false)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	for _, u := range users {
		admin := "No"
		if u.IsAdmin {
			admin = "Yes"
		}
		active := "No"
		if u.Active {
			active = "Yes"
		}
		table.Append([]string{u.ID, u.Name, u.Email, admin, active, u.CreatedAt})
	}
	table.Render()
}

func PrintProjectsTable(projects []models.Project) {
	if len(projects) == 0 {
		fmt.Println("No projects found.")
		return
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Name", "Description", "Owner ID", "Created At"})
	table.SetBorder(false)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	for _, p := range projects {
		table.Append([]string{p.ID, p.Name, p.Description, p.OwnerID, p.CreatedAt})
	}
	table.Render()
}

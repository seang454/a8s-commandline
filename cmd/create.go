package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yourname/a8s/internal/api"
	"github.com/yourname/a8s/internal/models"
)

var (
	userName  string
	userEmail string
	isAdmin   bool
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a resource",
}

var createUserCmd = &cobra.Command{
	Use:   "user",
	Short: "Create a new user",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := api.NewClient(
			viper.GetString("api_url"),
			viper.GetString("api_token"),
		)
		user := models.CreateUserRequest{
			Name:    userName,
			Email:   userEmail,
			IsAdmin: isAdmin,
		}
		created, err := client.CreateUser(user)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✓ User created: %s (ID: %s)\n", created.Name, created.ID)
		return nil
	},
}

func init() {
	createUserCmd.Flags().StringVar(&userName, "name", "", "user full name (required)")
	createUserCmd.Flags().StringVar(&userEmail, "email", "", "user email (required)")
	createUserCmd.Flags().BoolVar(&isAdmin, "admin", false, "grant admin privileges")
	createUserCmd.MarkFlagRequired("name")
	createUserCmd.MarkFlagRequired("email")
	createCmd.AddCommand(createUserCmd)
}

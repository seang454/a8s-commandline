package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yourname/a8s/internal/api"
)

var deleteID string

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a resource",
}

var deleteUserCmd = &cobra.Command{
	Use:   "user",
	Short: "Delete a user by ID",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := api.NewClient(
			viper.GetString("api_url"),
			viper.GetString("api_token"),
		)
		if err := client.DeleteUser(deleteID); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✓ User %s deleted successfully\n", deleteID)
		return nil
	},
}

func init() {
	deleteUserCmd.Flags().StringVar(&deleteID, "id", "", "user ID to delete (required)")
	deleteUserCmd.MarkFlagRequired("id")
	deleteCmd.AddCommand(deleteUserCmd)
}

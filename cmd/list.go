package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yourname/a8s/internal/api"
	"github.com/yourname/a8s/internal/output"
)

var listAll bool

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List resources",
}

var listUsersCmd = &cobra.Command{
	Use:   "users",
	Short: "List all users",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := api.NewClient(
			viper.GetString("api_url"),
			viper.GetString("api_token"),
		)
		users, err := client.ListUsers(listAll)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		switch outputFmt {
		case "json":
			output.PrintJSON(users)
		default:
			output.PrintUsersTable(users)
		}
		return nil
	},
}

var listProjectsCmd = &cobra.Command{
	Use:   "projects",
	Short: "List all projects",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := api.NewClient(
			viper.GetString("api_url"),
			viper.GetString("api_token"),
		)
		projects, err := client.ListProjects()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		switch outputFmt {
		case "json":
			output.PrintJSON(projects)
		default:
			output.PrintProjectsTable(projects)
		}
		return nil
	},
}

func init() {
	listUsersCmd.Flags().BoolVar(&listAll, "all", false, "include inactive users")
	listCmd.AddCommand(listUsersCmd)
	listCmd.AddCommand(listProjectsCmd)
}

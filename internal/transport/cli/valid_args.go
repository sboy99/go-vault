package cli

import (
	"context"

	"github.com/sboy99/go-vault/internal/app"
	"github.com/spf13/cobra"
)

func restoreValidArgs(backup *app.BackupService) func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		list, err := backup.ListBackups(context.Background(), 50, 0)
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}
		names := make([]string, 0, len(list)*2)
		for _, v := range list {
			names = append(names, v.BackupId, v.Name)
		}
		return names, cobra.ShellCompDirectiveNoFileComp
	}
}

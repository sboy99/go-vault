package cli

import (
	"github.com/sboy99/go-vault/internal/meta"
	"github.com/spf13/cobra"
)

func restoreBackupValidArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	backupMetaList, err := meta.ListBackupMeta(50, 0)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	names := make([]string, 0, len(backupMetaList)*2)
	for _, v := range backupMetaList {
		names = append(names, v.BackupId, v.Name)
	}
	return names, cobra.ShellCompDirectiveNoFileComp
}

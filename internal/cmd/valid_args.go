package cmd

import (
	"github.com/sboy99/go-vault/internal/meta"
	"github.com/spf13/cobra"
)

func restoreBackupValidArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	maxSize, offset := 15, 0
	backupMetaList, err := meta.ListBackupMeta(maxSize, offset)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	backupNames := make([]string, len(backupMetaList))
	for i, v := range backupMetaList {
		backupNames[i] = v.Name
	}
	return backupNames, cobra.ShellCompDirectiveNoFileComp
}

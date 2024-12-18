package pgdump

import "fmt"

func makeSqlComment(comment string) string {
	return fmt.Sprintf("--\n-- %s\n--\n", comment)
}

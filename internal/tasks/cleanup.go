package tasks

import (
	"context"
	"fmt"
)

func HandleCleanup(ctx context.Context, payload []byte) error {
	// TODO: implement cleanup logic (e.g. delete expired tokens, old logs, etc.)
	fmt.Println("[task] running system cleanup")
	return nil
}

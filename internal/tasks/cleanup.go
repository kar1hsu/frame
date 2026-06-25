package tasks

import (
	"context"

	"github.com/kar1hsu/frame/internal/app"
)

func HandleCleanup(ctx context.Context, payload []byte) error {
	// TODO: implement cleanup logic (e.g. delete expired tokens, old logs, etc.)
	app.Log.Info("[task] running system cleanup")
	return nil
}

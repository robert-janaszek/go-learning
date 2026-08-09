package hook

import "context"

func UseCancel() context.CancelFunc {
	return runtime.cancel
}

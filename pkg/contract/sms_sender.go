//go:generate mockery --name=SmsSender

package contract

import "context"

type SmsSender interface {
	Send(ctx context.Context, mobile, content string) error
}

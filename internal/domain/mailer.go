package domain

import (
	"context"

	"github.com/akfaiz/go-mailgen"
)

type Mailer interface {
	Send(ctx context.Context, msg *mailgen.Builder) error
}

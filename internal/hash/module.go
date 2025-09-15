package hash

import (
	"github.com/akfaiz/go-vue-starter-kit/internal/hash/argon2id"
	"github.com/akfaiz/go-vue-starter-kit/internal/hash/jwt"
	"go.uber.org/fx"
)

var Module = fx.Module("hash",
	fx.Provide(
		argon2id.NewHasher,
		jwt.NewJWTManager,
	),
)

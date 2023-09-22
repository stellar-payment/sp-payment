package scopeutil

import (
	"context"

	"github.com/stellar-payment/sp-payment/internal/commonkey"
	"github.com/stellar-payment/sp-payment/internal/indto"
	"github.com/stellar-payment/sp-payment/internal/util/ctxutil"
)

func ValidateScope(ctx context.Context, scope string) (ok bool) {
	usrScope := ctxutil.GetCtx[indto.UserScopeMap](ctx, commonkey.SCOPE_CTX_KEY)
	_, ok = usrScope[scope]
	return
}

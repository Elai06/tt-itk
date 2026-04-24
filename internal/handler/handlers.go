package handler

type Handlers struct {
	Wallet    WalletHandler
	Auth      AuthHandler
	Exchanger ExchangerHandler
}

func NewHandlers(h *WalletHandler, a *AuthHandler, e *ExchangerHandler) *Handlers {
	return &Handlers{Wallet: *h, Auth: *a, Exchanger: *e}
}

package handler

type Handlers struct {
	Wallet WalletHandler
	Auth   AuthHandler
}

func NewHandlers(h *WalletHandler, a *AuthHandler) *Handlers {
	return &Handlers{Wallet: *h, Auth: *a}
}

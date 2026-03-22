package dto

type WalletRequest struct {
	UUID          int64  `form:"valletId"`
	OperationType string `form:"operationType"`
	Amount        int64  `form:"amount"`
}

type WalletResponse struct {
	UUID   int64 `form:"walletId"`
	Amount int64 `form:"amount"`
}

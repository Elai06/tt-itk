package dto

type WalletRequest struct {
	UUID          int64  `form:"walletId"`
	OperationType string `form:"operationType"`
	Amount        int64  `form:"amount"`
}

type WalletResponse struct {
	UUID   int64 `json:"walletId"`
	Amount int64 `json:"amount"`
}
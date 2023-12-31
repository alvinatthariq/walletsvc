package entity

type Meta struct {
	Path       string `json:"path"`
	StatusCode int    `json:"status_code"`
	Status     string `json:"status"`
	Message    string `json:"message"`
	Timestamp  string `json:"timestamp"`
	Error      string `json:"error,omitempty"`
}

type HTTPEmptyResp struct {
	Meta Meta `json:"metadata"`
}

type HTTPInitWalletResp struct {
	Data   HTTPInitWalletData `json:"data"`
	Status string             `json:"status"`
}

type HTTPInitWalletData struct {
	Token string `json:"token"`
}

type HTTPWalletResp struct {
	Status string         `json:"status"`
	Data   HTTPWalletData `json:"data"`
}

type HTTPWalletData struct {
	Wallet Wallet `json:"wallet"`
}

type HTTPTransactionResp struct {
	Status string              `json:"status"`
	Data   HTTPTransactionData `json:"data"`
}

type HTTPTransactionData struct {
	Transactions []Transaction `json:"transactions"`
}

type HTTPDepositResp struct {
	Status string          `json:"status"`
	Data   HTTPDepositData `json:"data"`
}

type HTTPDepositData struct {
	Deposit Deposit `json:"deposit"`
}

type HTTPWithdrawResp struct {
	Status string           `json:"status"`
	Data   HTTPWithdrawData `json:"data"`
}

type HTTPWithdrawData struct {
	Withdraw Withdraw `json:"withdraw"`
}

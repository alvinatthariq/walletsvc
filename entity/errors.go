package entity

import "fmt"

const (
	// Code SQL Error from https://github.com/go-sql-driver/mysql/blob/master/errors.go
	CodeMySQLDuplicateEntry             = 1062
	CodeMySQLForeignKeyConstraintFailed = 1452
	CodeMySQLTableNotExist              = 1146
)

var (
	ErrorAccountWalletAlreadyExist      error = fmt.Errorf("Account Wallet Already Exist")
	ErrorWalletAlreadyExist             error = fmt.Errorf("Wallet Already Exist")
	ErrorWalletAlreadyEnabled           error = fmt.Errorf("Wallet Already Enabled")
	ErrorDepositReferenceIDMustBeUnique error = fmt.Errorf("Deposit Reference ID Must Be Unique")
	ErrorWalletDisabled                 error = fmt.Errorf("Wallet Disabled")
	ErrorWalletNotFound                 error = fmt.Errorf("Wallet Not Found")
	ErrorInvalidAuthToken               error = fmt.Errorf("Invalid Auth Token")
)

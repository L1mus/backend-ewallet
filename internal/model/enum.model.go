package model

type OauthProvider string

const (
	OauthGoogle   OauthProvider = "google"
	OauthFacebook OauthProvider = "facebook"
)

type SubscribeStatus string

const (
	StatusActive      SubscribeStatus = "active"
	StatusUnsubscribe SubscribeStatus = "unsubscribe"
)

type TypeTransaction string

const (
	TransactionIncome  TypeTransaction = "income"
	TransactionExpense TypeTransaction = "expense"
)

type TypeActivityTransaction string

const (
	ActivityTransfer TypeActivityTransaction = "transfer"
	ActivityTopup    TypeActivityTransaction = "topup"
)

type StatusTransaction string

const (
	StatusPending StatusTransaction = "pending"
	StatusSuccess StatusTransaction = "success"
	StatusFailed  StatusTransaction = "failed"
)

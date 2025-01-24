package expense

type TransactionType string

const (
	TransactionTypeExpense  TransactionType = "expense"
	TransactionTypeIncoming TransactionType = "income"
)

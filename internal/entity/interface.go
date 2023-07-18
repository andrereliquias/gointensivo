package entity

type OrderRepositoryInterface interface {
	Save(order *Order) error
	GetTotalTransactions(order *Order) (int, error)
}

package operationtype

type Service interface {
	FindAll() (Indexed, error)
}

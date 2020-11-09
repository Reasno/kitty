package contract

type Marshaller interface {
	Marshal() ([]byte, error)
}

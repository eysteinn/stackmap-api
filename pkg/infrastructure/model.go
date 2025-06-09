package infrastructure

type Infrastructure interface {
	Create() error
	Delete() error
}

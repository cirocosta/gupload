package core

type Server interface {
	Listen() (err error)
	Close()
}

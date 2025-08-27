package domain

type ConnectionRepo interface {
	NextID() string
	Save(connection *Connection) error
	List() ([]*Connection, error)
	FindByID(id string) (*Connection, error)
	DeleteByID(id string) error
}

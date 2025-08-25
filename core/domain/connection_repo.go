package domain

type ConnectionRepo interface {
	NextID() string
	Save(connection *Connection) error
}

package ports

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"jamo/backend/internal/core/domain/entity"
)

type Repository interface {
	CreateProduct(product entity.Product) (interface{}, error)
	GetProduct(amount int) (interface{}, error)
	SubscribeToNewsLetter(body entity.Subscriber) error
	GetProductById(id string) (interface{}, error)
	GetCartItems(ids []primitive.ObjectID) (interface{}, error)
	CreateOrder(order entity.Order) (interface{}, error)
	GetOrderById(id string) (interface{}, error)
	UpdateOrderPayment(id string) (interface{}, error)
	GetOrders(page string) (interface{}, error)
	GetDashBoardValues() (interface{}, error)
	CreateContactMessage(body entity.ContactMessage) error
	GetAdminMsgs(page string) (interface{}, error)
}

package services

import (
	"jamo/backend/internal/core/domain/entity"
	"jamo/backend/internal/core/helper"
	ports "jamo/backend/internal/port"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type backendService struct {
	Repository ports.Repository
}

func NewService(repository ports.Repository) *backendService {
	return &backendService{
		Repository: repository,
	}
}

func (s *backendService) GetProduct(amount int) (interface{}, error) {
	return s.Repository.GetProduct(amount)
}

func (s *backendService) CreateProduct(product entity.Product) (interface{}, error) {
	product.ProductID = primitive.NewObjectID()
	product.CreatedAt = helper.ParseTimeToString(time.Now())

	return s.Repository.CreateProduct(product)
}

func (s *backendService) SubscribeToNewsLetter(body entity.Subscriber) error {
	body.ID = primitive.NewObjectID()
	body.SubscribedAt = helper.ParseTimeToString(time.Now())
	return s.Repository.SubscribeToNewsLetter(body)
}

func (s *backendService) GetProductById(ref string) (interface{}, error) {
	return s.Repository.GetProductById(ref)
}

func (s *backendService) ContactAdmin(body entity.ContactMessage) error {
	body.To = helper.Config.SMTPUsername
	body.From = body.Email
	body.SentAt = helper.ParseTimeToString(time.Now())
	if err := helper.SendMail("contactmail.html", body); err != nil {
		helper.LogEvent("ERROR", err.Error())
	}

	return s.Repository.CreateContactMessage(body)
}

func (s *backendService) GetCartItems(ids []string) (interface{}, error) {
	idHexes := []primitive.ObjectID{}

	for _, id := range ids {
		idHex, _ := primitive.ObjectIDFromHex(id)
		idHexes = append(idHexes, idHex)
	}

	return s.Repository.GetCartItems(idHexes)
}

func (s *backendService) CreateOrder(order entity.Order) (interface{}, error) {
	order.ID = primitive.NewObjectID()
	order.CreatedAt = helper.ParseTimeToString(time.Now())

	if err := order.Validate(); err != nil {
		return nil, err
	}

	order.CartSubtotal += 1500 // add 1500 for delivery fee

	id, err := s.Repository.CreateOrder(order)
	if err != nil {
		return nil, err
	}

	err = helper.SendMail("ordermail.html", entity.ContactMessage{
		To:      order.DeliveryInfo.RecipientEmail,
		Message: id.(string),
	})

	if err != nil {
		helper.LogEvent("ERROR", "sending mail to client on successful order: "+err.Error())
	}

	return id, nil
}

func (s *backendService) UpdateOrderPayment(id string) (interface{}, error) {
	return s.Repository.UpdateOrderPayment(id)
}

func (s *backendService) GetOrders(page string) (interface{}, error) {
	return s.Repository.GetOrders(page)
}

func (s *backendService) GetDashBoardValues() (interface{}, error) {
	return s.Repository.GetDashBoardValues()
}

func (s *backendService) GetAdminMsgs(page string) (interface{}, error) {
	return s.Repository.GetAdminMsgs(page)
}

func (s *backendService) GetOrderById(id string) (interface{}, error) {
	return s.Repository.GetOrderById(id)
}

func (s *backendService) UpdateDeliveryStatus(id string) error {
	return s.Repository.UpdateDeliveryStatus(id)
}

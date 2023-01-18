package repository

import (
	"fmt"
	"jamo/backend/internal/core/domain/entity"
	"jamo/backend/internal/core/helper"

	"go.mongodb.org/mongo-driver/bson/primitive"

	// ports "jamo/backend/internal/port"
	"context"
	"errors"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type DatabaseInfra struct {
	ProductCollection    *mongo.Collection
	NewsletterCollection *mongo.Collection
	OrderCollection      *mongo.Collection
	MessagesCollection   *mongo.Collection
}

func NewInfra(ProductCollection, NewsletterCollection, OrderCollection, MessagesCollection *mongo.Collection) *DatabaseInfra {
	return &DatabaseInfra{
		ProductCollection:    ProductCollection,
		NewsletterCollection: NewsletterCollection,
		OrderCollection:      OrderCollection,
		MessagesCollection:   MessagesCollection,
	}
}

//Repo implements the repository.UserRepository interface
// var _ ports.Repository = &DatabaseInfra{}

func (r *DatabaseInfra) GetProduct(amount int) (interface{}, error) {
	cursor, err := r.ProductCollection.Aggregate(context.TODO(), []bson.M{
		{"$sample": bson.M{"size": amount}},
	})

	if err != nil {
		helper.LogEvent("ERROR", err.Error())
		return nil, errors.New("something went wrong")
	}
	defer cursor.Close(context.TODO())

	var products []entity.Product

	if err := cursor.All(context.TODO(), &products); err != nil {
		helper.LogEvent("ERROR", "decoding products in method getproduct: "+err.Error())
		return nil, errors.New("something went wrong")
	}

	return products, nil
}

func (r *DatabaseInfra) CreateProduct(product entity.Product) (interface{}, error) {
	_, err := r.ProductCollection.InsertOne(context.TODO(), product)

	if err != nil {
		helper.LogEvent("ERROR", err.Error())
		return nil, errors.New("something went wrong")
	}

	helper.LogEvent("INFO", "succesfully added new product to database")

	return product.ProductID.Hex(), nil
}

func (r *DatabaseInfra) SubscribeToNewsLetter(body entity.Subscriber) error {
	if singleResult := r.NewsletterCollection.FindOne(context.TODO(), bson.M{"email": body.Email}); singleResult.Err() == nil {
		log.Println("email found")
		helper.LogEvent("ERROR", "user already subscribed")
		return helper.USER_ALREADY_A_SUBSCRIBER
	}

	_, err := r.NewsletterCollection.InsertOne(context.TODO(), body)
	if err != nil {
		helper.LogEvent("ERROR", "inserting subscription email into database:"+err.Error())
		return errors.New("something went wrong")
	}

	helper.LogEvent("INFO", "successfully inserted subscriber's email into database")

	if err := helper.SendMail("newsletter.html", entity.ContactMessage{To: body.Email}); err != nil {
		helper.LogEvent("ERROR", "sending newsletter confirmation mail to client:"+err.Error())
		return helper.NEWSLETTER_MAIL_ERROR
	}

	helper.LogEvent("INFO", "successfully sent confirmation email to new subscriber")

	return nil
}

func (r *DatabaseInfra) GetProductByRef(ref string) (interface{}, error) {
	id, err := primitive.ObjectIDFromHex(ref)
	if err != nil {
		helper.LogEvent("ERROR", "invalid product reference:"+err.Error())
		return nil, errors.New("invalid product reference")
	}

	filter := bson.M{"_id": id}

	singleResult := r.ProductCollection.FindOne(context.TODO(), filter)
	if singleResult.Err() != nil {
		switch singleResult.Err() {
		case mongo.ErrNoDocuments:
			helper.LogEvent("ERROR", "no product found:"+err.Error())
			return nil, errors.New("no product found")

		default:
			helper.LogEvent("ERROR", "finding product by ref:"+err.Error())
			return nil, errors.New("something went wrong")
		}
	}

	product := entity.Product{}

	if err := singleResult.Decode(&product); err != nil {
		helper.LogEvent("ERROR", "decoding product from singleresult:"+err.Error())
		return nil, errors.New("something went wrong")
	}

	return product, nil
}

func (r *DatabaseInfra) GetCartItems(ids []primitive.ObjectID) (interface{}, error) {
	matchStage := bson.M{
		"$match": bson.M{
			"_id": bson.M{
				"$in": ids,
			},
		},
	}

	projectStage := bson.M{
		"$project": bson.M{
			"brand":      0,
			"category":   0,
			"created_at": 0,
			"desc":       0,
			"rating":     0,
		},
	}

	cursor, err := r.ProductCollection.Aggregate(context.TODO(), []bson.M{matchStage, projectStage})
	if err != nil {
		helper.LogEvent("ERROR", "getting cart items: "+err.Error())
		return nil, errors.New("something went wrong")
	}
	defer cursor.Close(context.TODO())

	var products []entity.Product

	if err := cursor.All(context.TODO(), &products); err != nil {
		helper.LogEvent("ERROR", "decoding products in method getcartitems: "+err.Error())
		return nil, errors.New("something went wrong")
	}

	return products, nil
}

func (r *DatabaseInfra) CreateOrder(order entity.Order) (interface{}, error) {
	_, err := r.OrderCollection.InsertOne(context.TODO(), order)
	if err != nil {
		helper.LogEvent("ERROR", "creating new order: "+err.Error())
		return nil, errors.New("something went wrong")
	}

	return order.ID.Hex(), nil
}

func (r *DatabaseInfra) UpdateOrderPayment(id string) (interface{}, error) {
	idHex, _ := primitive.ObjectIDFromHex(id)
	_, err := r.OrderCollection.UpdateOne(context.TODO(), bson.M{"_id": idHex}, bson.M{"$set": bson.M{"paymentStatus": "PAID"}})
	if err != nil {
		helper.LogEvent("ERROR", fmt.Sprintf("updating order payment with id{%v} not successful: %v", id, err.Error()))
		return nil, errors.New("something went wrong")
	}

	return id, nil
}

func (r *DatabaseInfra) CreateContactMessage(body entity.ContactMessage) error {
	_, err := r.MessagesCollection.InsertOne(context.TODO(), body)
	if err != nil {
		helper.LogEvent("ERROR", "adding message to database: "+err.Error())
		return errors.New("something went wrong")
	}

	return nil
}

func (r *DatabaseInfra) GetOrders(page string) (interface{}, error) {
	findOptions, err := GetPage(page)
	if err != nil {
		helper.LogEvent("ERROR", map[string]interface{}{"find options": err.Error()})
		return nil, errors.New("invalid page number")
	}

	findOptions = findOptions.SetSort(bson.M{"cartSubtotal": -1}).SetProjection(bson.M{
		"cartItems":  0,
		"created_at": 0,
	})

	cursor, err := r.OrderCollection.Find(context.TODO(), bson.M{"deliveryStatus": "UNDELIVERED"}, findOptions)
	if err != nil {
		helper.LogEvent("ERROR", map[string]interface{}{"find": err.Error()})
		return nil, errors.New("something went wrong")
	}
	defer cursor.Close(context.TODO())

	var orders []entity.Order

	if err := cursor.All(context.TODO(), &orders); err != nil {
		helper.LogEvent("ERROR", map[string]interface{}{"cursor.all": err.Error()})
		return nil, errors.New("something went wrong")
	}

	return orders, nil
}

func (r *DatabaseInfra) GetDashBoardValues() (interface{}, error) {
	values := entity.DashboardValues{}
	totalOrders, err := r.OrderCollection.CountDocuments(context.TODO(), bson.M{})
	if err != nil {
		helper.LogEvent("ERROR", "could not retrieve total orders count from database: "+err.Error())
		return nil, errors.New("something went wrong")
	}
	values.TotalOrders = int(totalOrders)

	pendingOrders, err := r.OrderCollection.CountDocuments(context.TODO(), bson.M{"deliveryStatus": "UNDELIVERED"})
	if err != nil {
		helper.LogEvent("ERROR", "could not retrieve pending orders count from database: "+err.Error())
		return nil, errors.New("something went wrong")
	}
	values.PendingOrders = int(pendingOrders)

	values.CompletedOrders = int(totalOrders) - int(pendingOrders)

	matchStage := bson.M{"$match": bson.M{"paymentStatus": "PAID"}}
	groupStage := bson.M{"$group": bson.M{"_id": "", "total": bson.M{"$sum": "$cartSubtotal"}}}
	projectStage := bson.M{"$project": bson.M{"_id": 0, "total": "$total"}}
	cursor, err := r.OrderCollection.Aggregate(context.TODO(), []bson.M{matchStage, groupStage, projectStage})
	if err != nil {
		helper.LogEvent("ERROR", "could not retrieve total revenue for pending orders from database:"+err.Error())
		return nil, errors.New("something went wrong")
	}
	defer cursor.Close(context.TODO())

	var data []primitive.M
	if err := cursor.All(context.TODO(), &data); err != nil {
		helper.LogEvent("ERROR", "cursor all:"+err.Error())
		return nil, errors.New("something went wrong")
	}
	values.TotalRevenue = data[0]["total"].(float64)

	values.NetProfit = values.TotalRevenue - helper.Config.TotalExpense

	return values, nil
}

func (r *DatabaseInfra) GetAdminMsgs(page string) (interface{}, error) {
	findOptions, err := GetPage(page)
	if err != nil {
		helper.LogEvent("ERROR", map[string]interface{}{"find options": err.Error()})
		return nil, errors.New("invalid page number")
	}

	findOptions = findOptions.SetSort(bson.M{"sent_at": 1}).SetProjection(bson.M{
		"to":   0,
		"from": 0,
	})

	cursor, err := r.MessagesCollection.Find(context.TODO(), bson.M{}, findOptions)
	if err != nil {
		helper.LogEvent("ERROR", map[string]interface{}{"find": err.Error()})
		return nil, errors.New("something went wrong")
	}
	defer cursor.Close(context.TODO())

	var messages []entity.ContactMessage

	if err := cursor.All(context.TODO(), &messages); err != nil {
		helper.LogEvent("ERROR", map[string]interface{}{"cursor.all": err.Error()})
		return nil, errors.New("something went wrong")
	}

	return messages, nil
}

func (r *DatabaseInfra) GetOrderById(id string) (interface{}, error) {
	idHex, _ := primitive.ObjectIDFromHex(id)
	singleRes := r.OrderCollection.FindOne(context.TODO(), bson.M{"_id": idHex})
	if singleRes.Err() != nil {
		helper.LogEvent("ERROR", "could not get order with specified hex: "+singleRes.Err().Error())
		return nil, errors.New("something went wrong")
	}

	var order entity.Order
	if err := singleRes.Decode(&order); err != nil {
		helper.LogEvent("ERROR", "decoding from single order: "+singleRes.Err().Error())
		return nil, errors.New("something went wrong")
	}

	return order, nil
}

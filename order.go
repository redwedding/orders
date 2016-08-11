package main

import (
	"github.com/satori/go.uuid"
	"log"
	"time"
	"errors"
)

type Order struct {
	Id           string        `json:"id"`
	Number       string        `json:"number"`
	Reference    string        `json:"reference"`
	Status       string        `json:"status"`
	CreatedAt    time.Time     `json:"createdAt"`
	UpdatedAt    time.Time     `json:"updatedAt"`
	Notes        string        `json:"notes"`
	Price        int           `json:"price"`
	Items        []OrderItem   `json:"items"`
	Transactions []Transaction `json:"transactions"`
}

func (o Order) ValidadeNewOrder() error {
	if o.Number == "" {
		return errors.New("Order number cannot be empty")
	}

	if o.Reference == "" {
		return errors.New("Order Reference cannot be empty")
	}

	if o.Status == "" {
		return errors.New("Order Status cannot be empty")
	}

	if len(o.Items) > 0 {
		for i := range o.Items {
			return o.Items[i].ValidadeNewOrderItem();
		}
	}

	if len(o.Transactions) > 0 {
		for k := range o.Transactions {
			return o.Transactions[k].ValidateNewTransaction()
		}
	}
	return nil
}

func (order *Order) CalculatePrice() {
	for _, item := range order.Items {
		order.Price += item.UnitPrice * item.Quantity
	}
}

func (order *Order) Save() error {
	order.Id = uuid.NewV4().String()
	order.Status = "DRAFT"
	order.CreatedAt = time.Now()

	err := order.ValidadeNewOrder()
	if err != nil {
		log.Print(err)
		return err;
	}

	err = session.Query("INSERT INTO orders (id,number,reference,status,created_at) VALUES (?,?,?,?,?)",
		order.Id, order.Number, order.Reference, order.Status, order.CreatedAt).Exec()

	if err != nil {
		log.Print(err)
	}

	return err
}

func (order *Order) FindId(id string) error {
	return session.Query("SELECT id FROM orders WHERE id = ? ", id).Scan(&order.Id)
}

func (order *Order) GetOrder(id string) error {

	err := session.Query("SELECT id, number, reference, status, notes, price, created_at, updated_at, items, transactions from orders WHERE id = ? ", id).
		 Scan(&order.Id, &order.Number, &order.Reference, &order.Status, &order.Notes, &order.Price, &order.CreatedAt,
			&order.UpdatedAt, &order.Items, &order.Transactions)

	order.CalculatePrice()

	if err != nil {
		log.Print(err)
	}

	return err
}


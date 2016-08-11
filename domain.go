package main

import (
	"fmt"
	"github.com/satori/go.uuid"
	"log"
	"time"
	"errors"
)

type Config struct {
	Scyllaclusters []string
	Serverport     int
}

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

type OrderItem struct {
	Sku       string `cql:"sku" json:"sku"`
	UnitPrice int    `cql:"unit_price" json:"unit_price"`
	Quantity  int    `cql:"quantity" json:"quantity"`
}


func (oi OrderItem) ValidadeNewOrderItem() error {
	if oi.Quantity <= 0 {
		return errors.New("OrderItem Quantity cannot be less or equal to 0")
	}

	if oi.Sku == "" {
		return errors.New("OrderItem Sku cannot be empty")
	}

	if oi.UnitPrice < 0 {
		return  errors.New("OrderItem UnitPrice cannot be less than 0")
	}

	return nil;
}


type Transaction struct {
	Id                string `json:"id" cql:"id"`
	ExternalId        string `json:"external_id" cql:"external_id"`
	Amount            int    `json:"amount" cql:"amount"`
	Type              string `json:"type" cql:"type"`
	AuthorizationCode string `json:"authorization_code" cql:"authorization_code"`
	CardBrand         string `json:"card_brand" cql:"card_brand"`
	CardBin           string `json:"card_bin" cql:"card_bin"`
	CardLast          string `json:"card_last" cql:"card_last"`
	OrderId           string `json:"order_id" cql:"order_id"`
}

func (t Transaction) ValidateNewTransaction() error {

	if t.ExternalId == "" {
		return errors.New("Transaction ExternalId cannot be empty")
	}

	if t.Amount <= 0 {
		return  errors.New("Transaction Amount cannot be less or equal to 0")
	}

	if t.AuthorizationCode == "" {
		return errors.New("Transaction Authorization code cannot be empty")
	}

	if t.Type == "" {
		return errors.New("Transaction Type cannot be empty")
	}

	if t.CardBin == "" {
		return  errors.New("Transaction Card Bin cannot be empty")
	}

	if t.CardBrand == "" {
		return errors.New("Transaction Card Brand cannot be empty")
	}

	if t.CardLast == "" {
		return errors.New("Transaction Card Last cannot be empty")
	}

	if t.OrderId == "" {
		return errors.New("Transaction OrderId cannot be empty")
	}

	return nil;
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

	if err != nil {
		log.Print(err)
	}

	return err
}

func (item *OrderItem) Save(order_id string) error {


	err := item.ValidadeNewOrderItem()
	if err != nil {
		log.Print(err)
		return err;
	}

	query := fmt.Sprintf("UPDATE orders SET updated_at = ?, items = items + [{sku: %v, unit_price: %v, quantity: %v}] WHERE id = %v",
		item.Sku, item.UnitPrice, item.Quantity, order_id)

	log.Print(query)
	err = session.Query(query, time.Now()).Exec()

	if err != nil {
		log.Print(err)
	}

	return err
}

func (tran *Transaction) Save(order_id string) error {
	tran.Id = uuid.NewV4().String()
	tran.OrderId = order_id

	err := tran.ValidateNewTransaction()
	if err != nil {
		log.Print(err)
		return err;
	}
	query := fmt.Sprintf("UPDATE orders SET updated_at = ?, transactions = transactions + [{id: %v, external_id: '%v', amount: %v, type: '%v', authorization_code: '%v', card_brand: '%v', card_bin: '%v', card_last: '%v'}] WHERE id = %v",
		tran.Id, tran.ExternalId, tran.Amount, tran.Type, tran.AuthorizationCode, tran.CardBrand, tran.CardBin, tran.CardLast, tran.OrderId)

	err = session.Query(query, time.Now()).Exec()

	if err != nil {
		log.Print(err)
	}

	return err
}


package main

import (
	"fmt"
	"log"
	"time"
	"errors"
)

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


func (item *OrderItem) Save(order_id string) error {


	err := item.ValidadeNewOrderItem()
	if err != nil {
		log.Print(err)
		return err;
	}

	query := fmt.Sprintf("UPDATE orders SET updated_at = ?, items = items + [{sku: %v, unit_price: %v, quantity: %v}] WHERE id = %v",
		item.Sku, item.UnitPrice, item.Quantity, order_id)

	err = session.Query(query, time.Now()).Exec()

	if err != nil {
		log.Print(err)
	}

	return err
}

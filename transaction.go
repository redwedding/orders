package main

import (
	"fmt"
	"github.com/satori/go.uuid"
	"log"
	"time"
	"errors"
)

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


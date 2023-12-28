package models

import (
	"time"

	"github.com/aikintech/scim-api/pkg/constants"
	nanoid "github.com/matoous/go-nanoid/v2"
	"gorm.io/gorm"
)

type Transaction struct {
	ID             string `gorm:"primaryKey;size:40"`
	UserID         string `gorm:"size:40;index;not null"`
	ReferenceID    string `gorm:"size:40"`
	Provider       string `gorm:"size:20"`
	IdempotencyKey string `gorm:"size:40;index;not null"`
	Currency       string `gorm:"size:3;not null"`
	Amount         int64  `gorm:"not null"`
	Type           string `gorm:"size:40;not null"`
	Method         string `gorm:"size:40;not null"`
	Status         string `gorm:"size:40;not null;default:'pending'"`
	Processed      bool
	Description    string    `gorm:"size:255"`
	CreatedAt      time.Time `gorm:"autoCreateTime;not null"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime;not null"`

	User *User `gorm:"foreignKey:UserID"`
}

type TransactionResource struct {
	ID          string    `json:"id"`
	Provider    string    `json:"provider"`
	Currency    string    `json:"currency"`
	Amount      float64   `json:"amount"`
	Type        string    `json:"type"`
	Method      string    `json:"method"`
	Status      string    `json:"status"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	User        UserRel   `json:"user"`
}

func (t *Transaction) BeforeCreate(db *gorm.DB) (err error) {
	t.ID = nanoid.MustGenerate(constants.ALPHABETS, 30)

	return
}

func TransactionToResource(t *Transaction) TransactionResource {
	return TransactionResource{
		ID:          t.ID,
		Provider:    t.Provider,
		Currency:    t.Currency,
		Amount:      float64(t.Amount / 100),
		Type:        t.Type,
		Method:      t.Method,
		Status:      t.Status,
		Description: t.Description,
		CreatedAt:   t.CreatedAt,
		User:        ToUserRelResource(t.User),
	}
}

func TransactionsToResourceCollection(transactions []*Transaction) []TransactionResource {
	var resources []TransactionResource

	for _, transaction := range transactions {
		resources = append(resources, TransactionToResource(transaction))
	}

	return resources
}

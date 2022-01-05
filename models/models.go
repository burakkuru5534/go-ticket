package models

import (
	"go.mongodb.org/mongo-driver/x/mongo/driver/uuid"
	"time"

	"gopkg.in/guregu/null.v4"
)

//log table model
type ArInternalMetadata struct {
	Key       string      `json:"Key"`
	Value     null.String `json:"Value"`
	CreatedAt time.Time   `json:"CreatedAt"`
	UpdatedAt time.Time   `json:"UpdatedAt"`
}

//purchases model
type Purchases struct {
	ID             uuid.UUID `json:"ID"`
	Quantity       null.Int  `json:"Quantity"`
	UserID         string    `json:"UserID"`
	TicketOptionID uuid.UUID `json:"TicketOptionID"`
	CreatedAt      time.Time `json:"CreatedAt"`
	UpdatedAt      time.Time `json:"UpdatedAt"`
}

//migration schema model
type SchemaMigrations struct {
	Version string `json:"Version"`
}

//ticket options model
type TicketOptions struct {
	ID         string      `json:"ID"`
	Name       null.String `json:"Name"`
	Desc       null.String `json:"Desc"`
	Allocation int64       `json:"Allocation"`
	CreatedAt  time.Time   `json:"CreatedAt"`
	UpdatedAt  time.Time   `json:"UpdatedAt"`
}

type Tickets struct {
	ID             string    `json:"ID"`
	TicketOptionID uuid.UUID `json:"TicketOptionID"`
	PurchaseID     uuid.UUID `json:"PurchaseID"`
	CreatedAt      time.Time `json:"CreatedAt"`
	UpdatedAt      time.Time `json:"UpdatedAt"`
}

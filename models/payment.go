package models

import (
	"time"

	"malakashuttle/entities"

	"gorm.io/gorm"
)

// PaymentStatus enum
type PaymentStatus string

const (
	PaymentStatusPending PaymentStatus = "pending"
	PaymentStatusSuccess PaymentStatus = "success"
	PaymentStatusFailed  PaymentStatus = "failed"
)

// Payment model untuk business logic layer
type Payment struct {
	ID            uint          `json:"id"`
	BookingID     uint          `json:"booking_id"`
	PaymentMethod string        `json:"payment_method"`
	PaymentStatus PaymentStatus `json:"payment_status"`
	PaymentDate   *time.Time    `json:"payment_date"`
	ProofImageURL string        `json:"proof_image_url"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}

// PaymentHandler untuk converter operations
type PaymentHandler struct{}

// NewPaymentHandler membuat instance baru dari PaymentHandler
func NewPaymentHandler() *PaymentHandler {
	return &PaymentHandler{}
}

// FromEntity mengkonversi entity.Payment ke model.Payment
func (h *PaymentHandler) FromEntity(entity *entities.Payment) *Payment {
	if entity == nil {
		return nil
	}

	return &Payment{
		ID:            entity.ID,
		BookingID:     entity.BookingID,
		PaymentMethod: entity.PaymentMethod,
		PaymentStatus: PaymentStatus(entity.PaymentStatus),
		PaymentDate:   entity.PaymentDate,
		ProofImageURL: entity.ProofImageURL,
		CreatedAt:     entity.CreatedAt,
		UpdatedAt:     entity.UpdatedAt,
	}
}

// ToEntity mengkonversi model.Payment ke entity.Payment
func (h *PaymentHandler) ToEntity(model *Payment) *entities.Payment {
	if model == nil {
		return nil
	}

	return &entities.Payment{
		Model: gorm.Model{
			ID:        model.ID,
			CreatedAt: model.CreatedAt,
			UpdatedAt: model.UpdatedAt,
		},
		BookingID:     model.BookingID,
		PaymentMethod: model.PaymentMethod,
		PaymentStatus: entities.PaymentStatus(model.PaymentStatus),
		PaymentDate:   model.PaymentDate,
		ProofImageURL: model.ProofImageURL,
	}
}

// FromEntityList mengkonversi slice entity.Payment ke slice model.Payment
func (h *PaymentHandler) FromEntityList(entities []*entities.Payment) []*Payment {
	if entities == nil {
		return nil
	}

	models := make([]*Payment, len(entities))
	for i, entity := range entities {
		models[i] = h.FromEntity(entity)
	}
	return models
}

// Global handler instance
var PaymentHandlerInstance = NewPaymentHandler()

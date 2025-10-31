package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Colaborador struct {
	ID             uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	Nome           string         `gorm:"not null" json:"nome" binding:"required"`
	CPF            string         `gorm:"uniqueIndex;not null;size:11" json:"cpf" binding:"required"`
	RG             *string        `gorm:"uniqueIndex;size:20" json:"rg,omitempty"`
	DepartamentoID uuid.UUID      `gorm:"type:uuid;not null" json:"departamento_id" binding:"required"`
	Departamento   *Departamento  `gorm:"foreignKey:DepartamentoID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"departamento,omitempty"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

func (c *Colaborador) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.Must(uuid.NewV7())
	}
	return nil
}

func (c *Colaborador) TableName() string {
	return "colaboradores"
}

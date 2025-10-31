package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Departamento struct {
	ID                     uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	Nome                   string         `gorm:"not null" json:"nome" binding:"required"`
	GerenteID              uuid.UUID      `gorm:"type:uuid;not null" json:"gerente_id" binding:"required"`
	Gerente                *Colaborador   `gorm:"foreignKey:GerenteID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"gerente,omitempty"`
	DepartamentoSuperiorID *uuid.UUID     `gorm:"type:uuid" json:"departamento_superior_id,omitempty"`
	DepartamentoSuperior   *Departamento  `gorm:"foreignKey:DepartamentoSuperiorID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"departamento_superior,omitempty"`
	Subdepartamentos       []Departamento `gorm:"foreignKey:DepartamentoSuperiorID" json:"subdepartamentos,omitempty"`
	Colaboradores          []Colaborador  `gorm:"foreignKey:DepartamentoID" json:"colaboradores,omitempty"`
	CreatedAt              time.Time      `json:"created_at"`
	UpdatedAt              time.Time      `json:"updated_at"`
	DeletedAt              gorm.DeletedAt `gorm:"index" json:"-"`
}

func (d *Departamento) BeforeCreate(tx *gorm.DB) error {
	if d.ID == uuid.Nil {
		d.ID = uuid.Must(uuid.NewV7())
	}
	return nil
}

func (d *Departamento) TableName() string {
	return "departamentos"
}

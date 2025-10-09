package models

import (
	"time"
)

type NodoDetalle struct {
	Id                 string    `bson:"_id" json: "_id, omitempty" `
	Subgrupo_id        string    `json:"subgrupo_id" `
	Nombre             string    `json:"nombre" `
	Descripcion        string    `json:"descripcion" `
	Dato               string    `json:"dato" `
	Activo             bool      `json:"activo" `
	Fecha_creacion     time.Time `json:"fecha_creacion" `
	Fecha_modificacion time.Time `json:"fecha_modificacion" `
}

type Dato struct {
	Type     string `json: 'type' `
	Required string `json: 'required' `
}

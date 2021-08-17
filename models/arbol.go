package models

import "time"

type Nodo struct {
	Id                 string    `json: "_id" `
	Nombre             string    `json: "nombre" `
	Descripcion        string    `json: "descripcion" `
	Hijos              []string  `json: "hijos" `
	Padre              string    `json: "padre" `
	Activo             bool      `json: "activo" `
	Fecha_creacion     time.Time `json: "fecha_creacion" `
	Fecha_modificacion time.Time `json: "fecha_modificacion" `
}

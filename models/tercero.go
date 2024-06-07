package models

import (
	"time"
)

type Vinculacion struct {
	Id                     int
	TerceroPrincipalId     *Tercero
	TerceroRelacionadoId   *Tercero
	TipoVinculacionId      int
	CargoId                int
	DependenciaId          int
	Soporte                int
	PeriodoId              int
	FechaInicioVinculacion time.Time
	FechaFinVinculacion    time.Time
	Activo                 bool
	FechaCreacion          string
	FechaModificacion      string
	Alternancia            bool
}

type Tercero struct {
	Id                int
	NombreCompleto    string
	PrimerNombre      string
	SegundoNombre     string
	PrimerApellido    string
	SegundoApellido   string
	LugarOrigen       int
	Activo            bool
	FechaCreacion     string
	FechaModificacion string
	UsuarioWSO2       string
}

type DatosIdentificacion struct {
	Activo             bool        `json:"Activo"`
	CiudadExpedicion   int64       `json:"CiudadExpedicion"`
	DigitoVerificacion int64       `json:"DigitoVerificacion"`
	DocumentoSoporte   int64       `json:"DocumentoSoporte"`
	FechaCreacion      string      `json:"FechaCreacion"`
	FechaExpedicion    interface{} `json:"FechaExpedicion"`
	FechaModificacion  string      `json:"FechaModificacion"`
	ID                 int64       `json:"Id"`
	Numero             string      `json:"Numero"`
	TerceroID          TerceroID   `json:"TerceroId"`
}

type TerceroID struct {
	Activo              bool                `json:"Activo"`
	FechaCreacion       string              `json:"FechaCreacion"`
	FechaModificacion   string              `json:"FechaModificacion"`
	FechaNacimiento     time.Time           `json:"FechaNacimiento"`
	ID                  int64               `json:"Id"`
	LugarOrigen         int64               `json:"LugarOrigen"`
	NombreCompleto      string              `json:"NombreCompleto"`
	PrimerApellido      string              `json:"PrimerApellido"`
	PrimerNombre        string              `json:"PrimerNombre"`
	SegundoApellido     string              `json:"SegundoApellido"`
	SegundoNombre       string              `json:"SegundoNombre"`
	TipoContribuyenteID TipoContribuyenteID `json:"TipoContribuyenteId"`
	UsuarioWSO2         string              `json:"UsuarioWSO2"`
}

type TipoContribuyenteID struct {
	Activo            bool   `json:"Activo"`
	CodigoAbreviacion string `json:"CodigoAbreviacion"`
	Descripcion       string `json:"Descripcion"`
	FechaCreacion     string `json:"FechaCreacion"`
	FechaModificacion string `json:"FechaModificacion"`
	ID                int64  `json:"Id"`
	Nombre            string `json:"Nombre"`
}

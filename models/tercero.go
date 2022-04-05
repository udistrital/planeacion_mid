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

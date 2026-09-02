package arbolHelper

import (
	"testing"

	"github.com/udistrital/planeacion_mid/models"
)

func TestBuildTreeConservaOrdenYExponePosicion(t *testing.T) {
	hijos := []models.Nodo{
		{Nombre: "Actividad", Activo: true},
		{Nombre: "Fecha de inicio", Activo: true},
		{Nombre: "Fecha de finalización", Activo: true},
	}
	ids := []map[string]interface{}{
		{"_id": "actividad"},
		{"_id": "fecha-inicio"},
		{"_id": "fecha-fin"},
	}

	arbol := BuildTree(hijos, ids)
	for posicion, nodo := range arbol {
		if nodo["id"] != ids[posicion]["_id"] {
			t.Fatalf("posición %d: se alteró el orden recibido", posicion)
		}
		if nodo["orden"] != posicion {
			t.Fatalf("posición %d: orden calculado = %v", posicion, nodo["orden"])
		}
	}
}

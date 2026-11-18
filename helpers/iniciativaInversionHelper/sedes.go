package iniciativaInversionHelper

import (
	"encoding/json"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Las sedes de la Universidad (Localización > Sedes donde se ejecuta la alternativa) no están en
// elemento_catalogo: salen de Oikos, de los espacios físicos de tipo SEDE. Los demás catálogos
// se leen del crud. Para el validador y el cliente es un catálogo más (/catalogos/sedes).
const (
	// CatalogoSedes es el nombre del catálogo que se lee de Oikos.
	CatalogoSedes = "sedes"
	// CodigoOtraSede es la opción que se agrega al final: no está en Oikos y es la que hace
	// aparecer «¿Cuál otra sede?» (solo_si en la plantilla).
	CodigoOtraSede = "Otra"
	// tipoEspacioSede es el tipo de espacio físico SEDE en Oikos.
	tipoEspacioSede = 38
)

// FuenteSedes entrega las sedes activas de la Universidad.
type FuenteSedes interface {
	SedesActivas() ([]ItemCatalogo, error)
}

// RepositorioConSedes lee el catálogo de sedes de Oikos y los demás del repositorio base.
type RepositorioConSedes struct {
	base  RepositorioElementosCatalogo
	sedes FuenteSedes
}

// NuevoRepositorioConSedes arma el repositorio.
func NuevoRepositorioConSedes(base RepositorioElementosCatalogo, sedes FuenteSedes) *RepositorioConSedes {
	return &RepositorioConSedes{base: base, sedes: sedes}
}

// ElementosDeCatalogo implementa RepositorioElementosCatalogo. Las sedes van por nombre y al
// final «Otra». Lo que se guarda en la iniciativa es el Id de Oikos: si allá se desactiva una
// sede, la validación la reporta como «ya no está en la lista».
func (r *RepositorioConSedes) ElementosDeCatalogo(nombre string) ([]ItemCatalogo, error) {
	if nombre != CatalogoSedes {
		return r.base.ElementosDeCatalogo(nombre)
	}
	sedes, err := r.sedes.SedesActivas()
	if err != nil {
		return nil, err
	}
	return append(sedes, ItemCatalogo{Codigo: CodigoOtraSede, Nombre: CodigoOtraSede}), nil
}

// ClienteSedesOikos implementa FuenteSedes con HTTP.
type ClienteSedesOikos struct {
	oikos string
	http  *http.Client
}

// NuevoClienteSedesOikos crea el cliente. El host va sin esquema, como en conf/app.conf
// (OikosService).
func NuevoClienteSedesOikos(hostOikos string) *ClienteSedesOikos {
	return &ClienteSedesOikos{
		oikos: "http://" + strings.TrimSuffix(hostOikos, "/"),
		http:  &http.Client{Timeout: 30 * time.Second},
	}
}

// SedesActivas consulta los espacios físicos de tipo SEDE y deja las activas, ordenadas por
// nombre. El código es el Id: el Codigo de Oikos se repite entre sedes.
func (c *ClienteSedesOikos) SedesActivas() ([]ItemCatalogo, error) {
	var espacios []struct {
		ID     int    `json:"Id"`
		Nombre string `json:"Nombre"`
		Estado string `json:"Estado"`
	}
	ruta := c.oikos + "/espacio_fisico?query=" + url.QueryEscape("tipo_espacio_fisico_id:"+strconv.Itoa(tipoEspacioSede)) + "&limit=0"
	respuesta, err := c.http.Get(ruta)
	if err != nil {
		return nil, errorServicio("no se pudo conectar con Oikos: %v", err)
	}
	defer respuesta.Body.Close()
	if respuesta.StatusCode != http.StatusOK {
		return nil, errorServicio("Oikos respondió HTTP %d en %s", respuesta.StatusCode, ruta)
	}
	if err := json.NewDecoder(respuesta.Body).Decode(&espacios); err != nil {
		return nil, errorServicio("respuesta inesperada de Oikos en %s: %v", ruta, err)
	}

	sedes := make([]ItemCatalogo, 0, len(espacios))
	vistos := map[int]bool{}
	for _, espacio := range espacios {
		nombre := strings.Join(strings.Fields(espacio.Nombre), " ")
		if espacio.ID == 0 || nombre == "" || vistos[espacio.ID] || !strings.EqualFold(strings.TrimSpace(espacio.Estado), "Activo") {
			continue
		}
		vistos[espacio.ID] = true
		sedes = append(sedes, ItemCatalogo{Codigo: strconv.Itoa(espacio.ID), Nombre: nombre})
	}
	sort.SliceStable(sedes, func(i, j int) bool { return sedes[i].Nombre < sedes[j].Nombre })
	return sedes, nil
}

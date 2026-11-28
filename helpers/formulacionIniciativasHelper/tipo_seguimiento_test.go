package formulacionIniciativasHelper

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/astaxie/beego"
)

const (
	idTipoIniciativas  = "6b08ffbe542b3dbdf40538cd"
	idTipoFormulacion  = "61f236f525e40c582a0840d0"
	cuerpoFechasValido = `[{"Id":14,"Nombre":"FACULTAD DE INGENIERIA"}]`
)

// crudDePrueba responde el tipo de seguimiento FIPI_SP, ningún registro de fechas y acepta lo que
// se guarde. Devuelve las peticiones que recibió («MÉTODO ruta»).
func crudDePrueba(t *testing.T) (*[]string, func()) {
	t.Helper()
	var mutex sync.Mutex
	peticiones := []string{}
	servidor := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mutex.Lock()
		peticiones = append(peticiones, r.Method+" "+r.URL.Path)
		mutex.Unlock()
		datos := interface{}([]interface{}{})
		switch {
		case strings.HasPrefix(r.URL.Path, "/tipo-seguimiento"):
			datos = []interface{}{map[string]interface{}{"_id": idTipoIniciativas, "codigo_abreviacion": "FIPI_SP"}}
		case r.Method == http.MethodPost:
			datos = map[string]interface{}{"_id": "nuevo"}
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"Success": true, "Status": "200", "Message": "ok", "Data": datos})
	}))
	anterior := beego.AppConfig.String("PlanesService")
	if err := beego.AppConfig.Set("PlanesService", servidor.URL); err != nil {
		t.Fatalf("no se pudo configurar PlanesService: %v", err)
	}
	return &peticiones, func() {
		servidor.Close()
		_ = beego.AppConfig.Set("PlanesService", anterior)
	}
}

func cuerpoFechas(tipoSeguimiento string) map[string]interface{} {
	return map[string]interface{}{
		"tipo_seguimiento_id": tipoSeguimiento,
		"fecha_inicio":        "2026-11-01",
		"fecha_fin":           "2026-11-30",
		"unidades_interes":    cuerpoFechasValido,
	}
}

func TestDefinirFechasRechazaOtroTipoDeSeguimientoSinEscribir(t *testing.T) {
	peticiones, cerrar := crudDePrueba(t)
	defer cerrar()

	var falla interface{}
	func() {
		defer func() { falla = recover() }()
		DefinirFechasIniciativasFormulacion(cuerpoFechas(idTipoFormulacion))
	}()

	detalle, ok := falla.(map[string]interface{})
	if !ok || detalle["status"] != "400" || !strings.Contains(detalle["err"].(string), "FIPI_SP") {
		t.Fatalf("se esperaba un 400 por el tipo de seguimiento, llegó %v", falla)
	}
	for _, peticion := range *peticiones {
		if strings.Contains(peticion, "periodo-seguimiento") {
			t.Fatalf("con otro tipo de seguimiento no debe leer ni escribir fechas: %v", *peticiones)
		}
	}
}

func TestDefinirFechasConElTipoDeIniciativasGuarda(t *testing.T) {
	peticiones, cerrar := crudDePrueba(t)
	defer cerrar()

	escritos := DefinirFechasIniciativasFormulacion(cuerpoFechas(idTipoIniciativas))

	if len(escritos) != 1 {
		t.Fatalf("se esperaba un registro nuevo, llegaron %v", escritos)
	}
	if ultima := (*peticiones)[len(*peticiones)-1]; ultima != "POST /periodo-seguimiento-iniciativas" {
		t.Fatalf("se esperaba que terminara creando el registro de fechas, peticiones: %v", *peticiones)
	}
}

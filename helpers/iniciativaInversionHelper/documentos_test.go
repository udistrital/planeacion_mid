package iniciativaInversionHelper

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"
)

// gestorMemoria imita al gestor documental: guarda lo que se sube y responde Id 101, 102… y
// Enlace uuid-1, uuid-2…; con err, falla.
type gestorMemoria struct {
	subidos []DocumentoGestor
	err     error
}

func (g *gestorMemoria) Subir(documento DocumentoGestor) (*ArchivoGestor, error) {
	if g.err != nil {
		return nil, g.err
	}
	g.subidos = append(g.subidos, documento)
	return &ArchivoGestor{ID: 100 + len(g.subidos), Enlace: fmt.Sprintf("uuid-%d", len(g.subidos))}, nil
}

// definicionConAnexos es la plantilla de prueba más un módulo Anexos como el de la plantilla: un
// documento fijo y la lista de otros anexos (nombre y archivo), en PDF o Word de máximo 1 MB.
func definicionConAnexos(t *testing.T) *Definicion {
	t.Helper()
	unMB := 1.0
	documento := Dato{Requerido: true, Formatos: []string{"pdf", "docx"}, TamanoMaximoMB: &unMB}
	submodulo := &Nodo{Nombre: "Anexos", Dato: Dato{Clave: "anexos"}, Hijos: []*Nodo{
		campoDePrueba("carta_presentacion", TipoDocumento, documento),
		campoDePrueba("otros_anexos", TipoRepetible, Dato{},
			campoDePrueba("nombre", TipoTexto, Dato{Requerido: true}),
			campoDePrueba("archivo", TipoDocumento, documento)),
	}}
	definicion := definicionDePrueba(t)
	definicion.Modulos = append(definicion.Modulos, &Nodo{Nombre: "Anexos", Dato: Dato{Clave: "anexos", Modulo: true}, Hijos: []*Nodo{submodulo}})
	return definicion
}

func servicioConAnexos(t *testing.T) (*Servicio, *repoMemoria, *gestorMemoria, string) {
	t.Helper()
	repo := nuevoRepoMemoria()
	gestor := &gestorMemoria{}
	servicio := conPersonasDePrueba(NuevoServicio(repo, definicionFija{definicionConAnexos(t)}, catalogosDePrueba(), gestor))
	return servicio, repo, gestor, crearIniciativaDePrueba(t, servicio).ID
}

func archivoDePrueba(bytes int) string {
	return base64.StdEncoding.EncodeToString(make([]byte, bytes))
}

// bytesMinimos es el archivo más pequeño que recibe el gestor: su base64 tiene 1000 caracteres.
const bytesMinimos = 748

// La carga solo sube el archivo y responde la referencia: el submódulo no cambia (ni sus datos ni
// su versión) hasta que se guarda.
func TestSubirDocumentoRespondeLaReferenciaSinCambiarElSubmodulo(t *testing.T) {
	servicio, _, gestor, iniciativaID := servicioConAnexos(t)
	antes, err := servicio.ObtenerSubmodulo(iniciativaID, "anexos")
	if err != nil {
		t.Fatal(err)
	}

	referencia, err := servicio.SubirDocumento(iniciativaID, "anexos", SolicitudDocumento{
		Clave: "otros_anexos.archivo", Nombre: " Cotización.PDF ", Archivo: archivoDePrueba(bytesMinimos), Actor: formulador(),
	})
	if err != nil {
		t.Fatal(err)
	}

	if referencia.ID != 101 || referencia.Enlace != "uuid-1" || referencia.Nombre != "Cotización.PDF" ||
		referencia.Usuario != responsableDePrueba || referencia.UsuarioNombre != personasDePrueba[responsableDePrueba] {
		t.Fatalf("referencia = %+v", referencia)
	}
	if _, err := time.Parse(time.RFC3339, referencia.Fecha); err != nil {
		t.Fatalf("fecha %q: %v", referencia.Fecha, err)
	}
	if len(gestor.subidos) != 1 || gestor.subidos[0].Nombre != "Cotización.PDF" || gestor.subidos[0].Contenido != archivoDePrueba(bytesMinimos) ||
		gestor.subidos[0].Descripcion != "archivo de la iniciativa de inversión "+iniciativaID {
		t.Fatalf("subidos al gestor = %+v", gestor.subidos)
	}
	despues, err := servicio.ObtenerSubmodulo(iniciativaID, "anexos")
	if err != nil {
		t.Fatal(err)
	}
	if despues.Submodulo.Version != antes.Submodulo.Version || !reflect.DeepEqual(despues.Submodulo.Datos, antes.Submodulo.Datos) {
		t.Fatalf("el submódulo cambió: antes %+v, después %+v", antes.Submodulo, despues.Submodulo)
	}
}

func TestSubirDocumentoRechazaSinLlamarAlGestor(t *testing.T) {
	casos := []struct {
		nombre    string
		solicitud SolicitudDocumento
		status    int
		mensaje   string
	}{
		{"campo que no es documento", SolicitudDocumento{Clave: "otros_anexos.nombre", Nombre: "a.pdf", Archivo: archivoDePrueba(1), Actor: formulador()}, http.StatusBadRequest, mensajeSolicitudInvalida},
		{"campo que no existe", SolicitudDocumento{Clave: "presupuesto", Nombre: "a.pdf", Archivo: archivoDePrueba(1), Actor: formulador()}, http.StatusBadRequest, mensajeSolicitudInvalida},
		{"sin usuario", SolicitudDocumento{Clave: "carta_presentacion", Nombre: "a.pdf", Archivo: archivoDePrueba(1)}, http.StatusBadRequest, "No se pudo identificar su usuario o sus roles. Vuelva a iniciar sesión e intente de nuevo."},
		{"sin nombre", SolicitudDocumento{Actor: formulador(), Clave: "carta_presentacion", Archivo: archivoDePrueba(1)}, http.StatusBadRequest, mensajeSolicitudInvalida},
		{"archivo que no es base64", SolicitudDocumento{Clave: "carta_presentacion", Nombre: "a.pdf", Archivo: "no es base64", Actor: formulador()}, http.StatusBadRequest, mensajeSolicitudInvalida},
		{"formato no permitido", SolicitudDocumento{Clave: "carta_presentacion", Nombre: "foto.png", Archivo: archivoDePrueba(1), Actor: formulador()}, http.StatusBadRequest, "«foto.png» no es un formato permitido. Use: pdf, docx."},
		{"sin extensión", SolicitudDocumento{Clave: "carta_presentacion", Nombre: "carta", Archivo: archivoDePrueba(1), Actor: formulador()}, http.StatusBadRequest, "«carta» no es un formato permitido. Use: pdf, docx."},
		{"archivo vacío", SolicitudDocumento{Clave: "carta_presentacion", Nombre: "a.pdf", Archivo: "", Actor: formulador()}, http.StatusBadRequest, "«a.pdf» está vacío."},
		{"archivo más pequeño de lo que recibe el gestor", SolicitudDocumento{Clave: "carta_presentacion", Nombre: "corto.pdf", Archivo: archivoDePrueba(bytesMinimos - 1), Actor: formulador()}, http.StatusBadRequest, "«corto.pdf» es demasiado pequeño para guardarlo. Use un archivo con más contenido."},
		{"archivo más grande que el máximo", SolicitudDocumento{Clave: "otros_anexos.archivo", Nombre: "grande.docx", Archivo: archivoDePrueba(1024*1024 + 1), Actor: formulador()}, http.StatusBadRequest, "«grande.docx» pasa de 1 MB."},
	}
	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			servicio, _, gestor, iniciativaID := servicioConAnexos(t)
			_, err := servicio.SubirDocumento(iniciativaID, "anexos", caso.solicitud)

			negocio := ComoErrorNegocio(err)
			if err == nil || negocio.Status != caso.status || negocio.Mensaje != caso.mensaje {
				t.Fatalf("error = %v (status %d, mensaje %q); se esperaba %d %q", err, negocio.Status, negocio.Mensaje, caso.status, caso.mensaje)
			}
			if len(gestor.subidos) != 0 {
				t.Fatalf("se subió al gestor: %+v", gestor.subidos)
			}
		})
	}
}

func TestSubirDocumentoEnIniciativaQueNoSePuedeEditar(t *testing.T) {
	servicio, repo, gestor, iniciativaID := servicioConAnexos(t)
	repo.planes[iniciativaID].EstadoPlanID = "ep-er" // enviada a revisión

	_, err := servicio.SubirDocumento(iniciativaID, "anexos", SolicitudDocumento{Clave: "carta_presentacion", Nombre: "carta.pdf", Archivo: archivoDePrueba(1), Actor: formulador()})

	if negocio := ComoErrorNegocio(err); err == nil || negocio.Status != http.StatusConflict || !strings.Contains(negocio.Mensaje, "no se puede editar") {
		t.Fatalf("error = %v; se esperaba 409", err)
	}
	if len(gestor.subidos) != 0 {
		t.Fatalf("se subió al gestor: %+v", gestor.subidos)
	}
}

// En Por corregir se suben anexos solo a los submódulos que no están aprobados.
func TestSubirDocumentoEnUnSubmoduloAprobado(t *testing.T) {
	servicio, repo, gestor, iniciativaID := servicioConAnexos(t)
	repo.planes[iniciativaID].EstadoPlanID = "ep-pc"
	for _, documento := range repo.submodulos {
		if documento.PlanID == iniciativaID && documento.Submodulo == "anexos" {
			documento.EstadoSubmoduloID = "es-apr"
		}
	}

	_, err := servicio.SubirDocumento(iniciativaID, "anexos", SolicitudDocumento{Actor: formulador(), Clave: "carta_presentacion", Nombre: "carta.pdf", Archivo: archivoDePrueba(bytesMinimos)})

	if negocio := ComoErrorNegocio(err); err == nil || negocio.Status != http.StatusConflict || !strings.Contains(negocio.Mensaje, "aprobado") {
		t.Fatalf("error = %v; se esperaba 409", err)
	}
	if len(gestor.subidos) != 0 {
		t.Fatalf("se subió al gestor: %+v", gestor.subidos)
	}
}

func TestSubirDocumentoSinNombreSiTercerosNoResponde(t *testing.T) {
	servicio, _, gestor, iniciativaID := servicioConAnexos(t)
	tercerosCaidos(servicio)

	referencia, err := servicio.SubirDocumento(iniciativaID, "anexos", SolicitudDocumento{Clave: "carta_presentacion", Nombre: "carta.pdf", Archivo: archivoDePrueba(bytesMinimos), Actor: formulador()})
	if err != nil {
		t.Fatal(err)
	}
	if len(gestor.subidos) != 1 || referencia.Usuario != responsableDePrueba || referencia.UsuarioNombre != "" {
		t.Fatalf("referencia = %+v", referencia)
	}
}

func TestSubirDocumentoCuandoFallaElGestor(t *testing.T) {
	servicio, _, gestor, iniciativaID := servicioConAnexos(t)
	gestor.err = errorServicio("el gestor documental respondió 500")

	_, err := servicio.SubirDocumento(iniciativaID, "anexos", SolicitudDocumento{Clave: "carta_presentacion", Nombre: "carta.pdf", Archivo: archivoDePrueba(bytesMinimos), Actor: formulador()})

	if negocio := ComoErrorNegocio(err); err == nil || negocio.Status != http.StatusBadGateway {
		t.Fatalf("error = %v; se esperaba 502", err)
	}
}

// Al guardar, la referencia queda solo con sus claves y el submódulo se valida como los demás.
func TestGuardarAnexosConReferencias(t *testing.T) {
	servicio, _, _, iniciativaID := servicioConAnexos(t)
	carta := map[string]interface{}{
		"id": 101, "enlace": "uuid-1", "nombre": "carta.pdf", "usuario": "1020", "usuario_nombre": "LAURA GOMEZ RIOS",
		"fecha": "2026-11-20T10:00:00Z", "descripcion": "sobra", "nuxeo_id": 101,
	}

	vista := guardar(t, servicio, iniciativaID, "anexos", Datos{
		"carta_presentacion": carta,
		"otros_anexos": []interface{}{
			map[string]interface{}{"id": "oa-1", "nombre": " Cotización ", "archivo": map[string]interface{}{"id": 102, "enlace": "uuid-2", "nombre": "cotizacion.docx"}},
		},
	})

	exigirSinErrores(t, "anexos", vista.Errores)
	if vista.Estado.Codigo != EstadoSubmoduloFormulado {
		t.Fatalf("anexos quedó %s", vista.Estado.Codigo)
	}
	esperada := map[string]interface{}{
		"id": float64(101), "enlace": "uuid-1", "nombre": "carta.pdf", "usuario": "1020", "usuario_nombre": "LAURA GOMEZ RIOS", "fecha": "2026-11-20T10:00:00Z",
	}
	if guardada := comoJSONDePrueba(t, vista.Datos["carta_presentacion"]); guardada != comoJSONDePrueba(t, esperada) {
		t.Fatalf("carta guardada = %s", guardada)
	}
	fila := filaConID(t, filasEn(t, vista.Datos, "otros_anexos"), "oa-1")
	if fila["nombre"] != "Cotización" {
		t.Fatalf("fila = %+v", fila)
	}
}

func TestValidarDocumento(t *testing.T) {
	definicion := definicionConAnexos(t)
	submodulo := definicion.Submodulo("anexos")
	validador := NuevoValidador(catalogosDePrueba())
	validar := func(datos Datos) []ErrorCampo {
		return validador.ValidarSubmodulo(submodulo, map[string]Datos{"anexos": NuevoNormalizador(catalogosDePrueba()).Normalizar(submodulo, datos)})
	}

	errores := validar(Datos{"otros_anexos": []interface{}{map[string]interface{}{"id": "oa-1", "nombre": "Plano"}}})
	exigirError(t, errores, "carta_presentacion", "carta_presentacion: es obligatorio.")
	exigirError(t, errores, "otros_anexos[oa-1].archivo", "archivo: es obligatorio.")

	errores = validar(Datos{
		"carta_presentacion": map[string]interface{}{"id": 1, "nombre": "carta.pdf"},
		"otros_anexos":       []interface{}{map[string]interface{}{"id": "oa-1", "nombre": "Plano", "archivo": map[string]interface{}{"id": 2, "enlace": "uuid-2", "nombre": "plano.dwg"}}},
	})
	exigirError(t, errores, "carta_presentacion", "carta_presentacion: el archivo no quedó bien adjuntado. Adjúntelo de nuevo.")
	exigirError(t, errores, "otros_anexos[oa-1].archivo", "archivo: «plano.dwg» no es un formato permitido. Use: pdf, docx.")

	errores = validar(Datos{"carta_presentacion": "carta.pdf"})
	exigirError(t, errores, "carta_presentacion", "carta_presentacion: el archivo no quedó bien adjuntado.")
}

func TestClienteGestorDocumental(t *testing.T) {
	var recibido []map[string]interface{}
	servidor := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/document/uploadAnyFormat" {
			t.Errorf("solicitud inesperada: %s %s", r.Method, r.URL.Path)
		}
		if err := json.NewDecoder(r.Body).Decode(&recibido); err != nil {
			t.Errorf("cuerpo inválido: %v", err)
		}
		switch recibido[0]["nombre"] {
		case "rechazado-en-el-sobre.pdf":
			json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400", "Error": "tipo de documento no válido"})
		case "rechazado.txt":
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"Success": false, "Status": "400", "Message": "El archivo no corresponde a su extensión"})
		case "gestor-caido.pdf":
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("<html>error</html>"))
		default:
			json.NewEncoder(w).Encode(map[string]interface{}{"Status": "200", "res": map[string]interface{}{"Id": 7, "Enlace": "uuid-7", "Nombre": "carta.pdf"}})
		}
	}))
	defer servidor.Close()
	cliente := NuevoClienteGestorDocumental(servidor.URL + "/")

	archivo, err := cliente.Subir(DocumentoGestor{Nombre: "carta.pdf", Descripcion: "Carta", Contenido: "JVBERi0="})
	if err != nil {
		t.Fatal(err)
	}
	if *archivo != (ArchivoGestor{ID: 7, Enlace: "uuid-7"}) {
		t.Fatalf("archivo = %+v", archivo)
	}
	esperado := map[string]interface{}{"IdTipoDocumento": float64(tipoDocumentoGestor), "nombre": "carta.pdf", "descripcion": "Carta", "metadatos": map[string]interface{}{}, "file": "JVBERi0="}
	if len(recibido) != 1 || !reflect.DeepEqual(recibido[0], esperado) {
		t.Fatalf("cuerpo recibido = %+v", recibido)
	}

	// Si el gestor no acepta el archivo, el usuario lo sabe y el log lleva la respuesta del gestor, venga la causa
	// en el campo que venga; si el gestor falla, el mensaje es el general de servicio.
	casos := []struct{ nombre, mensaje, enElDetalle string }{
		{"rechazado-en-el-sobre.pdf", mensajeArchivoRechazado, "tipo de documento no válido"},
		{"rechazado.txt", mensajeArchivoRechazado, "El archivo no corresponde a su extensión"},
		{"gestor-caido.pdf", mensajeServicio, "<html>error</html>"},
	}
	for _, caso := range casos {
		_, err := cliente.Subir(DocumentoGestor{Nombre: caso.nombre, Contenido: "JVBERi0="})
		negocio := ComoErrorNegocio(err)
		if negocio.Status != http.StatusBadGateway || negocio.Mensaje != caso.mensaje || !strings.Contains(negocio.Detalle, caso.enElDetalle) {
			t.Errorf("%s: error = %+v", caso.nombre, negocio)
		}
	}
}

func TestRecortarDejaLaRespuestaEnUnTamanoParaElLog(t *testing.T) {
	if got := recortar([]byte("  corto  ")); got != "corto" {
		t.Fatalf("recortar = %q", got)
	}
	largo := recortar([]byte(strings.Repeat("á", 600)))
	if len([]rune(largo)) != 501 || !strings.HasSuffix(largo, "…") {
		t.Fatalf("recortar dejó %d caracteres", len([]rune(largo)))
	}
}

func comoJSONDePrueba(t *testing.T, valor interface{}) string {
	t.Helper()
	texto, err := json.Marshal(valor)
	if err != nil {
		t.Fatal(err)
	}
	return string(texto)
}

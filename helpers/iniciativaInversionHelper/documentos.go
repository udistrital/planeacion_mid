package iniciativaInversionHelper

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"
)

// Los campos de tipo documento (Anexos) guardan la referencia de un archivo del gestor documental
// de OAS. El archivo se sube con SubirDocumento, que solo responde la referencia: esta queda en la
// iniciativa cuando se guarda el submódulo, como cualquier otro valor, y quitarla o eliminar la fila
// de un adicional también se aplica al guardar. El archivo sigue en el gestor.

// tipoDocumentoGestor es el tipo de documento con el que se suben los anexos: el 66, «Soportes
// planeación», el mismo de los proyectos de inversión (MACRO).
const tipoDocumentoGestor = 66

// minimoBase64Gestor es el menor base64 que recibe el gestor documental: con menos de 1000
// caracteres (archivos de menos de 748 bytes) responde «invalid file» (postAny, en
// gestor_documental_mid/controllers/document.py). Se revisa antes de llamarlo para avisar la causa.
const minimoBase64Gestor = 1000

// camposReferencia son las claves de una referencia; el normalizador descarta las demás.
var camposReferencia = []string{"id", "enlace", "nombre", "usuario", "usuario_nombre", "fecha"}

// ReferenciaDocumento es el valor de un campo documento: el archivo en el gestor (Id y Enlace, con
// el que el cliente lo vuelve a pedir), su nombre, quién lo subió y cuándo. Como en las
// observaciones, Usuario es el documento y UsuarioNombre su nombre en Terceros (vacío si no se obtuvo).
type ReferenciaDocumento struct {
	ID            int    `json:"id"`
	Enlace        string `json:"enlace"`
	Nombre        string `json:"nombre"`
	Usuario       string `json:"usuario"`
	UsuarioNombre string `json:"usuario_nombre,omitempty"`
	Fecha         string `json:"fecha"`
}

// SolicitudDocumento es el cuerpo de POST /iniciativas/:id/submodulos/:clave/documentos. Clave es
// la ruta del campo en el submódulo sin ids de fila ("carta_presentacion" u "otros_anexos.archivo");
// Archivo, el contenido en base64.
type SolicitudDocumento struct {
	Clave   string `json:"clave"`
	Nombre  string `json:"nombre"`
	Archivo string `json:"archivo"`
	Usuario string `json:"usuario"`
}

// DocumentoGestor es un archivo para subir al gestor documental (Contenido en base64).
type DocumentoGestor struct {
	Nombre      string
	Descripcion string
	Contenido   string
}

// ArchivoGestor identifica un archivo ya subido al gestor documental.
type ArchivoGestor struct {
	ID     int
	Enlace string
}

// GestorDocumental sube archivos al gestor documental de OAS.
type GestorDocumental interface {
	Subir(documento DocumentoGestor) (*ArchivoGestor, error)
}

// SubirDocumento revisa que la iniciativa se pueda editar y esté dentro de las fechas de
// formulación de su unidad, que clave sea un campo documento del submódulo y que el archivo cumpla
// sus formatos y su tamaño máximo; lo sube al gestor y responde la referencia, con el nombre en
// Terceros de quien lo sube. Si no se obtiene el nombre, la referencia va sin él y el motivo queda
// en el log. No cambia el submódulo.
func (s *Servicio) SubirDocumento(iniciativaID, claveSubmodulo string, solicitud SolicitudDocumento) (*ReferenciaDocumento, error) {
	nombre := strings.TrimSpace(solicitud.Nombre)
	usuario := strings.TrimSpace(solicitud.Usuario)
	if nombre == "" || usuario == "" {
		return nil, errorFormatoSolicitud("la carga de un documento necesita nombre y usuario (clave %q)", solicitud.Clave)
	}
	ctx, err := s.cargarIniciativa(iniciativaID)
	if err != nil {
		return nil, err
	}
	if !PermiteEditar(ctx.estado.Codigo) {
		return nil, errorConflicto("La iniciativa está «%s»; no se puede editar.", ctx.estado.Nombre)
	}
	if err := s.exigirFechasHabilitadas(ctx.plan.DependenciaID); err != nil {
		return nil, err
	}
	submodulo, _, err := ctx.submodulo(claveSubmodulo)
	if err != nil {
		return nil, err
	}
	campo := campoEnRuta(submodulo, solicitud.Clave)
	if campo == nil || campo.Dato.Tipo != TipoDocumento {
		return nil, errorFormatoSolicitud("%q no es un campo documento del submódulo %s", solicitud.Clave, claveSubmodulo)
	}
	if mensaje := formatoNoPermitido(campo, nombre); mensaje != "" {
		return nil, errorSolicitud("%s", mensaje)
	}
	contenido, err := base64.StdEncoding.DecodeString(solicitud.Archivo)
	if err != nil {
		return nil, errorFormatoSolicitud("el archivo de %q no está en base64: %v", solicitud.Clave, err)
	}
	if len(contenido) == 0 {
		return nil, errorSolicitud("«%s» está vacío.", nombre)
	}
	if len(solicitud.Archivo) < minimoBase64Gestor {
		return nil, errorSolicitud("«%s» es demasiado pequeño para guardarlo. Use un archivo con más contenido.", nombre)
	}
	if maximo := campo.Dato.TamanoMaximoMB; maximo != nil && float64(len(contenido)) > *maximo*1024*1024 {
		return nil, errorSolicitud("«%s» pasa de %v MB.", nombre, *maximo)
	}
	archivo, err := s.gestor.Subir(DocumentoGestor{
		Nombre:      nombre,
		Descripcion: fmt.Sprintf("%s de la iniciativa de inversión %s", campo.Nombre, iniciativaID),
		Contenido:   solicitud.Archivo,
	})
	if err != nil {
		return nil, err
	}
	usuarioNombre, err := s.nombreTercero(usuario)
	if err != nil {
		logs.Warn("iniciativaInversion: el documento %q de la iniciativa %s se sube sin el nombre de quien lo sube: %v", solicitud.Clave, iniciativaID, err)
	}
	return &ReferenciaDocumento{
		ID:            archivo.ID,
		Enlace:        archivo.Enlace,
		Nombre:        nombre,
		Usuario:       usuario,
		UsuarioNombre: usuarioNombre,
		Fecha:         time.Now().UTC().Format(time.RFC3339),
	}, nil
}

// campoEnRuta busca el campo de la ruta "grupo.campo" (sin ids de fila) entre los hijos del submódulo.
func campoEnRuta(submodulo *Nodo, ruta string) *Nodo {
	if strings.TrimSpace(ruta) == "" {
		return nil
	}
	nodo := submodulo
	for _, clave := range strings.Split(ruta, ".") {
		if nodo = nodo.Hijo(clave); nodo == nil {
			return nil
		}
	}
	return nodo
}

// formatoNoPermitido devuelve el mensaje para el usuario si la extensión del nombre no está en los
// formatos del campo; vacío si se permite (o si el campo no define formatos).
func formatoNoPermitido(campo *Nodo, nombre string) string {
	formatos := campo.Dato.Formatos
	if len(formatos) == 0 {
		return ""
	}
	extension := extensionDe(nombre)
	for _, formato := range formatos {
		if strings.EqualFold(formato, extension) {
			return ""
		}
	}
	return fmt.Sprintf("«%s» no es un formato permitido. Use: %s.", nombre, strings.Join(formatos, ", "))
}

// extensionDe devuelve la extensión del nombre en minúscula, sin el punto; vacía si no tiene.
func extensionDe(nombre string) string {
	punto := strings.LastIndex(nombre, ".")
	if punto < 0 {
		return ""
	}
	return strings.ToLower(nombre[punto+1:])
}

// referenciaLimpia deja solo las claves de camposReferencia, con los textos recortados. Lo que no es
// un objeto se conserva para que el validador lo reporte.
func referenciaLimpia(valor interface{}) interface{} {
	referencia, esObjeto := comoObjeto(valor)
	if !esObjeto {
		return valor
	}
	limpia := Datos{}
	for _, clave := range camposReferencia {
		if dato, existe := referencia[clave]; existe && dato != nil {
			limpia[clave] = limpiarEscalar(dato)
		}
	}
	return limpia
}

// ClienteGestorDocumental implementa GestorDocumental con HTTP: POST /document/uploadAnyFormat, que
// recibe cualquier formato (/document/upload solo recibe PDF).
type ClienteGestorDocumental struct {
	base string
	http *http.Client
}

// NuevoClienteGestorDocumental crea el cliente. El host incluye el esquema, como en conf/app.conf
// (GestorDocumental).
func NuevoClienteGestorDocumental(host string) *ClienteGestorDocumental {
	return &ClienteGestorDocumental{
		base: strings.TrimSuffix(host, "/"),
		http: &http.Client{Timeout: 60 * time.Second},
	}
}

// Subir implementa GestorDocumental. Si sube el archivo, el gestor responde { Status, res: { Id, Enlace, … } };
// si no, el detalle para el log lleva su respuesta completa (recortada), porque no siempre trae la causa en el mismo campo.
func (c *ClienteGestorDocumental) Subir(documento DocumentoGestor) (*ArchivoGestor, error) {
	cuerpo, err := json.Marshal([]map[string]interface{}{{
		"IdTipoDocumento": tipoDocumentoGestor,
		"nombre":          documento.Nombre,
		"descripcion":     documento.Descripcion,
		"metadatos":       map[string]interface{}{},
		"file":            documento.Contenido,
	}})
	if err != nil {
		return nil, fmt.Errorf("no se pudo serializar el documento %q: %w", documento.Nombre, err)
	}
	ruta := c.base + "/document/uploadAnyFormat"
	respuesta, err := c.http.Post(ruta, "application/json", bytes.NewReader(cuerpo))
	if err != nil {
		return nil, errorServicio("no se pudo conectar con el gestor documental (%s): %v", ruta, err)
	}
	defer respuesta.Body.Close()
	contenido, err := io.ReadAll(io.LimitReader(respuesta.Body, 1<<20))
	if err != nil {
		return nil, errorServicio("no se pudo leer la respuesta del gestor documental (%s): %v", ruta, err)
	}

	var sobre struct {
		Status interface{} `json:"Status"`
		Res    *struct {
			ID     int    `json:"Id"`
			Enlace string `json:"Enlace"`
		} `json:"res"`
	}
	if err := json.Unmarshal(contenido, &sobre); err != nil {
		return nil, errorServicio("el gestor documental respondió HTTP %d sin JSON válido (%s): %s", respuesta.StatusCode, ruta, recortar(contenido))
	}
	status := codigoReal(respuesta.StatusCode, sobre.Status)
	detalle := fmt.Sprintf("el gestor documental respondió %d al subir %q (%s): %s", status, documento.Nombre, ruta, recortar(contenido))
	switch {
	case status >= 400 && status < 500:
		// Rechazó el archivo (su contenido, su formato o el tipo de documento): el usuario puede probar con otro.
		return nil, nuevoError(http.StatusBadGateway, "%s", mensajeArchivoRechazado).conDetalle("%s", detalle)
	case status >= 400 || sobre.Res == nil || sobre.Res.Enlace == "":
		return nil, errorServicio("%s", detalle)
	}
	return &ArchivoGestor{ID: sobre.Res.ID, Enlace: sobre.Res.Enlace}, nil
}

// mensajeArchivoRechazado es lo que ve el usuario cuando el gestor documental no acepta el archivo.
const mensajeArchivoRechazado = "No se aceptó el archivo. Revise que su contenido corresponda a su formato o pruebe con otro archivo."

// recortar deja la respuesta de un servicio en un tamaño razonable para el log.
func recortar(contenido []byte) string {
	const maximo = 500
	texto := strings.TrimSpace(string(contenido))
	if len([]rune(texto)) <= maximo {
		return texto
	}
	return string([]rune(texto)[:maximo]) + "…"
}

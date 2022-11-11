package models

import (
	"github.com/astaxie/beego"
)

func RegistrarDoc(documento []map[string]interface{}) (status interface{}, outputError interface{}) {

	var resultadoRegistro map[string]interface{}

	errRegDoc := SendJson("http://"+beego.AppConfig.String("GestorDocumental")+"document/upload", "POST", &resultadoRegistro, documento)

	if resultadoRegistro["Status"].(string) == "200" && errRegDoc == nil {
		return resultadoRegistro["res"], nil
	} else {
		return nil, resultadoRegistro["Error"].(string)
	}

}

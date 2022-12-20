package models

import (
	"fmt"

	"github.com/astaxie/beego"
)

func RegistrarDoc(documento []map[string]interface{}) (status interface{}, outputError interface{}) {

	var resultadoRegistro map[string]interface{}
	fmt.Println(">>>>>>URL:", "http://"+beego.AppConfig.String("GestorDocumental")+"document/upload")
	fmt.Println(">>>>>>Data:", documento)
	errRegDoc := SendJson("http://"+beego.AppConfig.String("GestorDocumental")+"document/upload", "POST", &resultadoRegistro, documento)
	fmt.Println(">>>>>>Res:", resultadoRegistro)

	if resultadoRegistro["Status"].(string) == "200" && errRegDoc == nil {
		return resultadoRegistro["res"], nil
	} else {
		return nil, resultadoRegistro["Error"].(string)
	}

}

package formulacionhelper

// SinTipoPlan devuelve los planes que no son del tipo tipoPlanID, en el mismo orden. Con tipoPlanID
// vacío (p. ej. el tipo no existe en el ambiente) no hay nada que quitar y los devuelve todos.
func SinTipoPlan(planes []map[string]interface{}, tipoPlanID string) []map[string]interface{} {
	if tipoPlanID == "" {
		return planes
	}
	filtrados := make([]map[string]interface{}, 0, len(planes))
	for _, plan := range planes {
		if id, _ := plan["tipo_plan_id"].(string); id == tipoPlanID {
			continue
		}
		filtrados = append(filtrados, plan)
	}
	return filtrados
}

// IDPrimerTipoPlan devuelve el _id del primer tipo de plan de una consulta por código, o vacío si
// la consulta no trajo ninguno.
func IDPrimerTipoPlan(tipos []map[string]interface{}) string {
	if len(tipos) == 0 {
		return ""
	}
	id, _ := tipos[0]["_id"].(string)
	return id
}

package models

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

var global *context.Context

func SendJson(urlp string, trequest string, target interface{}, datajson interface{}) error {
	b := new(bytes.Buffer)
	if datajson != nil {
		if err := json.NewEncoder(b).Encode(datajson); err != nil {
			beego.Warn("Error codificando JSON en SendJson:", err)
		}
	}

	client := &http.Client{}
	req, err := http.NewRequest(trequest, urlp, b)
	if err != nil {
		beego.Error("Error creando request:", err)
		return err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	defer func() {
		if r := recover(); r != nil {
			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				beego.Error("Error leyendo respuesta en defer:", err)
				return
			}
			defer resp.Body.Close()

			if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
				beego.Warn("Error decodificando respuesta en defer:", err)
			}
		}
	}()

	header := GetHeader().Request.Header
	req.Header.Set("Authorization", header["Authorization"][0])

	resp, err := client.Do(req)
	if err != nil {
		beego.Error("Error leyendo respuesta:", err)
		return err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		beego.Warn("Error decodificando respuesta JSON:", err)
		return err
	}
	return nil
}

func SetHeader(ctx *context.Context) {
	global = ctx
}

func GetHeader() (ctx *context.Context) {
	return global
}

func GetJson(urlp string, target interface{}) error {
	req, err := http.NewRequest("GET", urlp, nil)
	if err != nil {
		beego.Error("Error creando request:", err)
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				beego.Error("Error leyendo respuesta en defer GetJson:", err)
				return
			}
			defer resp.Body.Close()

			if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
				beego.Warn("Error decodificando respuesta en defer GetJson:", err)
			}
		}
	}()

	header := GetHeader().Request.Header
	req.Header.Set("Authorization", header["Authorization"][0])

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		beego.Error("Error ejecutando request en GetJson:", err)
		return err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		beego.Warn("Error decodificando respuesta JSON en GetJson:", err)
		return err
	}
	return nil
}

func GetJsonTest(url string, target interface{}) (response *http.Response, err error) {
	client := &http.Client{Timeout: 10 * time.Second}
	response, err = client.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if err := json.NewDecoder(response.Body).Decode(target); err != nil {
		beego.Warn("Error decodificando respuesta JSON en GetJsonTest:", err)
		return response, err
	}
	return response, nil
}

func diff(a, b time.Time) (year, month, day int) {
	if a.Location() != b.Location() {
		b = b.In(a.Location())
	}
	if a.After(b) {
		a, b = b, a
	}
	y1, M1, d1 := a.Date()
	y2, M2, d2 := b.Date()

	year = int(y2 - y1)
	month = int(M2 - M1)
	day = int(d2 - d1)

	if day < 0 {
		t := time.Date(y1, M1, 32, 0, 0, 0, 0, time.UTC)
		day += 32 - t.Day()
		month--
	}
	if month < 0 {
		month += 12
		year--
	}
	return
}

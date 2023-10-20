package gcfbackend1214005

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/whatsauth/watoken"
)

func GCFHandler(MONGOCONNSTRINGENV, dbname, collectionname string) string {
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	datagedung := GetAllBangunanLineString(mconn, collectionname)
	return GCFReturnStruct(datagedung)
}

func GCFPostCoordinateLonLat(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	req := new(Credential)
	conn := SetConnection(MONGOCONNSTRINGENV, dbname)
	resp := new(CoorLonLatProperties)
	err := json.NewDecoder(r.Body).Decode(&resp)
	if err != nil {
		req.Status = false
		req.Message = "error parsing application/json: " + err.Error()
	} else {
		req.Status = true
		Ins := InsertDataLonlat(conn, collectionname,
			resp.Coordinates,
			resp.Name,
			resp.Volume,
			resp.Type)
		req.Message = fmt.Sprintf("%v:%v", "Berhasil Input data", Ins)
	}
	return GCFReturnStruct(req)
}

func GCFPostHandler(PASETOPRIVATEKEYENV, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	var Response Credential
	Response.Status = false
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	var datauser User
	err := json.NewDecoder(r.Body).Decode(&datauser)
	if err != nil {
		Response.Message = "error parsing application/json: " + err.Error()
	} else {
		if IsPasswordValid(mconn, collectionname, datauser) {
			Response.Status = true
			tokenstring, err := watoken.Encode(datauser.Username, os.Getenv(PASETOPRIVATEKEYENV))
			if err != nil {
				Response.Message = "Gagal Encode Token : " + err.Error()
			} else {
				Response.Message = "Selamat Datang"
				Response.Token = tokenstring
			}
		} else {
			Response.Message = "Password Salah"
		}
	}

	return GCFReturnStruct(Response)
}

func GCFReturnStruct(DataStuct any) string {
	jsondata, _ := json.Marshal(DataStuct)
	return string(jsondata)
}

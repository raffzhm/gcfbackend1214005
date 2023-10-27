package gcfbackend1214005

import (
	pasproj "github.com/e-dumas-sukasari/webpasetobackend"
	"github.com/petapedia/peda"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"go.mongodb.org/mongo-driver/mongo"

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

func SignUpGCF(PASETOPRIVATEKEYENV, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	resp := new(pasproj.Credential)
	userdata := new(User)
	resp.Status = false
	conn := SetConnection(MONGOCONNSTRINGENV, dbname)
	err := json.NewDecoder(r.Body).Decode(&userdata)
	if err != nil {
		resp.Message = "error parsing application/json: " + err.Error()
	} else {
		resp.Status = true
		hash, err := pasproj.HashPass(userdata.Password)
		if err != nil {
			resp.Message = "Gagal Hash Password" + err.Error()
		}
		InsertUserdata(conn, userdata.Username, hash)
		resp.Message = "Berhasil Input data"
	}
	response := pasproj.ReturnStringStruct(resp)
	return response
}

func SignInGCF(PASETOPRIVATEKEYENV, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	var resp pasproj.Credential
	mconn := pasproj.MongoCreateConnection(MONGOCONNSTRINGENV, dbname)
	var datauser peda.User
	err := json.NewDecoder(r.Body).Decode(&datauser)
	if err != nil {
		resp.Message = "error parsing application/json: " + err.Error()
	} else {
		if peda.IsPasswordValid(mconn, collectionname, datauser) {
			tokenstring, err := watoken.Encode(datauser.Username, os.Getenv(PASETOPRIVATEKEYENV))
			if err != nil {
				resp.Message = "Gagal Encode Token : " + err.Error()
			} else {
				resp.Status = true
				resp.Message = "Selamat Datang"
				resp.Token = tokenstring
			}
		} else {
			resp.Message = "Password Salah"
		}
	}
	return pasproj.ReturnStringStruct(resp)
}

func GCFReturnStruct(DataStuct any) string {
	jsondata, _ := json.Marshal(DataStuct)
	return string(jsondata)
}

func InsertUserdata(MongoConn *mongo.Database, username, password string) (InsertedID interface{}) {
	req := new(User)
	req.Username = username
	req.Password = password
	return pasproj.InsertOneDoc(MongoConn, "user", req)
}
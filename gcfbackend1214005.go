package gcfbackend1214005

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	pasproj "github.com/e-dumas-sukasari/webpasetobackend"
	"github.com/petapedia/peda"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/whatsauth/watoken"
)

func GCFHandler(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	resp := new(pasproj.Credential)
	tokenLogin := r.Header.Get("Login")

	if tokenLogin == "" {
		resp.Status = false
		resp.Message = "Header Login Not Exist"
		return pasproj.ReturnStringStruct(resp)
	}

	// Validate the token using your existing logic
	existing := IsExist(tokenLogin, os.Getenv(MONGOCONNSTRINGENV))

	if !existing {
		resp.Status = false
		resp.Message = "Kamu belum memiliki akun"
		return pasproj.ReturnStringStruct(resp)
	}

	// Proceed with your existing logic to get data
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	datagedung := GetAllBangunanLineString(mconn, collectionname)

	return GCFReturnStruct(datagedung)
}


func GCFPostCoordinateLonLat(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	req := new(Credential)
	tokenLogin := r.Header.Get("Login")

	if tokenLogin == "" {
		req.Status = false
		req.Message = "Header Login Not Exist"
		return GCFReturnStruct(req)
	}

	// Validate the token using your existing logic
	existing := IsExist(tokenLogin, os.Getenv(MONGOCONNSTRINGENV))

	if !existing {
		req.Status = false
		req.Message = "Kamu belum memiliki akun"
		return GCFReturnStruct(req)
	}

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

func SignUpGCF(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	resp := new(pasproj.Credential)
	userdata := new(RegistUser)
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
			tokenstring, err := watoken.Encode(datauser.Username, os.Getenv(MONGOCONNSTRINGENV))
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

func GCFUpdateGeo(MONGOCONNSTRINGENV, Mongostring, dbname, colname string, r *http.Request) string {
    req := new(Credential)
    resp := new(CoorLonLatProperties)

    tokenLogin := r.Header.Get("Login")
    if tokenLogin == "" {
        req.Status = false
        req.Message = "Header Login Not Exist"
        return GCFReturnStruct(req)
    }

    // Validate the token using your existing logic
    existing := IsExist(tokenLogin, os.Getenv(MONGOCONNSTRINGENV))

    if !existing {
        req.Status = false
        req.Message = "Kamu belum memiliki akun"
        return GCFReturnStruct(req)
    }

    conn := SetConnection(Mongostring, dbname)
    err := json.NewDecoder(r.Body).Decode(&resp)
    if err != nil {
        req.Status = false
        req.Message = "error parsing application/json: " + err.Error()
    } else {
        req.Status = true
        Ins := UpdateDataGeojson(conn, colname,
            resp.Name,
            resp.Volume,
            resp.Type)
        req.Message = fmt.Sprintf("%v:%v", "Berhasil Update data", Ins)
    }
    return GCFReturnStruct(req)
}



func GCFDelDataGeo(MONGOCONNSTRINGENV, Mongostring, dbname, colname string, r *http.Request) string {
    req := new(Credential)
    resp := new(CoorLonLatProperties)

    tokenLogin := r.Header.Get("Login")
    if tokenLogin == "" {
        req.Status = false
        req.Message = "Header Login Not Exist"
        return GCFReturnStruct(req)
    }

    // Validate the token using your existing logic
    existing := IsExist(tokenLogin, os.Getenv(MONGOCONNSTRINGENV))

    if !existing {
        req.Status = false
        req.Message = "Kamu belum memiliki akun"
        return GCFReturnStruct(req)
    }

    conn := SetConnection(Mongostring, dbname)
    err := json.NewDecoder(r.Body).Decode(&resp)
    if err != nil {
        req.Status = false
        req.Message = "error parsing application/json: " + err.Error()
    } else {
        req.Status = true
        delResult, delErr := DeleteDataGeojson(conn, colname, resp.Name)
        if delErr != nil {
            req.Status = false
            req.Message = "error deleting data: " + delErr.Error()
        } else {
            req.Message = fmt.Sprintf("Berhasil menghapus data. Jumlah data terhapus: %v", delResult.DeletedCount)
        }
    }
    return GCFReturnStruct(req)
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

func IsExist(Tokenstr, PublicKey string) bool {
	id := watoken.DecodeGetId(PublicKey, Tokenstr)
	if id == "" {
		return false
	}
	return true
}
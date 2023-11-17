package gcfbackend1214005

import (
	"os"
	"context"
	"github.com/aiteung/atdb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetConnection(MONGOCONNSTRINGENV, dbname string) *mongo.Database {
	var DBmongoinfo = atdb.DBInfo{
		DBString: os.Getenv(MONGOCONNSTRINGENV),
		DBName:   dbname,
	}
	return atdb.MongoConnect(DBmongoinfo)
}

func GetAllBangunanLineString(mongoconn *mongo.Database, collection string) []GeoJson {
	lokasi := atdb.GetAllDoc[[]GeoJson](mongoconn, collection)
	return lokasi
}

func IsPasswordValid(mongoconn *mongo.Database, collection string, userdata User) bool {
	filter := bson.M{"username": userdata.Username}
	res := atdb.GetOneDoc[User](mongoconn, collection, filter)
	return CheckPasswordHash(userdata.Password, res.Password)
}

func InsertDataLonlat(mongoconn *mongo.Database, collection string, coordinate [][]float64, name, volume, tipe string) (InsertedID interface{}) {
	req := new(CoorLonLatProperties)
	req.Type = tipe
	req.Coordinates = coordinate
	req.Name = name
	req.Volume = volume

	ins := atdb.InsertOneDoc(mongoconn, collection, req)
	return ins
}

func UpdateDataGeojson(mongoconn *mongo.Database, colname, name, newVolume, newTipe string) error {
    // Filter berdasarkan nama
    filter := bson.M{"name": name}

    // Update data yang akan diubah
    update := bson.M{
        "$set": bson.M{
            "volume": newVolume,
            "tipe":   newTipe,
        },
    }

    // Mencoba untuk mengupdate dokumen
    _, err := mongoconn.Collection(colname).UpdateOne(context.TODO(), filter, update)
    if err != nil {
        return err
    }

    return nil
}

func DeleteDataGeojson(mongoconn *mongo.Database, colname string, name string) (*mongo.DeleteResult, error) {
    filter := bson.M{"name": name}
    del, err := mongoconn.Collection(colname).DeleteOne(context.TODO(), filter)
    if err != nil {
        return nil, err
    }
    return del, nil
}

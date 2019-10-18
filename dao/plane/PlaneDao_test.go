package dao

import (
	"RESTApp/model"
	"RESTApp/utils/mongo"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2"
)

var gDB *mgo.Session

func TestPutPlane(t *testing.T) {
	testCases := []struct {
		url        string
		db         string
		collection string
		err        error
	}{
		{
			url:        "localhost:27017",
			db:         "testing",
			collection: "planes",
			err:        nil,
		},
	}

	gDB, _ := mongo.GetDataBaseSession(testCases[0].url)

	plane := model.Plane{
		Pid:      7,
		Name:     "MIG19",
		NoWheels: 6,
		Engines:  4,
		PType:    "Attack",
	}
	err := PutPlane(plane, gDB)
	assert.Equal(t, testCases[0].err, err, "Expected %v but got %v", testCases[0].err, err)
	defer gDB.Close()
}

func TestGetPlane(t *testing.T) {
	tc := []struct {
		Name string
		err  error
	}{
		{
			Name: "Airbus",
			err:  nil,
		},
	}

	gDB, _ := mongo.GetDataBaseSession("localhost:27017")

	p := GetPlane(tc[0].Name, gDB)
	assert.Equal(t, tc[0].Name, p.Name, "Expected %s but got %s", tc[0].Name, p.Name)

	// p1 := GetPlane(tc[1].Name, gDB)
	// assert.Equal
	defer gDB.Close()
}

func TestUpdatePlane(t *testing.T) {
	pl := model.Plane{
		Pid:      8,
		Name:     "Boeing 777",
		NoWheels: 24,
		Engines:  8,
		PType:    "Cargo",
	}

	gDB, _ := mongo.GetDataBaseSession("localhost:27017")

	defer gDB.Close()

	gP := GetPlane(pl.Name, gDB)
	//fmt.Println("...............................getPlane", gP)
	gP.Pid = pl.Pid
	gP.NoWheels = pl.NoWheels
	gP.Engines = pl.Engines
	gP.PType = pl.PType
	p, _ := UpdatePlane(gP, gDB)
	//fmt.Println("...............................updatePlane", p)
	assert.Equal(t, p.Pid, pl.Pid, "Exepected %s but got %v", p.Pid, pl.Pid)
}

func TestRemovePlane(t *testing.T) {
	gDB, _ := mongo.GetDataBaseSession("localhost:27017")
	defer gDB.Close()
	err := DeletePlane("MIG19", gDB)
	assert.Equalf(t, true, err, "Expected %s but got %s", true, err)
}

func TestRemovePlaneErr(t *testing.T) {
	gDB, _ := mongo.GetDataBaseSession("localhost:27017")
	defer gDB.Close()

	err := DeletePlane("name", gDB)
	assert.Error(t, errors.New("not found"), err, "..")
}

func TestGetAllPlanes(t *testing.T) {
	gDB, _ := mongo.GetDataBaseSession("localhost:27017")
	defer gDB.Close()

	planes, _ := GetAllPlanes(gDB)

	var p string
	for _, val := range planes {
		if val.Name == "MIG19" {
			p = val.Name
		}
		continue
	}
	//assert.Equal(t, []model.Plane, planes, "-------")
	assert.Equal(t, p, "MIG19")
}

func TestDeleteByID(t *testing.T) {
	gDB, _ := mongo.GetDataBaseSession("localhost:27017")

	err := DeletePlaneByID(7, gDB)
	assert.Equal(t, true, err)
}

func TestDeleteByIDErr(t *testing.T) {
	gDB, _ := mongo.GetDataBaseSession("localhost:27017")

	err := DeletePlaneByID(1221321343, gDB)
	assert.Error(t, errors.New("not found"), err, "...")
}

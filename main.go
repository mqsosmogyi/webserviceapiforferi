package main

import (
	"fmt"
	"net/http"

	"github.com/astaxie/beego/orm" //For SQL queries
	"github.com/gin-gonic/gin"     //Simple REST Api service
	_ "github.com/lib/pq"          //PostgreSQL Driver
)

var ormObject orm.Ormer
var ORM orm.Ormer

//Database table structure
type device struct {
	Id     int    `json:"id" orm:"auto"`
	Name   string `json:"name" orm:"size(128)"`
	Status string `json:"status" orm:"size(128)"`
}

//Connection to local Postgres DB
func ConnectToDb() {
	orm.RegisterDriver("postgres", orm.DRPostgres)
	orm.RegisterDataBase("default", "postgres", "user=postgres password=laszlo123 dbname=webapi host=localhost sslmode=disable")
	orm.RegisterModel(new(device))
	ormObject = orm.NewOrm()
}

func GetOrmObject() orm.Ormer {
	return ormObject
}

func init() {
	ConnectToDb()
	ORM = GetOrmObject()
}

//Get the device list from DB -> device table | example GET request http://localhost:8888/devices
func getDevices(c *gin.Context) {
	var name []device
	fmt.Println(ORM)
	_, err := ORM.QueryTable("device").All(&name)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "devices": &name})
	} else {
		c.JSON(http.StatusInternalServerError,
			gin.H{"status": http.StatusInternalServerError, "error": "Failed to get device list"})
	}
}

//Add a new device to DB -> device table with id/name/status parameters | example POST request http://localhost:8888/devices/Mobile | json should contain id/name/status fields
func createDevice(c *gin.Context) {
	var newDevice device
	c.BindJSON(&newDevice)
	test, err := ORM.Insert(&newDevice)
	if err != nil {

		c.JSON(http.StatusInternalServerError,
			gin.H{"status": http.StatusInternalServerError, "error": "Failed add a new device"})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"id":     newDevice.Id,
			"status": http.StatusOK,
			"name":   newDevice.Name})
	}
	fmt.Print(test)
}

//Remove a device from DB by name | example DELETE request: http://localhost:8888/devices/Mobile
func removeDevice(c *gin.Context) {
	var deleteDevice device
	c.BindJSON(&deleteDevice)
	_, err := ORM.QueryTable("device").Filter("name", deleteDevice.Name).Delete()
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
	} else {
		c.JSON(http.StatusInternalServerError,
			gin.H{"status": http.StatusInternalServerError, "error": "Failed to delete device"})
	}
}

func main() {
	router := gin.Default()
	router.GET("/devices", getDevices)
	router.POST("/devices", createDevice)
	router.DELETE("/devices/:name", removeDevice)
	router.Run("localhost:8888")
}

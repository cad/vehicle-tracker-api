package repository

import (
	"fmt"
	"log"
	//"github.com/jinzhu/gorm"
	"time"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

)


type Vehicle struct {
	ID        uint `json:"-" gorm:"primary_key"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"-"`
	PlateID string `json:"plate_id" gorm:"not null;unique_index"`
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}


func GetVehicleByID (vehicleId int) (Vehicle, error) {
	var vehicle Vehicle
	db.First(&vehicle, vehicleId)
	if vehicle != (Vehicle{}) {
		return vehicle, nil
	}
	return vehicle, &VehicleError{What: "Vehicle", Type: "Not-Found", Arg: fmt.Sprintf("%v", vehicleId)}
}


func GetVehicleByPlateID (plateID string) (Vehicle, error) {
	var vehicle Vehicle
	db.Where(&Vehicle{PlateID: plateID}).First(&vehicle)
	if vehicle != (Vehicle{}) {
		return vehicle, nil
	}
	return vehicle, &VehicleError{What: "Vehicle", Type: "Not-Found", Arg: plateID}
}


func GetAllVehicles () []Vehicle {
	var vehicles []Vehicle

	db.Find(&vehicles)

	return vehicles
}


func CreateVehicle (plateID string) {
	vehicle := Vehicle{
		PlateID: plateID,  // TODO(cad): sanitize `plateID`
	}
	db.Create(&vehicle)  // TODO(cad): check here if vehicle
	                     // created suc⎈cessfully.
}


func SyncVehicleByID(vehicleID int, lat float64, lon float64) error {
	var vehicle Vehicle
	vehicle, err := GetVehicleByID(vehicleID)
	if err != nil {
		return err
	}

	vehicle.Lat = lat
	vehicle.Lon = lon
	db.Save(&vehicle)
	return nil
}


func SyncVehicleByPlateID(plateID string, lat float64, lon float64) error {
	var vehicle Vehicle
	vehicle, err := GetVehicleByPlateID(plateID)
	if err != nil {
		return err
	}

	vehicle.Lat = lat
	vehicle.Lon = lon
	db.Save(&vehicle)
	return nil
}


func DeleteVehicleByID (id int) error{
	vehicle, err := GetVehicleByID(id)
	if err != nil {
		return err
	}

	log.Println(vehicle)
	db.Delete(&vehicle)
	return nil
}


func DeleteVehicleByPlateID (plateID string) error{
	vehicle, err := GetVehicleByPlateID(plateID)
	if err != nil {
		return err
	}

	log.Println(vehicle)
	db.Delete(&vehicle)
	return nil
}


type VehicleError struct {
	What string
	Type string
	Arg string
}


func (e VehicleError) Error() string {
	return fmt.Sprintf("%s-%s: %s", e.What, e.Arg, e.Type)
}

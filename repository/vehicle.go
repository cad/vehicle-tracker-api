package repository

import (
	"fmt"
	"log"
	"strconv"
	"time"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
//	"github.com/jinzhu/gorm"

)

const (
	SCHOOL_BUS = "SCHOOL-BUS"
	SOLAR_CAR = "SOLAR-CAR"
)

var VEHICLE_TYPES []string = []string{SCHOOL_BUS, SOLAR_CAR}

type Vehicle struct {
	ID        uint       `json:"-"           gorm:"primary_key"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-"`
	PlateID   string     `json:"plate_id"    gorm:"not null;unique_index"`
	Agent     Agent      `json:"agent"       gorm:"not null;ForeignKey:AgentID"`
	AgentID   uint       `json:"-"`
	Groups    []Group    `json:"groups"      gorm:"many2many:vehicle_group;"`
	Type      string     `json:"type"`
}

type Group struct {
	ID        uint       `json:"id"           gorm:"primary_key"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-"`
	Name      string     `json:"name"         gorm:"not null;unique_index"`
}


func GetVehicleByPlateID (plateID string) (Vehicle, error) {
	var vehicle Vehicle
	db.Preload("Groups").Preload("Agent").Where(&Vehicle{PlateID: plateID}).First(&vehicle)
	if vehicle.ID != 0 {
		return vehicle, nil
	}
	return vehicle, &VehicleError{What: "Vehicle", Type: "Not-Found", Arg: plateID}
}


func GetAllVehicles () []Vehicle {
	var vehicles []Vehicle

	db.Preload("Groups").Preload("Agent").Find(&vehicles)

	return vehicles
}


func FilterVehicles (types []string, groupids []uint) []Vehicle {
	var vehicles []Vehicle

	q := db.Preload("Groups").Preload("Agent")

	if len(types) > 1 {
		q = q.Where("type in (?)", types)
	} else if len(types) == 1 {
		q = q.Where("type = ?", types[0])
	}

	if len(groupids) > 1 {
		q = q.Where("groupid in (?)", groupids)
	} else if len(groupids) > 1 {
		q = q.Where("groupid = ?", groupids[0])
	}

	q.Find(&vehicles)

	return vehicles
}

func VehicleSetAgent(plateID,uUID string) error {
	vehicle, err := GetVehicleByPlateID(plateID)
	if err != nil {
		return err
	}
	vehicle.Agent, err = GetAgentByUUID(uUID)
	if err != nil {
		return err
	}

	db.Save(&vehicle)
	return nil
}

func CreateVehicle (plateID string, agentID string, groupIDs []int, vehicleType string) error {
	groups := make([]Group, 0)
	// Sanitize incoming
	for _, item := range groupIDs {
		var group Group
		db.First(&group, item)
		if group == (Group{}) {
			return &VehicleError{What: "Group", Type: "Not-Found", Arg: strconv.Itoa(item)}
		}
		groups = append(groups, group)
	}

	typeFound := false
	for _, item := range VEHICLE_TYPES {
		if vehicleType == item {
			typeFound = true
		}
	}

	if !typeFound {
		return &VehicleError{What: "VehicleType", Type: "Not-Found", Arg: vehicleType}
	}


	// Create Vehicle
	vehicle := Vehicle{
		PlateID: plateID,  // TODO(cad): sanitize `plateID`
		Type: vehicleType,
	}
	db.Where(&Agent{UUID: agentID}).First(&vehicle.Agent)
	if vehicle.Agent == (Agent{}) {
		return &VehicleError{What: "Agent", Type: "Not-Found", Arg: agentID}
	}

	vehicle.Groups = make([]Group,0)
	// Set groups if not empty
	if len(groups) > 0 {
		vehicle.Groups = groups
	}

	db.Create(&vehicle)  // TODO(cad): check here if vehicle
	if db.NewRecord(&vehicle) == true {
		return &VehicleError{What: "Vehicle.PlateID", Type: "Already-Exists", Arg: plateID}
	}
	// created suc⎈cessfully.
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


func CreateNewGroup(name string) (uint, error) {
	group := Group{Name: name}
	db.Create(&group)
	if db.NewRecord(&group) {
		return group.ID, &VehicleError{
			What: "Group.Name",
			Type: "Unknown-Error",
			Arg: name,
		}
	}
	return group.ID, nil
}


func GetAllGroups () []Group {
	var groups []Group

	db.Find(&groups)

	return groups
}

func GetGroupByID (iD uint) (Group, error) {
	var group Group
	db.Where(&Group{ID: iD}).First(&group)
	if group.ID != 0 {
		return group, nil
	}
	return group, &VehicleError{What: "VehicleGroup.ID", Type: "Not-Found", Arg: fmt.Sprintf("%d", iD)}
}

func GetGroupByName (name string) (Group, error) {
	var group Group
	db.Where(&Group{Name: name}).First(&group)
	if group.ID != 0 {
		return group, nil
	}
	return group, &VehicleError{What: "VehicleGroup.Name", Type: "Not-Found", Arg: name}
}

type VehicleError struct {
	What string
	Type string
	Arg string
}


func (e VehicleError) Error() string {
	return fmt.Sprintf("%s: <%s> %s", e.Type, e.What, e.Arg)
}

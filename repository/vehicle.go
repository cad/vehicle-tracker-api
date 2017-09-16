package repository

import (
	"fmt"
	"strconv"
	"time"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
	//	"github.com/jinzhu/gorm"
)

const (
	SCHOOL_BUS = "SCHOOL-BUS"
	SOLAR_CAR  = "SOLAR-CAR"
)

var VEHICLE_TYPES []string = []string{SCHOOL_BUS, SOLAR_CAR}

type Vehicle struct {
	ID        uint      `json:"-"           gorm:"primary_key"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"updated_at"`
	PlateID   string    `json:"plate_id"    gorm:"not null;unique_index"`
	Agent     *Agent    `json:"agent"       gorm:"ForeignKey:AgentID"`
	AgentID   uint      `json:"-"`
	Groups    []Group   `json:"groups"      gorm:"many2many:vehicle_group;"`
	Type      string    `json:"type"`
}

type Group struct {
	ID        uint       `json:"id"           gorm:"primary_key"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-"`
	Name      string     `json:"name"         gorm:"not null;unique_index"`
}

func GetVehicleByPlateID(plateID string) (Vehicle, error) {
	var vehicle Vehicle
	if plateID == "" {
		return vehicle, &VehicleError{What: "plateID", Type: "Empty", Arg: plateID}
	}
	db.Preload("Groups").Preload("Agent").Where(&Vehicle{PlateID: plateID}).First(&vehicle)
	if vehicle.ID != 0 {
		return vehicle, nil
	}
	return vehicle, &VehicleError{What: "Vehicle", Type: "Not-Found", Arg: plateID}
}

func GetVehicleByAgentUUID(agentUUID string) (Vehicle, error) {
	var vehicle Vehicle
	if agentUUID == "" {
		return vehicle, &VehicleError{What: "agentUUID", Type: "Empty", Arg: agentUUID}
	}
	agent, err := GetAgentByUUID(agentUUID)
	if err != nil {
		return vehicle, err
	}

	db.Preload("Groups").Preload("Agent").Where(&Vehicle{AgentID: agent.ID}).First(&vehicle)

	if db.NewRecord(&vehicle) {
		return vehicle, &VehicleError{What: "Vehicle.Agent", Type: "Not-Found", Arg: agentUUID}
	}
	return vehicle, nil
}

func GetAllVehicles() []Vehicle {
	var vehicles []Vehicle

	db.Preload("Groups").Preload("Agent").Find(&vehicles)

	return vehicles
}

func FilterVehicles(vehicleType string, groupID uint) []Vehicle {
	var vehicles []Vehicle

	q := db.Preload("Groups").Preload("Agent").Joins("JOIN vehicle_group ON vehicle_group.vehicle_id = vehicles.id")
	if groupID != *new(uint) {
		q = q.Where("vehicle_group.group_id = ?", groupID)
	}

	if vehicleType != *new(string) {
		q = q.Where("vehicles.type = ?", vehicleType)
	}

	q.Find(&vehicles)

	return vehicles
}

func VehicleSetAgent(plateID, uUID string) error {
	vehicle, err := GetVehicleByPlateID(plateID)
	if err != nil {
		return err
	}
	agent, err := GetAgentByUUID(uUID)
	if err != nil {
		return err
	}
	agentVehicle := agent.Vehicle()
	if agentVehicle != nil {
		agentVehicle.Agent = nil
		db.Save(&agentVehicle)
	}
	vehicle.Agent = &agent

	db.Save(&vehicle)
	return nil
}

func VehicleUnsetAgent(plateID string) error {
	vehicle, err := GetVehicleByPlateID(plateID)
	if err != nil {
		return err
	}

	vehicle.AgentID = 0
	vehicle.Agent = nil

	db.Save(&vehicle)
	return nil
}

func CreateVehicle(plateID string, agentUUID string, groupIDs []int, vehicleType string) error {
	if plateID == "" {
		return &VehicleError{What: "plateID", Type: "Empty", Arg: plateID}
	}
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
		PlateID: plateID, // TODO(cad): sanitize `plateID`
		Type:    vehicleType,
	}
	if agentUUID != "" {
		var a Agent
		db.Where(&Agent{UUID: agentUUID}).First(&a)
		vehicle.Agent = &a
		//vehicle.AgentID = a.ID
	}

	vehicle.Groups = make([]Group, 0)
	// Set groups if not empty
	if len(groups) > 0 {
		vehicle.Groups = groups
	}

	db.Create(&vehicle) // TODO(cad): check here if vehicle
	if db.NewRecord(&vehicle) == true {
		return &VehicleError{What: "Vehicle.PlateID", Type: "Already-Exists", Arg: plateID}
	}
	// created suc⎈cessfully.
	return nil
}

func DeleteVehicleByPlateID(plateID string) error {
	vehicle, err := GetVehicleByPlateID(plateID)
	if err != nil {
		return err
	}

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
			Arg:  name,
		}
	}
	return group.ID, nil
}

func DeleteGroup(groupID uint) error {
	var group Group

	db.First(&group, groupID)
	if db.NewRecord(&group) {
		return &VehicleError{
			What: "Group.ID",
			Type: "Unknown-Error",
			Arg:  fmt.Sprintf("%d", groupID),
		}
	}
	db.Delete(&group)
	return nil
}

func GetAllGroups() []Group {
	var groups []Group

	db.Find(&groups)

	return groups
}

func GetAllTypes() []string {
	return VEHICLE_TYPES
}

func GetGroupByID(iD uint) (Group, error) {
	var group Group
	db.Where(&Group{ID: iD}).First(&group)
	if group.ID != 0 {
		return group, nil
	}
	return group, &VehicleError{What: "VehicleGroup.ID", Type: "Not-Found", Arg: fmt.Sprintf("%d", iD)}
}

func GetGroupByName(name string) (Group, error) {
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
	Arg  string
}

func (e VehicleError) Error() string {
	return fmt.Sprintf("%s: <%s> %s", e.Type, e.What, e.Arg)
}

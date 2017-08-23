package endpoints

import (
	"github.com/cad/vehicle-tracker-api/repository"
)

// Returns a vehicle
// swagger:response
type VehicleSuccessVehicleResponse struct {
	// Vehicle
	// in: body
	Body repository.Vehicle
}

// Returns list of vehicles
// swagger:response
type VehicleSuccessVehiclesResponse struct {
	// Vehicles
	// in: body
	Body []repository.Vehicle
}

// Returns a vehicle group
// swagger:response
type VehicleSuccessVehicleGroupResponse struct {
	// Vehicle Group
	// in: body
	Body repository.Group
}

// Returns list of vehicle groups
// swagger:response
type VehicleSuccessVehicleGroupsResponse struct {
	// Vehicles
	// in: body
	Body []repository.Group
}

// Returns list of vehicle types
// swagger:response
type VehicleSuccessVehicleTypesResponse struct {
	// Vehicle types
	// in: body
	Body []string
}

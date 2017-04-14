package endpoints

import (
	"github.com/cad/vehicle-tracker-api/repository"
)


// Returns an agent
// swagger:response
type AgentSuccessAgentResponse struct {
	// Agent
	// in: body
	Body repository.Agent
}


// Returns empty object
// swagger:response
type AgentSuccessEmptyResponse struct {
}


// Returns list of agents
// swagger:response
type AgentSuccessAgentsResponse struct {
	// Agents
	// in: body
	Body []repository.Agent

}

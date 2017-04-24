package repository

import (
	"fmt"
//	"log"
	"time"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

)

type Agent struct {
	ID        uint       `json:"-"    gorm:"primary_key"`
	UUID      string     `json:"uuid" gorm:"not null;unique_index"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-"`

	Label     string     `json:"label"`

	Lat       string    `json:"lat"`
	Lon       string    `json:"lon"`
	TS        string  `json:"gps_ts"`
}

func GetAllAgents () []Agent {
	var agents []Agent

	db.Find(&agents)

	return agents
}

func GetAgentByUUID(uUID string) (Agent, error) {
	var agent Agent

	db.Where(&Agent{UUID: uUID}).First(&agent)
	if agent != (Agent{}) {
		return agent, nil
	}

	return agent, AgentError{
		What: "Agent",
		Type: "Not-Found",
		Arg: uUID,
	}
}

func CreateNewAgent(uUID string) (Agent, error) {
	agent := Agent{
		UUID: uUID,
	}
	db.Create(&agent)
	if agent == (Agent{}) {
		return agent, AgentError{
			What: "Agent",
			Type: "Can-Not-Create",
			Arg: uUID,
		}
	}
	return agent, nil
}

func SetLabelByUUID(uUID string, label string) error {
	var agent Agent
	agent, err := GetAgentByUUID(uUID)
	if err != nil {
		return err
	}

	agent.Label = label
	db.Save(&agent)
	return nil
}

func SyncAgentByUUID(uUID string, lat string, lon string, ts string) error {
	var agent Agent
	agent, err := GetAgentByUUID(uUID)
	if (err != nil) && (agent == Agent{}) {
		agent, err = CreateNewAgent(uUID)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	agent.Lat = lat
	agent.Lon = lon
	agent.TS = ts
	db.Save(&agent)
	return nil
}

type AgentError struct {
	What string
	Type string
	Arg string
}


func (e AgentError) Error() string {
	return fmt.Sprintf("%s: <%s> %s", e.Type, e.What, e.Arg)
}

package models

import "time"

type Agent struct {
	Id int64
	AgentId string `json:"agent_id"`
	PlatForm string `json:"platform"`
	Architecture string `json:"architecture"`
	UserName string `json:"user_name"`
	UserGUID string `json:"user_guid"`
	HostName string `json:"host_name"`
	Ips []string `json:"ips" xorm:"text"`
	Pid int	`json:"pid"`
	Debug bool `json:"debug"`
	Proto string `json:"proto"`
	UserAgent string `json:"user_agent"`
	Initial bool `json:"initial"`
	CreateTime time.Time `xorm:"created"`
	UpdateTime time.Time `xorm:"updated"`
	Version int	`xorm:version`
}

func (a *Agent) Insert() error {
	_, err := Engine.Insert(a)

	return err
}

func ListAgents() ([]Agent, error)  {
	agents := make([]Agent, 0)
	err := Engine.Find(&agents)

	return agents, err
}

func UpdateAgent(agentId string, curAgent Agent) error {
	agent := new(Agent)
	has, err := ExitAgentId(agentId)
	if err != nil {
		return err
	}
	if has {
		_, err = Engine.Id(agent.Id).Update(curAgent)
	}
	return err
}

func ExitAgentId(agentId string) (bool, error) {
	agent := new(Agent)
	has, err := Engine.Where("agent_id=?", agentId).Get(agent)

	return has, err
}

func RemoveAll() error {
	_, err := Engine.Exec("delete from agent")

	return err
}

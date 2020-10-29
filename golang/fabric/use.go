package fabric

import (
	"bytes"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/event"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/status"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"net/http"
	"time"
)

func init() {
	var err error
	sdk, err = fabsdk.New(config.FromFile(configPath))
	if err != nil {
		panic(err)
	}
}

func queryUser(ctx *gin.Context) {
	userId := ctx.Param("id")

	resp, err := channelQuery("queryUser", [][]byte{
		[]byte(userId),
	})

	if err != nil {
		ctx.String(http.StatusOK, err.Error())
		return
	}

	ctx.String(http.StatusOK, bytes.NewBuffer(resp.Payload).String())
}

func manageBlockChain() {
	ctx := sdk.Context(fabsdk.WithOrg(org), fabsdk.WithUser(user))

	cli, err := resmgmt.New(ctx)
	if err != nil {
		panic(err)
	}

	cli.SaveChannel(resmgmt.SaveChannelRequest{}, resmgmt.WithOrdererEndpoint(""), resmgmt.WithTargetEndpoints())
}

func queryBlockchain() {
	ctx := sdk.ChannelContext(channelName, fabsdk.WithOrg(org), fabsdk.WithUser(user))

	cli, err := ledger.New(ctx)
	if err != nil {
		panic(err)
	}

	resp, err := cli.QueryInfo(ledger.WithTargetEndpoints())
	cli.QueryBlockByHash(resp.BCI.CurrentBlockHash)

	for i := uint64(0); i <= resp.BCI.Height; i++ {
		cli.QueryBlock(i)
	}
}

func channelExecute(fcn string, args [][]byte) (channel.Response, error) {
	ctx := sdk.ChannelContext(channelName, fabsdk.WithOrg(org), fabsdk.WithUser(user))

	cli, err := channel.New(ctx)
	if err != nil {
		return channel.Response{}, err
	}

	resp, err := cli.Execute(channel.Request{
		ChaincodeID: chaincodeName,
		Fcn:         fcn,
		Args:        args,
	}, channel.WithTargetEndpoints("peer0.org1.imocc.com"))

	if err != nil {
		return channel.Response{}, err
	}

	go func() {
		reg, ccvert, err := cli.RegisterChaincodeEvent(chaincodeName, "eventname")
		if err != nil {
			return
		}
		defer cli.UnregisterChaincodeEvent(reg)

		timeoutctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		for {
			select {
			case evt := <-ccvert:
				fmt.Printf("received event of tx %s: %+v", resp.TransactionID, evt)
			case <-timeoutctx.Done():
				fmt.Println()
				return
			}
		}
	}()

	go func() {
		eventcli, err := event.New(ctx)
		if err != nil {
			return
		}

		reg, status, err := eventcli.RegisterTxStatusEvent(string(resp.TransactionID))
		defer eventcli.Unregister(reg)

		timeoutctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		for {
			select {
			case evt := <-status:
				fmt.Printf("received event of tx %s: %+v", resp.TransactionID, evt)
			case <-timeoutctx.Done():
				fmt.Println("event timeout, exit!")
				return
			}
		}
	}()
}

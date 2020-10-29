package main

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/hyperledger/fabric/core/ledger/util"
	"github.com/hyperledger/fabric/protos/common"
	"github.com/hyperledger/fabric/protos/utils"
)

func parseBlock(block *common.Block) error {
	var err error

	// Handle header
	fmt.Printf("Block: Number=[%d], CurrentBlockHash=[%s], PreviousBlockHash=[%s]\n",
		block.GetHeader().Number,
		base64.StdEncoding.EncodeToString(block.GetHeader().DataHash),
		base64.StdEncoding.EncodeToString(block.GetHeader().PreviousHash))
	fmt.Printf("%f", time.Now().UnixNano())
	// Handle transaction
	var tranNo int64 = -1
	txsFilter := util.TxValidationFlags(block.Metadata.Metadata[common.BlockMetadataIndex_TRANSACTIONS_FILTER])
	if len(txsFilter) == 0 {
		txsFilter = util.NewTxValidationFlags(len(block.Data.Data))
		block.Metadata.Metadata[common.BlockMetadataIndex_TRANSACTIONS_FILTER] = txsFilter
	}

	for _, envBytes := range block.Data.Data {
		tranNo++
		if txsFilter.IsInvalid(int(tranNo)) {
			fmt.Printf("    Transaction: No=[%d], Status=[INVALID]\n", tranNo)
			continue
		} else {
			fmt.Printf("    Transaction: No=[%d], Status=[VALID]\n",   tranNo)
		}

		var env *common.Envelope
		if env, err = utils.GetEnvelopeFromBlock(envBytes); err != nil {
			return err
		}

		var payload *common.Payload
		if payload, err = utils.GetPayload(env); err != nil {
			return err
		}

		var chdr *common.ChannelHeader
		if chdr, err = utils.UnmarshalChannelHeader(payload.Header.ChannelHeader); err != nil {
			return err
		}
		fmt.Printf("        txid=[%s], channel=[%s]\n", chdr.TxId, chdr.ChannelId)
	}

	return nil
}
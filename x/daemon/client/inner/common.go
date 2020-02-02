package inner

import (
	"fmt"
	// "github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	
	"github.com/rhizomata-io/dist-daemon-cosmos/node"
	"github.com/rhizomata-io/dist-daemon-cosmos/x/daemon"
	"github.com/rhizomata-io/dist-daemon-cosmos/x/daemon/internal/types"
)

func BroadcastTx(msg types.DaemonMsg) (err error) {
	err = msg.ValidateBasic()
	if err != nil {
		return err
	}
	
	cliCtx := node.GetCLIContext()
	
	txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cliCtx.Codec)).WithChainID(node.ChainID())
	
	sigMsg, err := txBldr.BuildSignMsg([]sdk.Msg{msg})
	
	if err != nil {
		// fmt.Println("---- ERROR : BuildSignMsg :: " , err)
		return err
	}
	
	bytes, err := cliCtx.Codec.MarshalJSON(sigMsg)
	
	if err != nil {
		// fmt.Println("---- ERROR : MarshalJSON :: " , err)
		return err
	} else {
		fmt.Println("---- MarshalJSON :: " , string(bytes))
	}
	
	// err = utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
	res, err := cliCtx.BroadcastTx(bytes)
	if err != nil {
		return err
	} else {
		err = cliCtx.PrintOutput(res)
	}
	
	return err
}

func Query(msg sdk.Msg) (err error) {
	cliCtx := node.GetCLIContext()
	res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", daemon.StoreKey, msg.Type()), nil)
	if err != nil {
		fmt.Printf("Could not Query - %s \n")
		return err
	}
	
	cliCtx.Codec.MustUnmarshalJSON(res, &msg)
	return nil
}

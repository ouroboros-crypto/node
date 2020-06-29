package rest

import (
	"encoding/base64"
	"fmt"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	ouroTypes "github.com/ouroboros-crypto/node/x/ouroboros/types"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/ouroboros/profile/{address}", getHandler(cliCtx)).Methods("GET")
	r.HandleFunc("/ouroboros/profile/{address}/{coin}", getHandlerCoin(cliCtx)).Methods("GET")

	r.HandleFunc("/ouroboros/encode", encodeTx(cliCtx)).Methods("POST")
	r.HandleFunc("/ouroboros/decode", decodeTx(cliCtx)).Methods("POST")

}

func getHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)

		if !ok {
			return
		}

		vars := mux.Vars(r)

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/ouroboros/profile/%s/ouro",  vars["address"]), nil)

		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponseBare(w, cliCtx, res)
	}
}

func getHandlerCoin(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)

		if !ok {
			return
		}

		vars := mux.Vars(r)

		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/ouroboros/profile/%s/%s", vars["address"], vars["coin"]), nil)

		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)

		rest.PostProcessResponse(w, cliCtx, res)
	}
}


func encodeTx(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.StdTx

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		err = cliCtx.Codec.UnmarshalJSON(body, &req)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// re-encode it via the Amino wire protocol
		txBytes, err := cliCtx.Codec.MarshalBinaryLengthPrefixed(req)

		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		// base64 encode the encoded tx bytes
		txBytesBase64 := base64.StdEncoding.EncodeToString(txBytes)

		response := ouroTypes.EncodeResp{Tx: txBytesBase64}
		rest.PostProcessResponseBare(w, cliCtx, response)
	}
}

func decodeTx(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ouroTypes.DecodeReq
		var resp types.StdTx

		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		err = cliCtx.Codec.UnmarshalJSON(body, &req)

		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		decodedTx, err := base64.StdEncoding.DecodeString(req.Tx)

		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// re-encode it via the Amino wire protocol
		err = cliCtx.Codec.UnmarshalBinaryLengthPrefixed(decodedTx, &resp)

		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		rest.PostProcessResponseBare(w, cliCtx, resp)
	}
}
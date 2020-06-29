package rest

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/emission/get", queryGetHandler(cliCtx)).Methods("GET")
	r.HandleFunc("/emission/get/{coin}", queryGetCoinHandler(cliCtx)).Methods("GET")
}

// Returns OURO emission
func queryGetHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, _, err := cliCtx.QueryWithData("custom/emission/get/ouro", nil)

		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

// Returns custom coin emission
func queryGetCoinHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)

		if !ok {
			return
		}

		vars := mux.Vars(r)

		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/emission/get/%s", vars["coin"]), nil)

		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)

		rest.PostProcessResponseBare(w, cliCtx, res)
	}
}

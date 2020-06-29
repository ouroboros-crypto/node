package rest

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/paramining/get/{address}", queryGetHandler(cliCtx)).Methods("GET")
	r.HandleFunc("/posmining/get/{address}/{coin}", queryGetWithCoinHandler(cliCtx)).Methods("GET")

}

func queryGetHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		vars := mux.Vars(r)

		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/posmining/get/%s/ouro", vars["address"]), nil)

		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)

		rest.PostProcessResponseBare(w, cliCtx, res)
	}
}


func queryGetWithCoinHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)

		if !ok {
			return
		}

		vars := mux.Vars(r)

		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/posmining/get/%s/%s", vars["address"], vars["coin"]), nil)

		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}
		cliCtx = cliCtx.WithHeight(height)

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

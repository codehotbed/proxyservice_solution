package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-openapi/runtime/middleware"

	"ProxyService/proxyservice/internal/queue"
	"ProxyService/proxyservice/openapi/autogen/proxyservice/server/operations/rpc"
	"ProxyService/proxyservice/utils"
)

var (
	clientsQueue  = "backend.clients"
	invoicesQueue = "backend.invoices"
)

// RPCApi ...
type RPCApi struct {
	queue *queue.Queue
}

// NewRPCApi RPCApi constructor
func NewRPCApi(q *queue.Queue) *RPCApi {
	return &RPCApi{queue: q}
}

// GetClients Returns all clients acquired from clients queue
func (api *RPCApi) GetClients(params rpc.GetClientsParams) middleware.Responder {
	qResp, err := api.queue.CallRPC(clientsQueue, nil)
	if err != nil {
		errResponse := utils.NewErrorResponse(http.StatusInternalServerError, err.Error())
		return rpc.NewGetClientsInternalServerError().WithPayload(errResponse)
	}

	return rpc.NewGetClientsOK().WithPayload(qResp)
}

// GetInvoices Returns all invoices acquired from invoices queue based on given client_id parameter
func (api *RPCApi) GetInvoices(params rpc.GetInvoicesParams) middleware.Responder {

	args, err := json.Marshal(params.ClientID)
	if err != nil {
		errResponse := utils.NewErrorResponse(http.StatusInternalServerError, err.Error())
		return rpc.NewGetInvoicesInternalServerError().WithPayload(errResponse)
	}

	qResp, err := api.queue.CallRPC(invoicesQueue, &args)
	if err != nil {
		errResponse := utils.NewErrorResponse(http.StatusInternalServerError, err.Error())
		return rpc.NewGetInvoicesInternalServerError().WithPayload(errResponse)
	}

	return rpc.NewGetInvoicesOK().WithPayload(qResp)
}

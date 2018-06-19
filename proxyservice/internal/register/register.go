package register

import (
	APIs "ProxyService/proxyservice/internal/api"
	"ProxyService/proxyservice/internal/queue"
	"ProxyService/proxyservice/openapi/autogen/proxyservice/server/operations"
	"ProxyService/proxyservice/openapi/autogen/proxyservice/server/operations/rpc"
)

// Router Holds all micro-service apis
type Router struct {
	RPCApi     *APIs.RPCApi
	SwaggerAPI *operations.ProxyServiceAPI
}

// NewRouter Router constructor
func NewRouter(swaggerAPI *operations.ProxyServiceAPI, queue *queue.Queue) *Router {
	return &Router{
		RPCApi:     APIs.NewRPCApi(queue),
		SwaggerAPI: swaggerAPI,
	}
}

// RegisterRoutes Registers all routes along with their handlers
func (router *Router) RegisterRoutes() {
	router.SwaggerAPI.RPCGetClientsHandler = rpc.GetClientsHandlerFunc(router.RPCApi.GetClients)
	router.SwaggerAPI.RPCGetInvoicesHandler = rpc.GetInvoicesHandlerFunc(router.RPCApi.GetInvoices)
}

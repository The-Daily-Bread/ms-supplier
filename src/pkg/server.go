package pkg

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/rpc"
	rpc_json "github.com/gorilla/rpc/json"
)

type Reply struct {
	Message string `json:"message"`
}

type Item struct {
	Name     string  `json:"name"`
	Price    float32 `json:"price"`
	Quantity int     `json:"quantity"`
}

type SupplierRequest struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Itens   []Item `json:"itens"`
}

type JSONServer struct{}
type Args struct {
	SupplierRequest
}

const url = "mongodb://root:password@localhost:27013/"

func (t *JSONServer) CreateSupplierRequest(r *http.Request, args *Args, reply *Reply) error {
	supplier := SupplierRequest{
		Name:    args.Name,
		Address: args.Address,
		Itens:   args.Itens,
	}

	client, err := ConnectMongoDb(url)
	if err != nil {
		return err
	}

	defer client.Disconnect(context.TODO())

	collection := client.Database("ms-supplier").Collection("supplier_request")

	_, err = collection.InsertOne(context.TODO(), supplier)
	if err != nil {
		return err
	}

	reply.Message = "Supplier created: " + supplier.Name

	return nil
}

func Start() {
	s := rpc.NewServer()
	s.RegisterCodec(rpc_json.NewCodec(), "application/json")
	s.RegisterService(new(JSONServer), "")

	r := mux.NewRouter()
	r.Handle("/rpc", s)

	http.ListenAndServe(":3003", r)
}

package pkg

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/rpc"
	rpc_json "github.com/gorilla/rpc/json"
)

type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data";omitempty`
}

type Supplier struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type JSONServer struct{}
type Args struct {
	Name string
}

func (t *JSONServer) CreateSupplierRequest(r *http.Request, args *Args, reply *Response) error {
	var supplier Supplier

	err := json.NewDecoder(r.Body).Decode(&supplier)
	if err != nil {
		*reply = Response{
			Message: err.Error(),
		}
	}

	*reply = Response{
		Message: "success",
		Data:    supplier,
	}

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

package transfer_microservice

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	"github.com/gorilla/mux"

	httptransport "github.com/go-kit/kit/transport/http"
)

var (
	ErrBadRouting = errors.New("inconsistent mapping between route and handler (programmer error)")
)

func MakeHTTPhander(s TransferService, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	e := MakeTransferEndpoints(s)
	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		httptransport.ServerErrorEncoder(encodeError),
	}

	r.Methods("GET").Path("/transfer/").Handler(httptransport.NewServer(
		e.GetTransferListEndpoint,
		decodeGetTransferList,
		encodeResponse,
		options...,
	))

	r.Methods("GET").Path("/transfer/waiting").Handler(httptransport.NewServer(
		e.GetWaitingTransferEndpoint,
		decodeGetWaitingTransferRequest,
		encodeResponse,
		options...,
	))

	r.Methods("POST").Path("/transfer/").Handler(httptransport.NewServer(
		e.CreateEndpoint,
		decodeCreateRequest,
		encodeResponse,
		options...,
	))

	r.Methods("POST").Path("/transfer/pay").Handler(httptransport.NewServer(
		e.PostTransferStatusEndpoint,
		decodePostTransferStatusEndpoint,
		encodeResponse,
		options...,
	))
	return r
}

type errorer interface {
	error() error
}

func decodePostTransferStatusEndpoint(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req PostTransferStatusRequest
	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		return nil, e
	}
	return req, nil
}

func decodeCreateRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req CreateRequest
	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		return nil, e
	}
	return req, nil
}

func decodeGetWaitingTransferRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req GetWaitingTransferRequest
	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		return nil, e
	}
	return req, nil
}

func decodeGetTransferList(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req GetTransferListRequest
	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		return nil, e
	}
	return req, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if e, ok := response.(errorer); ok && e.error() != nil {
		// Not a Go kit transport error, but a business-logic error.
		// Provide those as HTTP errors.
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(codeFrom(err))
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func codeFrom(err error) int {
	switch err {
	case ErrNotFound:
		return http.StatusNotFound
	case ErrNotAnId, ErrNotFound:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

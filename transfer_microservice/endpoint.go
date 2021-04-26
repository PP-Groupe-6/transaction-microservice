package transfer_microservice

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/endpoint"
)

const (
	PENDING = 0
	PAID    = 1
	EXPIRED = 2
)

type TransferEndpoints struct {
	GetTransferListEndpoint    endpoint.Endpoint
	GetWaitingTransferEndpoint endpoint.Endpoint
	CreateEndpoint             endpoint.Endpoint
	PostTransferStatusEndpoint endpoint.Endpoint
}

func MakeTransferEndpoints(s TransferService) TransferEndpoints {
	return TransferEndpoints{
		GetTransferListEndpoint:    MakeGetTransferListEndpoint(s),
		GetWaitingTransferEndpoint: MakeGetWaitingTransferEndpoint(s),
		CreateEndpoint:             MakeCreateEndpoint(s),
		PostTransferStatusEndpoint: MakePostTransferStatusEndpoint(s),
	}
}

type GetTransferListRequest struct {
	ClientID string
}

type FormatedTransfer struct {
	Type     string  `json:"type"`
	Role     string  `json:"role"`
	FullName string  `json:"name"`
	Amount   float64 `json:"transactionAmount"`
	Date     string  `json:"transactionDate"`
}

type GetTransferListResponse struct {
	Transfers []FormatedTransfer `json:"transfers"`
}

func MakeGetTransferListEndpoint(s TransferService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetTransferListRequest)
		transfers, err := s.GetTransferList(ctx, req.ClientID)
		if err != nil {
			return nil, err
		}
		accountInfo, err := s.GetAccountInformation(ctx, req.ClientID)
		if err != nil {
			return nil, err
		}
		formatedName := accountInfo.Surname + " " + accountInfo.Name
		response := make([]FormatedTransfer, 0)
		for _, transfer := range transfers {
			response = append(response, FormatedTransfer{
				Type:     "transfer",
				Amount:   transfer.Amount,
				FullName: formatedName,
				Date:     transfer.ExecutionDate,
			})

			if transfer.AccountPayerId == req.ClientID {
				response[len(response)-1].Role = "payer"
			} else if transfer.AccountReceiverId == req.ClientID {
				response[len(response)-1].Role = "receiver"
			}
		}
		return GetTransferListResponse{response}, err

	}
}

type GetWaitingTransferRequest struct {
	ClientID string
}

type FormatedWaitingTransfer struct {
	ID               string  `json:"transferId"`
	Mail             string  `json:"mailAdressTransferPayer"`
	Amount           float64 `json:"transferAmount"`
	ExecutionDate    string  `json:"executionTransferDate"`
	ReceiverQuestion string  `json:"receiverQuestion"`
	ReceiverAnswer   string  `json:"receiverAnswer"`
}

type GetWaitingTransferListResponse struct {
	Transfers []FormatedWaitingTransfer `json:"transfers"`
}

func MakeGetWaitingTransferEndpoint(s TransferService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetWaitingTransferRequest)
		transfers, err := s.GetTransferList(ctx, req.ClientID)
		if err != nil {
			return nil, err
		}
		accountInfo, err := s.GetAccountInformation(ctx, req.ClientID)
		if err != nil {
			return nil, err
		}
		response := make([]FormatedWaitingTransfer, 0)

		for _, transfer := range transfers {
			response = append(response, FormatedWaitingTransfer{
				ID:               transfer.ID,
				Mail:             accountInfo.Mail,
				Amount:           transfer.Amount,
				ExecutionDate:    transfer.ExecutionDate,
				ReceiverQuestion: transfer.ReceiverQuestion,
				ReceiverAnswer:   transfer.ReceiverAnswer,
			})
		}
		return GetWaitingTransferListResponse{response}, err

	}
}

type CreateRequest struct {
	MailAdressTransferPayer    string
	MailAdressTransferReceiver string
	TransferAmount             float64
	TransferType               string
	ReceiverQuestion           string
	ReceiverAnswer             string
	ExecutionTransferDate      string
}

type CreateResponse struct {
	Type                        string `json:"transfer_type,omitempty"`
	Amount                      string `json:"transfer_amount,omitempty"`
	EmailAdressTransferPayer    string `json:"transfer_payer_mail,omitempty"`
	EmailAdressTransferReceiver string `json:"transfer_receiver_mail,omitempty"`
	ReceiverQuestion            string `json:"receiver_question,omitempty"`
	ReceiverAnswer              string `json:"receiver_answer,omitempty"`
	ExecutionTransferDate       string `json:"executed_transfer_date,omitempty"`
}

func MakeCreateEndpoint(s TransferService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateRequest)
		fmt.Println(req.ExecutionTransferDate)
		idPayer, err := s.GetIdFromMail(ctx, req.MailAdressTransferPayer)
		if err != nil {
			fmt.Println("Payer ID not found")
			return nil, err
		}
		idReceiver, err := s.GetIdFromMail(ctx, req.MailAdressTransferReceiver)
		if err != nil {
			fmt.Print("Reciever ID not found")
			return nil, err
		}

		toAdd := Transfer{
			ID:                "",
			Type:              req.TransferType,
			State:             0,
			Amount:            req.TransferAmount,
			AccountPayerId:    idPayer,
			AccountReceiverId: idReceiver,
			ReceiverQuestion:  req.ReceiverQuestion,
			ReceiverAnswer:    req.ReceiverAnswer,
			ExecutionDate:     req.ExecutionTransferDate,
		}

		transfer, err := s.Create(ctx, toAdd)
		if (err == nil && transfer != Transfer{}) {
			return CreateResponse{
				transfer.Type,
				fmt.Sprint(transfer.Amount),
				req.MailAdressTransferPayer,
				req.MailAdressTransferReceiver,
				transfer.ReceiverQuestion,
				transfer.ReceiverAnswer,
				transfer.ExecutionDate,
			}, nil
		} else {
			return CreateResponse{}, err
		}
	}
}

type PostTransferStatusRequest struct {
	ID string `json:"transfer_id"`
}

type PostTransferStatusResponse struct {
	Done bool `json:"done"`
}

func MakePostTransferStatusEndpoint(s TransferService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(PostTransferStatusRequest)
		res, err := s.PostTransferStatus(ctx, req.ID)

		if err == nil && res {
			return PostTransferStatusResponse{res}, nil
		} else {
			return PostTransferStatusResponse{res}, err
		}

	}
}

func StateToString(stateID int) string {
	switch stateID {
	case PENDING:
		return "Pending"
	case PAID:
		return "Paid"
	case EXPIRED:
		return "Expired"
	}
	return ""
}

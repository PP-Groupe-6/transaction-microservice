package transfer_microservice

import (
	"context"
	"errors"
	"fmt"
	"strconv"
)

type TransferService interface {
	Create(ctx context.Context, transfer Transfer) (Transfer, error)
	Read(ctx context.Context, id string) (Transfer, error)
	Update(ctx context.Context, id string, transfer Transfer) (Transfer, error)
	Delete(ctx context.Context, id string) error
	GetWaitingTransfer(ctx context.Context, id string) ([]*Transfer, error)
	GetTransferList(ctx context.Context, id string) ([]*Transfer, error)
	UpdateTransferStatus(ctx context.Context, id string) error
}

var (
	ErrNotAnId         = errors.New("not an ID")
	ErrNotFound        = errors.New("transaction not found")
	ErrNoTransfer      = errors.New("transfer field is empty")
	ErrNoUpdate        = errors.New("could not update transfer")
	ErrNoDb            = errors.New("could not access database")
	ErrAlreadyExist    = errors.New("transfer id already exists")
	ErrNoInsert        = errors.New("insert did not go through")
	ErrInconsistentIDs = errors.New("could not access database")
	ErrNotEnoughMoney  = errors.New("payer account has not enough money")
)

type transferService struct {
	DbInfos dbConnexionInfo
}

func NewTransferService(dbinfos dbConnexionInfo) TransferService {
	return &transferService{
		DbInfos: dbinfos,
	}
}

func (s *transferService) UpdateTransferStatus(ctx context.Context, id string) error {
	transfer, err := s.Read(ctx, id)
	if err != nil {
		return err
	}
	db := GetDbConnexion(s.DbInfos)

	fmt.Println(transfer)
	var amount string
	err = db.Select(&amount, "SELECT account_amount FROM account WHERE client_id='"+transfer.AccountPayerId+"')")

	if err != nil {
		return err
	}

	fmt.Println(amount)

	if s, _ := strconv.ParseFloat(amount, 64); s < transfer.Amount {
		return ErrNotEnoughMoney
	}

	tx := db.MustBegin()
	tx.MustExec("UPDATE transfer SET transfer_type = '"+transfer.Type+"', transfer_state="+fmt.Sprint(transfer.State)+", transfer_amount ="+fmt.Sprint(transfer.Amount)+", account_transfer_payer_id = '"+transfer.AccountPayerId+"', account_transfer_receiver_id = '"+transfer.AccountReceiverId+"', receiver_question = '"+transfer.ReceiverQuestion+"', receiver_answer = '"+transfer.ReceiverAnswer+"', scheduled_transfer_date = '"+transfer.ScheduledDate+"', executed_transfer_date = '"+transfer.ExecutedDate+"' WHERE transfer_id=$1", id)
	tx.Commit()
	db.Close()

	return nil
}

func (s *transferService) GetTransferList(ctx context.Context, id string) ([]*Transfer, error) {
	db := GetDbConnexion(s.DbInfos)
	transfers := make([]*Transfer, 0)
	rows, err := db.Queryx("SELECT * FROM transfer WHERE account_transfer_payer_id=$1 OR account_transfer_receiver_id=$1", id)

	for rows.Next() {
		var t Transfer
		if err = rows.StructScan(&t); err != nil {
			return nil, err
		}
		transfers = append(transfers, &t)
	}
	if err != nil {
		return nil, err
	}
	return transfers, err
}

func (s *transferService) GetWaitingTransfer(ctx context.Context, id string) ([]*Transfer, error) {
	db := GetDbConnexion(s.DbInfos)
	transfers := make([]*Transfer, 0)

	rows, err := db.Queryx("SELECT * FROM transfer WHERE account_transfer_receiver_id=$1 AND transfer_state=0", id)

	for rows.Next() {
		var t Transfer
		if err = rows.StructScan(&t); err != nil {
			return nil, err
		}
		transfers = append(transfers, &t)
	}
	if err != nil {
		return nil, err
	}
	return transfers, err
}

func (s *transferService) Create(ctx context.Context, transfer Transfer) (Transfer, error) {
	if (transfer == Transfer{}) {
		return Transfer{}, ErrNoTransfer
	}

	if testID, _ := s.Read(ctx, transfer.ID); (testID != Transfer{}) {
		return Transfer{}, ErrAlreadyExist
	}
	db := GetDbConnexion(s.DbInfos)

	//validations

	tx := db.MustBegin()
	res := tx.MustExec("INSERT INTO transfer VALUES ('" + transfer.ID + "','" + transfer.Type + "'," + fmt.Sprint(transfer.State) + "," + fmt.Sprint(transfer.Amount) + ",'" + transfer.AccountPayerId + "','" + transfer.AccountReceiverId + "','" + transfer.ReceiverQuestion + "','" + transfer.ReceiverAnswer + "','" + transfer.ScheduledDate + "','" + transfer.ExecutedDate + "')")
	tx.Commit()
	db.Close()

	if nRows, err := res.RowsAffected(); nRows != 1 || err != nil {
		if err != nil {
			return Transfer{}, err
		}
		return Transfer{}, ErrNoInsert
	}

	return s.Read(ctx, transfer.ID)

}

func (s *transferService) Read(ctx context.Context, id string) (Transfer, error) {
	db := GetDbConnexion(s.DbInfos)

	res := Transfer{}
	err := db.Get(&res, "SELECT * FROM transfer WHERE transfer_id=$1", id)

	if err != nil {
		return Transfer{}, err
	}

	return res, nil
}

func (s *transferService) Update(ctx context.Context, id string, transfer Transfer) (Transfer, error) {
	if (transfer == Transfer{}) {
		return Transfer{}, ErrNoTransfer
	}

	if testID, _ := s.Read(ctx, id); (testID == Transfer{}) {
		return Transfer{}, ErrNotFound
	}

	db := GetDbConnexion(s.DbInfos)
	tx := db.MustBegin()
	res := tx.MustExec("UPDATE transfer SET transfer_type = '"+transfer.Type+"', transfer_state="+fmt.Sprint(transfer.State)+", transfer_amount ="+fmt.Sprint(transfer.Amount)+", account_transfer_payer_id = '"+transfer.AccountPayerId+"', account_transfer_receiver_id = '"+transfer.AccountReceiverId+"', receiver_question = '"+transfer.ReceiverQuestion+"', receiver_answer = '"+transfer.ReceiverAnswer+"', scheduled_transfer_date = '"+transfer.ScheduledDate+"', executed_transfer_date = '"+transfer.ExecutedDate+"' WHERE transfer_id=$1", id)
	tx.Commit()
	db.Close()

	if nRows, err := res.RowsAffected(); nRows != 1 || err != nil {
		if err != nil {
			return Transfer{}, err
		}
		return Transfer{}, ErrNoInsert
	}

	return s.Read(ctx, transfer.ID)
}

func (s *transferService) Delete(ctx context.Context, id string) error {

	if testID, _ := s.Read(ctx, id); (testID == Transfer{}) {
		return ErrNotFound
	}
	db := GetDbConnexion(s.DbInfos)
	tx := db.MustBegin()
	res := tx.MustExec("DELETE FROM transfer WHERE transfer_id=$1", id)

	if nRows, err := res.RowsAffected(); nRows != 1 || err != nil {
		if err != nil {
			return err
		}
	}
	tx.Commit()
	db.Close()

	return nil
}

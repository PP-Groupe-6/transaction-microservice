package transfer_microservice

import (
	"context"
	"errors"
	"fmt"

	"github.com/rs/xid"
)

type TransferService interface {
	Create(ctx context.Context, transfer Transfer) (Transfer, error)
	Read(ctx context.Context, id string) (Transfer, error)
	Update(ctx context.Context, id string, transfer Transfer) (Transfer, error)
	Delete(ctx context.Context, id string) error
	GetWaitingTransfer(ctx context.Context, id string) ([]*Transfer, error)
	GetTransferList(ctx context.Context, id string) ([]*Transfer, error)
	UpdateTransferStatus(ctx context.Context, id string) error
	GetAccountInformation(ctx context.Context, id string) (AccountInfo, error)
	GetIdFromMail(ctx context.Context, mail string) (string, error)
	PostTransferStatus(ctx context.Context, id string) (bool, error)
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
	DbInfos DbConnexionInfo
}

func NewTransferService(dbinfos DbConnexionInfo) TransferService {
	return &transferService{
		DbInfos: dbinfos,
	}
}
func (s *transferService) PostTransferStatus(ctx context.Context, id string) (bool, error) {
	if id == "" {
		return false, ErrNotAnId
	}

	TransferToPay, errorR := s.Read(ctx, id)

	if (TransferToPay == Transfer{} && errorR != nil) {
		return false, ErrNotFound
	}

	db := GetDbConnexion(s.DbInfos)

	// Dans un premier temps on récupère le solde du payeur
	payerBalance := float64(0.0)

	fmt.Println(TransferToPay.AccountPayerId)
	errPB := db.Get(&payerBalance, "SELECT account_amount FROM account WHERE client_id=$1", TransferToPay.AccountPayerId)

	// On récupère ensuite le solde du receveur
	recieverBalance := float64(0.0)
	errRB := db.Get(&recieverBalance, "SELECT account_amount FROM account WHERE client_id=$1", TransferToPay.AccountReceiverId)

	if errPB != nil {
		fmt.Println("Payer balance error")
		return false, ErrNotFound
	}
	if errRB != nil {
		fmt.Println("Reciever balance error")
		return false, ErrNotFound
	}

	// On regarde si le payeur a les fonds pour payer la facture
	if payerBalance < TransferToPay.Amount {
		return false, ErrNotEnoughMoney
	}

	tx := db.MustBegin()
	// On mets à jour le solde du payeur
	resPayer := tx.MustExec("UPDATE account SET account_amount = '"+fmt.Sprint(payerBalance-TransferToPay.Amount)+"' WHERE client_id=$1", TransferToPay.AccountPayerId)

	if rows, errUpdate := resPayer.RowsAffected(); rows != 1 {
		tx.Rollback()
		return false, errUpdate
	}

	// On mets à jour le solde du receveur
	resReciever := tx.MustExec("UPDATE account SET account_amount = '"+fmt.Sprint(recieverBalance+TransferToPay.Amount)+"' WHERE client_id=$1", TransferToPay.AccountPayerId)
	if rows, errUpdate := resReciever.RowsAffected(); rows != 1 {
		tx.Rollback()
		return false, errUpdate
	}

	//On change l'état de la facture a payer
	resInvoice := tx.MustExec("UPDATE transfer SET transfer_state = '"+fmt.Sprint(PAID)+"' WHERE transfer_id=$1", TransferToPay.ID)
	if rows, errUpdate := resInvoice.RowsAffected(); rows != 1 {
		tx.Rollback()
		return false, errUpdate
	}

	tx.Commit()
	db.Close()

	return true, nil
}

func (s *transferService) GetIdFromMail(ctx context.Context, mail string) (string, error) {
	db := GetDbConnexion(s.DbInfos)

	var res string

	err := db.Get(&res, "SELECT client_id FROM account where mail_adress=$1", mail)
	if err != nil {
		return "", ErrNotFound
	}

	return res, err
}

func (s *transferService) GetAccountInformation(ctx context.Context, id string) (AccountInfo, error) {
	db := GetDbConnexion(s.DbInfos)

	res := AccountInfo{}
	err := db.Get(&res, "SELECT name, surname, mail_adress, account_amount FROM account where client_id=$1", id)

	if err != nil {
		return AccountInfo{}, err
	}
	return res, err
}

func (s *transferService) UpdateTransferStatus(ctx context.Context, id string) error {
	transfer, err := s.Read(ctx, id)
	if err != nil {
		return err
	}
	db := GetDbConnexion(s.DbInfos)

	tx := db.MustBegin()
	tx.MustExec("UPDATE transfer SET transfer_type = '"+transfer.Type+"', transfer_state="+fmt.Sprint(transfer.State)+", transfer_amount ="+fmt.Sprint(transfer.Amount)+", account_transfer_payer_id = '"+transfer.AccountPayerId+"', account_transfer_receiver_id = '"+transfer.AccountReceiverId+"', receiver_question = '"+transfer.ReceiverQuestion+"', receiver_answer = '"+transfer.ReceiverAnswer+"', executed_transfer_date = '"+transfer.ExecutionDate+"' WHERE transfer_id=$1", id)
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
	id := xid.New()
	db := GetDbConnexion(s.DbInfos)

	//validations

	tx := db.MustBegin()
	res := tx.MustExec("INSERT INTO transfer VALUES ('" + id.String() + "','" + transfer.Type + "'," + fmt.Sprint(transfer.State) + "," + fmt.Sprint(transfer.Amount) + ",'" + transfer.AccountPayerId + "','" + transfer.AccountReceiverId + "','" + transfer.ReceiverQuestion + "','" + transfer.ReceiverAnswer + "','" + transfer.ExecutionDate + "')")
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
	res := tx.MustExec("UPDATE transfer SET transfer_type = '"+transfer.Type+"', transfer_state="+fmt.Sprint(transfer.State)+", transfer_amount ="+fmt.Sprint(transfer.Amount)+", account_transfer_payer_id = '"+transfer.AccountPayerId+"', account_transfer_receiver_id = '"+transfer.AccountReceiverId+"', receiver_question = '"+transfer.ReceiverQuestion+"', receiver_answer = '"+transfer.ReceiverAnswer+"', executed_transfer_date = '"+transfer.ExecutionDate+"' WHERE transfer_id=$1", id)
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

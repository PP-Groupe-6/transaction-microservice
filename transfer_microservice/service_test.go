package transfer_microservice

import (
	"context"
	"testing"
)

type TestData struct {
	s             TransferService
	mockTransfer  Transfer
	otherTransfer Transfer
}

func NewTestData() TestData {
	info := dbConnexionInfo{
		"postgre://",
		"5432",
		"prix_banque_test",
		"admin",
		"secret",
	}

	s := NewTransferService(info)

	mockTransfer := Transfer{
		"QItqpLGDZJHu4JqGcvMnCBNUZJSOPpY1A3GK6bkNjDpEkyKc1nQj22SEMrq7VR146ApOBwP6X6NAfTtG9yiORGx0pOt5diWDMAe940JGqqiFeQFov2xoKcX7shkee5UEcLnXhhkijFTGcCMrfXcs4QoXUKLLSwa91a0xjfZz9lvag2pyI78R1GQZ1HqJp9yUbwgsmQeou7fm1C2WFrqobtwVI2rPPxU9wdeEFimAPOTuAZewhUTp236SP8q",
		"instant",
		0,
		667,
		"sIowRDsqanK3vj0jfRVn1i8yLrmJfu93qDDlZwkeHFl4td0W2czjJbutqwibI8iaQJ7skSHtLpWHUtfN7gFQ0f40e6J1Fie4LeuRrmLHkxfpr6bv5VOYpwGvDyoux7Zus0fw2R2IRWEr3CqKtrohdX8t9pf37I17WoSVFg83hrb18BoKD3h989i3I36GAjXGLyEWbj6RsD6lt5TEQOjwJEZDZTeBOUOq0fNOUFmEW47cEgQ2R4DvIj5AN2iPDsv",
		"fErnq0RHXlGBI6DuQ88O3T6BCIizTwgt1YtNte1lAUE1uqoJOVDxHUihPSTTf57GVORuB8XFT1f8lUASGP8p0Fzj69wGDOv1tzsnwSlHbPp4M2fggbiNItk0w10E7Ro3sZ0V77osOzXU43pLHZ53gFLDOrG8NVzUTr0FM33ySDa5f53KTJ7AUfTujnbiVwiwIWWCS10YBOKcMqGvJ5s48AkThnzqfSCIRM2Omh0xeJvn4RSYSROfsi8ol3iqbPa",
		"Ça va, tout va bien ?",
		"Oui",
		"2021-04-29T00:00:00Z",
		"2021-04-29T00:00:00Z",
	}

	otherTransfer := Transfer{
		"QItqpLGDZJHu4JqGcvMnCBNUZJSOPpY1A3GK6bkNjDpEkyKc1nQj22SEMrq7VR146ApOBwP6X6NAfTtG9yiORGx0pOt5diWDMAe940JGqqiFeQFov2xoKcX7shkee5UEcLnXhhkijFTGcCMrfXcs4QoXUKLLSwa91a0xjfZz9lvag2pyI78R1GQZ1HqJp9yUbwgsmQeou7fm1C2WFrqobtwVI2rPPxU9wdeEFimAPOTuAZewhUTp236SP8q",
		"scheduled",
		0,
		50,
		"7xZnb9WK362TUHQkkLyCAnaaLiF5b55OQX77nRyh4kUGuFq17z3Cn4LKfKN2sD108L79knYWu8O5VvMpq5ei5beoZsOJq0qtj2fBl7R1kc6UdNHAcDAnpWvklEyhk9u39hGzDUx7dqCX9Rd1mEDMvhrFq5Dt5DDzAUWI6Sr1z9LVVSeu4T8gOZSt9EFxAX4OWLxAVKK6PNv3D77SOYunRk5CUggH9GYWjDJ8O1C2lUICOKjDd4QRyyK7Ovcs9Dh",
		"sIowRDsqanK3vj0jfRVn1i8yLrmJfu93qDDlZwkeHFl4td0W2czjJbutqwibI8iaQJ7skSHtLpWHUtfN7gFQ0f40e6J1Fie4LeuRrmLHkxfpr6bv5VOYpwGvDyoux7Zus0fw2R2IRWEr3CqKtrohdX8t9pf37I17WoSVFg83hrb18BoKD3h989i3I36GAjXGLyEWbj6RsD6lt5TEQOjwJEZDZTeBOUOq0fNOUFmEW47cEgQ2R4DvIj5AN2iPDsv",
		"Ça va, tout va bien ?",
		"NON",
		"2021-04-29T00:00:00Z",
		"2021-04-29T00:00:00Z",
	}

	return TestData{
		s,
		mockTransfer,
		otherTransfer,
	}
}
func TestDelete(t *testing.T) {
	testData := NewTestData()

	errEmptyID := testData.s.Delete(context.TODO(), "")
	if errEmptyID == nil {
		t.Errorf("Passed an empty ID, method should have raised an error")
	}

	errInvalidID := testData.s.Delete(context.TODO(), "lmao")
	if errInvalidID == nil {
		t.Errorf("Passed wrong ID, should have raised an error")
	}

	err := testData.s.Delete(context.TODO(), testData.mockTransfer.ID)
	if err != nil {
		t.Errorf("Passed a valid ID, should not have raised an error")
	}
}

func TestGetTransferList(t *testing.T) {
	testData := NewTestData()
	result, err := testData.s.GetWaitingTransfer(context.TODO(), "sIowRDsqanK3vj0jfRVn1i8yLrmJfu93qDDlZwkeHFl4td0W2czjJbutqwibI8iaQJ7skSHtLpWHUtfN7gFQ0f40e6J1Fie4LeuRrmLHkxfpr6bv5VOYpwGvDyoux7Zus0fw2R2IRWEr3CqKtrohdX8t9pf37I17WoSVFg83hrb18BoKD3h989i3I36GAjXGLyEWbj6RsD6lt5TEQOjwJEZDZTeBOUOq0fNOUFmEW47cEgQ2R4DvIj5AN2iPDsv")

	if err != nil {
		t.Errorf("Method should not have failed : " + err.Error())
	}
	for _, transfer := range result {
		t.Log(transfer)
	}
}
func TestCreate(t *testing.T) {
	testData := NewTestData()

	_, err := testData.s.Create(context.TODO(), Transfer{})

	if err == nil {
		t.Errorf("Passed empty transfer field to create function, should have raised an error")
	}

	result, err := testData.s.Create(context.TODO(), testData.mockTransfer)

	if err != nil {
		t.Errorf("Valid transfer, method shound not raise an error : " + err.Error())

	}

	if result != testData.mockTransfer {
		t.Errorf("Returned transfer is not the same as the one created : " + testData.mockTransfer.ID + " got : " + result.ID)
	}
}

func TestRead(t *testing.T) {
	testData := NewTestData()

	_, err := testData.s.Read(context.TODO(), "")

	if err == nil {
		t.Errorf("Passed empty transfer id, should have raised an error")
	}

	result, err := testData.s.Read(context.TODO(), testData.mockTransfer.ID)

	if err != nil {
		t.Errorf("Valid ID, method should not fail : " + err.Error())
	}

	if result.ID != testData.mockTransfer.ID {
		t.Errorf("Returned transfer is not the same as the one specified")
	}

}

func TestUpdate(t *testing.T) {
	testData := NewTestData()

	_, errEmptyID := testData.s.Update(context.TODO(), "", testData.mockTransfer)

	if errEmptyID == nil {
		t.Errorf("Passed empty id, should have raised an error")
	}

	_, errEmptyTransfer := testData.s.Update(context.TODO(), testData.mockTransfer.ID, Transfer{})

	if errEmptyTransfer == nil {
		t.Errorf("Passed empty transfer, should have raised an error")
	}

	_, errInconsistentIDs := testData.s.Update(context.TODO(), "lmao", testData.mockTransfer)

	if errInconsistentIDs == nil {
		t.Errorf("Passed inconsistent IDs, should have raised an error ")
	}

	_, err := testData.s.Update(context.TODO(), testData.mockTransfer.ID, testData.otherTransfer)
	if err != nil {
		t.Errorf("Valid transfer ID, method should not have raised an error : " + err.Error())
	}

	dbResult, err := testData.s.Read(context.TODO(), testData.mockTransfer.ID)

	if err != nil {
		t.Errorf("Error during read")
	}

	if dbResult != testData.otherTransfer && dbResult == testData.mockTransfer {
		t.Errorf("Update did not work")
	}

	if dbResult != testData.otherTransfer && dbResult != testData.mockTransfer {
		t.Errorf("Fetched result is not the test transfer we asked for")
	}

}

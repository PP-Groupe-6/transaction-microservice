# transfer-microservice

## Comment lancer et compiler le microservice

Pour compiler, aller à la racine du projet et utiliser la commande 
```powershell
go build
```
Cette commande produit un exécutable qu'il suffit de lancer pour que le microservice soit actif.

## Comment accéder au microservice

Ce microservice se lance sur localhost:8001 par défaut. Pour en changer la configuration, modifiez le fichier main.go à la ligne 29 :
```go
err := http.ListenAndServe(":<port>", accountService.MakeHTTPHandler(service, logger))
```

Pour tester le microservice nous conseillons l'outil [Postman](https://www.postman.com) et [la collection fournie avec le microservice](https://github.com/PP-Groupe-6/transfer-microservice/blob/master/Transfer.postman_collection.json).

La liste des Url est la suivante :
| URL                     | Méthode           | Param (JSON dans le body) | Retour               |
| ----------------------- |:-----------------:| :------------------------:| :-------------------:|
| localhost:8001/transfer/  | GET             | {"ClientID": "\<ID\>"}      |{"transfers": [{"type": "\<type\>","role": "\<role\>","name": "\<fullname\>","transactionAmount": \<amount\>,"transactionDate": "\<date\>"}, ... ]}|
| localhost:8001/transfer/waiting   | GET     | {"ClientID": "\<ID\>"}      |{"transfers": [{"transferId": "\<transferId\>","mailAdressTransferPayer": "\<mailAdressTransferPayer\>","transferAmount": \<transferAmount\>,"executionTransferDate": "\<executionTransferDate\>","receiverQuestion": "\<receiverQuestion\>"},"receiverAnswer": "\<receiverAnswer\>"}, ... ]}  |
| localhost:8001/transfer/   | POST              | {"EmailAdressTransferPayer":"\<EmailAdressTransferPayer\>","EmailAdressTransferReceiver": "\<EmailAdressTransferReceiver\>","TransferAmount":\<TransferAmount\>,"TransferType":"\<TransferType\>","ReceiverQuestion":"\<ReceiverQuestion\>","ReceiverAnswer":"\<ReceiverAnswer\>","ExecutionTransferDate":"\<ExecutionTransferDate\>"}| |
| localhost:8001/transfer/pay  | POST             | {"transfer_id": "\<transfer_id\>"}      |{"done": \<bool\>}|

package main

import (
	"database/sql"
	"fmt"

	"github.com/Sans-arch/fc-walletcore/internal/database"
	"github.com/Sans-arch/fc-walletcore/internal/event"
	"github.com/Sans-arch/fc-walletcore/internal/usecase/create_account"
	"github.com/Sans-arch/fc-walletcore/internal/usecase/create_client"
	"github.com/Sans-arch/fc-walletcore/internal/usecase/create_transaction"
	"github.com/Sans-arch/fc-walletcore/pkg/events"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", "root", "root", "localhost", "3306", "wallet"))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	eventDispatcher := events.NewEventDispatcher()
	transactionCreatedEvent := event.NewTransactionCreated()
	// eventDispatcher.Register("TransactionCreated", handler)

	clientDb := database.NewClientDB(db)
	accountDb := database.NewAccountDB(db)
	transactionDb := database.TransactionDB(db)

	createClientUsecase := create_client.NewCreateClientUsecase(clientDb)
	createAccountUsecase := create_account.NewCreateAccountUsecase(accountDb, clientDb)
	createTransactionUsecase := create_transaction.NewTransactionUsecase(transactionDb, accountDb, eventDispatcher, transactionCreatedEvent)
}

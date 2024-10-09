package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Sans-arch/fc-walletcore/internal/database"
	"github.com/Sans-arch/fc-walletcore/internal/event"
	"github.com/Sans-arch/fc-walletcore/internal/usecase/create_account"
	"github.com/Sans-arch/fc-walletcore/internal/usecase/create_client"
	"github.com/Sans-arch/fc-walletcore/internal/usecase/create_transaction"
	"github.com/Sans-arch/fc-walletcore/internal/web"
	"github.com/Sans-arch/fc-walletcore/internal/web/webserver"
	"github.com/Sans-arch/fc-walletcore/pkg/events"
	"github.com/Sans-arch/fc-walletcore/pkg/uow"
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

	ctx := context.Background()
	uow := uow.NewUow(ctx, db)

	uow.Register("AccountDB", func(tx *sql.Tx) interface{} {
		return database.NewAccountDB(db)
	})

	uow.Register("TransactionDB", func(tx *sql.Tx) interface{} {
		return database.NewTransactionDB(db)
	})

	createClientUsecase := create_client.NewCreateClientUsecase(clientDb)
	createAccountUsecase := create_account.NewCreateAccountUsecase(accountDb, clientDb)
	createTransactionUsecase := create_transaction.NewTransactionUsecase(uow, eventDispatcher, transactionCreatedEvent)

	webserver := webserver.NewWebServer(":3000")

	clientHandler := web.NewWebClientHandler(*createClientUsecase)
	accountHandler := web.NewWebAccountHandler(*createAccountUsecase)
	transactionHandler := web.NewWebTransactionHandler(*createTransactionUsecase)

	webserver.AddHandler("/clients", clientHandler.CreateClient)
	webserver.AddHandler("/accounts", accountHandler.CreateAccount)
	webserver.AddHandler("/transactions", transactionHandler.CreateTransaction)

	webserver.Start()
}

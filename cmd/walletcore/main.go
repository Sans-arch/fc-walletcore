package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Sans-arch/fc-walletcore/internal/database"
	"github.com/Sans-arch/fc-walletcore/internal/event"
	"github.com/Sans-arch/fc-walletcore/internal/event/handler"
	"github.com/Sans-arch/fc-walletcore/internal/usecase/create_account"
	"github.com/Sans-arch/fc-walletcore/internal/usecase/create_client"
	"github.com/Sans-arch/fc-walletcore/internal/usecase/create_transaction"
	"github.com/Sans-arch/fc-walletcore/internal/web"
	"github.com/Sans-arch/fc-walletcore/internal/web/webserver"
	"github.com/Sans-arch/fc-walletcore/pkg/events"
	"github.com/Sans-arch/fc-walletcore/pkg/kafka"
	"github.com/Sans-arch/fc-walletcore/pkg/uow"
	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", "root", "root", "mysql", "3306", "wallet"))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	configMap := ckafka.ConfigMap{
		"bootstrap.servers": "kafka:29092",
		"group.id":          "wallet",
	}
	kafkaProducer := kafka.NewKafkaProducer(&configMap)

	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("TransactionCreated", handler.NewTransactionCreatedKafkaHandler(kafkaProducer))
	eventDispatcher.Register("BalanceUpdated", handler.NewUpdateBalanceKafkaHandler(kafkaProducer))
	transactionCreatedEvent := event.NewTransactionCreated()
	balanceUpdatedEvent := event.NewBalanceUpdated()

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
	createTransactionUsecase := create_transaction.NewTransactionUsecase(uow, eventDispatcher, transactionCreatedEvent, balanceUpdatedEvent)

	webserver := webserver.NewWebServer(":8080")

	clientHandler := web.NewWebClientHandler(*createClientUsecase)
	accountHandler := web.NewWebAccountHandler(*createAccountUsecase)
	transactionHandler := web.NewWebTransactionHandler(*createTransactionUsecase)

	webserver.AddHandler("/clients", clientHandler.CreateClient)
	webserver.AddHandler("/accounts", accountHandler.CreateAccount)
	webserver.AddHandler("/transactions", transactionHandler.CreateTransaction)

	fmt.Println("Server is running")
	webserver.Start()
}

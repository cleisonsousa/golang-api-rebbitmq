package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-lang-api/internal/order/infra/database"
	"github.com/go-lang-api/internal/order/usecase"
	"github.com/go-lang-api/pkg/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	db, err := sql.Open("sqlite3", "./orders.db")
	if err != nil {
		panic(err)
	}

	defer db.Close()
	re := database.NewOrderRepository(db)
	uc := usecase.CalculateFinalPriceUseCase{OrderRepository: re}

	ch, err := rabbitmq.OpenChannel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()
	out := make(chan amqp.Delivery)
	go rabbitmq.Consume(ch, out)
	// go kafka.Consume(ch, out)
	// go sqs.Consume(ch, out)

	// for msg := range out {
	// 	var inputDTO usecase.OrderInputDTO
	// 	err := json.Unmarshal(msg.Body, &inputDTO)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	outputDTO, err := uc.Execute(inputDTO)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	fmt.Println(outputDTO)
	// 	msg.Ack(false)

	// }

	qtdWorkers := 150
	for i:=1; i <= qtdWorkers; i++ {
		go worker(out, &uc, i)
	}

	http.HandleFunc("/total", func(w http.ResponseWriter, r *http.Request){
		getTotalUC := usecase.GetTotalUseCase{OrderRepository: re}
		total, err := getTotalUC.Execute()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
		json.NewEncoder(w).Encode(total)
	})

	http.ListenAndServe(":8080", nil) // Chama serve HTTP, ria uma thread
	// go http.ListenAndServe(":8080", nil)
}

func worker(deliveryMessage <-chan amqp.Delivery, uc *usecase.CalculateFinalPriceUseCase, workerID int) {
	println()
	for msg := range deliveryMessage {
		var inputDTO usecase.OrderInputDTO
		err := json.Unmarshal(msg.Body, &inputDTO)
		if err != nil {
			panic(err)
		}
		outputDTO, err := uc.Execute(inputDTO)
		if err != nil{
			panic(err)
		}
		fmt.Printf("worker %d has processed order %s\n", workerID, outputDTO.ID)
		msg.Ack(false)
		time.Sleep(1* time.Second)
	}
}

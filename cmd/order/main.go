package main

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/andrereliquias/gointensivo/internal/infra/database"
	"github.com/andrereliquias/gointensivo/internal/usecase"
	"github.com/andrereliquias/gointensivo/pkg/rabbitmq"
	_ "github.com/mattn/go-sqlite3"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Car struct {
	Model string
	Color string
}

// metodo
func (c Car) Start() {
	println(c.Model + " has been started")
}

func (c *Car) ChangeColor(color string) {
	c.Color = color
	println("new color: " + c.Color)
}

// funcao
func soma(x, y int) int {
	return x + y
}

func main() {

	db, err := sql.Open("sqlite3", "db.sqlite3")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	orderRepository := database.NewOrderRepository(db)

	uc := usecase.NewCalculateFinalPrice(orderRepository)

	// abre canal com rabbitmq
	ch, err := rabbitmq.OpenChannel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	// abre canal de comunicacao entre as threads
	msgRabbitmqChannel := make(chan amqp.Delivery)

	// thread pra ler os dados do rabbitmq e jogar no canal msgRabbitmqChannel
	go rabbitmq.Consume(ch, msgRabbitmqChannel) // processo continuo que trava a thread

	rabbitmqWorker(msgRabbitmqChannel, uc) // le os dados do canal para executar o useCase

}

func rabbitmqWorker(msgChann chan amqp.Delivery, uc *usecase.CalculateFinalPrice) {
	fmt.Println("Starting rabbitmq")

	for msg := range msgChann {
		var input usecase.OrderInput
		err := json.Unmarshal(msg.Body, &input)

		if err != nil {
			panic(err)
		}

		output, err := uc.Execute(input)

		if err != nil {
			panic(err)
		}

		msg.Ack(false)
		fmt.Println("Mensagem processada e salva no banco", output)

	}

}

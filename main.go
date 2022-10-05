package main

import (
	"fmt"

	"github.com/architecture-it/go-platform/AMQStream"
	"github.com/architecture-it/integracion-schemas-event-go/Onboarding/Events"
	TestEvents "github.com/architecture-it/integracion-schemas-event-go/Test"
)

func main() {
	config, err := AMQStream.AddKafka()

	if err != nil {
		fmt.Printf("Revento todo corran por sus vidas %v", err)
	}

	event2 := Events.Pedido{}
	event := TestEvents.KafkaDemo{}

	topics := []string{"KafkaDemoTest"}

	config.ToProducer(&event2, topics)
	config.ToProducer(&event, topics)

	ejecutar()
}

func ejecutar() {

	team := TestEvents.NewUnionNullTeam()
	boss := TestEvents.NewUnionNullString()

	boss.String = "pepe"

	team.Team.Boss = boss
	team.Team.Tl = nil
	team.Team.Members = []string{"NC", "NZ", "LL", "LO", "GZ"}

	evento := TestEvents.KafkaDemo{
		Id:               1,
		DidYouUnderstand: true,
		Attendance:       200,
		Time:             20221008101213,
		Title:            "Que se yo",
		Me: TestEvents.Person{
			Name:      "Lucas",
			Surname:   "lucero",
			OnSite:    true,
			Seniority: "Sr.",
			Team:      team,
		},
	}

	evento2 := Events.Pedido{
		Id:                      "1",
		NumeroDePedido:          23,
		CicloDelPedido:          "32",
		CodigoDeContratoInterno: 12,
		EstadoDelPedido:         "na",
		CuentaCorriente:         32,
		Cuando:                  "123",
	}

	AMQStream.To(&evento, "key-go")
	AMQStream.To(&evento2, "key-pedido")

}

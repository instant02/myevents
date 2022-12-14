package amqp

import (
	"encoding/json"
	"fmt"

	"github.com/martin-helmich/cloudnativego-backend/src/lib/msgqueue"
	"github.com/streadway/amqp"
)

type amqpEventEmitter struct{
	connection *amqp.Connection
}

func (a *amqpEventEmitter) setup()error{
	channel,err := a.connection.Channel()
	if err != nil{
		return err
	}

	defer channel.Close()

	return channel.ExchangeDeclare("events","topic",true,false,false,false,nil)
}

func NewAMQPEventEmitter(conn *amqp.Connection)(*amqpEventEmitter,error){
	emitter := amqpEventEmitter{
		connection: conn,
	}

	err := emitter.setup()
	if  err != nil {
		return nil,err
	}

	return &emitter,nil
}

func(a *amqpEventEmitter)Emit(event msgqueue.Event) error{
	channel,err := a.connection.Channel()
	if err != nil{
		return err
	}

	defer channel.Close()

	jsonBody,err := json.Marshal(event)
	if err != nil{
		return fmt.Errorf("could not JSON-serialize event: %s",err)
	}

	msg := amqp.Publishing{
		Headers: amqp.Table{"x-event-name":event.EventName()},
		ContentType: "application/json",
		Body: jsonBody,
	}

	return channel.Publish("events", event.EventName(),false,false,msg)
}
# amqp091-comp

This package can convert github.com/Azure/go-amqp (amqp 1.0) to github.com/rabbitmq/amqp091-go (amqp 0.9.1) messages and vice versa.
It is based on the 

#Usage: 
Q: When you should use this package?

A: When you have two clients communicating through RabbitMQ, one (Bob) is using amqp 1.0 but the other one (Alice) wants to use amqp 0.9.1.
    The second one would use this package to translate messages from amqp 1.0 format to amqp 0.9.1 and vice versa, while keeping the original message structure and data.

If Bob is the producer (amqp 1.0 sender), Bob will send a message (Azure/go-amqp) to RabbitMQ and Alice (amqp 0.9.1 consumer) would get a message (rabbitmq/amqp091-go).
 In order for Alice to decode the message into the original message Bob sent, Alice would use the ConvertTo10 method to convert the message.
 
If Alice is the producer (amqp 0.9.1 sender), Alice will send a message (rabbitmq/amqp091-go) to RabbitMQ and Bob (amqp 1.0 consumer) would get a message (Azure/go-amqp).
 In order to enable Bob to read the message, Alice would create a Azure/go-amqp message and encode it into a rabbitmq/amqp091-go message using the ConvertTo091, method before sending.

Q: Wh

#Note:
Although it has a similar purpose to RabbitMQ Shovel plugin (see https://www.rabbitmq.com/docs/shovel),
 this package only handles conversion of messages from amqp 1.0 format to amqp 0.9.1 format and vice versa.
 This enables sending messages using amqp 0.9.1 that can be read using amqp 1.0 and sending messages using amqp 1.0 that can be read using amqp 0.9.1.
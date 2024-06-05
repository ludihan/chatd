const WebSocket = require('ws')
const amqp = require('amqplib/callback_api')

const wss = new WebSocket.Server({ port: 3000 })

wss.on('connection', ws => {
  ws.on('message', message => {
    console.log(JSON.parse(message.toString()))
    amqp.connect('amqps://drydhblv:nI2ZVgy8acFxj73dR7s8tGFa4zJnENZ7@prawn.rmq.cloudamqp.com/drydhblv', function (error0, connection) {
      if (error0) {
        throw error0;
      }
      connection.createChannel(function (error1, channel) {
        if (error1) {
          throw error1;
        }
        var exchange = JSON.parse(message);

        channel.assertExchange(exchange, 'fanout', {
          durable: false,
          autoDelete: true
        });

        channel.assertQueue('', {
          exclusive: true,
          deletWhenUnused: true
        }, function (error2, q) {
          if (error2) {
            throw error2;
          }
          console.log(" [*] Waiting for messages in %s. To exit press CTRL+C", q.queue);
          channel.bindQueue(q.queue, exchange, '');

          channel.consume(q.queue, function (msg) {
            if (msg.content) {
              console.log(" [x] %s", msg.content.toString());
              
              let content = JSON.parse(msg.content.toString())
              ws.send(JSON.stringify({ type: 'chat', message: content }))
            }
          }, {
            noAck: true
          });
        });
      });
    });
  })
})
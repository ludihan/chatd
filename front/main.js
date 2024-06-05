document.addEventListener("DOMContentLoaded", function () {
    var joinButton = document.getElementById("join-user");
    var sendButton = document.getElementById("send-message");
    var messageInput = document.getElementById("message-input");
    var messagesContainer = document.querySelector(".messages");

    joinButton.addEventListener("click", function () {
        var usernameInput = document.getElementById("username");
        var username = usernameInput.value;

        if (username.trim() !== "") {
            var joinScreen = document.querySelector(".join-screen");
            var chatScreen = document.querySelector(".chat-screen");
            joinScreen.classList.remove("active");
            chatScreen.classList.add("active");
        } else {
            alert("Por favor, insira um nome de usuário.");
        }
    });

    sendButton.addEventListener("click", function () {
        var message = messageInput.value.trim();
        var username = document.getElementById("username").value.trim();
        var token = document.getElementById("token").value.trim();

        if (message !== "") {
            messageInput.value = "";

            var messageObject = {
                exchange: token,
                body: message,
                userId: username,
            };

            var jsonMessage = JSON.stringify(messageObject);
            // Você pode enviar o jsonMessage para onde precisar aqui
            fetch("http://localhost:8080/publish", {
                method: "POST",
                body: jsonMessage,
                headers: {
                    "Content-type": "application/json; charset=UTF-8",
                },
            }).then(res => {
                renderMessage("my", {
                    userId: username,
                    body: message
                })
            });
        } else {
            alert("Por favor, insira uma mensagem.");
        }
    });
});

const app = document.querySelector(".app");
var amqp = require('amqplib/callback_api');

function configAmqp() {
    amqp.connect('amqps://drydhblv:nI2ZVgy8acFxj73dR7s8tGFa4zJnENZ7@prawn.rmq.cloudamqp.com/drydhblv', function (error0, connection) {
        if (error0) {
            throw error0;
        }
        connection.createChannel(function (error1, channel) {
            if (error1) {
                throw error1;
            }
            var exchange = token;

            channel.assertExchange(exchange, 'fanout', {
                durable: false,
            });

            channel.assertQueue('', {
                exclusive: true
            }, function (error2, q) {
                if (error2) {
                    throw error2;
                }
                console.log(" [*] Waiting for messages in %s. To exit press CTRL+C", q.queue);
                channel.bindQueue(q.queue, exchange, '');

                channel.consume(q.queue, function (msg) {
                    if (msg.content) {
                        console.log(" [x] %s", msg.content.toString());
                    }
                }, {
                    noAck: true
                });
            });
        });
    });
}

function renderMessage(type, message) {
    let messageContainer = app.querySelector(".chat-screen .messages")
    if (type == "my") {
        let el = document.createElement("div")
        el.setAttribute("class", "message my-message")
        el.innerHTML = `
                <div>
                    <div class = "name">Você</div>
                    <div class = "name"> ${message.body}</div>
                </div>
            
            `
        messageContainer.appendChild(el)
    } else if (type == "other") {
        let el = document.createElement("div")
        el.setAttribute("class", "message other-message")
        el.innerHTML = `
                <div>
                    <div class = "name">${message.userId}</div>
                    <div class = "name"> ${message.body}</div>
                </div>
            
            `
        messageContainer.appendChild(el)
    } else if (type == "update") {
        let el = document.createElement("div")
        el.setAttribute("class", "update")
        el.innerText = message
        messageContainer.appendChild(el)
    }
    messageContainer.scrollTop = messageContainer.scrollHeight - messageContainer.clientHeight;
}

configAmqp()
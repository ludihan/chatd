document.addEventListener("DOMContentLoaded", function () {
    var joinButton = document.getElementById("join-user");
    var sendButton = document.getElementById("send-message");
    var messageInput = document.getElementById("message-input");
    var messagesContainer = document.querySelector(".messages");

    joinButton.addEventListener("click", function () {
        var usernameInput = document.getElementById("username");
        var username = usernameInput.value;
        var token = document.getElementById("token").value.trim();

        if (username.trim() !== "") {
            var joinScreen = document.querySelector(".join-screen");
            var chatScreen = document.querySelector(".chat-screen");
            joinScreen.classList.remove("active");
            chatScreen.classList.add("active");
        } else {
            alert("Por favor, insira um nome de usuário.");
        }

        var socket = new WebSocket("ws://localhost:3000")

        socket.onopen = function (e) {
            console.log("Conexão estabelecida")
            socket.send(JSON.stringify(token))
        }

        socket.onclose = function (event) {
            console.log('conexão fechada')
        }

        socket.onmessage = function (event) {
            let data = JSON.parse(event.data)
            switch (data.type) {
                case 'chat':
                    console.log(data.message.body)
                    // console.log(data.message.content)
                    renderMessage("other", {
                        userId: data.message.userId,
                        body: data.message.body
                    })
                    break
                default:
                    console.log("quebrou")
            }
        }

        socket.onerror = function (error) {
            console.log(`Erro ${error.message}`)
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
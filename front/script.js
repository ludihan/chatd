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
            console.log(jsonMessage)
            // Você pode enviar o jsonMessage para onde precisar aqui
            fetch("http://localhost:8080/publish", {
                method: "POST",
                body: jsonMessage,
                headers: {
                    "Content-type": "application/json; charset=UTF-8",
                },
            }).then(res => {
                console.log("Request sent:", res)
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
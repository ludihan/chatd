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
            });
        } else {
            alert("Por favor, insira uma mensagem.");
        }
    });
});

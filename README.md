# rabbitmq-wrapper
Um projeto de chat online com base em mensageria

O nome rabbitmq-""wrapper"" vem do fato de que qualquer mensagem precisa
passar por uma api antes, funcionando como um bloqueador de mensagens maliciosas
(ao contrário de entrar diretamente na fila do rabbitmq)

# Build
## Requisitos
- Go 1.22
- Node
## Comando:
(execute na pasta raiz do projeto)
```sh
go mod tidy && \
go build -o start-server ./server && \
go build -o chat ./tui-client && \
cd front && \
npm install && \
cd ..
```

# Como executar
## Servidor
Forneça um arquivo toml com as informações adequadas
```sh
./start-server server-config.toml
```

## Front
### Requisitos: servidor rodando
Execute o arquivo "server.js" na pasta "front" com o node
```sh
node server.js
```

## Cliente de terminal
### Requisitos: servidor rodando
Execute o arquivo fornecendo seu nome, nome da exchange, url do rabbitmq e url da api
```sh
./chat paulo xyz amqp://algumacoisa@algumacoisa http/localhost:8080/publish
```

# Como o servidor funciona?
O servidor fornece dois endpoints:
- "/": o frontend web (precisa da api e serviço do node)
- "/publish": a api wrapper do rabbitmq

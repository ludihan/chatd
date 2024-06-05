# rabbitmq-wrapper
Isto é um projeto de chat online com base em mensageria

O nome rabbitmq-""wrapper"" vem do fato de que qualquer mensagem precisa
passar por uma api antes que funciona como um bloqueador de mensagens maliciosas
(ao contrário de entrar direto na fila do rabbitmq diretamente)

# Como buildar
## Requisitos
- Go 1.22
## Comando:
```sh
go mod tidy && go build -o start-server ./server && go build -o chat ./tui-client
```

# Como executar
## Servidor
Forneça um arquivo toml com as informações adequadas
```sh
./start-server server-config.toml
```
## Cliente de terminal
Forneça seu nome, nome da exchange, url do rabbitmq e url da api

!!! O servidor precisa estar rodando !!!
```sh
./chat paulo 777 amqp://algumacoisa@algumacoisa http/localhost:8080/publish
```

# Como o servidor funciona?
O servidor fornece dois endpoints:
- "/": o frontend da coisa
- "/publish": a api wrapper do rabbitmq

Para acessar a aplicação acesse "localhost:8080" no seu browser (8080 é a porta padrão)

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
go build -o start-server ./server
```

# Como executar
Forneça um arquivo toml com as informações adequadas para o executável
## Comando:
```sh
./start-server server-config.toml
```

# Como usar
A aplicação fornece dois endpoints:
- "/": o frontend da coisa
- "/publish": a api wrapper do rabbitmq

Para acessar a aplicação acesse "localhost:8080" no seu browser (8080 é a porta padrão)

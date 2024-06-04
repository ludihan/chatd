# rabbitmq-wrapper
Este projeto é uma camada de abstração por cima do serviço de mensageria "rabbitmq".

Ele provê a capacidade de filtrar mensagens da fila com base no conteúdo

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
O resto de como a aplicação funciona se descreve por si só

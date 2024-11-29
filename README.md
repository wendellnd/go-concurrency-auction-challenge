# go-concurrency-auction-challenge

#### Fork

`https://github.com/devfullcycle/labs-auction-goexpert`

#### Execução

1. Execute o comando `docker compose up --build` para subir a API e o mongodb

2. Crie um novo leilão utilizando o arquivo [./api/auction.http](./api/auction.http) a requisição vai retornar um novo leilão aberto.
   **Exemplo:**

```
POST http://localhost:8080/auction
Content-Type: application/jsons

{
    "product_name": "Carro",
    "category": "Carro",
    "description": "Carro de testes",
    "condition": 1
}
```

**Resposta**

```
HTTP/1.1 201 Created
Content-Type: application/json; charset=utf-8
Date: Fri, 29 Nov 2024 22:05:38 GMT
Content-Length: 189
Connection: close

{
  "id": "6e4b3448-044e-43d4-9faa-6be516c02568",
  "product_name": "Carro",
  "category": "Carro",
  "description": "Carro de testes",
  "condition": 1,
  "status": 0,
  "timestamp": "2024-11-29T22:05:38.813149751Z"
}
```

3. Execute a requisição `GET /auction/:id`, utilizando o id gerado na primeira requisição.
   - O status inicial do leilão deve ser 0(Aberto).
   - Após 20 segundos o status será atualizado para 1(Completo/Fechado) automaticamente, utilizando os conceitos de concorrência.

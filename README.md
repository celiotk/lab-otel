# Lab Observabilidade & Open Telemetry
Este projeto implementa tracing distribuído entre dois serviços: um para validação de CEP e outro para consulta de temperatura.

## Configuração
* Configure o campo `WEATHER_API_KEY` no arquivo `.env` com a sua chave de acesso do [www.weatherapi.com](https://www.weatherapi.com/).

## Passos para execução
* Inicie os serviços com o comando:
  ```bash
  docker-compose up -d
  ```
* Aguarde a inicialização completa da aplicação.

## Como usar
 * Faça uma requisição POST para o serviço principal (serviço A):
### Exemplo de requisição:
  ```
  POST http://localhost:8181/temperature HTTP/1.1
  Content-Type: application/json

  {
    "cep": "01001000"
  }
  ```
  Substitua `01001000` pelo CEP do local que deseja consultar.
## Observando os traces
* Para visualizar o trace das requisições, acesse a interface do [Zipkin](http://localhost:9411/)
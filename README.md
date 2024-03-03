# O que é isso?

Este é um serviço recebe um CEP brasileiro e retorna temperaturas (celsius, fahrenheit e kelvin)

Nesta abordagem são executados dois serviços onde o primeiro recebe um cep e o válida este sendo válido requisita ao segundo a temperatura.

Para medirmos os tracings distrbuidos foi utilizado o OpenTelemetry.

# Como executá-lo?

Antes de executar o docker compose verifique se as portas 8080 e 9411 estão livres para não sofrer com conflitos

```bash
docker compose up
```

Execute a solicitação atraves do endpoint https://localhost/ utilizando o Verbo POST com um body {"cep" : "01153000"}

```bash
curl --request POST --url 'http://localhost:8080' -H "Content-Type: application/json" -d '{"cep" : "01153000"}'
```
# Como visualizar 

Acesse o Zipkin UI para visualizar o tracing da aplicação no endpoint `http://localhost:9411/`:

# Exemplo de saida esperado
 https://drive.google.com/file/d/1kXIdkTgbUaeoRdEfeG-lDCU505nZkFi0/view?usp=drive_link





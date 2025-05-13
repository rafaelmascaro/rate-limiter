# rate-limiter

Rate Limiter em Go que pode ser configurado para limitar o número máximo de requisições por segundo com base no **endereço IP** ou **token de acesso**. O Rate Limiter trabalha como um middleware injetado ao servidor web (para essa aplicação, possui um endpoint que retorna um "Hello, World!", rodando na porta 8080`).

## Funcionamento do Rate Limiter
O Rate Limiter controla o tráfego de requisições para um serviço web. Limitando o número de requisições recebidas por segundo a um **limite** definido, e, caso o limite seja atingido, bloqueia o tráfego dessa origem por um **intervalo de tempo** definido. O limite de requisições e tempo de bloqueio podem ser definidos das seguintes formas:
1. **Token de acesso único**: especificado no header API_KEY da requisição.
2. **Endereço IP**: valores exclusivos para o IP.
3. **Default**: valida o IP, com valores default da aplicação.

A definição do limiter é feita na ordem **Token > IP > Default**, ou seja, caso o Token e o IP da requisição tenham limiters pré configurados, o Rate Limiter se baseará nas configurações do Token.

## Configuração do Rate Limiter
Configurações de limiter:
* Limit: número máximo de requisições por segundo.
* Time Block: tempo de bloqueio em caso de atingimento do limite (em segundos).

Armazenamento das configurações por tipo de validação:
* Token de acesso único: Configurações armazenadas por Token no banco de dados.
* Endereço IP: Configurações armazenadas por IP no banco de dados.
* Default: Configurado no arquivo .env (parametros RATE_LIMIT_DEFAULT e TIME_BLOCK_DEFAULT)

Sobre as configurações armazenadas em banco de dados:
* Nessa aplicação estamos usando um banco de dados Redis, sendo o db 0 para a aplicação e o 4 para testes automatizados.
* A aplicação conta com limiters pré-configurados que sobem pelo docker-compose (service redis_init).
* Também é possível criar um endpoint para o método AddHash de repository, que insere configurações de limiter no banco de dados.

## Arquitetura do Rate Limiter
* /internal/repository/: diretório que faz a comunicação com o banco de dados (no caso dessa aplicação, Redis).
* /internal/entity/interfaces.go: possui a interface com o repository para o resto da aplicação, facilitando o uso de diferentes bancos de dados ou troca.
* /internal/usecase/rate_limiter.go: possui toda a lógica de controle de tráfego e orquestração do banco de dados.
* /internal/middleware/rate_limiter.go: middleware injetado aos serviços web, que orquestra o usecase.

## Executando a aplicação em ambiente local
1. Certifique-se de ter o Docker instalado.
2. Suba os containers necessários executando o comando:
    ```bash
    docker-compose up --build app
    ```
3. Aguarde até que a mensagem de que a aplicação está rodando na porta :8080 seja exibida nos logs.
4. Pronto! O serviço esta disponível no ambiente local. Pode ser consumido usando o modelo disponível em `api/rate_limiter.http` ou pela curl abaixo (ajustar o API KEY):
    ```bash
    curl http://localhost:8080/
    ```

## Testes automatizados
A aplicação conta com testes automatizados na camada de middleware, que demonstram o funcionamento de todo o rate limiter. Os testes contam com os seguintes cenários:
1. **Default**: Demonstra o funcionamento para uma requisição sem token e sem IP pré-configurado (`limit: 10` e `time block: 20s`). 
2. **Token**: Token e IP pré-configurados, demonstrando também a soberania da configuração do token (`limit: 10` e `time block: 10s`)
3. **IP**: Adiciona configuração do IP local no banco de dados antes de demonstrar o funcionamento (`limit: 5` e `time block: 15s`)

Para cada um dos cenários, o teste executa o seguinte fluxo:
1. Valida se todas as requisições retornam o código http 200 até atingir o limit.
2. Valida se a próxima requisição retorna o código http 429.
3. Aguarda somente 5 segundos e faz uma nova requisição, que deve retornar o código http 429 por não ter corrido todo o tempo de bloqueio.
4. Aguarda o restante do tempo de bloqueio e faz uma nova requisição, que deve retornar o código http 200.

Passo a passo para execução dos testes automatizados:
1. Certifique-se de ter o Docker instalado.
2. Dispare os testes executando o comando:
    ```bash
    docker-compose up test
    ```
3. Aguarde até que a execução seja finalizada e avalie os logs exibidos no terminal.
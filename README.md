# rate-limiter

Rate Limiter em Go que pode ser configurado para limitar o número máximo de requisições por segundo com base no **endereço IP** ou **token de acesso**. O Rate Limiter funciona como um middleware integrado ao servidor web (nesta aplicação, possui um endpoint que retorna um "Hello, World!", rodando na porta `:8080`).

## Funcionamento do Rate Limiter
O Rate Limiter gerencia o tráfego de requisições para um serviço web, limitando o número de requisições recebidas por segundo a um **limite** definido. Se o limite for atingido, ele bloqueia o tráfego dessa origem por um **intervalo de tempo** especificado. O limite de requisições e o tempo de bloqueio podem ser configurados das seguintes maneiras:
1. **Token de acesso único**: especificado no cabeçalho API_KEY da requisição.
2. **Endereço IP**: valores exclusivos para o IP.
3. **Padrão**: valida o IP com valores padrão da aplicação.

A definição do limiter segue a ordem **Token > IP > Padrão**, ou seja, se o Token e o IP da requisição tiverem limiters pré-configurados, o Rate Limiter usará as configurações do Token.

## Configuração do Rate Limiter
Configurações do limiter:
* Limite: número máximo de requisições por segundo.
* Tempo de Bloqueio: tempo de bloqueio em caso de atingimento do limite (em segundos).

Armazenamento das configurações por tipo de validação:
* Token de acesso único: Configurações armazenadas por Token no banco de dados.
* Endereço IP: Configurações armazenadas por IP no banco de dados.
* Padrão: Configurado no arquivo .env (parâmetros RATE_LIMIT_DEFAULT e TIME_BLOCK_DEFAULT)

Sobre as configurações armazenadas no banco de dados:
* Nesta aplicação, estamos usando um banco de dados Redis, sendo o db 0 para a aplicação e o 4 para testes automatizados.
* A aplicação possui limiters pré-configurados que são carregados pelo docker-compose (serviço redis_init).
* Também é possível criar um endpoint para o método AddHash do repositório, que insere configurações de limiter no banco de dados.

## Arquitetura do Rate Limiter
* /internal/repository/: diretório que faz a comunicação com o banco de dados (nesta aplicação, Redis).
* /internal/entity/interfaces.go: contém a interface com o repositório para o restante da aplicação, facilitando o uso de diferentes bancos de dados ou troca.
* /internal/usecase/rate_limiter.go: contém toda a lógica de controle de tráfego e orquestração do banco de dados.
* /internal/middleware/rate_limiter.go: middleware integrado aos serviços web, que orquestra o caso de uso.

## Executando a aplicação em ambiente local
1. Certifique-se de ter o Docker instalado.
2. Suba os containers necessários executando o comando:
    ```bash
    docker-compose up --build app
    ```
3. Aguarde até que a mensagem de que a aplicação está rodando na porta :8080 seja exibida nos logs.
4. O serviço estará disponível no ambiente local. Pode ser consumido usando o modelo disponível em `api/rate_limiter.http` ou pelo curl abaixo (ajustar o API KEY):
    ```bash
    curl http://localhost:8080/
    ```

## Testes automatizados
A aplicação possui testes automatizados na camada de middleware, que demonstram o funcionamento de todo o rate limiter. Os testes cobrem os seguintes cenários:
1. **Padrão**: Demonstra o funcionamento para uma requisição sem token e sem IP pré-configurado (`limite: 10` e `tempo de bloqueio: 20s`). 
2. **Token**: Token e IP pré-configurados, demonstrando também a prioridade da configuração do token (`limite: 10` e `tempo de bloqueio: 10s`)
3. **IP**: Adiciona configuração do IP local no banco de dados antes de demonstrar o funcionamento (`limite: 5` e `tempo de bloqueio: 15s`)

Para cada um dos cenários, o teste executa o seguinte fluxo:
1. Valida se todas as requisições retornam o código http 200 até atingir o limite.
2. Valida se a próxima requisição retorna o código http 429.
3. Aguarda apenas 5 segundos e faz uma nova requisição, que deve retornar o código http 429 por não ter passado todo o tempo de bloqueio.
4. Aguarda o restante do tempo de bloqueio e faz uma nova requisição, que deve retornar o código http 200.

Passo a passo para execução dos testes automatizados:
1. Certifique-se de ter o Docker instalado.
2. Execute os testes com o comando:
    ```bash
    docker-compose up test
    ```
3. Aguarde até que a execução seja finalizada e avalie os logs exibidos no terminal.
# MQTT Broker (Go)

Este broker atua como a primeira camada de recepção, gerenciando a telemetria enviada pelos sensores de campo.

## Justificativa Técnica das Ferramentas

| Tecnologia | Justificativa |
| :--- | :--- |
| **Go (Golang)** | Alta performance em concorrência (*goroutines*), essencial para lidar com múltiplos sensores simultâneos com baixo consumo de recursos e binários leves. |
| **Mochi MQTT v2** | Biblioteca de broker *embeddable* que permite customização total via **Hooks**, permitindo injetar lógica de monitoramento e integração diretamente no motor do servidor. |

## Arquitetura de Software (Hooks)

O projeto utiliza o conceito de **Hooks** para interceptar pacotes sem interromper o fluxo principal do protocolo:
1. **Auth Hook:** Define a política de acesso (atualmente em modo *Allow All* para desenvolvimento).
2. **Custom Log Hook:** Intercepta o evento `OnPublished` para exibir dados em tempo real no terminal, garantindo rastreabilidade durante os testes.

---

## Como Executar (Localmente)

1. **Instalar dependências:** 
```bash
   go mod tidy
```

2. **Rodar o projeto:**

```bash
go run main.go
```
## Como Executar (via Docker)
Para facilitar o deploy e garantir isolamento, utilize o Docker:

1. **Construir a imagem:**

```bash
docker build -t mqtt-broker .
```
2. **Rodar o container:**

```bash
docker run -d -p 1883:1883 --name broker-container mqtt-broker
```
3. **Verificar os logs (mensagens dos sensores):**

```bash
docker logs -f broker-container
```

## Testando a Conexão
Conecte qualquer cliente MQTT (MQTTX, MQTT Explorer ou ESP32) ao endereço:

- **Host**: ``localhost``

- **Porta**: ``1883``

Ao publicar uma mensagem, o terminal (ou logs do docker) exibirá:
``Mensagem no tópico [fazenda/sensor1]: {"umidade": 45.2}``


### O que isso resolve?
* **Reprodutibilidade:** Consegue rodar o broker com um único comando Docker, sem precisar instalar o Go na máquina.
* **Profissionalismo:** O binário estático no `alpine` resulta em uma imagem de apenas alguns megabytes, mostrando eficiência.


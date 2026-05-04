# MQTT Broker (Go)

Este broker atua como a primeira camada de recepção, gerenciando a telemetria enviada pelos sensores de campo.

## Justificativa Técnica das Ferramentas

| Tecnologia | Justificativa |
| :--- | :--- |
| **Go (Golang)** | Alta performance em concorrência (*goroutines*), essencial para lidar com múltiplos sensores simultâneos com baixo consumo de recursos e binários leves. |
| **Mochi MQTT v2** | Biblioteca de broker *embeddable* que permite customização total via **Hooks**, permitindo injetar lógica de monitoramento e integração diretamente no motor do servidor. |

### O que isso resolve?
* **Reprodutibilidade:** Consegue rodar o broker com um único comando Docker, sem precisar instalar o Go na máquina.
* **Profissionalismo:** O binário estático no `alpine` resulta em uma imagem de apenas alguns megabytes, mostrando eficiência.



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
ou ```docker compose up -d```

3. **Verificar os logs (mensagens dos sensores):**

```bash
docker logs -f broker-container
```

## Testando a Conexão
Conecte qualquer cliente MQTT (MQTTX, MQTT Explorer ou ESP32) ao endereço:

- **Host**: ``localhost``

- **Porta**: ``1883``

-----
-----
_______

## O Broker recebe o seguinte padrao de mensagem
`` mosquitto_pub -h <ip_do_broker> -i "<ID_de_quem_envia podendo ser o deviceId>" -u "<id_do_campo>" -P "<senha_do_user>" -t "campo/<id_do_campo>/<tipo_de_remetente>/<deviceId>/dados" -m {payload:"payload"} -r`` 

(tipo_de_rementente === sensor)

``mosquitto_pub -h localhost -i "1e23456" -u "fazenda1" -P "pass" -t "campo/fazenda1/sensor/1e23456/dados" -m {payload:"payload"} -r``

(-i === deviceId) <br>
(-u === userName/id do campo) <br>
(-P === senha do usuario) <br>
(-t === Topico) <br>
(-m === payload) <br>
(-r === armazenar o último valor válido para aquele tópico.) *persistencia* <br> 

### Topic 
``campo/{campoId}/sensor/{deviceId}/dados``

##### Exemplo:

``campo/123/sensor/1e23456/dados``

### Payload
##### Ideia de chaves exemplo:

`` {timestamp, temp_ar, hum_ar, hum_solo, luminosidade, bateria (?), lat, long, flag_validade}``



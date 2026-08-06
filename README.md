# MQTT Broker (Go)

Este broker atua como a primeira camada de recepção, gerenciando a telemetria enviada pelos sensores de userId.

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
`` mosquitto_pub -h <ip_do_broker> -i "<deviceId>" -u "<userId>" -P "<senha_do_user>" -t "userId/<userId>/<tipo_de_remetente>/<deviceId>/dados" -m {payload:"payload"} -r`` 

(tipo_de_rementente === sensor)

``mosquitto_pub -h localhost -i "1e23456" -u "fazenda1" -P "pass" -t "userId/fazenda1/sensor/1e23456/dados" -m '{"name":"batata", "soil_temperature": 12, "soil_moisture": 123, "air_humidity": 124, "luminosity": 125, "air_temperature": 126, "battery": 127, "latitude": "-23.3454252", "longitude": "-46.123231","timestamp":1778275967}' -r``


(-i === deviceId) <br>
(-u === userId) <br>
(-P === senha do usuario) <br>
(-t === Topico) <br>
(-m === payload) <br>
(-r === armazenar o último valor válido para aquele tópico.) *persistencia* <br> 

### Topic 
``userId/{userId}/sensor/{deviceId}/dados``

##### Exemplo:

``userId/fazenda1/sensor/1e23456/dados``

### Payload
##### Ideia de chaves exemplo:

`` {timestamp, temp_ar, hum_ar, hum_solo, luminosidade, bateria (?), lat, long, flag_validade}``







mosquitto_pub -h localhost -i "5e76ce4fd99eefe3" -u "public" -P "public" -t "application/c914c505-9a26-4c4b-acf1-400fa51514a0/device/5e76ce4fd99eefe3/event/up" -m '{
    "deduplicationId": "9b434ab2-b582-495f-bd01-98daa78e9bd8",
    "time": "2026-07-08T20:10:19.414+00:00",
    "deviceInfo": {
        "tenantId": "52f14cd4-c6f1-4fbd-8f87-4025e1d49242",
        "tenantName": "IMT",
        "applicationId": "c914c505-9a26-4c4b-acf1-400fa51514a0",
        "applicationName": "IMT-Students",
        "deviceProfileId": "5311144d-4ae6-4f50-8907-8e9359f0e18a",
        "deviceProfileName": "1.0.3-ABP",
        "deviceName": "2026-tcc-cmd03",
        "devEui": "5e76ce4fd99eefe3",
        "deviceClassEnabled": "CLASS_A",
        "tags": {}
    },
    "devAddr": "d99eefe3",
    "adr": true,
    "dr": 0,
    "fCnt": 0,
    "fPort": 1,
    "confirmed": false,
    "data": "AAAAAwn2F3ARlB9AJxAB",
    "rxInfo": [
        {
            "gatewayId": "c0ba1ffffe007566",
            "uplinkId": 2377,
            "gwTime": "2026-07-08T20:10:19.414638+00:00",
            "nsTime": "2026-07-08T20:10:17.929726028+00:00",
            "timeSinceGpsEpoch": "1467576637.414s",
            "rssi": -57,
            "snr": 7.5,
            "channel": 6,
            "rfChain": 1,
            "location": {
                "latitude": -23.65001966566475,
                "longitude": -46.574355343130826
            },
            "context": "58BCzA==",
            "metadata": {
                "region_config_id": "au915_0",
                "region_common_name": "AU915"
            },
            "crcStatus": "CRC_OK"
        }
    ],
    "txInfo": {
        "frequency": 916400000,
        "modulation": {
            "lora": {
                "bandwidth": 125000,
                "spreadingFactor": 12,
                "codeRate": "CR_4_5"
            }
        }
    }
}' -r



mosquitto_pub -h localhost -i "5e76ce4fd99eefe3" -u "public" -P "public" -t "application/c914c505-9a26-4c4b-acf1-400fa51514a0/device/5e76ce4fd99eefe3/event/up" -m '{
    "deduplicationId": "75643b5e-84d0-4094-80ce-af89a5547d3c",
    "time": "2026-07-08T19:53:55.352+00:00",
    "deviceInfo": {
        "tenantId": "52f14cd4-c6f1-4fbd-8f87-4025e1d49242",
        "tenantName": "IMT",
        "applicationId": "c914c505-9a26-4c4b-acf1-400fa51514a0",
        "applicationName": "IMT-Students",
        "deviceProfileId": "5311144d-4ae6-4f50-8907-8e9359f0e18a",
        "deviceProfileName": "1.0.3-ABP",
        "deviceName": "2026-tcc-cmd03",
        "devEui": "5e76ce4fd99eefe3",
        "deviceClassEnabled": "CLASS_A",
        "tags": {}
    },
    "devAddr": "d99eefe3",
    "adr": true,
    "dr": 0,
    "fCnt": 0,
    "fPort": 1,
    "confirmed": false,
    "data": "AAABlQn2F3ARlB9AJxAB",
    "rxInfo": [
        {
            "gatewayId": "c0ba1ffffe007566",
            "uplinkId": 12697,
            "gwTime": "2026-07-08T19:53:55.352280+00:00",
            "nsTime": "2026-07-08T19:53:53.890121503+00:00",
            "timeSinceGpsEpoch": "1467575653.352s",
            "rssi": -47,
            "snr": 8.8,
            "channel": 2,
            "location": {
                "latitude": -23.65001966566475,
                "longitude": -46.574355343130826
            },
            "context": "rRipNQ==",
            "metadata": {
                "region_common_name": "AU915",
                "region_config_id": "au915_0"
            },
            "crcStatus": "CRC_OK"
        }
    ],
    "txInfo": {
        "frequency": 915600000,
        "modulation": {
            "lora": {
                "bandwidth": 125000,
                "spreadingFactor": 12,
                "codeRate": "CR_4_5"
            }
        }
    }
}' -r
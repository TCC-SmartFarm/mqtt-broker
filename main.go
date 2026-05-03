package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"strings"

	"github.com/mochi-mqtt/server/v2"
	// "github.com/mochi-mqtt/server/v2/hooks/auth"
	"github.com/mochi-mqtt/server/v2/hooks/storage/badger"
	"github.com/mochi-mqtt/server/v2/listeners"
	"github.com/mochi-mqtt/server/v2/packets" // Importante para reconhecer o tipo 'packets.Packet'
)

func main() {
	server := mqtt.New(nil)

	err := server.AddHook(new(AuthHook), nil)
    if err != nil {
        log.Fatal(err)
    }

	err = server.AddHook(new(badger.Hook), &badger.Options{
		Path: "./data/db", // Onde os dados (inclusive mensagens retidas) serão salvos
	})
    if err != nil {
        log.Fatal(err)
    }

	// Hook de log direto no main para não precisar de outro arquivo
	_ = server.AddHook(&LogHook{}, nil)

	tcp := listeners.NewTCP(listeners.Config{
		ID:      "mqtt-tcp",
		Address: ":1883",
	})
	
	if err := server.AddListener(tcp); err != nil {
		log.Fatal(err)
	}

	go func() {
		err := server.Serve()
		if err != nil {
			log.Fatal(err)
		}
	}()

	log.Println("Broker rodando e printando mensagens...")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig
	server.Close()
}

// HOOK para log

type LogHook struct {
	mqtt.HookBase
}

func (h *LogHook) ID() string {
	return "log-hook"
}

func (h *LogHook) Provides(b byte) bool {
	return b == mqtt.OnPublished
}

func (h *LogHook) OnPublished(cl *mqtt.Client, pk packets.Packet) {
	log.Printf("Mensagem no tópico [%s]: %s\n", pk.TopicName, string(pk.Payload))
}


// HOOK para autenticação

type AuthHook struct {
	mqtt.HookBase
}

func (h *AuthHook) ID() string {
	return "mqtt-broker-auth-hook"
}

// este Hook cuida da Conexão (Auth) e Acesso a Tópicos (ACL)
func (h *AuthHook) Provides(b byte) bool {
	return b == mqtt.OnConnectAuthenticate || b == mqtt.OnACLCheck
}

// AUTENTICAÇÃO: Verifica Usuário e Senha
func (h *AuthHook) OnConnectAuthenticate(cl *mqtt.Client, pk packets.Packet) bool {
	// Vou validar com variáveis de ambiente 
	user := string(pk.Connect.Username)
	pass := string(pk.Connect.Password)

	if user == "admin" && pass == "admin" {
		return true // Admin tem acesso total
	}

	if user == "sensor" && pass == "123" {
		return true
	}

	if user == "mqtt_sub" && pass == "mqtt_sub" {
		return true
	}

	log.Printf("Conexão negada: Usuário %s incorreto", user)
	return false
}

// ACL: Verifica se o cliente pode publicar/assinar em um tópico
func (h *AuthHook) OnACLCheck(cl *mqtt.Client, topic string, write bool) bool {
	// Regra de Ouro: Admin pode tudo
	if string(cl.Properties.Username) == "admin" {
		return true
	}

	// Regra para Sensores: Só podem publicar no seu próprio prefixo
	// Exemplo: O cliente "esp32" só publica em "campo/+/sensor/<cliente_id>/dados"
	if string(cl.Properties.Username) == "sensor" {
		if write {
			topic = strings.TrimSuffix(topic, "/")

			parts := strings.Split(topic, "/")

			log.Printf("DEBUG ACL -> topic=%s parts=%v", topic, parts)

			// Esperado: campo/<algo>/sensor/<clientID>/dados
			if len(parts) != 5 {
				return false
			}

			if parts[0] != "campo" {
				return false
			}

			if parts[2] != "sensor" {
				return false
			}

			if parts[3] != cl.ID {
				return false
			}

			if parts[4] != "dados" {
				return false
			}

			log.Printf("WRITE OK sensor [%s] topic [%s]", cl.ID, topic)
			return true
		}


		// leitura bloqueada
		if !write {
			return true
		}
	}

	// Regra de Leitura (Subscribe): Impedir que sensores leiam dados de outros
	return false 
}
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

	if user == "fazenda0" && pass == "pass" {
		return true
	}

	if user == "fazenda1" && pass == "pass" {
		return true
	}

	if user == "fazenda2" && pass == "pass" {
		return true
	}

	if user == "fazenda3" && pass == "pass" {
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
	userId := string(cl.Properties.Username)
	
	// Regra de Ouro: Admin pode tudo
	if userId == "admin" {
		return true
	}

	// mqtt_sub pode ler TUDO
    if userId == "mqtt_sub" {
        return !write // Permite apenas leitura (Subscribe)
    }

    // Regra para Usuários de Fazendas (Sensores)
    if write { // Tentativa de Publicação
        topic = strings.TrimSuffix(topic, "/")
        parts := strings.Split(topic, "/")

        // Validação: userId / {userId} / sensor / {clientID} / dados
        // Exemplo: userId/fazenda1/sensor/sensor01/dados
        if len(parts) == 5 && 
           parts[0] == "userId" && 
           parts[1] == userId && // O segundo nível TEM que ser o nome do usuário
           parts[2] == "sensor" && 
           parts[3] == cl.ID && 
           parts[4] == "dados" {
            return true
        }
        
        log.Printf("ACL NEGADA: Usuário %s tentou publicar em tópico proibido: %s", userId, topic)
        return false
    }

	// Regra de Leitura (Subscribe): Impedir que sensores leiam dados de outros
	return false 
}
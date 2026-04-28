package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/hooks/auth"
	"github.com/mochi-mqtt/server/v2/listeners"
	"github.com/mochi-mqtt/server/v2/packets" // Importante para reconhecer o tipo 'packets.Packet'
)

func main() {
	server := mqtt.New(nil)
	_ = server.AddHook(new(auth.AllowHook), nil) // qualquer cliente pode se conectar, sem autenticação // lembrar de trancar essas portas depois que tiver o banco e o esp funcionando

	// Hook de log direto no main para não precisar de outro arquivo
	_ = server.AddHook(&LogHook{}, nil)

	tcp := listeners.NewTCP(listeners.Config{
		ID:      "t1",
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

// --- HOOK para log ---

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
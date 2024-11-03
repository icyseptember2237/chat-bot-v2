package hook

import "chatbot/msg"

type Hook func(message *msg.ReceiveMessage) bool

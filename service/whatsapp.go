package service

import (
	"context"
	waProto "go.mau.fi/whatsmeow/binary/proto"

)

type WhatsApp interface {
	SendMessage(ctx context.Context, recipientNumber string, message *waProto.Message) error
}

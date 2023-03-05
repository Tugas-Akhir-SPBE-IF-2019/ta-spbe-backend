package whatsapp

import (
	"context"

	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
	// "google.golang.org/protobuf/proto"
)

type WhatsMeow struct {
	Client *whatsmeow.Client
}

func (wm *WhatsMeow) SendMessage(ctx context.Context, recipientNumber string, message *waProto.Message) error {
	const initialMessage = `*[OTOMATISASI PENILAIAN SPBE]*

` + "```" + `Terima kasih telah menggunakan Aplikasi Otomatisasi Penilaian SPBE. Hasil penilaian anda akan keluar dalam beberapa saat lagi.` + "```"
		// initialTemplateMessage := &waProto.Message{
		// 	TemplateMessage: &waProto.TemplateMessage{
		// 		HydratedTemplate: &waProto.TemplateMessage_HydratedFourRowTemplate{
		// 			Title: &waProto.TemplateMessage_HydratedFourRowTemplate_HydratedTitleText{
		// 				HydratedTitleText: "The Title",
		// 			},
		// 			TemplateId:          proto.String("template-id"),
		// 			HydratedContentText: proto.String("The Content"),
		// 			HydratedFooterText:  proto.String("The Footer"),
		// 			HydratedButtons: []*waProto.HydratedTemplateButton{

		// 				// This for URL button
		// 				{
		// 					Index: proto.Uint32(1),
		// 					HydratedButton: &waProto.HydratedTemplateButton_UrlButton{
		// 						UrlButton: &waProto.HydratedTemplateButton_HydratedURLButton{
		// 							DisplayText: proto.String("The Link"),
		// 							Url:         proto.String("https://fb.me/this"),
		// 						},
		// 					},
		// 				},

		// 				// This for call button
		// 				{
		// 					Index: proto.Uint32(2),
		// 					HydratedButton: &waProto.HydratedTemplateButton_CallButton{
		// 						CallButton: &waProto.HydratedTemplateButton_HydratedCallButton{
		// 							DisplayText: proto.String("Call us"),
		// 							PhoneNumber: proto.String("1234567890"),
		// 						},
		// 					},
		// 				},

		// 				// This is just a quick reply
		// 				{
		// 					Index: proto.Uint32(3),
		// 					HydratedButton: &waProto.HydratedTemplateButton_QuickReplyButton{
		// 						QuickReplyButton: &waProto.HydratedTemplateButton_HydratedQuickReplyButton{
		// 							DisplayText: proto.String("Quick reply"),
		// 							Id:          proto.String("quick-id"),
		// 						},
		// 					},
		// 				},
		// 			},
		// 		},
		// 	},
		// }

	recipient := types.NewJID(recipientNumber, "s.whatsapp.net")
	_, err := wm.Client.SendMessage(context.Background(), recipient, message)

	return err
}

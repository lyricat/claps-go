package util

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"claps-test/model"

	"github.com/fox-one/mixin-sdk-go"
	"github.com/gofrs/uuid"
)

const (
	MessageCategoryPlainText    = "PLAIN_TEXT"
	MessageCategoryPlainImage   = "PLAIN_IMAGE"
	MessageCategoryPlainVideo   = "PLAIN_VIDEO"
	MessageCategoryPlainData    = "PLAIN_DATA"
	MessageCategoryPlainSticker = "PLAIN_STICKER"
	MessageCategoryPlainContact = "PLAIN_CONTACT"
	MessageCategoryPlainAudio   = "PLAIN_AUDIO"
	MessageCategoryAppCard      = "APP_CARD"
	MessageCategoryAppButtons   = "APP_BUTTON_GROUP"
	MessageCategoryPlainPost    = "PLAIN_POST"
	MessageCategoryPlainLive    = "PLAIN_LIVE"

	// @TODO message state management
	MessageStateInit      = "init"
	MessageStateSending   = "sending"
	MessageStateDelivered = "delivered"
)

// CreateConversation create a conversation with a specified user. In practice, save the conversation id for further using
func CreateConversation(ctx context.Context, client *mixin.Client, userID string) (*mixin.Conversation, error) {
	conversation, err := client.CreateContactConversation(ctx, userID)
	if err != nil {
		return nil, err
	}
	return conversation, nil
}

// SendMessage send a message to a specified user.
func SendMessage(ctx context.Context, client *mixin.Client, conversation *mixin.Conversation, category string, data []byte) error {
	message := &mixin.MessageRequest{
		ConversationID: conversation.ConversationID,
		MessageID:      uuid.Must(uuid.NewV4()).String(),
		Category:       MessageCategoryAppCard,
		Data:           base64.StdEncoding.EncodeToString(data),
	}

	err := client.SendMessage(ctx, message)
	if err != nil {
		return err
	}
	return nil
}

// SendTransferNotification send a notification when a transfer delivered.
func SendTransferNotification(ctx context.Context, client *mixin.Client, conversation *mixin.Conversation, assetID, amount, traceID string) error {
	asset, err := (&model.Asset{}).GetAssetById(assetID)
	if err != nil {
		return err
	}

	card, err := json.Marshal(map[string]string{
		"icon_url":    asset.IconUrl,
		"title":       amount,
		"description": asset.Name,
		"action":      fmt.Sprintf("mixin://snapshots?trace=%s", traceID),
	})
	if err != nil {
		return err
	}
	return SendMessage(ctx, client, conversation, MessageCategoryAppCard, card)
}

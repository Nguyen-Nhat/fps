package slack

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"

	config "git.teko.vn/loyalty-system/loyalty-file-processing/configs"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/common/constant"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/common/errorz"
	"go.tekoapis.com/tekone/library/nap"
)

const (
	fieldMaxLength     = 1950
	defaultHttpTimeout = 5 // seconds
)

// Client ...
type Client interface {
	Send(ctx context.Context, payload *Payload) error
	SendInfo(ctx context.Context, message string, req interface{}, args map[string]string, mentionUsers ...string)
	SendError(ctx context.Context, message string, req interface{}, error error, args map[string]string, mentionUsers ...string)
	SendWarning(ctx context.Context, message string, req interface{}, args map[string]string, mentionUsers ...string)
}

type clientImpl struct {
	config                *config.SlackWebhook
	httpClient            *nap.Nap
	defaultMentionUserIds []string
}

type AttachmentField struct {
	Title string `json:"title"`
	Value string `json:"value"`
}

type Attachment struct {
	Color      string            `json:"color,omitempty"`
	Text       string            `json:"text,omitempty"`
	Fields     []AttachmentField `json:"fields,omitempty"`
	MarkdownIn []string          `json:"markdown_in,omitempty"`
}

type Payload struct {
	Username    string        `json:"username,omitempty"` // use for directed message
	IconEmoji   string        `json:"icon_emoji,omitempty"`
	Channel     string        `json:"channel,omitempty"`
	Text        string        `json:"text,omitempty"`
	Attachments []*Attachment `json:"attachments,omitempty"`
	Markdown    bool          `json:"markdown,omitempty"`
}

func NewSlackClient(cfg config.SlackWebhook) Client {
	if cfg.Enable && len(cfg.Url) == 0 {
		panic("Slack webhook url is required")
	}

	return &clientImpl{
		config:                &cfg,
		httpClient:            nap.New().Base(cfg.Url),
		defaultMentionUserIds: strings.Split(cfg.MentionUserIds, constant.SplitByComma),
	}
}

func (c *clientImpl) Send(ctx context.Context, payload *Payload) error {
	if !c.config.Enable {
		ctxzap.Extract(ctx).Info("Slack is disabled")
		return nil
	}
	ctx, cancel := context.WithTimeout(ctx, defaultHttpTimeout*time.Second)
	defer cancel()
	resp, err := c.httpClient.SetContext(ctx).Post("").BodyJSON(payload).Receive(nil, nil)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return errorz.SlackSendingErr(resp.Status)
	}
	return nil
}

func (c *clientImpl) SendInfo(ctx context.Context, message string, req interface{}, args map[string]string, mentionUsers ...string) {
	if !c.config.Enable {
		ctxzap.Extract(ctx).Info("Slack is disabled")
		return
	}

	attachmentFields := make([]AttachmentField, 0)
	if c.config.Environment != "" {
		attachmentFields = append(attachmentFields, AttachmentField{
			Title: "Environment:",
			Value: fmt.Sprintf("`%s`", c.config.Environment),
		})
	}
	if req != nil {
		valueBytes, err := json.Marshal(req)
		if err != nil {
			ctxzap.Extract(ctx).Error("Slack.SendError | Marshalling req got error", zap.Error(err))
			return
		}
		attachmentFields = append(attachmentFields, AttachmentField{
			Title: "Request:",
			Value: fmt.Sprintf("```%s```", string(valueBytes)),
		})
	}
	for key, value := range args {
		attachmentFields = append(attachmentFields, AttachmentField{
			Title: key,
			Value: value,
		})
	}

	payload := &Payload{
		Channel: c.config.AlertChannel,
		Text:    fmt.Sprintf("%s %s %s", IconBell, getMessageWithServiceName(message), c.convertToMentionText(mentionUsers)),
		Attachments: []*Attachment{
			{
				Color:      ColorInfo,
				Fields:     attachmentFields,
				MarkdownIn: []string{"text", "fields"},
			},
		},
		Markdown: true,
	}

	err := c.Send(ctx, payload)
	if err != nil {
		ctxzap.Extract(ctx).Info("Can't send slack message", zap.Any("payload", payload), zap.Any("err", err))
	}
}

func (c *clientImpl) SendError(ctx context.Context, message string, req interface{}, error error, args map[string]string, mentionUsers ...string) {
	if !c.config.Enable {
		ctxzap.Extract(ctx).Info("Slack is disabled")
		return
	}

	attachmentFields := make([]AttachmentField, 0)
	if c.config.Environment != "" {
		attachmentFields = append(attachmentFields, AttachmentField{
			Title: "Environment:",
			Value: fmt.Sprintf("`%s`", c.config.Environment),
		})
	}
	if req != nil {
		valueBytes, err := json.Marshal(req)
		if err != nil {
			ctxzap.Extract(ctx).Error("Slack.SendError | Marshalling req got error", zap.Error(err))
			return
		}
		reqMessage := string(valueBytes)
		if len(reqMessage) > fieldMaxLength {
			reqMessage = reqMessage[:fieldMaxLength] + "..."
		}
		attachmentFields = append(attachmentFields, AttachmentField{
			Title: "Request:",
			Value: fmt.Sprintf("```%s```", reqMessage),
		})
	}
	attachmentFields = append(attachmentFields, AttachmentField{
		Title: "Error:",
		Value: fmt.Sprintf("```%v```", error),
	})
	for key, value := range args {
		attachmentFields = append(attachmentFields, AttachmentField{
			Title: key,
			Value: value,
		})
	}

	payload := &Payload{
		Channel: c.config.AlertChannel,
		Text:    fmt.Sprintf("%s %s %s", IconNoEntry, getMessageWithServiceName(message), c.convertToMentionText(mentionUsers)),
		Attachments: []*Attachment{
			{
				Color:      ColorErr,
				Fields:     attachmentFields,
				MarkdownIn: []string{"text", "fields"},
			},
		},
		Markdown: true,
	}

	err := c.Send(ctx, payload)
	if err != nil {
		ctxzap.Extract(ctx).Info("Can't send slack message", zap.Any("payload", payload), zap.Any("err", err))
	}
}

func (c *clientImpl) SendWarning(ctx context.Context, message string, req interface{}, args map[string]string, mentionUsers ...string) {
	if !c.config.Enable {
		ctxzap.Extract(ctx).Info("Slack is disabled")
		return
	}

	attachmentFields := make([]AttachmentField, 0)
	if c.config.Environment != "" {
		attachmentFields = append(attachmentFields, AttachmentField{
			Title: "Environment:",
			Value: fmt.Sprintf("`%s`", c.config.Environment),
		})
	}
	if req != nil {
		valueBytes, err := json.Marshal(req)
		if err != nil {
			ctxzap.Extract(ctx).Error("Slack.SendError | Marshalling req got error", zap.Error(err))
			return
		}
		attachmentFields = append(attachmentFields, AttachmentField{
			Title: "Request:",
			Value: fmt.Sprintf("```%s```", string(valueBytes)),
		})
	}
	for key, value := range args {
		attachmentFields = append(attachmentFields, AttachmentField{
			Title: key,
			Value: value,
		})
	}

	payload := &Payload{
		Channel: c.config.AlertChannel,
		Text:    fmt.Sprintf("%s %s %s", IconWarning, getMessageWithServiceName(message), c.convertToMentionText(mentionUsers)),
		Attachments: []*Attachment{
			{
				Color:      ColorWarning,
				Fields:     attachmentFields,
				MarkdownIn: []string{"text", "fields"},
			},
		},
		Markdown: true,
	}

	err := c.Send(ctx, payload)
	if err != nil {
		ctxzap.Extract(ctx).Info("Can't send slack message", zap.Any("payload", payload), zap.Any("err", err))
	}
}

func (c *clientImpl) convertToMentionText(mentionUsers []string) string {
	if len(mentionUsers) == 0 {
		mentionUsers = c.defaultMentionUserIds
	}
	var result strings.Builder
	for _, m := range mentionUsers {
		if m != constant.EmptyString {
			result.WriteString(fmt.Sprintf("<@%s> ", m))
		}
	}

	return result.String()
}

func getMessageWithServiceName(message string) string {
	return fmt.Sprintf("%s %s", ServiceName, message)
}

package pkgses

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"mime/multipart"
	"mime/quotedprintable"
	"net/textproto"
	"os"
	"path/filepath"

	cfg "github.com/aws/aws-sdk-go-v2/config"
	awsSes "github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
)

func SendEmailWithMultipleAttachmentsSES(from, to, subject, body string, filePaths []string) error {
	ctx := context.TODO()

	cfg, err := cfg.LoadDefaultConfig(ctx, cfg.WithRegion("ap-south-1"))
	if err != nil {
		return fmt.Errorf("failed to load AWS config: %v", err)
	}

	client := awsSes.NewFromConfig(cfg)

	var emailRaw bytes.Buffer
	writer := multipart.NewWriter(&emailRaw)
	boundary := writer.Boundary()

	// MIME headers
	fmt.Fprintf(&emailRaw, "From: %s\r\n", from)
	fmt.Fprintf(&emailRaw, "To: %s\r\n", to)
	fmt.Fprintf(&emailRaw, "Subject: %s\r\n", subject)
	fmt.Fprintf(&emailRaw, "MIME-Version: 1.0\r\n")
	fmt.Fprintf(&emailRaw, "Content-Type: multipart/mixed; boundary=\"%s\"\r\n\r\n", boundary)

	// Body (text or HTML)
	bodyPart, _ := writer.CreatePart(textproto.MIMEHeader{
		"Content-Type":              {"text/html; charset=\"utf-8\""},
		"Content-Transfer-Encoding": {"quoted-printable"},
	})
	qp := quotedprintable.NewWriter(bodyPart)
	qp.Write([]byte(body))
	qp.Close()

	// Loop through multiple attachments
	for _, filePath := range filePaths {
		if err := addAttachment(writer, filePath); err != nil {
			return err
		}
	}

	writer.Close()

	rawInput := &awsSes.SendRawEmailInput{
		RawMessage: &types.RawMessage{
			Data: emailRaw.Bytes(),
		},
	}

	_, err = client.SendRawEmail(ctx, rawInput)
	if err != nil {
		return fmt.Errorf("failed to send raw email: %v", err)
	}

	return nil
}

func addAttachment(writer *multipart.Writer, filePath string) error {
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("cannot read file %s: %v", filePath, err)
	}

	encoded := base64.StdEncoding.EncodeToString(fileBytes)
	filename := filepath.Base(filePath)

	header := textproto.MIMEHeader{}
	header.Add("Content-Type", fmt.Sprintf("%s; name=\"%s\"", detectMimeType(filename), filename))
	header.Add("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	header.Add("Content-Transfer-Encoding", "base64")

	part, err := writer.CreatePart(header)
	if err != nil {
		return fmt.Errorf("failed to create MIME part for %s: %v", filename, err)
	}

	for i := 0; i < len(encoded); i += 76 {
		end := i + 76
		if end > len(encoded) {
			end = len(encoded)
		}
		part.Write([]byte(encoded[i:end] + "\r\n"))
	}

	return nil
}

func detectMimeType(filename string) string {
	ext := filepath.Ext(filename)
	switch ext {
	case ".pdf":
		return "application/pdf"
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".txt":
		return "text/plain"
	default:
		return "application/octet-stream"
	}
}

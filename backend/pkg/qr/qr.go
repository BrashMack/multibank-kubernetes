package qr

import (
	"encoding/json"
	"fmt"
	"image"
	"log"
	"os"

	"github.com/liyue201/goqr"
	"github.com/skip2/go-qrcode"
)

// GenerateQRCode принимает map[string]string, сериализует в JSON и сохраняет QR-код в файл
func GenerateQRCode(data map[string]string, filename string, size int) error {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("json.Marshal: %w", err)
	}

	return qrcode.WriteFile(string(jsonBytes), qrcode.Medium, size, filename)
}

// DecodeQRCode читает QR-код из файла и возвращает данные как map[string]string
func DecodeQRCode(filename string) (map[string]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("os.Open: %w", err)
	}

	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("image.Decode: %w", err)
	}

	qrCodes, err := goqr.Recognize(img)
	if err != nil {
		return nil, fmt.Errorf("goqr.Recognize: %w", err)
	}

	if len(qrCodes) == 0 {
		return nil, ErrQrCodeNotFound
	}

	var data map[string]string

	if err = json.Unmarshal(qrCodes[0].Payload, &data); err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %w", err)
	}

	return data, nil
}

func example() {
	data := map[string]string{
		"accountIdTo": "40890c53c73e34742d",
		"amount":      "1000",
	}

	filename := "qr.png"

	if err := GenerateQRCode(data, filename, 256); err != nil {
		log.Fatalf("GenerateQRCode: %v", err)
	}

	fmt.Println("QR-код сохранён в", filename)

	decoded, err := DecodeQRCode(filename)
	if err != nil {
		log.Fatalf("DecodeQRCode: %v", err)
	}

	fmt.Println("Расшифрованные данные:", decoded)
}

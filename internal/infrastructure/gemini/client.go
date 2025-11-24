package gemini

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"sandbox/internal/domain/repository"
	transactionDTO "sandbox/internal/usecase/transaction"
)

const (
	geminiAPIURL = "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.5-flash:generateContent"
)

type Client struct {
	apiKey     string
	httpClient *http.Client
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 300 * time.Second,
		},
	}
}

func (c *Client) ExtractFromDocuments(ctx context.Context, documents []repository.Document, promptType string) (interface{}, error) {
	if len(documents) == 0 {
		return nil, errors.New("no documents provided")
	}

	if c.apiKey == "" {
		return nil, errors.New("GEMINI_API_KEY is not configured. Please set the environment variable and restart the application")
	}

	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("context cancelled before starting API call: %w", err)
	}

	prompt := c.definePrompt(promptType)
	parts := []map[string]interface{}{
		{"text": prompt},
	}

	for _, doc := range documents {
		base64Content := base64.StdEncoding.EncodeToString(doc.Content)
		parts = append(parts, map[string]interface{}{
			"inline_data": map[string]interface{}{
				"mime_type": doc.MimeType,
				"data":      base64Content,
			},
		})
	}

	body := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": parts,
			},
		},
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.getAPIURL(), bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			return nil, fmt.Errorf("context cancelled during API call: %w", ctx.Err())
		}
		return nil, fmt.Errorf("failed to call Gemini API: %w", err)
	}
	defer resp.Body.Close()

	bodyResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("gemini api error (status %d): %s", resp.StatusCode, string(bodyResp))
	}

	return c.parseResponse(bodyResp)
}

func (c *Client) getAPIURL() string {
	return fmt.Sprintf("%s?key=%s", geminiAPIURL, c.apiKey)
}

func (c *Client) definePrompt(promptType string) string {
	if promptType == "scanAssigneeTransaction" {
		return c.scanAssigneeTransaction()
	} else {
		return c.scanBusinessTripDocs()
	}
}

func (c *Client) scanBusinessTripDocs() string {
	return `Baca semua dokumen berikut (gambar atau PDF).
Ekstrak setiap transaksi dan tampilkan dalam format JSON valid berikut ini:

{
  "startDate": "YYYY-MM-DD", -> ambil dari file surat tugas
  "endDate": "YYYY-MM-DD", -> ambil dari file surat tugas
  "activityPurpose": "TUJUAN_AKTIVITAS", -> ambil dari file surat tugas
  "destinationCity": "KOTA_TUJUAN", -> ambil dari file surat tugas
  "spdDate": "YYYY-MM-DD", -> ambil dari file surat tugas
  "departureDate": "YYYY-MM-DD", -> ambil dari file surat tugas
  "returnDate": "YYYY-MM-DD", -> ambil dari file surat tugas
  "assignees": [
    {
      "name": "NAMA_PEGAWAI", -> ambil dari file surat tugas
      "spd_number": "NOMOR_SPD", -> ambil dari file surat tugas
      "employee_number": "NIP_PEGAWAI", -> ambil dari file surat tugas
      "position": "JABATAN_PEGAWAI", -> ambil dari file surat tugas
      "rank": "GOLONGAN_PEGAWAI", -> ambil dari file surat tugas
      "transactions": [
        {
          "name": "NAMA_PEMESAN_TRANSAKSI",
          "type": "accommodation | transport | other | allowance",
          "subtype": "hotel | flight | train | taxi | daily_allowance",
          "amount": number,
          "total_night": number,
          "subtotal": number, -> hasil amount*total_night kalo dia accomodation tapi kalo selain itu langsung ambil dari amount aja
	      "description" : string, -> ini adalah keterangan transaksi ini transaksi apa, misalkan gojek dari alamat1 ke alamat2, kalo hotel jelasin juga hotelnya
	      "transport_detail" : string, -> ini terisi hanya jika dia transport darat ya (pesawat tidak termasuk) 1.jika dia dari bandara soetta atau tujuannya ke bandara soetta maka valuenya menjadi "transport_asal" atau kalau dia transportasinya di jakarta juga masuk trasnport asal 2.jika mengandung bandara lain selain soetta maka valuenya adalah "transport_daerah"
        }
      ]
    }
  ]
}

- Kembalikan hasil hanya dalam JSON valid (tanpa teks tambahan).
- Jangan bungkus JSON dengan tanda kutip atau karakter escape.
- Jika total_night tidak ada, field tersebut boleh dihapus.
- Pastikan angka hanya berupa digit (tanpa simbol mata uang).
- Untuk data transaksi, nama yang digunakan harus sesuai dengan nama yang tercantum di surat tugas. Harap lakukan pengecekan dan pencocokan dengan surat tugas.
- Jika nama pemesan di transaksi tersebut tidak tercantum di surat tugas, mohon assign ke salah satu nama yang ada di surat tugas.
- Jangan menggunakan nama driver sebagai nama transaksi — gunakan nama pemesan.
- Group semua transaksi di bawah setiap assignee.

di bawah ini data uang harian aku minta untuk ambil datanya untuk di masukkan ke transactions sesuai dengan kota tujuannya yang ada di surat tugas misalnya dia di surabaya maka dia akan mengambil data jawa timur karena surabaya terletak di jawa timur dan jadikan datanya sebagai allowance
NO,PROVINSI,SATUAN,LUAR KOTA,DALAM KOTA LEBIH DARI 8 JAM,DIKLAT
1,ACEH,OH,Rp360.000,Rp140.000,Rp110.000
2,SUMATRA UTARA,OH,Rp370.000,Rp150.000,Rp110.000
3,RIAU,OH,Rp370.000,Rp150.000,Rp110.000
4,KEPULAUAN RIAU,OH,Rp370.000,Rp150.000,Rp110.000
5,JAMBI,OH,Rp370.000,Rp150.000,Rp110.000
6,SUMATRA BARAT,OH,Rp380.000,Rp150.000,Rp110.000
7,SUMATRA SELATAN,OH,Rp380.000,Rp150.000,Rp110.000
8,LAMPUNG,OH,Rp380.000,Rp150.000,Rp110.000
9,BENGKULU,OH,Rp380.000,Rp150.000,Rp110.000
10,BANGKA BELITUNG,OH,Rp410.000,Rp160.000,Rp120.000
11,BANTEN,OH,Rp370.000,Rp150.000,Rp110.000
12,JAWA BARAT,OH,Rp430.000,Rp170.000,Rp130.000
13,D.K.I. JAKARTA,OH,Rp530.000,Rp210.000,Rp160.000
14,JAWA TENGAH,OH,Rp370.000,Rp150.000,Rp110.000
15,D.I. YOGYAKARTA,OH,Rp420.000,Rp170.000,Rp130.000
16,JAWA TIMUR,OH,Rp410.000,Rp160.000,Rp120.000
17,BALI,OH,Rp480.000,Rp190.000,Rp140.000
18,NUSA TENGGARA BARAT,OH,Rp440.000,Rp190.000,Rp130.000
19,NUSA TENGGARA TIMUR,OH,Rp430.000,Rp170.000,Rp130.000
20,KALIMANTAN BARAT,OH,Rp380.000,Rp150.000,Rp110.000
21,KALIMANTAN TENGAH,OH,Rp360.000,Rp140.000,Rp110.000
22,KALIMANTAN SELATAN,OH,Rp380.000,Rp150.000,Rp110.000
23,KALIMANTAN TIMUR,OH,Rp430.000,Rp170.000,Rp130.000
24,KALIMANTAN UTARA,OH,Rp430.000,Rp170.000,Rp130.000
25,SULAWESI UTARA,OH,Rp370.000,Rp150.000,Rp110.000
26,GORONTALO,OH,Rp370.000,Rp150.000,Rp110.000
27,SULAWESI BARAT,OH,Rp410.000,Rp160.000,Rp120.000
28,SULAWESI SELATAN,OH,Rp430.000,Rp170.000,Rp130.000
29,SULAWESI TENGAH,OH,Rp370.000,Rp150.000,Rp110.000
30,SULAWESI TENGGARA,OH,Rp380.000,Rp150.000,Rp110.000
31,MALUKU,OH,Rp380.000,Rp150.000,Rp110.000
32,MALUKU UTARA,OH,Rp430.000,Rp170.000,Rp130.000
33,PAPUA,OH,Rp580.000,Rp230.000,Rp170.000
34,PAPUA BARAT,OH,Rp480.000,Rp190.000,Rp140.000
35,PAPUA BARAT DAYA,OH,Rp480.000,Rp190.000,Rp140.000
36,PAPUA TENGAH,OH,Rp580.000,Rp230.000,Rp170.000
37,PAPUA SELATAN,OH,Rp580.000,Rp230.000,Rp170.000
38,PAPUA PEGUNUNGAN,OH,Rp580.000,Rp230.000,Rp170.000
`
}

func (c *Client) scanAssigneeTransaction() string {
	return `Baca semua dokumen berikut (gambar atau PDF).
Ekstrak setiap transaksi dan tampilkan dalam format JSON valid berikut ini:

{
  "startDate": "YYYY-MM-DD", -> ambil dari file surat tugas
  "endDate": "YYYY-MM-DD", -> ambil dari file surat tugas
  "activityPurpose": "TUJUAN_AKTIVITAS", -> ambil dari file surat tugas
  "destinationCity": "KOTA_TUJUAN", -> ambil dari file surat tugas
  "spdDate": "YYYY-MM-DD", -> ambil dari file surat tugas
  "departureDate": "YYYY-MM-DD", -> ambil dari file surat tugas
  "returnDate": "YYYY-MM-DD", -> ambil dari file surat tugas
  "assignees": [
    {
      "name": "NAMA_PEGAWAI", -> ambil dari file surat tugas
      "spd_number": "NOMOR_SPD", -> ambil dari file surat tugas
      "employee_number": "NIP_PEGAWAI", -> ambil dari file surat tugas
      "position": "JABATAN_PEGAWAI", -> ambil dari file surat tugas
      "rank": "GOLONGAN_PEGAWAI", -> ambil dari file surat tugas
      "transactions": [
        {
          "name": "NAMA_PEMESAN_TRANSAKSI",
          "type": "accommodation | transport | other | allowance",
          "subtype": "hotel | flight | train | taxi | daily_allowance",
          "amount": number,
          "total_night": number,
          "subtotal": number, -> hasil amount*total_night kalo dia accomodation tapi kalo selain itu langsung ambil dari amount aja
	      "description" : string, -> ini adalah keterangan transaksi ini transaksi apa, misalkan gojek dari alamat1 ke alamat2, kalo hotel jelasin juga hotelnya
	      "transport_detail" : string, -> ini terisi hanya jika dia transport darat ya (pesawat tidak termasuk) 1.jika dia dari bandara soetta atau tujuannya ke bandara soetta maka valuenya menjadi "transport_asal" atau kalau dia transportasinya di jakarta juga masuk trasnport asal 2.jika mengandung bandara lain selain soetta maka valuenya adalah "transport_daerah"
        }
      ]
    }
  ]
}

- Kembalikan hasil hanya dalam JSON valid (tanpa teks tambahan).
- Jangan bungkus JSON dengan tanda kutip atau karakter escape.
- Jika total_night tidak ada, field tersebut boleh dihapus.
- Pastikan angka hanya berupa digit (tanpa simbol mata uang).
- Untuk data transaksi, nama yang digunakan harus sesuai dengan nama yang tercantum di surat tugas. Harap lakukan pengecekan dan pencocokan dengan surat tugas.
- Jika nama pemesan di transaksi tersebut tidak tercantum di surat tugas, mohon assign ke salah satu nama yang ada di surat tugas.
- Jangan menggunakan nama driver sebagai nama transaksi — gunakan nama pemesan.
- Group semua transaksi di bawah setiap assignee.
`
}

func (c *Client) parseResponse(bodyResp []byte) (*transactionDTO.RecapReportDTO, error) {
	var geminiAPIResponse geminiResponse
	if err := json.Unmarshal(bodyResp, &geminiAPIResponse); err != nil {
		return nil, fmt.Errorf("failed to parse Gemini API response wrapper: %w", err)
	}

	if len(geminiAPIResponse.Candidates) == 0 || len(geminiAPIResponse.Candidates[0].Content.Parts) == 0 {
		return nil, errors.New("empty response candidates or parts from Gemini API")
	}

	rawText := geminiAPIResponse.Candidates[0].Content.Parts[0].Text
	cleanJSON := c.cleanJSON(rawText)

	var geminiRawReport geminiReportResponse
	if err := json.Unmarshal([]byte(cleanJSON), &geminiRawReport); err != nil {
		return nil, fmt.Errorf("failed to parse Gemini report content: %w (raw: %s)", err, cleanJSON)
	}

	geminiRawReport.ReceiptSignatureDate = time.Now().Format("2006-01-02")

	assignees := make([]transactionDTO.AssigneeDTO, 0, len(geminiRawReport.Assignees))
	for _, rawAssignee := range geminiRawReport.Assignees {
		transactionsDTO := make([]transactionDTO.TransactionDTO, 0, len(rawAssignee.Transactions))
		for _, rawTx := range rawAssignee.Transactions {
			transactionsDTO = append(transactionsDTO, transactionDTO.TransactionDTO{
				Name:            rawTx.Name,
				Type:            rawTx.Type,
				Subtype:         rawTx.Subtype,
				Amount:          rawTx.Amount,
				TotalNight:      rawTx.TotalNight,
				Subtotal:        rawTx.Subtotal,
				PaymentType:     "", // Assuming default empty, needs to be derived if applicable
				Description:     rawTx.Description,
				TransportDetail: rawTx.TransportDetail,
			})
		}

		assignees = append(assignees, transactionDTO.AssigneeDTO{
			Name:           rawAssignee.Name,
			SpdNumber:      rawAssignee.SpdNumber,
			EmployeeID:     "", // Will be populated from external API later
			EmployeeNumber: rawAssignee.EmployeeNumber,
			Position:       rawAssignee.Position,
			Rank:           rawAssignee.Rank,
			Transactions:   transactionsDTO,
		})
	}

	return &transactionDTO.RecapReportDTO{
		StartDate:            geminiRawReport.StartDate,
		EndDate:              geminiRawReport.EndDate,
		ActivityPurpose:      geminiRawReport.ActivityPurpose,
		DestinationCity:      geminiRawReport.DestinationCity,
		SpdDate:              geminiRawReport.SpdDate,
		DepartureDate:        geminiRawReport.DepartureDate,
		ReturnDate:           geminiRawReport.ReturnDate,
		ReceiptSignatureDate: geminiRawReport.ReceiptSignatureDate,
		Assignees:            assignees,
	}, nil
}

func (c *Client) cleanJSON(s string) string {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "```json")
	s = strings.TrimPrefix(s, "```JSON")
	s = strings.TrimSuffix(s, "```")
	return strings.TrimSpace(s)
}

type geminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

type geminiReportResponse struct {
	StartDate            string                `json:"startDate"`
	EndDate              string                `json:"endDate"`
	ActivityPurpose      string                `json:"activityPurpose"`
	DestinationCity      string                `json:"destinationCity"`
	SpdDate              string                `json:"spdDate"`
	DepartureDate        string                `json:"departureDate"`
	ReturnDate           string                `json:"returnDate"`
	ReceiptSignatureDate string                `json:"receiptSignatureDate"`
	Assignees            []rawAssigneeResponse `json:"assignees"`
}

type rawAssigneeResponse struct {
	Name           string           `json:"name"`
	SpdNumber      string           `json:"spd_number"`
	EmployeeNumber string           `json:"employee_number"`
	Position       string           `json:"position"`
	Rank           string           `json:"rank"`
	Transactions   []rawTransaction `json:"transactions"`
}

type rawTransaction struct {
	Name            string `json:"name"`
	Type            string `json:"type"`
	Subtype         string `json:"subtype"`
	Amount          int32  `json:"amount"`
	TotalNight      *int32 `json:"total_night,omitempty"`
	Subtotal        int32  `json:"subtotal"`
	Description     string `json:"description"`
	TransportDetail string `json:"transport_detail"`
}

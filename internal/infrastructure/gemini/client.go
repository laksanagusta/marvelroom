package gemini

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"sandbox/internal/domain/repository"
	"sandbox/internal/infrastructure/llm"
	transactionDTO "sandbox/internal/usecase/transaction"
)

// geminiModels uses the shared model list from llm package
var geminiModels = llm.GeminiModels

type Client struct {
	apiKey       string
	openAIAPIKey string
	httpClient   *http.Client
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 300 * time.Second,
		},
	}
}

// NewClientWithFallback creates a new client with OpenAI fallback support
func NewClientWithFallback(geminiAPIKey, openAIAPIKey string) *Client {
	return &Client{
		apiKey:       geminiAPIKey,
		openAIAPIKey: openAIAPIKey,
		httpClient: &http.Client{
			Timeout: 300 * time.Second,
		},
	}
}

// SetOpenAIAPIKey sets the OpenAI API key for fallback
func (c *Client) SetOpenAIAPIKey(apiKey string) {
	c.openAIAPIKey = apiKey
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

	// Try each Gemini model in order
	var lastError error
	for _, model := range geminiModels {
		log.Printf("[GeminiClient] Trying model: %s", model)

		result, err := c.tryGeminiModel(ctx, documents, prompt, model)
		if err != nil {
			log.Printf("[GeminiClient] Model %s failed: %v", model, err)
			lastError = err
			continue
		}

		log.Printf("[GeminiClient] Successfully used model: %s", model)
		return result, nil
	}

	// If all Gemini models failed, try OpenAI GPT-4o-mini as last resort
	if c.openAIAPIKey != "" {
		log.Printf("[GeminiClient] All Gemini models failed, falling back to GPT-4o-mini")

		result, err := c.tryOpenAI(ctx, documents, prompt)
		if err != nil {
			log.Printf("[GeminiClient] GPT-4o-mini also failed: %v", err)
			return nil, fmt.Errorf("all LLM models failed (last error from GPT-4o-mini: %w)", err)
		}

		log.Printf("[GeminiClient] Successfully used GPT-4o-mini as fallback")
		return result, nil
	}

	return nil, fmt.Errorf("all Gemini models failed (OpenAI fallback not configured): %w", lastError)
}

// tryGeminiModel attempts to use a specific Gemini model
func (c *Client) tryGeminiModel(ctx context.Context, documents []repository.Document, prompt string, model string) (interface{}, error) {
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

	apiURL := c.getGeminiModelURL(model)
	req, err := http.NewRequestWithContext(ctx, "POST", apiURL, bytes.NewBuffer(jsonBody))
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

// tryOpenAI attempts to use OpenAI GPT-4o-mini as fallback
func (c *Client) tryOpenAI(ctx context.Context, documents []repository.Document, prompt string) (interface{}, error) {
	// Build content array for OpenAI
	content := []map[string]interface{}{
		{"type": "text", "text": prompt},
	}

	// Add documents as images for vision API (only image types)
	for _, doc := range documents {
		if strings.HasPrefix(doc.MimeType, "image/") {
			base64Content := base64.StdEncoding.EncodeToString(doc.Content)
			dataURL := fmt.Sprintf("data:%s;base64,%s", doc.MimeType, base64Content)
			content = append(content, map[string]interface{}{
				"type": "image_url",
				"image_url": map[string]string{
					"url": dataURL,
				},
			})
		} else {
			// For non-image files, add a note
			content = append(content, map[string]interface{}{
				"type": "text",
				"text": fmt.Sprintf("\n[Dokumen: %s (tipe: %s, ukuran: %d bytes)]", doc.MimeType, doc.MimeType, len(doc.Content)),
			})
		}
	}

	body := map[string]interface{}{
		"model": "gpt-4o-mini",
		"messages": []map[string]interface{}{
			{
				"role":    "user",
				"content": content,
			},
		},
		"max_tokens": 4096,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal OpenAI request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", llm.OpenAIAPIURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create OpenAI request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.openAIAPIKey))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call OpenAI API: %w", err)
	}
	defer resp.Body.Close()

	bodyResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read OpenAI response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OpenAI API error (status %d): %s", resp.StatusCode, string(bodyResp))
	}

	return c.parseOpenAIResponse(bodyResp)
}

// parseOpenAIResponse parses the OpenAI response format
func (c *Client) parseOpenAIResponse(bodyResp []byte) (*transactionDTO.RecapReportDTO, error) {
	var openAIResp struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
		Error *struct {
			Message string `json:"message"`
		} `json:"error,omitempty"`
	}

	if err := json.Unmarshal(bodyResp, &openAIResp); err != nil {
		return nil, fmt.Errorf("failed to parse OpenAI response: %w", err)
	}

	if openAIResp.Error != nil {
		return nil, fmt.Errorf("OpenAI API error: %s", openAIResp.Error.Message)
	}

	if len(openAIResp.Choices) == 0 {
		return nil, errors.New("empty response from OpenAI API")
	}

	rawText := openAIResp.Choices[0].Message.Content
	cleanJSON := c.cleanJSON(rawText)

	var geminiRawReport geminiReportResponse
	if err := json.Unmarshal([]byte(cleanJSON), &geminiRawReport); err != nil {
		return nil, fmt.Errorf("failed to parse report content from OpenAI: %w (raw: %s)", err, cleanJSON)
	}

	// Reuse the same conversion logic as parseResponse
	return c.convertToRecapReport(&geminiRawReport), nil
}

// getGeminiModelURL returns the API URL for a specific Gemini model
func (c *Client) getGeminiModelURL(model string) string {
	return fmt.Sprintf("%s/%s:generateContent?key=%s", llm.GeminiAPIBaseURL, model, c.apiKey)
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
  "assignmentLetterNumber": "NOMOR_SURAT_TUGAS", -> ambil nomor surat tugas dari file surat tugas (biasanya ada di bagian atas dokumen atau di header surat)
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
          "amount": number, -> jika dia uang harian maka akan dikali 80%
          "total_night": number, -> jika dia subtypenya daily_allowance/hotel, jika hotel maka akan dihitung total malamnya, jika daily_allowance maka akan dihitung total hari
          "subtotal": number, -> hasil amount*total_night kalo dia subtypenya daily_allowance/hotel tapi kalo selain itu langsung ambil dari amount aja
	      "description" : string, -> ini adalah keterangan transaksi ini transaksi apa, misalkan gojek dari alamat1 ke alamat2, kalo hotel jelasin juga hotelnya
	      "transport_detail" : string, -> ini terisi hanya jika dia transport darat ya (pesawat tidak termasuk) 1.jika dia dari bandara soetta atau tujuannya ke bandara soetta maka valuenya menjadi "transport_asal" atau kalau dia transportasinya di jakarta juga masuk trasnport asal 2.jika mengandung bandara lain selain soetta maka valuenya adalah "transport_daerah"
	      "is_valid": boolean -> PENTING: Validasi keaslian dokumen transaksi ini (lihat instruksi validasi keaslian dokumen di bawah)
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

VALIDASI KEASLIAN DOKUMEN (DOCUMENT AUTHENTICITY VALIDATION):
SANGAT PENTING: Untuk setiap transaksi, lakukan analisis untuk mendeteksi potensi dokumen palsu/fraud:

1. VERIFIKASI VISUAL DOKUMEN:
   - Periksa font dan formatting: dokumen palsu sering menggunakan font yang tidak konsisten atau salah
   - Periksa logo dan watermark: pastikan logo resmi perusahaan/hotel/maskapai terlihat asli dan tidak cacat
   - Periksa alignment dan spacing: dokumen palsu sering memiliki masalah alignment atau spacing yang tidak konsisten

2. VERIFIKASI DATA TRANSAKSI:
   - Periksa apakah tanggal transaksi masuk akal dan sesuai dengan periode perjalanan dinas
   - Periksa apakah harga wajar untuk jenis layanan dan lokasi yang disebutkan
   - Periksa apakah ada inkonsistensi data (misalnya: hotel di kota A tapi invoice menyebut kota B)
   - Periksa nomor invoice/booking reference: dokumen palsu sering menggunakan nomor yang tidak valid atau format yang salah

3. INDIKATOR DOKUMEN PALSU:
   - Font atau ukuran text yang tidak konsisten dalam satu dokumen
   - Kesalahan ejaan pada nama perusahaan resmi (hotel, maskapai, perusahaan, dll)
   - Format tanggal yang tidak standar atau inkonsisten
   - Harga yang sangat tidak wajar (terlalu murah atau terlalu mahal)
   - Logo yang buram, terpotong, atau terlihat hasil edit
   - Kualitas gambar dokumen yang sangat rendah atau tampak di-scan ulang berkali-kali
   - Informasi yang bertentangan dalam satu dokumen (misalnya check-in dan check-out date yang tidak masuk akal)
   - Dokumen yang terlihat seperti hasil edit software (misalnya ada jejak copy-paste, layer yang terlihat)

4. SET is_valid FIELD:
   - is_valid: true  -> Jika dokumen terlihat asli dan tidak ada indikator fraud
   - is_valid: false -> Jika ditemukan indikator dokumen palsu/fraud, atau ada kejanggalan sedikit saja

Cek dengan super teliti jika Anda merasa dokumen mencurigakan atau tidak yakin, lebih baik set is_valid: false untuk kehati-hatian.

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
14,JAWA TENGAH,OH,Rp370.000,Rp150.000Rp110.000
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
  "assignmentLetterNumber": "NOMOR_SURAT_TUGAS", -> ambil nomor surat tugas dari file surat tugas (biasanya ada di bagian atas dokumen atau di header surat)
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

// ExtractVaccineRecommendations extracts vaccine information from CDC HTML
func (c *Client) ExtractVaccineRecommendations(ctx context.Context, htmlContent string) (map[string]interface{}, error) {
	if c.apiKey == "" {
		return nil, errors.New("GEMINI_API_KEY is not configured")
	}

	prompt := c.getVaccineExtractionPrompt()
	fullPrompt := prompt + fmt.Sprintf("\n\nHTML Content to analyze:\n\n%s", htmlContent)

	// Try each Gemini model in order
	var lastError error
	for _, model := range geminiModels {
		log.Printf("[GeminiClient] ExtractVaccineRecommendations trying model: %s", model)

		result, err := c.tryGeminiVaccineModel(ctx, fullPrompt, model)
		if err != nil {
			log.Printf("[GeminiClient] Model %s failed for vaccine extraction: %v", model, err)
			lastError = err
			continue
		}

		log.Printf("[GeminiClient] Successfully used model %s for vaccine extraction", model)
		return result, nil
	}

	// If all Gemini models failed, try OpenAI GPT-4o-mini as last resort
	if c.openAIAPIKey != "" {
		log.Printf("[GeminiClient] All Gemini models failed for vaccine extraction, falling back to GPT-4o-mini")

		result, err := c.tryOpenAIVaccine(ctx, fullPrompt)
		if err != nil {
			log.Printf("[GeminiClient] GPT-4o-mini also failed for vaccine extraction: %v", err)
			return nil, fmt.Errorf("all LLM models failed for vaccine extraction (last error from GPT-4o-mini: %w)", err)
		}

		log.Printf("[GeminiClient] Successfully used GPT-4o-mini as fallback for vaccine extraction")
		return result, nil
	}

	return nil, fmt.Errorf("all Gemini models failed for vaccine extraction (OpenAI fallback not configured): %w", lastError)
}

// tryGeminiVaccineModel attempts vaccine extraction with a specific Gemini model
func (c *Client) tryGeminiVaccineModel(ctx context.Context, fullPrompt string, model string) (map[string]interface{}, error) {
	parts := []map[string]interface{}{
		{"text": fullPrompt},
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

	apiURL := c.getGeminiModelURL(model)
	req, err := http.NewRequestWithContext(ctx, "POST", apiURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
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

	return c.parseVaccineResponse(bodyResp)
}

// tryOpenAIVaccine attempts vaccine extraction with OpenAI GPT-4o-mini
func (c *Client) tryOpenAIVaccine(ctx context.Context, fullPrompt string) (map[string]interface{}, error) {
	body := map[string]interface{}{
		"model": "gpt-4o-mini",
		"messages": []map[string]interface{}{
			{
				"role":    "user",
				"content": fullPrompt,
			},
		},
		"max_tokens": 4096,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal OpenAI request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", llm.OpenAIAPIURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create OpenAI request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.openAIAPIKey))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call OpenAI API: %w", err)
	}
	defer resp.Body.Close()

	bodyResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read OpenAI response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OpenAI API error (status %d): %s", resp.StatusCode, string(bodyResp))
	}

	return c.parseOpenAIVaccineResponse(bodyResp)
}

// parseOpenAIVaccineResponse parses OpenAI response for vaccine extraction
func (c *Client) parseOpenAIVaccineResponse(bodyResp []byte) (map[string]interface{}, error) {
	var openAIResp struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
		Error *struct {
			Message string `json:"message"`
		} `json:"error,omitempty"`
	}

	if err := json.Unmarshal(bodyResp, &openAIResp); err != nil {
		return nil, fmt.Errorf("failed to parse OpenAI response: %w", err)
	}

	if openAIResp.Error != nil {
		return nil, fmt.Errorf("OpenAI API error: %s", openAIResp.Error.Message)
	}

	if len(openAIResp.Choices) == 0 {
		return nil, errors.New("empty response from OpenAI API")
	}

	rawText := openAIResp.Choices[0].Message.Content
	cleanJSON := c.cleanJSON(rawText)

	var result map[string]interface{}
	if err := json.Unmarshal([]byte(cleanJSON), &result); err != nil {
		return nil, fmt.Errorf("failed to parse vaccine JSON content from OpenAI: %w (raw: %s)", err, cleanJSON)
	}

	return result, nil
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

	return c.convertToRecapReport(&geminiRawReport), nil
}

// convertToRecapReport converts geminiReportResponse to RecapReportDTO
func (c *Client) convertToRecapReport(geminiRawReport *geminiReportResponse) *transactionDTO.RecapReportDTO {
	geminiRawReport.ReceiptSignatureDate = time.Now().Format("2006-01-02")

	assignees := make([]transactionDTO.AssigneeDTO, 0, len(geminiRawReport.Assignees))
	for _, rawAssignee := range geminiRawReport.Assignees {
		transactionsDTO := make([]transactionDTO.TransactionDTO, 0, len(rawAssignee.Transactions))
		for _, rawTx := range rawAssignee.Transactions {
			// Default isValid to true if not provided by AI
			isValid := true
			if rawTx.IsValid != nil {
				isValid = *rawTx.IsValid
			}

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
				IsValid:         isValid,
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
		StartDate:              geminiRawReport.StartDate,
		EndDate:                geminiRawReport.EndDate,
		ActivityPurpose:        geminiRawReport.ActivityPurpose,
		DestinationCity:        geminiRawReport.DestinationCity,
		SpdDate:                geminiRawReport.SpdDate,
		DepartureDate:          geminiRawReport.DepartureDate,
		ReturnDate:             geminiRawReport.ReturnDate,
		AssignmentLetterNumber: geminiRawReport.AssignmentLetterNumber,
		ReceiptSignatureDate:   geminiRawReport.ReceiptSignatureDate,
		Assignees:              assignees,
	}
}

func (c *Client) getVaccineExtractionPrompt() string {
	return `You are analyzing PRE-EXTRACTED vaccine-related HTML content from a CDC travel destination page. The content has been filtered to include only relevant vaccine tables and health sections.

Extract all vaccine and health information from this optimized content and return it as a valid JSON object with this structure:

{
  "countryName": "Country name from the page",
  "requiredVaccines": [
    {
      "name": "Vaccine name (e.g., Yellow Fever)",
      "description": "Why it's required or any details",
      "forWho": "Who needs it (e.g., All travelers, travelers from certain countries)"
    }
  ],
  "recommendedVaccines": [
    {
      "name": "Vaccine name (e.g., Hepatitis A)",
      "description": "Why it's recommended",
      "forWho": "Who should get it"
    }
  ],
  "considerVaccines": [
    {
      "name": "Vaccine name (e.g., Japanese Encephalitis)",
      "description": "When to consider it",
      "forWho": "Specific groups who should consider"
    }
  ],
  "malariaInfo": {
    "risk": "None/Low/Moderate/High - describe malaria risk",
    "prophylaxis": "Whether malaria prophylaxis is recommended and details"
  },
  "healthNotice": "Any important health notices or warnings",
  "lastUpdated": "Last updated date if mentioned"
}

CRITICAL GUIDELINES:
- Return ONLY valid JSON, no additional text
- Do not wrap JSON in quotes or markdown code blocks
- The HTML content is already filtered for relevance - focus on extracting structured data
- If a section has no information, return empty arrays []
- Use standard vaccine names: "Hepatitis A", "Hepatitis B", "Typhoid", "MMR", "Japanese encephalitis", "Yellow Fever", "Rabies", "Influenza", "COVID-19", "Polio", "Tetanus", "Diphtheria", "Pertussis", "Chickenpox (Varicella)", "Cholera"
- Convert "required", "mandatory", "for entry" to requiredVaccines
- Convert "recommended", "CDC recommends" to recommendedVaccines
- Convert "consider", "some travelers", "certain travelers" to considerVaccines

TABLE EXTRACTION INSTRUCTIONS:
- Look specifically for tables with class="disease" and id starting with "dest-vm-"
- Extract data from table rows: first column = disease/vaccine name, second column = recommendations, third column = clinical guidance
- For "Routine vaccines" in first column, extract the individual vaccines listed in the recommendations
- Pay attention to "Vaccine is not recommended" vs "Vaccine is not required" for proper categorization
- Extract malaria transmission areas, drug resistance, species, and recommended chemoprophylaxis when present`
}

func (c *Client) parseVaccineResponse(bodyResp []byte) (map[string]interface{}, error) {
	var geminiAPIResponse geminiResponse
	if err := json.Unmarshal(bodyResp, &geminiAPIResponse); err != nil {
		return nil, fmt.Errorf("failed to parse Gemini API response wrapper: %w", err)
	}

	if len(geminiAPIResponse.Candidates) == 0 || len(geminiAPIResponse.Candidates[0].Content.Parts) == 0 {
		return nil, errors.New("empty response candidates or parts from Gemini API")
	}

	rawText := geminiAPIResponse.Candidates[0].Content.Parts[0].Text
	cleanJSON := c.cleanJSON(rawText)

	var result map[string]interface{}
	if err := json.Unmarshal([]byte(cleanJSON), &result); err != nil {
		return nil, fmt.Errorf("failed to parse vaccine JSON content: %w (raw: %s)", err, cleanJSON)
	}

	return result, nil
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
	StartDate              string                `json:"startDate"`
	EndDate                string                `json:"endDate"`
	ActivityPurpose        string                `json:"activityPurpose"`
	DestinationCity        string                `json:"destinationCity"`
	SpdDate                string                `json:"spdDate"`
	DepartureDate          string                `json:"departureDate"`
	ReturnDate             string                `json:"returnDate"`
	AssignmentLetterNumber string                `json:"assignmentLetterNumber"`
	ReceiptSignatureDate   string                `json:"receiptSignatureDate"`
	Assignees              []rawAssigneeResponse `json:"assignees"`
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
	IsValid         *bool  `json:"is_valid,omitempty"` // Document authenticity validation from AI
}

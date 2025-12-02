package work_paper

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"sort"
	"strings"
	"time"

	"sandbox/internal/domain/entity"
	"sandbox/internal/domain/service"

	"github.com/fumiama/go-docx"
)

type GenerateWorkPaperDocxUseCase struct {
	deskService service.DeskService
}

func NewGenerateWorkPaperDocxUseCase(deskService service.DeskService) *GenerateWorkPaperDocxUseCase {
	return &GenerateWorkPaperDocxUseCase{
		deskService: deskService,
	}
}

func (uc *GenerateWorkPaperDocxUseCase) Execute(ctx context.Context, workPaperID string) ([]byte, error) {
	// 1. Fetch Data
	wp, err := uc.deskService.GetWorkPaper(ctx, workPaperID)
	if err != nil {
		return nil, fmt.Errorf("failed to get work paper: %w", err)
	}

	notes, err := uc.deskService.GetWorkPaperNotes(ctx, workPaperID)
	if err != nil {
		return nil, fmt.Errorf("failed to get notes: %w", err)
	}

	signatures, err := uc.deskService.GetWorkPaperSignatures(ctx, workPaperID)
	if err != nil {
		return nil, fmt.Errorf("failed to get signatures: %w", err)
	}

	// 2. Group Notes by Type
	groupedNotes := make(map[string][]*entity.WorkPaperNote)
	for _, note := range notes {
		if note.MasterItem != nil {
			groupedNotes[note.MasterItem.Type] = append(groupedNotes[note.MasterItem.Type], note)
		}
	}

	// // Sort types (A, B, C...)
	var types []string
	for t := range groupedNotes {
		types = append(types, t)
	}
	sort.Strings(types)

	// 3. Generate DOCX
	f := docx.New().WithDefaultTheme()

	// Title
	p := f.AddParagraph().Justification("center")
	p.AddText("KEMENTERIAN KESEHATAN RI").Bold().Size("24")

	p = f.AddParagraph().Justification("center")
	p.AddText("DIREKTORAT JENDERAL P2P").Bold().Size("24")

	p = f.AddParagraph().Justification("center")
	p.AddText("CATATAN HASIL REVIU (CHR)").Bold().Size("24")

	p = f.AddParagraph().Justification("center")
	p.AddText("LAPORAN KINERJA INSTANSI PEMERINTAH").Bold().Size("24")

	f.AddParagraph() // Spacer

	// Intro
	now := time.Now()
	dateStr := dateToIndonesianText(now)

	p = f.AddParagraph()
	p.AddText(fmt.Sprintf("Pada hari ini %s, kami yang bertanda tangan di bawah ini telah melakukan reviu terhadap Laporan Kinerja Semester %d tahun %d Unit Kerja/Satuan Kerja Ditjen P2P sebagai berikut:", dateStr, wp.Semester, wp.Year))

	f.AddParagraph()

	// // Entity Name
	orgName := "Unknown Organization"
	if wp.Organization != nil {
		orgName = wp.Organization.Name
	}

	p = f.AddParagraph()
	p.AddText("Entitas Akuntabilitas : " + orgName).Bold()

	f.AddParagraph()

	// Calculate rows for Notes Table
	rowCount := 2
	for _, t := range types {
		rowCount += 1 // Section Header
		rowCount += len(groupedNotes[t])
	}

	// Table
	table := f.AddTable(rowCount, 1, 9000, nil)

	// Ensure all cells have paragraphs
	for _, row := range table.TableRows {
		for _, cell := range row.TableCells {
			if len(cell.Paragraphs) == 0 {
				cell.AddParagraph()
			}
		}
	}

	currentRow := 0

	// Header
	cell := table.TableRows[currentRow].TableCells[0]
	// // cell.Shade("clear", "auto", "E0E0E0") // Disable shade for debugging
	// Use existing paragraph if available
	if len(cell.Paragraphs) > 0 {
		p = cell.Paragraphs[0]
	} else {
		p = cell.AddParagraph()
	}
	p.AddText("Uraian Catatan Hasil").Bold()
	currentRow++

	// // Intro in table
	cell = table.TableRows[currentRow].TableCells[0]
	if len(cell.Paragraphs) > 0 {
		p = cell.Paragraphs[0]
	} else {
		p = cell.AddParagraph()
	}
	p.AddText("Pelaksanaan Reviu Laporan Kinerja Instansi Pemerintah berdasarkan pada Surat Sesditjen...")
	currentRow++

	// Loop Types
	for _, t := range types {
		// Section Header
		cell = table.TableRows[currentRow].TableCells[0]
		// cell.Shade("clear", "auto", "E0E0E0") // Disable shade
		if len(cell.Paragraphs) > 0 {
			p = cell.Paragraphs[0]
		} else {
			p = cell.AddParagraph()
		}
		sectionTitle := fmt.Sprintf("%s. %s", t, getSectionTitle(t))
		p.AddText(sectionTitle).Bold()
		currentRow++

		// Notes
		typeNotes := groupedNotes[t]
		sort.Slice(typeNotes, func(i, j int) bool {
			return typeNotes[i].MasterItem.SortOrder < typeNotes[j].MasterItem.SortOrder
		})

		for _, note := range typeNotes {
			cell = table.TableRows[currentRow].TableCells[0]
			if len(cell.Paragraphs) > 0 {
				p = cell.Paragraphs[0]
			} else {
				p = cell.AddParagraph()
			}
			noteText := note.GetNotes()
			if noteText == "" {
				noteText = "Belum ada catatan."
			}
			p.AddText("- " + noteText)
			currentRow++
		}
	}

	f.AddParagraph() // Spacer

	// Signatures
	// 3 rows, 2 columns
	sigTable := f.AddTable(3, 2, 9000, nil)

	// Initialize all cells with an empty paragraph to avoid corruption
	for _, row := range sigTable.TableRows {
		for _, cell := range row.TableCells {
			if len(cell.Paragraphs) == 0 {
				cell.AddParagraph()
			}
		}
	}

	// Header Row
	sigTable.TableRows[0].TableCells[0].Paragraphs[0].AddText("Petugas Satker").Bold()
	sigTable.TableRows[0].TableCells[1].Paragraphs[0].AddText("Petugas Reviu").Bold()

	// Helper to add signature
	addSignatureToCell := func(cell *docx.WTableCell, sig *entity.WorkPaperSignature) {
		if sig == nil {
			return
		}

		// Image
		data := sig.GetSignatureData()
		if data.SignatureImage != "" {
			b64 := data.SignatureImage
			if idx := strings.Index(b64, ","); idx != -1 {
				b64 = b64[idx+1:]
			}
			imgBytes, err := base64.StdEncoding.DecodeString(b64)
			if err == nil {
				// Add new paragraph for image
				p := cell.AddParagraph()
				_, err := p.AddInlineDrawing(imgBytes)
				if err != nil {
					p.AddText("[Error adding image]")
				}
				// // Placeholder for debugging
				// p := cell.AddParagraph()
				// p.AddText("[Signature Image Placeholder]")
				// _ = imgBytes // suppress unused error
			}
		}

		// Name
		p := cell.AddParagraph()
		p.AddText(sig.UserName)

		// NIP
		p = cell.AddParagraph()
		p.AddText("NIP " + sig.UserID)
	}

	if len(signatures) > 0 {
		addSignatureToCell(sigTable.TableRows[1].TableCells[0], signatures[0])
	}
	if len(signatures) > 1 {
		addSignatureToCell(sigTable.TableRows[1].TableCells[1], signatures[1])
	}

	// Bottom Signature (Ketua Tim)
	if len(signatures) > 2 {
		cell := sigTable.TableRows[2].TableCells[0]
		cell.Paragraphs[0].AddText("Ketua Tim Kerja Informasi dan Kerjasama").Bold()
		addSignatureToCell(cell, signatures[2])
	}

	// Write to buffer
	var buf bytes.Buffer
	_, err = f.WriteTo(&buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func dateToIndonesianText(t time.Time) string {
	days := []string{"Minggu", "Senin", "Selasa", "Rabu", "Kamis", "Jumat", "Sabtu"}
	months := []string{"", "Januari", "Februari", "Maret", "April", "Mei", "Juni", "Juli", "Agustus", "September", "Oktober", "November", "Desember"}

	dayName := days[t.Weekday()]
	monthName := months[t.Month()]

	return fmt.Sprintf("%s tanggal %d, bulan %s, tahun %d", dayName, t.Day(), monthName, t.Year())
}

func getSectionTitle(t string) string {
	switch t {
	case "A":
		return "Eksistensi Laporan Kinerja"
	case "B":
		return "Format Laporan Kinerja"
	case "C":
		return "Mekanisme Penyusunan Laporan Kinerja"
	case "D":
		return "Substansi Laporan Kinerja"
	case "E":
		return "Pemanfaatan Laporan Kinerja"
	default:
		return "Lain-lain"
	}
}

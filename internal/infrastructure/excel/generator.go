package excel

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"sandbox/pkg/constants"
	"sandbox/utils"

	"github.com/xuri/excelize/v2"
)

const (
	workUnit            = "Setditjen Penanggulangan Penyakit"
	commitmentOfficer   = "Ruly Wahyuni, SE, MKM"
	commitmentOfficerId = "197508142000032001"

	expenditureTreasurer   = "Fatmawati Husain, SE, M.Ak"
	expenditureTreasurerId = "198608202005012002"

	payer   = "Marsaulina Siahaan, SE"
	payerId = "197101261997032002"
)

type Generator struct{}

func NewGenerator() *Generator {
	return &Generator{}
}

func (g *Generator) generateTitle(f *excelize.File, sheetName string) error {
	if err := f.MergeCell(sheetName, "A2", "U2"); err != nil {
		return err
	}
	if err := f.MergeCell(sheetName, "A3", "U3"); err != nil {
		return err
	}
	if err := f.MergeCell(sheetName, "A4", "U4"); err != nil {
		return err
	}

	titleStyle, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Font: &excelize.Font{
			Bold:   true,
			Size:   12,
			Family: "Tahoma",
		},
	})
	if err != nil {
		return err
	}
	subTitleStyle, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Font: &excelize.Font{
			Bold:   true,
			Size:   12,
			Family: "Tahoma",
		},
	})
	if err != nil {
		return err
	}
	accountStyle, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Font: &excelize.Font{
			Bold:   true,
			Size:   12,
			Family: "Tahoma",
		},
	})
	if err != nil {
		return err
	}

	if sheetName == "PEMANTAUAN REKAP RAMPUNG" {
		f.SetCellValue(sheetName, "A2", "Rekapitulasi Biaya Perjalanan Dinas Rampung")
	} else {
		f.SetCellValue(sheetName, "A2", "Rekapitulasi Uang Muka Biaya Perjalanan Dinas")
	}

	f.SetCellValue(sheetName, "A3", "Rekapitulasi Biaya Perjalanan Dinas dalam Rangka Pemantauan dan Evaluasi Pelaksanaan Program di Daerah")
	f.SetCellValue(sheetName, "A4", "AKUN : 4815.EBD.953.501.B.524111")

	if err := f.SetCellStyle(sheetName, "A2", "U2", titleStyle); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "A3", "U3", subTitleStyle); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "A4", "U4", accountStyle); err != nil {
		return err
	}
	return nil
}

func (g *Generator) generateTableHeader(f *excelize.File, sheetName string) error {
	headerStyle, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
			WrapText:   true,
		},
		Font: &excelize.Font{
			Bold:   true,
			Size:   10,
			Family: "Tahoma",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 2},
			{Type: "top", Color: "000000", Style: 2},
			{Type: "bottom", Color: "000000", Style: 2},
			{Type: "right", Color: "000000", Style: 2},
		},
	})
	if err != nil {
		return err
	}

	if err := f.SetCellValue(sheetName, "A8", "No"); err != nil {
		return err
	}
	if err := f.MergeCell(sheetName, "A8", "A10"); err != nil {
		return err
	}

	if err := f.SetCellValue(sheetName, "B8", "Nama"); err != nil {
		return err
	}
	if err := f.MergeCell(sheetName, "B8", "B10"); err != nil {
		return err
	}

	if err := f.SetCellValue(sheetName, "C8", "NIP"); err != nil {
		return err
	}
	if err := f.MergeCell(sheetName, "C8", "C10"); err != nil {
		return err
	}

	if err := f.SetCellValue(sheetName, "D8", "Jabatan"); err != nil {
		return err
	}
	if err := f.MergeCell(sheetName, "D8", "D10"); err != nil {
		return err
	}

	if err := f.SetCellValue(sheetName, "E8", "Gol"); err != nil {
		return err
	}
	if err := f.MergeCell(sheetName, "E8", "E10"); err != nil {
		return err
	}

	if err := f.SetCellValue(sheetName, "F8", "Tujuan"); err != nil {
		return err
	}
	if err := f.MergeCell(sheetName, "F8", "F10"); err != nil {
		return err
	}

	if err := f.SetCellValue(sheetName, "G8", "Tanggal"); err != nil {
		return err
	}
	if err := f.MergeCell(sheetName, "G8", "G10"); err != nil {
		return err
	}

	if err := f.SetCellValue(sheetName, "H8", "Uang Harian"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "H9", "Konstanta : 008448"); err != nil {
		return err
	}
	if err := f.MergeCell(sheetName, "H9", "K9"); err != nil {
		return err
	}

	if err := f.SetCellValue(sheetName, "L8", "Penginapan"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "L9", "Konstanta : 008447"); err != nil {
		return err
	}
	if err := f.MergeCell(sheetName, "L9", "O9"); err != nil {
		return err
	}

	if err := f.SetCellValue(sheetName, "P8", "Transport"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "P9", "Konstanta : 008446"); err != nil {
		return err
	}
	if err := f.MergeCell(sheetName, "P9", "T9"); err != nil {
		return err
	}

	if sheetName == "PEMANTAUAN REKAP RAMPUNG" {
		if err := f.SetCellValue(sheetName, "U8", "Jumlah SPJ Rampung (Rp)"); err != nil {
			return err
		}
		if err := f.SetCellValue(sheetName, "V8", "Jumlah SPJ Uang Muka (Rp)"); err != nil {
			return err
		}
		if err := f.SetCellValue(sheetName, "W8", "Jumlah Dibayarkan (Rp)"); err != nil {
			return err
		}
	} else {
		if err := f.SetCellValue(sheetName, "U8", "Jumlah Dibayarkan (Rp)"); err != nil {
			return err
		}
	}

	if err := f.SetCellValue(sheetName, "AA8", "No SPD"); err != nil {
		return err
	}

	// Sub-headers for Uang Harian
	if err := f.SetCellValue(sheetName, "H10", "Jml Hari"); err != nil {
		return err
	}
	if err := f.MergeCell(sheetName, "H10", "I10"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "J10", "Perhari"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "K10", "Jumlah"); err != nil {
		return err
	}

	// Sub-headers for Penginapan
	if err := f.SetCellValue(sheetName, "L10", "Jml Hari"); err != nil {
		return err
	}
	if err := f.MergeCell(sheetName, "L10", "M10"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "N10", "Perhari"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "O10", "Jumlah"); err != nil {
		return err
	}

	if err := f.SetCellValue(sheetName, "P10", "Tiket Pesawat"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "Q10", "Transport Asal"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "R10", "Transport Daerah"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "S10", "Transport Darat"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "T10", "Jumlah"); err != nil {
		return err
	}

	if err := f.MergeCell(sheetName, "A8", "A9"); err != nil {
		return err
	}
	if err := f.MergeCell(sheetName, "B8", "B9"); err != nil {
		return err
	}
	if err := f.MergeCell(sheetName, "C8", "C9"); err != nil {
		return err
	}
	if err := f.MergeCell(sheetName, "D8", "D9"); err != nil {
		return err
	}
	if err := f.MergeCell(sheetName, "E8", "E9"); err != nil {
		return err
	}
	if err := f.MergeCell(sheetName, "F8", "F9"); err != nil {
		return err
	}
	if err := f.MergeCell(sheetName, "G8", "G9"); err != nil {
		return err
	}

	if err := f.MergeCell(sheetName, "H8", "K8"); err != nil {
		return err
	}
	if err := f.MergeCell(sheetName, "L8", "O8"); err != nil {
		return err
	}
	if err := f.MergeCell(sheetName, "P8", "T8"); err != nil {
		return err
	}

	if err := f.MergeCell(sheetName, "A9", "A9"); err != nil {
		return err
	}

	if err := f.MergeCell(sheetName, "U8", "U10"); err != nil {
		return err
	}

	if sheetName == "PEMANTAUAN REKAP RAMPUNG" {
		if err := f.MergeCell(sheetName, "W8", "W10"); err != nil {
			return err
		}

		if err := f.MergeCell(sheetName, "V8", "V10"); err != nil {
			return err
		}
	}

	if sheetName == "PEMANTAUAN REKAP UANG MUKA" {
		if err := f.SetCellStyle(sheetName, "A8", "U10", headerStyle); err != nil {
			return err
		}
	} else {
		if err := f.SetCellStyle(sheetName, "A8", "W10", headerStyle); err != nil {
			return err
		}
	}

	if err := f.SetCellStyle(sheetName, "H9", "T9", g.dynamicStyle(f, []string{"top", "bottom", "left", "right"}, true, false, 2, "left", 0, false, 0)); err != nil {
		return err
	}

	for i := 8; i < 11; i++ {
		if err := f.SetRowHeight(sheetName, i, 28); err != nil {
			return err
		}
	}

	return nil
}

type PersonRecap struct {
	Name                    string
	NIP                     string
	Jabatan                 string
	Gol                     string
	Tujuan                  string
	Tanggal                 string
	NoSpd                   string
	UMUangHarianJmlHari     int32
	UMUangHarianPerhari     int32
	UMUangHarianJumlah      int32
	UMPenginapanJmlHari     int32
	UMPenginapanPerhari     int32
	UMPenginapanJumlah      int32
	UMTransportTiketPesawat int32
	UMTransportAsal         int32
	UMTransportDaerah       int32
	UMTransportDarat        int32
	UMTransportJumlah       int32
	UMTotalDibayarkan       int32

	RUangHarianJmlHari     int32
	RUangHarianPerhari     int32
	RUangHarianJumlah      int32
	RPenginapanJmlHari     int32
	RPenginapanPerhari     int32
	RPenginapanJumlah      int32
	RTransportTiketPesawat int32
	RTransportAsal         int32
	RTransportDaerah       int32
	RTransportDarat        int32
	RTransportJumlah       int32
	RTotalDibayarkan       int32
}

func (g *Generator) generateTableData(f *excelize.File, sheetName string, req RecapReport, currentRow int) (int, error) {
	textStyle, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Vertical: "center",
		},
		Font: &excelize.Font{
			Size:   10,
			Family: "Tahoma",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 2},
			{Type: "top", Color: "000000", Style: 2},
			{Type: "bottom", Color: "000000", Style: 2},
			{Type: "right", Color: "000000", Style: 2},
		},
	})
	if err != nil {
		return currentRow, err
	}

	numberStyle, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Vertical: "center",
		},
		Font: &excelize.Font{
			Size:   10,
			Family: "Tahoma",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 2},
			{Type: "top", Color: "000000", Style: 2},
			{Type: "bottom", Color: "000000", Style: 2},
			{Type: "right", Color: "000000", Style: 2},
		},
		NumFmt: 3,
	})
	if err != nil {
		return currentRow, err
	}

	personData := make(map[string]*PersonRecap)
	constUangHarianJmlHari := int32(2)
	constUangHarianPerhari := int32(688000)
	constUangHarianJumlah := constUangHarianJmlHari * constUangHarianPerhari

	jsn, _ := json.Marshal(req)
	fmt.Println(string(jsn))

	for _, assignee := range req.Assignees {
		// Use employee_number (NIP) as the primary identifier, fallback to employee_id
		employeeIdentifier := assignee.EmployeeNumber
		if employeeIdentifier == "" {
			employeeIdentifier = assignee.EmployeeID
		}

		if employeeIdentifier == "" {
			continue
		}

		data, exists := personData[employeeIdentifier]
		if !exists {
			data = &PersonRecap{
				Name:                assignee.Name,
				NIP:                 employeeIdentifier, // Use employee_number (NIP) for NIP field
				Jabatan:             assignee.Position,
				Gol:                 assignee.Rank,
				Tujuan:              req.DestinationCity,
				Tanggal:             req.DepartureDate,
				NoSpd:               assignee.SpdNumber,
				UMUangHarianJmlHari: constUangHarianJmlHari,
				UMUangHarianPerhari: constUangHarianPerhari,
				UMUangHarianJumlah:  constUangHarianJumlah,
			}
			personData[employeeIdentifier] = data
		}

		for _, tx := range assignee.Transactions {
			// Skip if transaction has zero amount
			if tx.Subtotal <= 0 {
				continue
			}
			switch constants.TransactionType(strings.ToLower(tx.Type)) {
			case constants.TransactionTypeAccommodation:
				if tx.PaymentType == "uang muka" {
					if tx.TotalNight != nil && *tx.TotalNight > 0 {
						data.UMPenginapanJmlHari += *tx.TotalNight
						data.RPenginapanJmlHari += *tx.TotalNight
					}
					if tx.Amount > 0 {
						data.UMPenginapanPerhari = tx.Amount
						data.RPenginapanPerhari = tx.Amount
					}
					data.UMPenginapanJumlah += tx.Subtotal
					data.RPenginapanJumlah += tx.Subtotal
				} else {
					if tx.TotalNight != nil && *tx.TotalNight > 0 {
						data.RPenginapanJmlHari += *tx.TotalNight
					}
					if tx.Amount > 0 {
						data.RPenginapanPerhari = tx.Amount
					}
					data.RPenginapanJumlah += tx.Subtotal
				}
			case constants.TransactionTypeTransport:
				if tx.PaymentType == "uang muka" {
					if strings.ToLower(tx.Subtype) == "flight" {
						data.UMTransportTiketPesawat += tx.Subtotal
						data.RTransportTiketPesawat += tx.Subtotal
					}

					if strings.ToLower(tx.Subtype) == "taxi" {
						switch strings.ToLower(tx.TransportDetail) {
						case "transport_darat":
							data.UMTransportDarat += tx.Subtotal
							data.RTransportDarat += tx.Subtotal
						case "transport_asal":
							data.UMTransportAsal += tx.Subtotal
							data.RTransportAsal += tx.Subtotal
						case "transport_daerah":
							data.UMTransportDaerah += tx.Subtotal
							data.RTransportDaerah += tx.Subtotal
						}
					}

					data.UMTransportJumlah += tx.Subtotal
					data.RTransportJumlah += tx.Subtotal
				} else {
					if strings.ToLower(tx.Subtype) == "flight" {
						data.RTransportTiketPesawat += tx.Subtotal
					}

					if strings.ToLower(tx.Subtype) == "taxi" {
						switch strings.ToLower(tx.TransportDetail) {
						case "transport_darat":
							data.RTransportDarat += tx.Subtotal
						case "transport_asal":
							data.RTransportAsal += tx.Subtotal
						case "transport_daerah":
							data.RTransportDaerah += tx.Subtotal
						}
					}

					data.RTransportJumlah += tx.Subtotal
				}
			case constants.TransactionTypeOther:
				if tx.PaymentType == "uang muka" {
					data.UMTotalDibayarkan += tx.Subtotal
				} else {
					data.RTotalDibayarkan += tx.Subtotal
				}
			}

			data.UMTotalDibayarkan = data.UMUangHarianJumlah + data.UMPenginapanJumlah + data.UMTransportJumlah
		}
	}

	personNo := 1
	for _, data := range personData {
		if err := f.SetCellValue(sheetName, fmt.Sprintf("A%d", currentRow), personNo); err != nil {
			return currentRow, err
		}
		if err := f.SetCellValue(sheetName, fmt.Sprintf("B%d", currentRow), data.Name); err != nil {
			return currentRow, err
		}
		if err := f.SetCellValue(sheetName, fmt.Sprintf("C%d", currentRow), data.NIP); err != nil {
			return currentRow, err
		}
		if err := f.SetCellValue(sheetName, fmt.Sprintf("D%d", currentRow), data.Jabatan); err != nil {
			return currentRow, err
		}
		if err := f.SetCellValue(sheetName, fmt.Sprintf("E%d", currentRow), data.Gol); err != nil {
			return currentRow, err
		}
		if err := f.SetCellValue(sheetName, fmt.Sprintf("F%d", currentRow), data.Tujuan); err != nil {
			return currentRow, err
		}
		if err := f.SetCellValue(sheetName, fmt.Sprintf("G%d", currentRow), data.Tanggal); err != nil {
			return currentRow, err
		}
		if err := f.SetCellValue(sheetName, fmt.Sprintf("I%d", currentRow), "Hari"); err != nil {
			return currentRow, err
		}
		if err := f.SetCellValue(sheetName, fmt.Sprintf("M%d", currentRow), "Hari"); err != nil {
			return currentRow, err
		}

		if err := f.SetCellValue(sheetName, fmt.Sprintf("H%d", currentRow), data.UMUangHarianJmlHari); err != nil {
			return currentRow, err
		}
		if err := f.SetCellValue(sheetName, fmt.Sprintf("J%d", currentRow), data.UMUangHarianPerhari); err != nil {
			return currentRow, err
		}
		if err := f.SetCellValue(sheetName, fmt.Sprintf("K%d", currentRow), data.UMUangHarianJumlah); err != nil {
			return currentRow, err
		}

		if sheetName == "PEMANTAUAN REKAP UANG MUKA" {
			if err := f.SetCellValue(sheetName, fmt.Sprintf("L%d", currentRow), data.UMPenginapanJmlHari); err != nil {
				return currentRow, err
			}
			if err := f.SetCellValue(sheetName, fmt.Sprintf("N%d", currentRow), data.UMPenginapanPerhari); err != nil {
				return currentRow, err
			}
			if err := f.SetCellValue(sheetName, fmt.Sprintf("O%d", currentRow), data.UMPenginapanJumlah); err != nil {
				return currentRow, err
			}
		} else {
			if err := f.SetCellValue(sheetName, fmt.Sprintf("L%d", currentRow), data.RPenginapanJmlHari); err != nil {
				return currentRow, err
			}
			if err := f.SetCellValue(sheetName, fmt.Sprintf("N%d", currentRow), data.RPenginapanPerhari); err != nil {
				return currentRow, err
			}
			if err := f.SetCellValue(sheetName, fmt.Sprintf("O%d", currentRow), data.RPenginapanJumlah); err != nil {
				return currentRow, err
			}
		}

		if sheetName == "PEMANTAUAN REKAP RAMPUNG" {
			if err := f.SetCellValue(sheetName, fmt.Sprintf("P%d", currentRow), data.RTransportTiketPesawat); err != nil {
				return currentRow, err
			}
			if err := f.SetCellValue(sheetName, fmt.Sprintf("Q%d", currentRow), data.RTransportAsal); err != nil {
				return currentRow, err
			}
			if err := f.SetCellValue(sheetName, fmt.Sprintf("R%d", currentRow), data.RTransportDaerah); err != nil {
				return currentRow, err
			}
			if err := f.SetCellValue(sheetName, fmt.Sprintf("S%d", currentRow), data.RTransportDarat); err != nil {
				return currentRow, err
			}
		} else {
			if err := f.SetCellValue(sheetName, fmt.Sprintf("P%d", currentRow), data.UMTransportTiketPesawat); err != nil {
				return currentRow, err
			}
			if err := f.SetCellValue(sheetName, fmt.Sprintf("Q%d", currentRow), data.UMTransportAsal); err != nil {
				return currentRow, err
			}
			if err := f.SetCellValue(sheetName, fmt.Sprintf("R%d", currentRow), data.UMTransportDaerah); err != nil {
				return currentRow, err
			}
			if err := f.SetCellValue(sheetName, fmt.Sprintf("S%d", currentRow), data.UMTransportDarat); err != nil {
				return currentRow, err
			}
		}

		if err := f.SetCellFormula(sheetName, fmt.Sprintf("T%d", currentRow), fmt.Sprintf("=P%d+Q%d+R%d+S%d", currentRow, currentRow, currentRow, currentRow)); err != nil {
			return currentRow, err
		}

		if sheetName == "PEMANTAUAN REKAP RAMPUNG" {
			if err := f.SetCellFormula(sheetName, fmt.Sprintf("U%d", currentRow), fmt.Sprintf("=T%d+O%d+K%d", currentRow, currentRow, currentRow)); err != nil {
				return currentRow, err
			}
			if err := f.SetCellFormula(sheetName, fmt.Sprintf("V%d", currentRow), fmt.Sprintf("='PEMANTAUAN REKAP UANG MUKA'!U%d", currentRow)); err != nil {
				return currentRow, err
			}
			if err := f.SetCellFormula(sheetName, fmt.Sprintf("W%d", currentRow), fmt.Sprintf("=U%d-V%d", currentRow, currentRow)); err != nil {
				return currentRow, err
			}
		} else {
			if err := f.SetCellFormula(sheetName, fmt.Sprintf("U%d", currentRow), fmt.Sprintf("=T%d+O%d+K%d", currentRow, currentRow, currentRow)); err != nil {
				return currentRow, err
			}
		}

		if err := f.SetCellValue(sheetName, fmt.Sprintf("AA%d", currentRow), data.NoSpd); err != nil {
			return currentRow, err
		}

		if sheetName == "PEMANTAUAN REKAP RAMPUNG" {
			if err := f.SetCellStyle(sheetName, fmt.Sprintf("A%d", currentRow), fmt.Sprintf("G%d", currentRow), textStyle); err != nil {
				return currentRow, err
			}
			if err := f.SetCellStyle(sheetName, fmt.Sprintf("I%d", currentRow), fmt.Sprintf("I%d", currentRow), textStyle); err != nil {
				return currentRow, err
			}
			if err := f.SetCellStyle(sheetName, fmt.Sprintf("M%d", currentRow), fmt.Sprintf("M%d", currentRow), textStyle); err != nil {
				return currentRow, err
			}
			if err := f.SetCellStyle(sheetName, fmt.Sprintf("AA%d", currentRow), fmt.Sprintf("AA%d", currentRow), textStyle); err != nil {
				return currentRow, err
			}

			if err := f.SetCellStyle(sheetName, fmt.Sprintf("H%d", currentRow), fmt.Sprintf("W%d", currentRow), numberStyle); err != nil {
				return currentRow, err
			}
		} else {
			if err := f.SetCellStyle(sheetName, fmt.Sprintf("A%d", currentRow), fmt.Sprintf("U%d", currentRow), textStyle); err != nil {
				return currentRow, err
			}
			if err := f.SetCellStyle(sheetName, fmt.Sprintf("AA%d", currentRow), fmt.Sprintf("AA%d", currentRow), textStyle); err != nil {
				return currentRow, err
			}

			if err := f.SetCellStyle(sheetName, fmt.Sprintf("H%d", currentRow), fmt.Sprintf("P%d", currentRow), numberStyle); err != nil {
				return currentRow, err
			}
			if err := f.SetCellStyle(sheetName, fmt.Sprintf("T%d", currentRow), fmt.Sprintf("U%d", currentRow), numberStyle); err != nil {
				return currentRow, err
			}

		}

		if err := f.SetRowHeight(sheetName, currentRow, 28); err != nil {
			return currentRow, err
		}

		currentRow++
		personNo++
	}
	return currentRow, nil
}

func (g *Generator) generateSummaryRow(f *excelize.File, sheetName string, currentRow int) error {
	summaryStyle, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "right",
			Vertical:   "center",
		},
		Font: &excelize.Font{
			Bold:   true,
			Size:   10,
			Family: "Tahoma",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 2},
			{Type: "top", Color: "000000", Style: 2},
			{Type: "bottom", Color: "000000", Style: 2},
			{Type: "right", Color: "000000", Style: 2},
		},
		NumFmt: 3,
	})

	summaryNumberStyle, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "right",
			Vertical:   "center",
		},
		Font: &excelize.Font{
			Bold:   true,
			Size:   10,
			Family: "Tahoma",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 2},
			{Type: "top", Color: "000000", Style: 2},
			{Type: "bottom", Color: "000000", Style: 2},
			{Type: "right", Color: "000000", Style: 2},
		},
		NumFmt: 3,
	})
	if err != nil {
		return err
	}

	totalRow := currentRow

	if sheetName == "PEMANTAUAN REKAP RAMPUNG" {
		if err := f.SetCellValue(sheetName, fmt.Sprintf("B%d", totalRow+1), "JUMLAH SPJ RAMPUNG"); err != nil {
			return err
		}

		if err := f.SetCellValue(sheetName, fmt.Sprintf("B%d", totalRow+2), "JUMLAH SPJ UANG MUKA"); err != nil {
			return err
		}

		if err := f.SetCellValue(sheetName, fmt.Sprintf("B%d", totalRow+3), "JUMLAH DIBAYARKAN"); err != nil {
			return err
		}

		if err := f.SetCellFormula(sheetName, fmt.Sprintf("K%d", totalRow+2), fmt.Sprintf("='PEMANTAUAN REKAP UANG MUKA'!K%d", totalRow)); err != nil {
			return err
		}
		if err := f.SetCellFormula(sheetName, fmt.Sprintf("K%d", totalRow+1), fmt.Sprintf("=SUM(K11:K%d)", totalRow-1)); err != nil {
			return err
		}

		if err := f.SetCellFormula(sheetName, fmt.Sprintf("O%d", totalRow+2), fmt.Sprintf("='PEMANTAUAN REKAP UANG MUKA'!O%d", totalRow)); err != nil {
			return err
		}
		if err := f.SetCellFormula(sheetName, fmt.Sprintf("O%d", totalRow+1), fmt.Sprintf("=SUM(O11:O%d)", totalRow-1)); err != nil {
			return err
		}

		if err := f.SetCellFormula(sheetName, fmt.Sprintf("P%d", totalRow+2), fmt.Sprintf("='PEMANTAUAN REKAP UANG MUKA'!P%d", totalRow)); err != nil {
			return err
		}
		if err := f.SetCellFormula(sheetName, fmt.Sprintf("P%d", totalRow+1), fmt.Sprintf("=SUM(P11:P%d)", totalRow-1)); err != nil {
			return err
		}

		if err := f.SetCellFormula(sheetName, fmt.Sprintf("Q%d", totalRow+2), fmt.Sprintf("='PEMANTAUAN REKAP UANG MUKA'!Q%d", totalRow)); err != nil {
			return err
		}
		if err := f.SetCellFormula(sheetName, fmt.Sprintf("Q%d", totalRow+1), fmt.Sprintf("=SUM(Q11:Q%d)", totalRow-1)); err != nil {
			return err
		}

		if err := f.SetCellFormula(sheetName, fmt.Sprintf("R%d", totalRow+2), fmt.Sprintf("='PEMANTAUAN REKAP UANG MUKA'!R%d", totalRow)); err != nil {
			return err
		}
		if err := f.SetCellFormula(sheetName, fmt.Sprintf("R%d", totalRow+1), fmt.Sprintf("=SUM(R11:R%d)", totalRow-1)); err != nil {
			return err
		}

		if err := f.SetCellFormula(sheetName, fmt.Sprintf("S%d", totalRow+2), fmt.Sprintf("='PEMANTAUAN REKAP UANG MUKA'!S%d", totalRow)); err != nil {
			return err
		}
		if err := f.SetCellFormula(sheetName, fmt.Sprintf("S%d", totalRow+1), fmt.Sprintf("=SUM(S11:S%d)", totalRow-1)); err != nil {
			return err
		}

		if err := f.SetCellFormula(sheetName, fmt.Sprintf("T%d", totalRow+2), fmt.Sprintf("='PEMANTAUAN REKAP UANG MUKA'!T%d", totalRow)); err != nil {
			return err
		}
		if err := f.SetCellFormula(sheetName, fmt.Sprintf("T%d", totalRow+1), fmt.Sprintf("=SUM(T11:T%d)", totalRow-1)); err != nil {
			return err
		}

		if err := f.SetCellFormula(sheetName, fmt.Sprintf("U%d", totalRow+1), fmt.Sprintf("=SUM(U11:U%d)", totalRow-1)); err != nil {
			return err
		}
		if err := f.SetCellFormula(sheetName, fmt.Sprintf("V%d", totalRow+2), fmt.Sprintf("=SUM(V11:V%d)", totalRow-1)); err != nil {
			return err
		}
		if err := f.SetCellFormula(sheetName, fmt.Sprintf("W%d", totalRow+3), fmt.Sprintf("=SUM(W11:W%d)", totalRow-1)); err != nil {
			return err
		}

		if err := f.SetCellFormula(sheetName, fmt.Sprintf("K%d", totalRow+3), fmt.Sprintf("=K%d-K%d", totalRow+1, totalRow+2)); err != nil {
			return err
		}
		if err := f.SetCellFormula(sheetName, fmt.Sprintf("O%d", totalRow+3), fmt.Sprintf("=O%d-O%d", totalRow+1, totalRow+2)); err != nil {
			return err
		}
		if err := f.SetCellFormula(sheetName, fmt.Sprintf("P%d", totalRow+3), fmt.Sprintf("=P%d-P%d", totalRow+1, totalRow+2)); err != nil {
			return err
		}
		if err := f.SetCellFormula(sheetName, fmt.Sprintf("Q%d", totalRow+3), fmt.Sprintf("=Q%d-Q%d", totalRow+1, totalRow+2)); err != nil {
			return err
		}
		if err := f.SetCellFormula(sheetName, fmt.Sprintf("R%d", totalRow+3), fmt.Sprintf("=R%d-R%d", totalRow+1, totalRow+2)); err != nil {
			return err
		}
		if err := f.SetCellFormula(sheetName, fmt.Sprintf("S%d", totalRow+3), fmt.Sprintf("=S%d-S%d", totalRow+1, totalRow+2)); err != nil {
			return err
		}
		if err := f.SetCellFormula(sheetName, fmt.Sprintf("T%d", totalRow+3), fmt.Sprintf("=T%d-T%d", totalRow+1, totalRow+2)); err != nil {
			return err
		}

		if err := f.SetCellStyle(sheetName, fmt.Sprintf("A%d", totalRow), fmt.Sprintf("W%d", totalRow+3), summaryStyle); err != nil {
			return err
		}

		if err := f.SetCellStyle(sheetName, fmt.Sprintf("J%d", totalRow+1), fmt.Sprintf("W%d", totalRow+3), summaryNumberStyle); err != nil {
			return err
		}

		if err := f.SetCellStyle(sheetName, fmt.Sprintf("B%d", totalRow+1), fmt.Sprintf("B%d", totalRow+3), g.dynamicStyle(f, []string{"top", "right", "bottom", "left"}, true, false, 2, "left", 0, false, 0)); err != nil {
			return err
		}

		for i := currentRow; i < currentRow+4; i++ {
			if err := f.SetRowHeight(sheetName, i, 28); err != nil {
				return err
			}
		}
	} else {
		if err := f.SetCellValue(sheetName, fmt.Sprintf("B%d", totalRow), "JUMLAH"); err != nil {
			return err
		}

		if err := f.SetCellFormula(sheetName, fmt.Sprintf("K%d", totalRow), fmt.Sprintf("=SUM(K11:K%d)", totalRow-1)); err != nil {
			return err
		}
		if err := f.SetCellFormula(sheetName, fmt.Sprintf("O%d", totalRow), fmt.Sprintf("=SUM(O11:O%d)", totalRow-1)); err != nil {
			return err
		}
		if err := f.SetCellFormula(sheetName, fmt.Sprintf("P%d", totalRow), fmt.Sprintf("=SUM(P11:P%d)", totalRow-1)); err != nil {
			return err
		}
		if err := f.SetCellFormula(sheetName, fmt.Sprintf("Q%d", totalRow), fmt.Sprintf("=SUM(Q11:Q%d)", totalRow-1)); err != nil {
			return err
		}
		if err := f.SetCellFormula(sheetName, fmt.Sprintf("R%d", totalRow), fmt.Sprintf("=SUM(R11:R%d)", totalRow-1)); err != nil {
			return err
		}
		if err := f.SetCellFormula(sheetName, fmt.Sprintf("S%d", totalRow), fmt.Sprintf("=SUM(S11:S%d)", totalRow-1)); err != nil {
			return err
		}
		if err := f.SetCellFormula(sheetName, fmt.Sprintf("T%d", totalRow), fmt.Sprintf("=SUM(T11:T%d)", totalRow-1)); err != nil {
			return err
		}
		if err := f.SetCellFormula(sheetName, fmt.Sprintf("U%d", totalRow), fmt.Sprintf("=SUM(U11:U%d)", totalRow-1)); err != nil {
			return err
		}

		if err := f.SetCellStyle(sheetName, fmt.Sprintf("A%d", totalRow), fmt.Sprintf("U%d", totalRow), summaryStyle); err != nil {
			return err
		}

		if err := f.SetCellStyle(sheetName, fmt.Sprintf("K%d", totalRow), fmt.Sprintf("U%d", totalRow), summaryNumberStyle); err != nil {
			return err
		}

		for i := currentRow; i < currentRow+1; i++ {
			if err := f.SetRowHeight(sheetName, i, 28); err != nil {
				return err
			}
		}
	}

	return nil
}

func (g *Generator) generateSignatureRow(f *excelize.File, sheetName string, currentRow int) error {
	nameStyle, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Vertical: "center",
		},
		Font: &excelize.Font{
			Bold:   true,
			Size:   10,
			Family: "Tahoma",
		},
	})

	totalRow := currentRow + 3
	if sheetName == "PEMANTAUAN REKAP RAMPUNG" {
		totalRow += 3
	}

	err := f.SetCellStyle(sheetName, fmt.Sprintf("B%d", totalRow), fmt.Sprintf("L%d", totalRow+2), g.dynamicStyle(f, []string{"top", "right", "bottom", "left"}, false, false, 0, "left", 0, false, 0))
	if err != nil {
		return err
	}

	if err := f.SetCellValue(sheetName, fmt.Sprintf("B%d", totalRow), "Mengetahui/Menyetujui"); err != nil {
		return err
	}

	if err := f.SetCellValue(sheetName, fmt.Sprintf("D%d", totalRow), "Setuju/Lunas dibayar"); err != nil {
		return err
	}

	totalRow += 1
	if err := f.SetCellValue(sheetName, fmt.Sprintf("B%d", totalRow), "Pejabat Pembuat Komitmen II"); err != nil {
		return err
	}

	currentDate := time.Now().Format("2 January 2006")
	if err := f.SetCellValue(sheetName, fmt.Sprintf("D%d", totalRow), fmt.Sprintf("Tanggal: %s", currentDate)); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, fmt.Sprintf("I%d", totalRow), "Yang Membayarkan"); err != nil {
		return err
	}

	totalRow += 1
	if err := f.SetCellValue(sheetName, fmt.Sprintf("B%d", totalRow), fmt.Sprintf("Unit Kerja %s", workUnit)); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, fmt.Sprintf("D%d", totalRow), "Bendahara Pengeluaran"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, fmt.Sprintf("I%d", totalRow), "Pemegang Uang Muka Tim Kerja Manajemen Risiko, Reformasi Birokrasi dan Monitoring Evaluasi"); err != nil {
		return err
	}

	totalRow += 5
	if err := f.SetCellValue(sheetName, fmt.Sprintf("B%d", totalRow), commitmentOfficer); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, fmt.Sprintf("D%d", totalRow), expenditureTreasurer); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, fmt.Sprintf("I%d", totalRow), payer); err != nil {
		return err
	}

	totalRow += 1
	if err := f.SetCellValue(sheetName, fmt.Sprintf("B%d", totalRow), commitmentOfficerId); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, fmt.Sprintf("D%d", totalRow), expenditureTreasurerId); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, fmt.Sprintf("I%d", totalRow), payerId); err != nil {
		return err
	}

	if sheetName == "PEMANTAUAN REKAP RAMPUNG" {
		err := f.SetCellStyle(sheetName, fmt.Sprintf("B%d", totalRow), fmt.Sprintf("L%d", totalRow), nameStyle)
		if err != nil {
			return err
		}
	} else {
		err := f.SetCellStyle(sheetName, fmt.Sprintf("B%d", totalRow), fmt.Sprintf("L%d", totalRow), nameStyle)
		if err != nil {
			return err
		}
	}

	landscape := "landscape"
	excelizeOne := 1
	excelizeTrue := true
	if err := f.SetPageLayout(sheetName, &excelize.PageLayoutOptions{
		Size:        utils.ToPtr(14),
		Orientation: &landscape,
		FitToWidth:  &excelizeOne,
		FitToHeight: &excelizeOne,
	}); err != nil {
		panic(err)
	}

	_ = f.SetPageMargins(sheetName, &excelize.PageLayoutMarginsOptions{
		Horizontally: &excelizeTrue,
	})

	zoomScale := 112.0
	_ = f.SetSheetView(sheetName, 0, &excelize.ViewOptions{
		ZoomScale: &zoomScale,
	})

	_ = f.DeleteDefinedName(&excelize.DefinedName{
		Name:  "_xlnm.Print_Area",
		Scope: sheetName,
	})

	ref := ""
	if sheetName == "PEMANTAUAN REKAP UANG MUKA" {
		ref = fmt.Sprintf("'%s'!$A$1:$U$24", sheetName)
	} else {
		ref = fmt.Sprintf("'%s'!$A$1:$W$26", sheetName)
	}

	if err := f.SetDefinedName(&excelize.DefinedName{
		Name:     "_xlnm.Print_Area",
		RefersTo: ref,
		Scope:    sheetName,
	}); err != nil {
		return err
	}

	return nil
}

func (g *Generator) setColumnWidths(f *excelize.File, sheetName string) error {
	if err := f.SetColWidth(sheetName, "A", "A", 5); err != nil {
		return err
	}
	if err := f.SetColWidth(sheetName, "B", "B", 30); err != nil {
		return err
	}
	if err := f.SetColWidth(sheetName, "C", "C", 20); err != nil {
		return err
	}
	if err := f.SetColWidth(sheetName, "D", "D", 25); err != nil {
		return err
	}
	if err := f.SetColWidth(sheetName, "E", "E", 25); err != nil {
		return err
	}
	if err := f.SetColWidth(sheetName, "F", "F", 30); err != nil {
		return err
	}
	if err := f.SetColWidth(sheetName, "G", "G", 25); err != nil {
		return err
	}

	if err := f.SetColWidth(sheetName, "H", "H", 10); err != nil {
		return err
	}
	if err := f.SetColWidth(sheetName, "I", "I", 15); err != nil {
		return err
	}
	if err := f.SetColWidth(sheetName, "J", "J", 15); err != nil {
		return err
	}

	if err := f.SetColWidth(sheetName, "K", "K", 15); err != nil {
		return err
	}

	if err := f.SetColWidth(sheetName, "O", "O", 15); err != nil {
		return err
	}

	if err := f.SetColWidth(sheetName, "L", "L", 10); err != nil {
		return err
	}
	if err := f.SetColWidth(sheetName, "M", "M", 15); err != nil {
		return err
	}
	if err := f.SetColWidth(sheetName, "N", "N", 15); err != nil {
		return err
	}

	if err := f.SetColWidth(sheetName, "P", "P", 15); err != nil {
		return err
	}
	if err := f.SetColWidth(sheetName, "Q", "S", 15); err != nil {
		return err
	}
	if err := f.SetColWidth(sheetName, "T", "T", 15); err != nil {
		return err
	}

	if err := f.SetColWidth(sheetName, "U", "U", 20); err != nil {
		return err
	}

	return nil
}

func (g *Generator) generateKw(f *excelize.File, sheetName string, req RecapReport) error {
	if _, err := f.NewSheet(sheetName); err != nil {
		return err
	}

	if err := f.SetCellValue(sheetName, "T1", 1); err != nil {
		return err
	}

	titleStyle, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Font: &excelize.Font{
			Bold:   true,
			Size:   10,
			Family: "Tahoma",
		},
	})

	if err := f.SetCellValue(sheetName, "A1", "KEMENTERIAN KESEHATAN REPUBLIK INDONESIA"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "A2", "DIREKTORAT JENDERAL"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "A3", "PENANGGULANAN PENYAKIT"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "A5", "J A K A R T A"); err != nil {
		return err
	}

	if err := f.MergeCell(sheetName, "A1", "L1"); err != nil {
		return err
	}
	if err := f.MergeCell(sheetName, "A2", "L2"); err != nil {
		return err
	}
	if err := f.MergeCell(sheetName, "A3", "L3"); err != nil {
		return err
	}
	if err := f.MergeCell(sheetName, "A5", "L5"); err != nil {
		return err
	}

	if err := f.SetCellStyle(sheetName, "A1", "L5", titleStyle); err != nil {
		return err
	}

	basicStyle, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Vertical: "center",
		},
		Font: &excelize.Font{
			Size:   10,
			Family: "Tahoma",
		},
	})

	currencyStyle, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Vertical: "center",
		},
		Font: &excelize.Font{
			Size:   10,
			Family: "Tahoma",
		},
		NumFmt: 3,
	})

	if err := f.SetCellValue(sheetName, "M1", "Tahun Anggaran"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "M2", "No Bukti"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "M3", "Akun"); err != nil {
		return err
	}

	if err := f.SetCellValue(sheetName, "O1", ":"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "O2", ":"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "O3", ":"); err != nil {
		return err
	}

	if err := f.SetCellValue(sheetName, "P1", "2025"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "P3", "024.05.WA.4815.EBD.953."); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "P4", "501.B.524111"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "T1", 1); err != nil {
		return err
	}

	if err := f.SetCellStyle(sheetName, "M1", "P4", basicStyle); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "M1", "P4", basicStyle); err != nil {
		return err
	}

	if err := f.SetCellValue(sheetName, "A6", "RINCIAN BIAYA PERJALANAN DINAS"); err != nil {
		return err
	}
	if err := f.MergeCell(sheetName, "A6", "S6"); err != nil {
	}
	if err := f.SetCellStyle(sheetName, "A6", "S6", titleStyle); err != nil {
		return err
	}

	if err := f.SetCellValue(sheetName, "A7", "Lampiran SPD Nomor"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "F7", ":"); err != nil {
		return err
	}
	if err := f.SetCellFormula(sheetName, "G7", "=VLOOKUP($T$1,'PEMANTAUAN REKAP UANG MUKA'!$A$11:$AA$100,27,FALSE)"); err != nil {
		return err
	}

	if err := f.SetCellValue(sheetName, "A8", "Tanggal"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "F8", ":"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "G8", req.StartDate); err != nil {
		return err
	}

	if err := f.SetCellStyle(sheetName, "A7", "S9", basicStyle); err != nil {
		return err
	}

	// header
	if err := f.SetCellValue(sheetName, "A10", "NO"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "C10", "PERINCIAN BIAYA"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "L10", "JUMLAH"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "O10", "KETERANGAN"); err != nil {
		return err
	}

	if err := f.MergeCell(sheetName, "A10", "B11"); err != nil {
		return err
	}
	if err := f.MergeCell(sheetName, "C10", "K11"); err != nil {
		return err
	}
	if err := f.MergeCell(sheetName, "L10", "N11"); err != nil {
		return err
	}
	if err := f.MergeCell(sheetName, "O10", "S11"); err != nil {
		return err
	}

	kwHeaderStyle, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Font: &excelize.Font{
			Bold:   true,
			Size:   10,
			Family: "Tahoma",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 2},
			{Type: "top", Color: "000000", Style: 2},
			{Type: "bottom", Color: "000000", Style: 2},
			{Type: "right", Color: "000000", Style: 2},
		},
	})

	if err := f.SetCellStyle(sheetName, "A10", "B11", kwHeaderStyle); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "C10", "K11", kwHeaderStyle); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "L10", "N11", kwHeaderStyle); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "O10", "S11", kwHeaderStyle); err != nil {
		return err
	}

	if err := f.SetCellValue(sheetName, "A13", "1"); err != nil {
		return err
	}

	if err := f.SetCellValue(sheetName, "C13", "Uang harian :"); err != nil {
		return err
	}

	if sheetName == "KW RAMPUNG" {
		if err := f.SetCellFormula(sheetName, "C14", "=VLOOKUP($T$1,'PEMANTAUAN REKAP RAMPUNG'!$A$11:$AC$100,8,FALSE)"); err != nil {
			return err
		}
	} else {
		if err := f.SetCellFormula(sheetName, "C14", "=VLOOKUP($T$1,'PEMANTAUAN REKAP UANG MUKA'!$A$11:$AC$100,8,FALSE)"); err != nil {
			return err
		}
	}

	if err := f.SetCellValue(sheetName, "D14", "hr"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "E14", "x"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "F14", "Rp."); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "H14", "-"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "L14", "Rp."); err != nil {
		return err
	}

	if sheetName == "KW RAMPUNG" {
		if err := f.SetCellFormula(sheetName, "M14", "=VLOOKUP($T$1,'PEMANTAUAN REKAP RAMPUNG'!$A$11:$AC$100,11,FALSE)"); err != nil {
			return err
		}
	} else {
		if err := f.SetCellFormula(sheetName, "M14", "=VLOOKUP($T$1,'PEMANTAUAN REKAP UANG MUKA'!$A$11:$AC$100,11,FALSE)"); err != nil {
			return err
		}
	}

	if err := f.SetCellStyle(sheetName, "M14", "M14", currencyStyle); err != nil {
		return err
	}

	if err := f.SetCellValue(sheetName, "A15", "2"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "C15", "Transport"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "C16", "a."); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "D16", "Tiket :"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "D17", fmt.Sprintf("- Pesawat Jakarta - %s (PP)", req.DestinationCity)); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "L17", "Rp."); err != nil {
		return err
	}
	if sheetName == "KW RAMPUNG" {
		if err := f.SetCellFormula(sheetName, "M17", "=VLOOKUP($T$1,'PEMANTAUAN REKAP RAMPUNG'!$A$11:$AC$100,20,FALSE)"); err != nil {
			return err
		}
	} else {
		if err := f.SetCellFormula(sheetName, "M17", "=VLOOKUP($T$1,'PEMANTAUAN REKAP UANG MUKA'!$A$11:$AC$100,20,FALSE)"); err != nil {
			return err
		}
	}

	if err := f.SetCellStyle(sheetName, "M17", "M17", currencyStyle); err != nil {
		return err
	}

	if err := f.SetCellValue(sheetName, "C18", "b."); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "D18", "Transport (PP):"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "D19", "- Transport Jakarta - Bandara Soetta (PP)"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "L19", "Rp."); err != nil {
		return err
	}

	if sheetName == "KW RAMPUNG" {
		if err := f.SetCellFormula(sheetName, "M19", "=VLOOKUP($T$1,'PEMANTAUAN REKAP RAMPUNG'!$A$11:$AC$100,17,FALSE)"); err != nil {
			return err
		}
	} else {
		if err := f.SetCellFormula(sheetName, "M19", "=VLOOKUP($T$1,'PEMANTAUAN REKAP UANG MUKA'!$A$11:$AC$100,17,FALSE)"); err != nil {
			return err
		}
	}

	if err := f.SetCellStyle(sheetName, "M19", "M19", currencyStyle); err != nil {
		return err
	}

	if err := f.SetCellValue(sheetName, "D20", fmt.Sprintf("- Transport Dareah %s (PP)", req.DestinationCity)); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "L20", "Rp."); err != nil {
		return err
	}

	if sheetName == "KW RAMPUNG" {
		if err := f.SetCellFormula(sheetName, "M20", "=VLOOKUP($T$1,'PEMANTAUAN REKAP RAMPUNG'!$A$11:$AC$100,18,FALSE)"); err != nil {
			return err
		}
	} else {
		if err := f.SetCellFormula(sheetName, "M20", "=VLOOKUP($T$1,'PEMANTAUAN REKAP UANG MUKA'!$A$11:$AC$100,18,FALSE)"); err != nil {
			return err
		}
	}

	if err := f.SetCellStyle(sheetName, "M20", "M20", currencyStyle); err != nil {
		return err
	}

	if err := f.SetCellValue(sheetName, "A21", "3"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "C21", "Biaya Penginapan : "); err != nil {
		return err
	}

	if sheetName == "KW RAMPUNG" {
		if err := f.SetCellFormula(sheetName, "C22", "=VLOOKUP($T$1,'PEMANTAUAN REKAP RAMPUNG'!$A$11:$AC$100,12,FALSE)"); err != nil {
			return err
		}
		if err := f.SetCellFormula(sheetName, "D22", "=VLOOKUP($T$1,'PEMANTAUAN REKAP RAMPUNG'!$A$11:$AC$100,13,FALSE)"); err != nil {
			return err
		}
	} else {
		if err := f.SetCellFormula(sheetName, "C22", "=VLOOKUP($T$1,'PEMANTAUAN REKAP UANG MUKA'!$A$11:$AC$100,12,FALSE)"); err != nil {
			return err
		}
		if err := f.SetCellFormula(sheetName, "D22", "=VLOOKUP($T$1,'PEMANTAUAN REKAP UANG MUKA'!$A$11:$AC$100,13,FALSE)"); err != nil {
			return err
		}
	}

	if err := f.SetCellValue(sheetName, "E22", "x"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "F22", "Rp."); err != nil {
		return err
	}
	if err := f.SetCellFormula(sheetName, "H22", "=VLOOKUP($T$1,'PEMANTAUAN REKAP UANG MUKA'!$A$11:$AC$100,14,FALSE)"); err != nil {
		return err
	}

	if err := f.SetCellStyle(sheetName, "H22", "H22", currencyStyle); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "L22", "Rp."); err != nil {
		return err
	}
	if err := f.SetCellFormula(sheetName, "M22", "=VLOOKUP($T$1,'PEMANTAUAN REKAP UANG MUKA'!$A$11:$AC$100,15,FALSE)"); err != nil {
		return err
	}

	if err := f.SetCellStyle(sheetName, "M22", "M22", currencyStyle); err != nil {
		return err
	}

	if err := f.SetCellStyle(sheetName, "C12", "C24", g.dynamicStyle(f, []string{"left"}, false, false, 2, "left", 0, false, 0)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "L12", "L24", g.dynamicStyle(f, []string{"left"}, false, false, 2, "left", 0, false, 0)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "O12", "O24", g.dynamicStyle(f, []string{"left"}, false, false, 2, "left", 0, false, 0)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "S12", "S24", g.dynamicStyle(f, []string{"right"}, false, false, 2, "left", 0, false, 0)); err != nil {
		return err
	}

	if err := f.MergeCell(sheetName, "A13", "B13"); err != nil {
		return err
	}
	if err := f.MergeCell(sheetName, "A15", "B15"); err != nil {
		return err
	}
	if err := f.MergeCell(sheetName, "A21", "B21"); err != nil {
		return err
	}

	noStyle, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Vertical:   "center",
			Horizontal: "left",
		},
		Font: &excelize.Font{
			Size:   10,
			Family: "Tahoma",
		},
	})

	if err := f.SetCellStyle(sheetName, "A13", "A21", noStyle); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "C24", "J U M L A H"); err != nil {
		return err
	}
	if err := f.MergeCell(sheetName, "C24", "K24"); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "C24", "C24", kwHeaderStyle); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "A24", "B24", g.dynamicStyle(f, []string{"top", "bottom"}, false, false, 2, "left", 0, false, 0)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "D24", "K24", g.dynamicStyle(f, []string{"top", "bottom"}, false, false, 2, "left", 0, false, 0)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "M24", "N24", g.dynamicStyle(f, []string{"top", "bottom"}, false, false, 2, "left", 3, false, 0)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "L24", "L24", g.dynamicStyle(f, []string{"top", "bottom", "left"}, false, false, 2, "left", 0, false, 0)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "P24", "R24", g.dynamicStyle(f, []string{"top", "bottom"}, false, false, 2, "left", 0, false, 0)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "O24", "O24", g.dynamicStyle(f, []string{"top", "bottom", "left"}, false, false, 2, "left", 0, false, 0)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "S24", "S24", g.dynamicStyle(f, []string{"top", "bottom", "right"}, false, false, 2, "left", 0, false, 0)); err != nil {
		return err
	}

	if err := f.SetCellValue(sheetName, "L24", "Rp."); err != nil {
		return err
	}
	if err := f.SetCellFormula(sheetName, "M24", "=SUM(M12:M23)"); err != nil {
		return err
	}

	if err := f.SetCellStyle(sheetName, "M24", "M24", currencyStyle); err != nil {
		return err
	}

	if err := f.SetCellStyle(sheetName, "M24", "M24", g.dynamicStyle(f, []string{"top", "bottom"}, true, false, 2, "left", 0, false, 0)); err != nil {
		return err
	}

	if err := f.SetCellStyle(sheetName, "L24", "L24", g.dynamicStyle(f, []string{"top", "bottom", "left"}, true, false, 2, "left", 0, false, 0)); err != nil {
		return err
	}

	if err := f.SetCellValue(sheetName, "A25", "TERBILANG:"); err != nil {
		return err
	}

	terbilang, err := f.CalcCellValue(sheetName, "M24")
	if err != nil {
		return err
	}

	valInt, err := strconv.ParseInt(terbilang, 10, 64)
	if err != nil {
		return err
	}
	terbilang = Terbilang(valInt)

	if err := f.SetCellValue(sheetName, "D25", terbilang); err != nil {
		return err
	}

	if err := f.SetCellValue(sheetName, "O27", fmt.Sprintf("Jakarta, %s", req.ReceiptSignatureDate)); err != nil {
		return err
	}

	if err := f.SetCellValue(sheetName, "A28", "Telah dibayar sejumlah"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "O28", "Telah menerima jumlah uang sebesar"); err != nil {
		return err
	}

	if err := f.SetCellValue(sheetName, "A29", "Rp."); err != nil {
		return err
	}
	if err := f.SetCellFormula(sheetName, "C29", "=M24"); err != nil {
		return err
	}

	if err := f.SetCellValue(sheetName, "O29", "Rp."); err != nil {
		return err
	}

	if err := f.SetCellFormula(sheetName, "P29", "=M24"); err != nil {
		return err
	}

	if err := f.MergeCell(sheetName, "C29", "F29"); err != nil {
		return err
	}

	if err := f.SetCellValue(sheetName, "A30", "Bendahara Pengeluaran Pembantu"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "I30", "PUM Timker"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "O30", "Yang Menerima"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "A31", "Unit Kerja Setditjen Penanggulangan Penyakit"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "A34", expenditureTreasurer); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "I34", payer); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "A35", "NIP"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "B35", expenditureTreasurerId); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "I35", "NIP"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "J35", payerId); err != nil {
		return err
	}
	if err := f.SetCellFormula(sheetName, "O34", "=VLOOKUP($T$1,'PEMANTAUAN REKAP UANG MUKA'!$A$11:$AC$100,2,FALSE)"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "O35", "NIP"); err != nil {
		return err
	}
	if err := f.SetCellFormula(sheetName, "P35", "=VLOOKUP($T$1,'PEMANTAUAN REKAP UANG MUKA'!$A$11:$AC$100,3,FALSE)"); err != nil {
		return err
	}

	if err := f.SetCellStyle(sheetName, "A28", "S35", g.dynamicStyle(f, []string{}, false, false, 2, "left", 0, false, 0)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "A28", "S35", g.dynamicStyle(f, []string{}, false, false, 2, "left", 0, false, 0)); err != nil {
		return err
	}

	if err := f.SetCellStyle(sheetName, "A29", "C29", g.dynamicStyle(f, []string{}, true, false, 2, "left", 3, false, 0)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "M24", "M24", g.dynamicStyle(f, []string{"top", "bottom"}, true, false, 2, "right", 3, false, 0)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "O29", "P29", g.dynamicStyle(f, []string{}, true, false, 2, "left", 3, false, 0)); err != nil {
		return err
	}

	if err := f.SetCellStyle(sheetName, "A35", "S35", g.dynamicStyle(f, []string{"bottom"}, false, false, 2, "left", 0, false, 0)); err != nil {
		return err
	}

	if err := f.SetCellValue(sheetName, "A36", "PERHITUNGAN SPD UANG MUKA"); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "A36", "A36", g.dynamicStyle(f, []string{}, true, false, 0, "center", 0, false, 0)); err != nil {
		return err
	}

	if err := f.MergeCell(sheetName, "A36", "S36"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "A37", "Ditetapkan sejumlah"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "I37", ":"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "J37", "Rp"); err != nil {
		return err
	}
	if err := f.SetCellFormula(sheetName, "K37", "=M24"); err != nil {
		return err
	}
	if err := f.MergeCell(sheetName, "K37", "M37"); err != nil {
		return err
	}

	if err := f.SetCellValue(sheetName, "A38", "Yang telah dibayar semula"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "I38", ":"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "J38", "Rp"); err != nil {
		return err
	}
	if err := f.SetCellFormula(sheetName, "K38", "=K37"); err != nil {
		return err
	}

	if err := f.MergeCell(sheetName, "K38", "M38"); err != nil {
		return err
	}

	if err := f.SetCellValue(sheetName, "A39", "Sisa Kurang/Lebih"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "I39", ":"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "J39", "Rp"); err != nil {
		return err
	}
	if err := f.SetCellFormula(sheetName, "K39", "=K37"); err != nil {
		return err
	}

	if err := f.MergeCell(sheetName, "K39", "M39"); err != nil {
		return err
	}

	if err := f.SetCellStyle(sheetName, "A37", "L40", g.dynamicStyle(f, []string{}, true, false, 0, "left", 0, false, 0)); err != nil {
		return err
	}

	if err := f.SetCellStyle(sheetName, "K37", "L39", g.dynamicStyle(f, []string{}, true, false, 0, "left", 3, false, 0)); err != nil {
		return err
	}

	if err := f.SetCellValue(sheetName, "M40", "an.  Kuasa Pengguna Anggaran"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "M41", "      Pejabat Pembuat Komitmen,"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "M42", "      Satker Kantor Pusat Ditjen Penanggulangan Penyakit"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "M43", "      Unit Kerja Setditjen Penanggulangan Penyakit"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "M47", fmt.Sprintf("      %s", expenditureTreasurer)); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "M48", fmt.Sprintf("      NIP  %s", expenditureTreasurer)); err != nil {
		return err
	}

	if err := f.SetCellStyle(sheetName, "A50", "S50", g.dynamicStyle(f, []string{"bottom"}, false, false, 8, "left", 0, false, 0)); err != nil {
		return err
	}

	if err := f.SetCellValue(sheetName, "A53", "DAFTAR PENELUARAN RIIL"); err != nil {
		return err
	}

	if err := f.MergeCell(sheetName, "A53", "S53"); err != nil {
		return err
	}

	if err := f.SetCellStyle(sheetName, "A53", "A53", g.dynamicStyle(f, []string{}, true, false, 0, "center", 0, false, 0)); err != nil {
		return err
	}

	if err := f.SetCellValue(sheetName, "A55", "Yang bertanda tangan dibawah ini  :"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "A57", "Nama"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "C57", ":"); err != nil {
		return err
	}
	if err := f.SetCellFormula(sheetName, "D57", "=VLOOKUP($T$1,'PEMANTAUAN REKAP UANG MUKA'!$A$11:$AC$12,2,FALSE)"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "A58", "NIP"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "C58", ":"); err != nil {
		return err
	}
	if err := f.SetCellFormula(sheetName, "D58", "=VLOOKUP($T$1,'PEMANTAUAN REKAP UANG MUKA'!$A$11:$AC$12,3,FALSE)"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "A59", "Jabatan"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "C59", ":"); err != nil {
		return err
	}
	if err := f.SetCellFormula(sheetName, "D59", "=VLOOKUP($T$1,'PEMANTAUAN REKAP UANG MUKA'!$A$11:$AC$12,4,FALSE)"); err != nil {
		return err
	}

	if err := f.SetCellFormula(sheetName, "A61", `="Berdasarkan Surat Perjalanan Dinas ( SPD ) Nomor "&VLOOKUP($T$1,'PEMANTAUAN REKAP UANG MUKA'!$A$11:$AJ$100,22,FALSE)`); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "N61", fmt.Sprintf("tanggal %s", req.SpdDate)); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "A62", "dengan ini kami menyatakan dengan sesungguhnya bahwa :"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "A64", "1.  Biaya transport pegawai dan / atau biaya penginapan dibawah ini yang tidak dapat "); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "A65", "     diperoleh bukti-bukti pengeluarannya, meliputi  :"); err != nil {
		return err
	}

	if err := f.SetCellValue(sheetName, "A67", "No"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "D67", "U r a i a n"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "O67", "Jumlah"); err != nil {
		return err
	}

	if err := f.MergeCell(sheetName, "A67", "C67"); err != nil {
		return err
	}
	if err := f.MergeCell(sheetName, "D67", "N67"); err != nil {
		return err
	}
	if err := f.MergeCell(sheetName, "O67", "S67"); err != nil {
		return err
	}

	if err := f.SetCellStyle(sheetName, "A67", "S67", g.dynamicStyle(f, []string{"top", "right", "bottom", "left"}, true, false, 2, "center", 0, false, 0)); err != nil {
		return err
	}

	if err := f.SetCellValue(sheetName, "B69", "1"); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "B69", "B69", g.dynamicStyle(f, []string{}, false, false, 0, "center", 0, false, 0)); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "D69", "Transport :"); err != nil {
		return err
	}
	if err := f.SetCellFormula(sheetName, "E70", "=D19"); err != nil {
		return err
	}
	if err := f.SetCellFormula(sheetName, "E71", "=D20"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "O70", "Rp."); err != nil {
		return err
	}

	if err := f.SetCellStyle(sheetName, "O70", "O70", currencyStyle); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "O71", "Rp."); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "O71", "O71", currencyStyle); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "D74", "JUMLAH"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "O74", "Rp."); err != nil {
		return err
	}
	if err := f.MergeCell(sheetName, "D74", "N74"); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "D68", "D73", g.dynamicStyle(f, []string{"left"}, false, false, 2, "left", 0, false, 0)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "O68", "O73", g.dynamicStyle(f, []string{"left"}, false, false, 2, "left", 0, false, 0)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "S68", "S73", g.dynamicStyle(f, []string{"right"}, false, false, 2, "left", 0, false, 0)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "A74", "S74", g.dynamicStyle(f, []string{"top", "bottom"}, false, false, 2, "left", 0, false, 0)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "D74", "D74", g.dynamicStyle(f, []string{"top", "bottom", "left"}, true, false, 2, "center", 0, false, 0)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "O74", "O74", g.dynamicStyle(f, []string{"top", "bottom", "left"}, true, false, 2, "left", 0, false, 0)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "S74", "S74", g.dynamicStyle(f, []string{"top", "bottom", "right"}, false, false, 2, "left", 0, false, 0)); err != nil {
		return err
	}

	if err := f.SetCellValue(sheetName, "A76", "2.  Jumlah uang tersebut pada angka 1 diatas benar-benar dikeluarkan untuk pelaksanaan "); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "A77", "     perjalanan dinas dimaksud dan apabila kemudian hari terdapat kelebihan atas pembayaran"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "A78", "     kami bersedia untuk menyetorkan kelebihan tersebut ke Kas Negara."); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "A80", "Demikian pernyataan ini kami buat dengan sebenarnya, untuk dipergunakan sebagaimana mestinya."); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "A83", "Mengetahui / Menyetujui"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "N83", fmt.Sprintf("Jakarta, %s", req.ReceiptSignatureDate)); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "N84", "Pelaksana SPD,"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "A84", "Pejabat Pembuat Komitmen II,"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "A85", "Satker Kantor Pusat Ditjen Penanggulangan Penyakit"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "A86", "Unit Kerja Setditjen Penanggulangan Penyakit"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "A89", expenditureTreasurer); err != nil {
		return err
	}
	if err := f.SetCellFormula(sheetName, "N89", "=VLOOKUP($T$1,'PEMANTAUAN REKAP UANG MUKA'!$A$11:$AC$100,2,FALSE)"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "A90", "NIP"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "B90", expenditureTreasurerId); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "N90", "NIP"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "O90", expenditureTreasurerId); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "A89", "S89", g.dynamicStyle(f, []string{}, true, false, 0, "left", 0, false, 0)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "A5", "L5", g.dynamicStyle(f, []string{"bottom"}, true, false, 1, "center", 0, false, 0)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "A55", "P65", g.dynamicStyle(f, []string{}, false, false, 0, "left", 0, false, 0)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "A76", "P86", g.dynamicStyle(f, []string{}, false, false, 0, "left", 0, false, 0)); err != nil {
		return err
	}

	if err := f.SetCellStyle(sheetName, "A12", "A23", g.dynamicStyle(f, []string{"left"}, false, false, 2, "center", 0, false, 0)); err != nil {
		return err
	}

	if err := f.SetCellStyle(sheetName, "A14", "A14", g.dynamicStyle(f, []string{"left"}, false, false, 2, "center", 0, false, 0)); err != nil {
		return err
	}

	if err := f.SetCellStyle(sheetName, "A22", "A22", g.dynamicStyle(f, []string{"left"}, false, false, 2, "left", 0, false, 0)); err != nil {
		return err
	}

	if err := f.SetCellStyle(sheetName, "A68", "A73", g.dynamicStyle(f, []string{"left"}, false, false, 2, "left", 0, false, 0)); err != nil {
		return err
	}

	if err := f.SetCellStyle(sheetName, "A74", "A74", g.dynamicStyle(f, []string{"left", "top", "bottom"}, false, false, 2, "left", 0, false, 0)); err != nil {
		return err
	}

	if err := f.SetCellStyle(sheetName, "A24", "A24", g.dynamicStyle(f, []string{"left", "top", "bottom"}, false, false, 2, "left", 0, false, 0)); err != nil {
		return err
	}

	if err := f.SetCellStyle(sheetName, "A25", "D25", g.dynamicStyle(f, []string{}, false, true, 2, "left", 0, false, 0)); err != nil {
		return err
	}

	if err := f.SetColWidth(sheetName, "A", "A", 4); err != nil {
		return err
	}
	if err := f.SetColWidth(sheetName, "B", "B", 3.4); err != nil {
		return err
	}
	if err := f.SetColWidth(sheetName, "C", "C", 5.2); err != nil {
		return err
	}
	if err := f.SetColWidth(sheetName, "D", "D", 6); err != nil {
		return err
	}
	if err := f.SetColWidth(sheetName, "E", "E", 1.8); err != nil {
		return err
	}
	if err := f.SetColWidth(sheetName, "F", "F", 6.8); err != nil {
		return err
	}
	if err := f.SetColWidth(sheetName, "G", "G", 2); err != nil {
		return err
	}
	if err := f.SetColWidth(sheetName, "H", "H", 13.1); err != nil {
		return err
	}
	if err := f.SetColWidth(sheetName, "I", "I", 3.6); err != nil {
		return err
	}
	if err := f.SetColWidth(sheetName, "J", "J", 3.5); err != nil {
		return err
	}
	if err := f.SetColWidth(sheetName, "K", "K", 3.6); err != nil {
		return err
	}
	if err := f.SetColWidth(sheetName, "L", "L", 4.8); err != nil {
		return err
	}
	if err := f.SetColWidth(sheetName, "M", "M", 11.8); err != nil {
		return err
	}
	if err := f.SetColWidth(sheetName, "N", "N", 3.6); err != nil {
		return err
	}
	if err := f.SetColWidth(sheetName, "O", "O", 3.5); err != nil {
		return err
	}
	if err := f.SetColWidth(sheetName, "P", "P", 20); err != nil {
		return err
	}
	if err := f.SetColWidth(sheetName, "Q", "Q", 4.3); err != nil {
		return err
	}
	if err := f.SetColWidth(sheetName, "R", "R", 2); err != nil {
		return err
	}
	if err := f.SetColWidth(sheetName, "S", "S", 4.75); err != nil {
		return err
	}

	for i := 1; i <= 100; i++ {
		f.SetRowHeight(sheetName, i, 13)
	}

	portrait := "portrait"
	excelizeOne := 1
	excelizeTrue := true
	if err := f.SetPageLayout(sheetName, &excelize.PageLayoutOptions{
		Size:        utils.ToPtr(14),
		Orientation: &portrait,
		FitToWidth:  &excelizeOne,
		FitToHeight: &excelizeOne,
	}); err != nil {
		panic(err)
	}

	_ = f.SetPageMargins(sheetName, &excelize.PageLayoutMarginsOptions{
		Left:         utils.ToPtr(0.2362204),
		Right:        utils.ToPtr(0.2362204),
		Top:          utils.ToPtr(0.3543307),
		Bottom:       utils.ToPtr(0.3543307),
		Header:       utils.ToPtr(0.3149606),
		Footer:       utils.ToPtr(0.3149606),
		Horizontally: &excelizeTrue,
		Vertically:   &excelizeTrue,
	})

	zoomScale := 154.0
	_ = f.SetSheetView(sheetName, 0, &excelize.ViewOptions{
		ZoomScale: &zoomScale,
	})

	_ = f.DeleteDefinedName(&excelize.DefinedName{
		Name:  "_xlnm.Print_Area",
		Scope: sheetName,
	})

	ref := fmt.Sprintf("'%s'!$A$1:$S$91", sheetName)

	if err := f.SetDefinedName(&excelize.DefinedName{
		Name:     "_xlnm.Print_Area",
		RefersTo: ref,
		Scope:    sheetName,
	}); err != nil {
		return err
	}

	return nil
}

func (g *Generator) generateSppd(f *excelize.File, sheetName string, req RecapReport) error {
	if _, err := f.NewSheet(sheetName); err != nil {
		return err
	}

	if err := f.SetCellValue(sheetName, "F1", "LAMPIRAN I"); err != nil {
		return err
	}

	if err := f.SetCellValue(sheetName, "F2", "PERATURAN MENTERI KEUANGAN REPUBLIK INDONESIA NOMOR 113/PMK.05/2012"); err != nil {
		return err
	}

	if err := f.SetCellValue(sheetName, "F3", "Tentang"); err != nil {
		return err
	}

	if err := f.SetCellValue(sheetName, "F4", "PERJALANAN DINAS JABATAN DALAM NEGERI BAGI PEJABAT NEGARA, PEGAWAI NEGERI, DAN PEGAWAI TIDAK TETAP"); err != nil {
		return err
	}

	if err := f.SetCellValue(sheetName, "A8", "KEMENTERIAN KESEHATAN REPUBLIK INDONNESIA"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "A9", "DIREKTORAT JENDERAL"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "A10", "PENANGGULANGAN PENYAKIT"); err != nil {
		return err
	}

	if err := f.SetCellValue(sheetName, "F10", "Lembar ke"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "G10", ": ………………………………………"); err != nil {
		return err
	}

	if err := f.SetCellValue(sheetName, "F11", "Kode No"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "G11", ": ………………………………………"); err != nil {
		return err
	}

	if err := f.SetCellValue(sheetName, "F12", "Nomor"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "G12", ":"); err != nil {
		return err
	}
	if err := f.SetCellFormula(sheetName, "H12", "=VLOOKUP($K$20,'PEMANTAUAN REKAP RAMPUNG'!$A$8:$AC$100,28,FALSE)"); err != nil {
		return err
	}

	if err := f.SetCellValue(sheetName, "A12", "J A K A R T A"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "A14", "SURAT PERJALANAN DINAS (SPD)"); err != nil {
		return err
	}

	if err := f.SetCellValue(sheetName, "A16", "1."); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "B16", "Pejabat Pembuat Komitmen I"); err != nil {
		return err
	}
	if err := f.SetCellFormula(sheetName, "D16", "=F63"); err != nil {
		return err
	}

	if err := f.SetCellValue(sheetName, "A19", "2."); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "B19", "Nama/NIP Pegawai yang melaksanakan"); err != nil {
		return err
	}
	if err := f.SetCellFormula(sheetName, "D19", "=VLOOKUP($K$20,'PEMANTAUAN REKAP RAMPUNG'!$A$8:$AC$100,2,FALSE)"); err != nil {
		return err
	}

	if err := f.SetCellValue(sheetName, "B20", "perjalanan dinas"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "D20", "NIP"); err != nil {
		return err
	}
	if err := f.SetCellFormula(sheetName, "E20", "=VLOOKUP($K$20,'PEMANTAUAN REKAP RAMPUNG'!$A$8:$AC$100,3,FALSE)"); err != nil {
		return err
	}

	if err := f.SetCellValue(sheetName, "B22", "a."); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "C22", "Pangkat dan Golongan"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "D22", "a."); err != nil {
		return err
	}

	if err := f.SetCellFormula(sheetName, "E22", "=VLOOKUP($K$20,'PEMANTAUAN REKAP RAMPUNG'!$A$8:$AG$100,33,FALSE)"); err != nil {
		return err
	}
	if err := f.SetCellFormula(sheetName, "F22", "=VLOOKUP($K$20,'PEMANTAUAN REKAP RAMPUNG'!$A$8:$AG$34,5,FALSE)"); err != nil {
		return err
	}

	if err := f.SetCellValue(sheetName, "A23", "3."); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "B23", "b."); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "C23", "Jabatan/Instansi"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "D23", "b."); err != nil {
		return err
	}
	if err := f.SetCellFormula(sheetName, "E23", "=VLOOKUP($K$20,'PEMANTAUAN REKAP RAMPUNG'!$A$8:$AC$100,4,FALSE)"); err != nil {
		return err
	}

	if err := f.SetCellValue(sheetName, "B24", "c."); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "C24", "Tingkat Biaya Perjalanan Dinas"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "D24", "c."); err != nil {
		return err
	}

	if err := f.SetCellValue(sheetName, "A27", "4."); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "B27", "Maksud perjalanan dinas"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "D26", req.ActivityPurpose); err != nil {
		return err
	}

	if err := f.SetCellValue(sheetName, "A29", "5."); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "B29", "Alat angkut yang dipergunakan"); err != nil {
		return err
	}
	if err := f.SetCellFormula(sheetName, "D29", "=VLOOKUP($K$20,'PEMANTAUAN REKAP RAMPUNG'!$A$8:$AJ$100,34,FALSE)"); err != nil {
		return err
	}

	if err := f.SetCellValue(sheetName, "A31", "6."); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "B31", "a."); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "C31", "Tempat Berangkat"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "D31", "a."); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "E31", "Jakarta"); err != nil {
		return err
	}

	if err := f.SetCellValue(sheetName, "B32", "b."); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "C32", "Tempat Tujuan"); err != nil {
		return err
	}
	if err := f.SetCellFormula(sheetName, "E32", "=VLOOKUP($K$20,'PEMANTAUAN REKAP RAMPUNG'!$A$8:$AC$100,6,FALSE)"); err != nil {
		return err
	}

	if err := f.SetCellValue(sheetName, "B35", "a."); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "C35", "Lamanya perjalanan dinas"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "D35", "a."); err != nil {
		return err
	}
	if err := f.SetCellFormula(sheetName, "E35", `=VLOOKUP(K20,'PEMANTAUAN REKAP RAMPUNG'!$A$8:$Y$100,8,0)&" hari"`); err != nil {
		return err
	}

	if err := f.SetCellValue(sheetName, "A36", "7."); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "B36", "b."); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "C36", "Tanggal Berangkat"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "D36", "b."); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "E36", req.DepartureDate); err != nil {
		return err
	}

	if err := f.SetCellValue(sheetName, "B37", "c."); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "C37", "Tanggal harus Kembali/tiba di"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "D37", "c."); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "E37", req.ReturnDate); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "C38", " tempat baru *)"); err != nil {
		return err
	}

	if err := f.SetCellValue(sheetName, "B39", "Pengikut  :           N  a  m  a"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "E39", "Tanggal Lahir"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "H39", "Keterangan"); err != nil {
		return err
	}

	if err := f.SetCellValue(sheetName, "B39", "Pengikut  :           N  a  m  a           "); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "B40", 1); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "B41", 2); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "B42", 3); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "A42", "8."); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "B43", 4); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "B44", 5); err != nil {
		return err
	}

	if err := f.SetCellValue(sheetName, "A47", "9."); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "B46", "Pembebanan Anggaran"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "B47", "a."); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "C47", "Instansi"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "D47", "a."); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "E47", "Sekretariat Ditjen Pencegahan dan Pengendalian Penyakit"); err != nil {
		return err
	}

	if err := f.SetCellValue(sheetName, "B48", "b."); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "C48", "Akun"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "D48", "b."); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "E48", "024.05.WA.4815.EBD.953.501.B.524111"); err != nil {
		return err
	}

	if err := f.SetCellValue(sheetName, "A51", "10."); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "B51", "Keterangan lain-lain"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "B53", "*coret yang tidak perlu"); err != nil {
		return err
	}

	if err := f.SetCellValue(sheetName, "F55", "Dikeluarkan di"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "G55", ":"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "H55", "Jakarta"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "F56", "Tanggal"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "G56", ":"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "H56", req.ReceiptSignatureDate); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "F58", "Pejabat Pembuat Komitmen II"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "F59", "Unit Kerja Setditjen Penanggulangan Penyakit"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "F63", "Ruly Wahyuni, SE, MKM"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetName, "F64", "NIP 197508142000032001"); err != nil {
		return err
	}

	if err := f.SetCellValue(sheetName, "K20", 1); err != nil {
		return err
	}

	if err := f.SetColWidth(sheetName, "A", "A", 4.5); err != nil {
		return err
	}

	if err := f.SetColWidth(sheetName, "B", "B", 2.6); err != nil {
		return err
	}

	if err := f.SetColWidth(sheetName, "C", "C", 31.6); err != nil {
		return err
	}
	if err := f.SetColWidth(sheetName, "D", "D", 3.5); err != nil {
		return err
	}
	if err := f.SetColWidth(sheetName, "E", "E", 19.6); err != nil {
		return err
	}
	if err := f.SetColWidth(sheetName, "F", "F", 13.6); err != nil {
		return err
	}
	if err := f.SetColWidth(sheetName, "G", "G", 1.2); err != nil {
		return err
	}
	if err := f.SetColWidth(sheetName, "H", "H", 13.3); err != nil {
		return err
	}
	if err := f.SetColWidth(sheetName, "I", "I", 16.8); err != nil {
		return err
	}

	if err := f.SetCellStyle(sheetName, "A1", "I63", g.dynamicStyle(f, []string{}, false, false, 0, "left", 0, false, 0)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "A8", "A14", g.dynamicStyle(f, []string{}, true, false, 0, "left", 0, false, 0)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "A14", "A14", g.dynamicStyle(f, []string{}, true, false, 0, "center", 0, false, 0)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "A16", "I53", g.dynamicStyle(f, []string{}, false, false, 0, "left", 0, false, 9)); err != nil {
		return err
	}

	if err := f.MergeCell(sheetName, "A14", "I14"); err != nil {
		return err
	}

	if err := f.SetCellStyle(sheetName, "A16", "A16", g.dynamicStyle(f, []string{"left", "top", "right"}, false, false, 1, "center", 0, false, 9)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "B16", "C16", g.dynamicStyle(f, []string{"top"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "D16", "D16", g.dynamicStyle(f, []string{"left", "top"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "E16", "H16", g.dynamicStyle(f, []string{"top"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "I16", "I16", g.dynamicStyle(f, []string{"top", "right"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}

	if err := f.SetCellStyle(sheetName, "I17", "I17", g.dynamicStyle(f, []string{"bottom", "right"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "B17", "C17", g.dynamicStyle(f, []string{"bottom"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "D17", "D17", g.dynamicStyle(f, []string{"bottom", "left"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "E17", "H17", g.dynamicStyle(f, []string{"bottom"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "I17", "I17", g.dynamicStyle(f, []string{"bottom", "right"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "A17", "A17", g.dynamicStyle(f, []string{"bottom", "right", "left"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}

	if err := f.SetCellStyle(sheetName, "A18", "B20", g.dynamicStyle(f, []string{"left"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "B18", "B20", g.dynamicStyle(f, []string{"left"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "D18", "D20", g.dynamicStyle(f, []string{"left"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "I18", "I20", g.dynamicStyle(f, []string{"right"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "A19", "A19", g.dynamicStyle(f, []string{"left"}, false, false, 1, "center", 0, false, 9)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "A21", "A21", g.dynamicStyle(f, []string{"left", "bottom", "right"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "B21", "C21", g.dynamicStyle(f, []string{"bottom"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "D21", "D21", g.dynamicStyle(f, []string{"left", "bottom"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "E21", "H21", g.dynamicStyle(f, []string{"bottom"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "I21", "I21", g.dynamicStyle(f, []string{"bottom", "right"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "A22", "A24", g.dynamicStyle(f, []string{"left"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "A23", "A23", g.dynamicStyle(f, []string{"left"}, false, false, 1, "center", 0, false, 9)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "B22", "A24", g.dynamicStyle(f, []string{"left"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "D22", "D24", g.dynamicStyle(f, []string{"left"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "I22", "I24", g.dynamicStyle(f, []string{"right"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}

	if err := f.SetCellStyle(sheetName, "A25", "A25", g.dynamicStyle(f, []string{"bottom", "left", "right"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "B25", "C25", g.dynamicStyle(f, []string{"bottom"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "D25", "D25", g.dynamicStyle(f, []string{"bottom", "left"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "E25", "H25", g.dynamicStyle(f, []string{"bottom"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "I25", "I25", g.dynamicStyle(f, []string{"bottom", "right"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}

	if err := f.SetCellStyle(sheetName, "A26", "A26", g.dynamicStyle(f, []string{"left"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "B26", "B26", g.dynamicStyle(f, []string{"left"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "D26", "D27", g.dynamicStyle(f, []string{"left", "right", "bottom"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "D28", "D28", g.dynamicStyle(f, []string{"left", "bottom"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "A27", "A27", g.dynamicStyle(f, []string{"left"}, false, false, 1, "center", 0, false, 9)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "B27", "B27", g.dynamicStyle(f, []string{"left"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}
	if err := f.MergeCell(sheetName, "D26", "I28"); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "A28", "B28", g.dynamicStyle(f, []string{"left", "bottom"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "C28", "C28", g.dynamicStyle(f, []string{"bottom"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "E28", "H28", g.dynamicStyle(f, []string{"bottom"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "I26", "H27", g.dynamicStyle(f, []string{"right"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "I28", "H28", g.dynamicStyle(f, []string{"right", "bottom"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}

	if err := f.SetCellStyle(sheetName, "A29", "B29", g.dynamicStyle(f, []string{"left", "bottom"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "C29", "C29", g.dynamicStyle(f, []string{"left", "bottom"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "D29", "D29", g.dynamicStyle(f, []string{"left", "bottom"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "E29", "H29", g.dynamicStyle(f, []string{"bottom"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "I29", "I29", g.dynamicStyle(f, []string{"bottom", "right"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}

	if err := f.SetCellStyle(sheetName, "A30", "A32", g.dynamicStyle(f, []string{"left"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "B30", "B32", g.dynamicStyle(f, []string{"left"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "D30", "D32", g.dynamicStyle(f, []string{"left"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "I30", "I32", g.dynamicStyle(f, []string{"right"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "A33", "B33", g.dynamicStyle(f, []string{"left", "bottom"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "C33", "C33", g.dynamicStyle(f, []string{"bottom"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "D33", "D33", g.dynamicStyle(f, []string{"left", "bottom"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "E33", "H33", g.dynamicStyle(f, []string{"bottom"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "I33", "I33", g.dynamicStyle(f, []string{"bottom", "right"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}

	if err := f.SetCellStyle(sheetName, "A34", "B34", g.dynamicStyle(f, []string{"left"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "D34", "D34", g.dynamicStyle(f, []string{"left"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "I34", "I34", g.dynamicStyle(f, []string{"right"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "A35", "A37", g.dynamicStyle(f, []string{"left"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "B35", "B37", g.dynamicStyle(f, []string{"left"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "D35", "D37", g.dynamicStyle(f, []string{"left"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "I35", "I37", g.dynamicStyle(f, []string{"right"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "A38", "B38", g.dynamicStyle(f, []string{"left", "bottom"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "C38", "C38", g.dynamicStyle(f, []string{"bottom"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "D38", "D38", g.dynamicStyle(f, []string{"left", "bottom"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "E38", "H38", g.dynamicStyle(f, []string{"bottom"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "I38", "I38", g.dynamicStyle(f, []string{"bottom", "right"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}

	if err := f.SetCellStyle(sheetName, "A39", "B39", g.dynamicStyle(f, []string{"left"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "D39", "D39", g.dynamicStyle(f, []string{"left"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "I39", "I39", g.dynamicStyle(f, []string{"right"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "A40", "A44", g.dynamicStyle(f, []string{"left"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "B40", "B44", g.dynamicStyle(f, []string{"left"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "D40", "D44", g.dynamicStyle(f, []string{"left"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "I40", "I44", g.dynamicStyle(f, []string{"right"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "A45", "B45", g.dynamicStyle(f, []string{"left", "bottom"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "C45", "C45", g.dynamicStyle(f, []string{"bottom"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "D45", "D45", g.dynamicStyle(f, []string{"left", "bottom"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "E45", "H45", g.dynamicStyle(f, []string{"bottom"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "I45", "I45", g.dynamicStyle(f, []string{"bottom", "right"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}

	if err := f.SetCellStyle(sheetName, "A46", "B46", g.dynamicStyle(f, []string{"left"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "D46", "D46", g.dynamicStyle(f, []string{"left"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "I46", "I46", g.dynamicStyle(f, []string{"right"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "A47", "A48", g.dynamicStyle(f, []string{"left"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "B47", "B48", g.dynamicStyle(f, []string{"left"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "D47", "D48", g.dynamicStyle(f, []string{"left"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "I47", "I48", g.dynamicStyle(f, []string{"right"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "A49", "B49", g.dynamicStyle(f, []string{"left", "bottom"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "C49", "C49", g.dynamicStyle(f, []string{"bottom"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "D49", "D49", g.dynamicStyle(f, []string{"left", "bottom"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "E49", "H49", g.dynamicStyle(f, []string{"bottom"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "I49", "I49", g.dynamicStyle(f, []string{"bottom", "right"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}

	if err := f.SetCellStyle(sheetName, "A50", "B50", g.dynamicStyle(f, []string{"left"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "D50", "D50", g.dynamicStyle(f, []string{"left"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "I50", "I50", g.dynamicStyle(f, []string{"right"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "A51", "A51", g.dynamicStyle(f, []string{"left"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "B51", "B51", g.dynamicStyle(f, []string{"left"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "D51", "D51", g.dynamicStyle(f, []string{"left"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "I51", "I51", g.dynamicStyle(f, []string{"right"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "A52", "B52", g.dynamicStyle(f, []string{"left", "bottom"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "C52", "C52", g.dynamicStyle(f, []string{"bottom"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "D52", "D52", g.dynamicStyle(f, []string{"left", "bottom"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "E52", "H52", g.dynamicStyle(f, []string{"bottom"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "I52", "I52", g.dynamicStyle(f, []string{"bottom", "right"}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}

	if err := f.SetCellStyle(sheetName, "A31", "A31", g.dynamicStyle(f, []string{"left"}, false, false, 1, "center", 0, false, 9)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "A29", "A29", g.dynamicStyle(f, []string{"left", "bottom"}, false, false, 1, "center", 0, false, 9)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "A36", "A36", g.dynamicStyle(f, []string{"left"}, false, false, 1, "center", 0, false, 9)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "A42", "A42", g.dynamicStyle(f, []string{"left"}, false, false, 1, "center", 0, false, 9)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "A47", "A47", g.dynamicStyle(f, []string{"left"}, false, false, 1, "center", 0, false, 9)); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetName, "A51", "A51", g.dynamicStyle(f, []string{"left"}, false, false, 1, "center", 0, false, 9)); err != nil {
		return err
	}

	if err := f.SetCellStyle(sheetName, "F1", "F4", g.dynamicStyle(f, []string{}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}

	if err := f.MergeCell(sheetName, "F2", "I2"); err != nil {
		return err
	}

	if err := f.MergeCell(sheetName, "F4", "I4"); err != nil {
		return err
	}

	if err := f.SetCellStyle(sheetName, "F10", "H12", g.dynamicStyle(f, []string{}, false, false, 1, "left", 0, false, 12)); err != nil {
		return err
	}

	if err := f.SetCellStyle(sheetName, "F55", "H64", g.dynamicStyle(f, []string{}, false, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}

	if err := f.SetCellStyle(sheetName, "F58", "H63", g.dynamicStyle(f, []string{}, true, false, 1, "left", 0, false, 9)); err != nil {
		return err
	}

	if err := f.SetCellStyle(sheetName, "D26", "D26", g.dynamicStyle(f, []string{}, false, false, 1, "left", 0, true, 9)); err != nil {
		return err
	}

	if err := f.SetRowHeight(sheetName, 2, 24.8); err != nil {
		return err
	}
	if err := f.SetRowHeight(sheetName, 4, 36); err != nil {
		return err
	}

	if err := f.SetRowHeight(sheetName, 22, 22); err != nil {
		return err
	}
	if err := f.SetRowHeight(sheetName, 23, 22); err != nil {
		return err
	}
	if err := f.SetRowHeight(sheetName, 24, 22); err != nil {
		return err
	}

	if err := f.SetRowHeight(sheetName, 31, 20); err != nil {
		return err
	}
	if err := f.SetRowHeight(sheetName, 32, 20); err != nil {
		return err
	}

	if err := f.SetRowHeight(sheetName, 35, 18); err != nil {
		return err
	}
	if err := f.SetRowHeight(sheetName, 36, 18); err != nil {
		return err
	}
	if err := f.SetRowHeight(sheetName, 37, 18); err != nil {
		return err
	}

	portrait := "portrait"
	excelizeOne := 1
	excelizeTrue := true
	if err := f.SetPageLayout(sheetName, &excelize.PageLayoutOptions{
		Size:        utils.ToPtr(14),
		Orientation: &portrait,
		FitToWidth:  &excelizeOne,
		FitToHeight: &excelizeOne,
	}); err != nil {
		panic(err)
	}

	_ = f.SetPageMargins(sheetName, &excelize.PageLayoutMarginsOptions{
		Left:         utils.ToPtr(0.2362204),
		Right:        utils.ToPtr(0.2362204),
		Top:          utils.ToPtr(0.3543307),
		Bottom:       utils.ToPtr(0.3543307),
		Header:       utils.ToPtr(0.3149606),
		Footer:       utils.ToPtr(0.3149606),
		Horizontally: &excelizeTrue,
		Vertically:   &excelizeTrue,
	})

	zoomScale := 136.0
	_ = f.SetSheetView(sheetName, 0, &excelize.ViewOptions{
		ZoomScale: &zoomScale,
	})

	_ = f.DeleteDefinedName(&excelize.DefinedName{
		Name:  "_xlnm.Print_Area",
		Scope: sheetName,
	})

	ref := fmt.Sprintf("'%s'!$A$1:$I$64", sheetName)

	if err := f.SetDefinedName(&excelize.DefinedName{
		Name:     "_xlnm.Print_Area",
		RefersTo: ref,
		Scope:    sheetName,
	}); err != nil {
		return err
	}

	return nil
}

func (g *Generator) GenerateRecapExcel(req RecapReport) (*bytes.Buffer, error) {
	if len(req.Assignees) == 0 {
		return nil, fmt.Errorf("no assignees provided")
	}

	f := excelize.NewFile()
	defer f.Close()

	var err error
	sheetName := "PEMANTAUAN REKAP UANG MUKA"
	firstSheet := f.GetSheetName(f.GetActiveSheetIndex())
	if err := f.SetSheetName(firstSheet, sheetName); err != nil {
		return nil, err
	}

	if err = g.generateTitle(f, sheetName); err != nil {
		return nil, err
	}

	if err = g.generateTableHeader(f, sheetName); err != nil {
		return nil, err
	}

	currentRow := 11
	if currentRow, err = g.generateTableData(f, sheetName, req, currentRow); err != nil {
		return nil, err
	}

	if err = g.generateSummaryRow(f, sheetName, currentRow); err != nil {
		return nil, err
	}

	if err = g.generateSignatureRow(f, sheetName, currentRow); err != nil {
		return nil, err
	}

	if err = g.setColumnWidths(f, sheetName); err != nil {
		return nil, err
	}

	kwUangMuka := "KW UANG MUKA"

	err = g.generateKw(f, kwUangMuka, req)
	if err != nil {
		return nil, err
	}

	sheetName = "PEMANTAUAN REKAP RAMPUNG"
	if _, err = f.NewSheet(sheetName); err != nil {
		return nil, err
	}

	if err = g.generateTitle(f, sheetName); err != nil {
		return nil, err
	}

	if err = g.generateTableHeader(f, sheetName); err != nil {
		return nil, err
	}

	if currentRow, err = g.generateTableData(f, sheetName, req, 11); err != nil {
		return nil, err
	}

	if err = g.generateSummaryRow(f, sheetName, currentRow); err != nil {
		return nil, err
	}

	if err = g.generateSignatureRow(f, sheetName, currentRow); err != nil {
		return nil, err
	}

	if err = g.setColumnWidths(f, sheetName); err != nil {
		return nil, err
	}

	kwRampung := "KW RAMPUNG"
	err = g.generateKw(f, kwRampung, req)
	if err != nil {
		return nil, err
	}

	sppd := "SPPD"
	err = g.generateSppd(f, sppd, req)
	if err != nil {
		return nil, err
	}

	var b bytes.Buffer
	if err = f.Write(&b); err != nil {
		return nil, fmt.Errorf("failed to write excel to buffer: %w", err)
	}

	return &b, nil
}

func (g *Generator) dynamicStyle(f *excelize.File, borderTypes []string, bold bool, italic bool, borderStyle int, horizontal string, numFmt int, wrapText bool, fontSize float64) int {
	var borders []excelize.Border

	for _, side := range borderTypes {
		switch side {
		case "top":
			borders = append(borders, excelize.Border{
				Type:  "top",
				Color: "000000",
				Style: borderStyle,
			})
		case "bottom":
			borders = append(borders, excelize.Border{
				Type:  "bottom",
				Color: "000000",
				Style: borderStyle,
			})
		case "left":
			borders = append(borders, excelize.Border{
				Type:  "left",
				Color: "000000",
				Style: borderStyle,
			})
		case "right":
			borders = append(borders, excelize.Border{
				Type:  "right",
				Color: "000000",
				Style: borderStyle,
			})
		}
	}

	var fontSizeVal float64 = 10
	if fontSize > 0 {
		fontSizeVal = fontSize
	}

	styleModel := &excelize.Style{
		Alignment: &excelize.Alignment{
			WrapText:   wrapText,
			Horizontal: horizontal,
			Vertical:   "center",
		},
		Font: &excelize.Font{
			Italic: italic,
			Bold:   bold,
			Size:   fontSizeVal,
			Family: "Tahoma",
		},
		Border: borders,
	}

	if numFmt > 0 {
		styleModel.NumFmt = numFmt
	}

	style, _ := f.NewStyle(styleModel)

	return style
}

var angka = []string{
	"", "Satu", "Dua", "Tiga", "Empat", "Lima", "Enam", "Tujuh", "Delapan", "Sembilan", "Sepuluh", "Sebelas",
}

func Terbilang(n int64) string {
	if n < 12 {
		return angka[n]
	} else if n < 20 {
		return Terbilang(n-10) + " Belas"
	} else if n < 100 {
		return Terbilang(n/10) + " Puluh " + Terbilang(n%10)
	} else if n < 200 {
		return "Seratus " + Terbilang(n-100)
	} else if n < 1000 {
		return Terbilang(n/100) + " Ratus " + Terbilang(n%100)
	} else if n < 2000 {
		return "Seribu " + Terbilang(n-1000)
	} else if n < 1000000 {
		return Terbilang(n/1000) + " Ribu " + Terbilang(n%1000)
	} else if n < 1000000000 {
		return Terbilang(n/1000000) + " Juta " + Terbilang(n%1000000)
	} else if n < 1000000000000 {
		return Terbilang(n/1000000000) + " Miliar " + Terbilang(n%1000000000)
	} else if n < 1000000000000000 {
		return Terbilang(n/1000000000000) + " Triliun " + Terbilang(n%1000000000000)
	}
	return ""
}

func CleanSpace(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

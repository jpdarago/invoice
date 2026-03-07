package main

import (
	"fmt"
	"image"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/signintech/gopdf"
)

const (
	quantityColumnOffset = 360
	rateColumnOffset     = 405
	amountColumnOffset   = 480
)

const (
	subtotalLabel = "Subtotal"
	discountLabel = "Discount"
	taxLabel      = "Tax"
	totalLabel    = "Total"
)

func writeLogo(pdf *gopdf.GoPdf, logo string, from string) error {
	if logo != "" {
		width, height, err := getImageDimension(logo)
		if err != nil {
			return fmt.Errorf("unable to read logo %s: %w", logo, err)
		}
		scaledWidth := 100.0
		scaledHeight := float64(height) * scaledWidth / float64(width)
		if err := pdf.Image(logo, pdf.GetX(), pdf.GetY(), &gopdf.Rect{W: scaledWidth, H: scaledHeight}); err != nil {
			return err
		}
		pdf.Br(scaledHeight + 24)
	}
	pdf.SetTextColor(55, 55, 55)

	formattedFrom := strings.ReplaceAll(from, `\n`, "\n")
	fromLines := strings.Split(formattedFrom, "\n")

	for i := 0; i < len(fromLines); i++ {
		if i == 0 {
			if err := pdf.SetFont("Inter", "", 12); err != nil {
				return err
			}
			if err := pdf.Cell(nil, fromLines[i]); err != nil {
				return err
			}
			pdf.Br(18)
		} else {
			if err := pdf.SetFont("Inter", "", 10); err != nil {
				return err
			}
			if err := pdf.Cell(nil, fromLines[i]); err != nil {
				return err
			}
			pdf.Br(15)
		}
	}
	pdf.Br(21)
	pdf.SetStrokeColor(225, 225, 225)
	pdf.Line(pdf.GetX(), pdf.GetY(), 260, pdf.GetY())
	pdf.Br(36)
	return nil
}

func writeTitle(pdf *gopdf.GoPdf, title, id, date string) error {
	if err := pdf.SetFont("Inter-Bold", "", 24); err != nil {
		return err
	}
	pdf.SetTextColor(0, 0, 0)
	if err := pdf.Cell(nil, title); err != nil {
		return err
	}
	pdf.Br(36)
	if err := pdf.SetFont("Inter", "", 12); err != nil {
		return err
	}
	pdf.SetTextColor(100, 100, 100)
	if err := pdf.Cell(nil, "#"); err != nil {
		return err
	}
	if err := pdf.Cell(nil, id); err != nil {
		return err
	}
	pdf.SetTextColor(150, 150, 150)
	if err := pdf.Cell(nil, "  ·  "); err != nil {
		return err
	}
	pdf.SetTextColor(100, 100, 100)
	if err := pdf.Cell(nil, date); err != nil {
		return err
	}
	pdf.Br(48)
	return nil
}

func writeDueDate(pdf *gopdf.GoPdf, due string) error {
	if err := pdf.SetFont("Inter", "", 9); err != nil {
		return err
	}
	pdf.SetTextColor(75, 75, 75)
	pdf.SetX(rateColumnOffset)
	if err := pdf.Cell(nil, "Due Date"); err != nil {
		return err
	}
	pdf.SetTextColor(0, 0, 0)
	if err := pdf.SetFontSize(11); err != nil {
		return err
	}
	pdf.SetX(amountColumnOffset - 15)
	if err := pdf.Cell(nil, due); err != nil {
		return err
	}
	pdf.Br(12)
	return nil
}

func writeBillTo(pdf *gopdf.GoPdf, to string) error {
	pdf.SetTextColor(75, 75, 75)
	if err := pdf.SetFont("Inter", "", 9); err != nil {
		return err
	}
	if err := pdf.Cell(nil, "BILL TO"); err != nil {
		return err
	}
	pdf.Br(18)
	pdf.SetTextColor(75, 75, 75)

	formattedTo := strings.ReplaceAll(to, `\n`, "\n")
	toLines := strings.Split(formattedTo, "\n")

	for i := 0; i < len(toLines); i++ {
		if i == 0 {
			if err := pdf.SetFont("Inter", "", 15); err != nil {
				return err
			}
			if err := pdf.Cell(nil, toLines[i]); err != nil {
				return err
			}
			pdf.Br(20)
		} else {
			if err := pdf.SetFont("Inter", "", 10); err != nil {
				return err
			}
			if err := pdf.Cell(nil, toLines[i]); err != nil {
				return err
			}
			pdf.Br(15)
		}
	}
	pdf.Br(64)
	return nil
}

func writeHeaderRow(pdf *gopdf.GoPdf) error {
	if err := pdf.SetFont("Inter", "", 9); err != nil {
		return err
	}
	pdf.SetTextColor(55, 55, 55)
	if err := pdf.Cell(nil, "ITEM"); err != nil {
		return err
	}
	pdf.SetX(quantityColumnOffset)
	if err := pdf.Cell(nil, "QTY"); err != nil {
		return err
	}
	pdf.SetX(rateColumnOffset)
	if err := pdf.Cell(nil, "RATE"); err != nil {
		return err
	}
	pdf.SetX(amountColumnOffset)
	if err := pdf.Cell(nil, "AMOUNT"); err != nil {
		return err
	}
	pdf.Br(24)
	return nil
}

func writeNotes(pdf *gopdf.GoPdf, notes string) error {
	pdf.SetY(600)

	if err := pdf.SetFont("Inter", "", 9); err != nil {
		return err
	}
	pdf.SetTextColor(55, 55, 55)
	if err := pdf.Cell(nil, "NOTES"); err != nil {
		return err
	}
	pdf.Br(18)
	if err := pdf.SetFont("Inter", "", 9); err != nil {
		return err
	}
	pdf.SetTextColor(0, 0, 0)

	formattedNotes := strings.ReplaceAll(notes, `\n`, "\n")
	notesLines := strings.Split(formattedNotes, "\n")

	for i := 0; i < len(notesLines); i++ {
		if err := pdf.Cell(nil, notesLines[i]); err != nil {
			return err
		}
		pdf.Br(15)
	}

	pdf.Br(48)
	return nil
}

func writeMetadata(pdf *gopdf.GoPdf, metadata map[string]string) error {
	pdf.SetY(690)

	if err := pdf.SetFont("Inter", "", 9); err != nil {
		return err
	}
	pdf.SetTextColor(55, 55, 55)
	if err := pdf.Cell(nil, "DETAILS"); err != nil {
		return err
	}
	pdf.Br(18)

	keys := make([]string, 0, len(metadata))
	for k := range metadata {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		if err := pdf.SetFont("Inter-Bold", "", 9); err != nil {
			return err
		}
		pdf.SetTextColor(55, 55, 55)
		if err := pdf.Cell(nil, k+": "); err != nil {
			return err
		}
		if err := pdf.SetFont("Inter", "", 9); err != nil {
			return err
		}
		pdf.SetTextColor(0, 0, 0)
		if err := pdf.Cell(nil, metadata[k]); err != nil {
			return err
		}
		pdf.Br(15)
	}

	return nil
}

func writeFooter(pdf *gopdf.GoPdf, id string) error {
	pdf.SetY(800)

	if err := pdf.SetFont("Inter", "", 10); err != nil {
		return err
	}
	pdf.SetTextColor(55, 55, 55)
	if err := pdf.Cell(nil, id); err != nil {
		return err
	}
	pdf.SetStrokeColor(225, 225, 225)
	pdf.Line(pdf.GetX()+10, pdf.GetY()+6, 550, pdf.GetY()+6)
	pdf.Br(48)
	return nil
}

func writeRow(pdf *gopdf.GoPdf, item string, quantity int, rate float64, currencySymbol string) error {
	if err := pdf.SetFont("Inter", "", 11); err != nil {
		return err
	}
	pdf.SetTextColor(0, 0, 0)

	total := float64(quantity) * rate
	amount := strconv.FormatFloat(total, 'f', 2, 64)

	if err := pdf.Cell(nil, item); err != nil {
		return err
	}
	pdf.SetX(quantityColumnOffset)
	if err := pdf.Cell(nil, strconv.Itoa(quantity)); err != nil {
		return err
	}
	pdf.SetX(rateColumnOffset)
	if err := pdf.Cell(nil, currencySymbol+strconv.FormatFloat(rate, 'f', 2, 64)); err != nil {
		return err
	}
	pdf.SetX(amountColumnOffset)
	if err := pdf.Cell(nil, currencySymbol+amount); err != nil {
		return err
	}
	pdf.Br(24)
	return nil
}

func writeTotals(pdf *gopdf.GoPdf, subtotal float64, tax float64, discount float64, currencySymbol string) error {
	pdf.SetY(600)

	if err := writeTotal(pdf, subtotalLabel, subtotal, currencySymbol); err != nil {
		return err
	}
	if tax > 0 {
		if err := writeTotal(pdf, taxLabel, tax, currencySymbol); err != nil {
			return err
		}
	}
	if discount > 0 {
		if err := writeTotal(pdf, discountLabel, discount, currencySymbol); err != nil {
			return err
		}
	}
	return writeTotal(pdf, totalLabel, subtotal+tax-discount, currencySymbol)
}

func writeTotal(pdf *gopdf.GoPdf, label string, total float64, currencySymbol string) error {
	if err := pdf.SetFont("Inter", "", 9); err != nil {
		return err
	}
	pdf.SetTextColor(75, 75, 75)
	pdf.SetX(rateColumnOffset)
	if err := pdf.Cell(nil, label); err != nil {
		return err
	}
	pdf.SetTextColor(0, 0, 0)
	if err := pdf.SetFontSize(12); err != nil {
		return err
	}
	pdf.SetX(amountColumnOffset - 15)
	if label == totalLabel {
		if err := pdf.SetFont("Inter-Bold", "", 11.5); err != nil {
			return err
		}
	}
	if err := pdf.Cell(nil, currencySymbol+strconv.FormatFloat(total, 'f', 2, 64)); err != nil {
		return err
	}
	pdf.Br(24)
	return nil
}

func getImageDimension(imagePath string) (int, int, error) {
	f, err := os.Open(imagePath)
	if err != nil {
		return 0, 0, err
	}
	defer f.Close()

	img, _, err := image.DecodeConfig(f)
	if err != nil {
		return 0, 0, err
	}
	return img.Width, img.Height, nil
}

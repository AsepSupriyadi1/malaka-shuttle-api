package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"malakashuttle/dto"

	"github.com/jung-kurt/gofpdf/v2"
)

// PDFReceiptGenerator handles PDF receipt generation
type PDFReceiptGenerator struct{}

// NewPDFReceiptGenerator creates a new PDF receipt generator
func NewPDFReceiptGenerator() *PDFReceiptGenerator {
	return &PDFReceiptGenerator{}
}

// GenerateBookingReceipt generates a PDF receipt for a booking
func (p *PDFReceiptGenerator) GenerateBookingReceipt(booking *dto.BookingResponse, outputPath string) error {
	// Create new PDF
	pdf := gofpdf.New(gofpdf.OrientationPortrait, gofpdf.UnitMillimeter, gofpdf.PageSizeA4, "")
	pdf.AddPage()

	// Set font
	pdf.SetFont("Arial", "B", 16)

	// Company Header
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(0, 10, "MALAKA SHUTTLE", "", 1, "C", false, 0, "")

	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(0, 6, "Jasa Transportasi Terpercaya", "", 1, "C", false, 0, "")
	pdf.Ln(10)

	// Title
	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(0, 8, "BOOKING RECEIPT", "", 1, "C", false, 0, "")
	pdf.Ln(12)

	// Booking Information
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(40, 6, "Booking ID:", "", 0, "L", false, 0, "")
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(0, 6, fmt.Sprintf("#%d", booking.ID), "", 1, "L", false, 0, "")

	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(40, 6, "Booking Date:", "", 0, "L", false, 0, "")
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(0, 6, booking.CreatedAt.Format("02 January 2006, 15:04 WIB"), "", 1, "L", false, 0, "")

	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(40, 6, "Status:", "", 0, "L", false, 0, "")
	pdf.SetFont("Arial", "", 10)
	statusColor := getStatusColor(string(booking.Status))
	pdf.SetTextColor(statusColor.R, statusColor.G, statusColor.B)
	pdf.CellFormat(0, 6, string(booking.Status), "", 1, "L", false, 0, "")
	pdf.SetTextColor(0, 0, 0) // Reset to black
	pdf.Ln(10)

	// Schedule Information
	if booking.Schedule != nil {
		pdf.SetFont("Arial", "B", 12)
		pdf.CellFormat(0, 8, "SCHEDULE DETAILS", "", 1, "L", false, 0, "")
		pdf.Ln(8)

		pdf.SetFont("Arial", "B", 10)
		pdf.CellFormat(40, 6, "Route:", "", 0, "L", false, 0, "")
		pdf.SetFont("Arial", "", 10)
		pdf.CellFormat(0, 6, fmt.Sprintf("%s → %s", booking.Schedule.Origin, booking.Schedule.Destination), "", 1, "L", false, 0, "")

		pdf.SetFont("Arial", "B", 10)
		pdf.CellFormat(40, 6, "Departure:", "", 0, "L", false, 0, "")
		pdf.SetFont("Arial", "", 10)
		pdf.CellFormat(0, 6, booking.Schedule.DepartureTime, "", 1, "L", false, 0, "")

		pdf.SetFont("Arial", "B", 10)
		pdf.CellFormat(40, 6, "Arrival:", "", 0, "L", false, 0, "")
		pdf.SetFont("Arial", "", 10)
		pdf.CellFormat(0, 6, booking.Schedule.ArrivalTime, "", 1, "L", false, 0, "")

		pdf.SetFont("Arial", "B", 10)
		pdf.CellFormat(40, 6, "Duration:", "", 0, "L", false, 0, "")
		pdf.SetFont("Arial", "", 10)
		pdf.CellFormat(0, 6, booking.Schedule.Duration, "", 1, "L", false, 0, "")
		pdf.Ln(10)
	}

	// Passenger Details
	if len(booking.Passengers) > 0 {
		pdf.SetFont("Arial", "B", 12)
		pdf.CellFormat(0, 8, "PASSENGER DETAILS", "", 1, "L", false, 0, "")
		pdf.Ln(8)

		// Table header
		pdf.SetFont("Arial", "B", 10)
		pdf.SetFillColor(240, 240, 240)
		pdf.CellFormat(10, 8, "No", "1", 0, "C", true, 0, "")
		pdf.CellFormat(80, 8, "Passenger Name", "1", 0, "C", true, 0, "")
		pdf.CellFormat(30, 8, "Seat", "1", 0, "C", true, 0, "")
		pdf.CellFormat(40, 8, "Price", "1", 1, "C", true, 0, "")

		// Table content
		pdf.SetFont("Arial", "", 10)
		pdf.SetFillColor(255, 255, 255)

		for i, passenger := range booking.Passengers {
			pdf.CellFormat(10, 8, fmt.Sprintf("%d", i+1), "1", 0, "C", false, 0, "")
			pdf.CellFormat(80, 8, passenger.PassengerName, "1", 0, "L", false, 0, "")
			pdf.CellFormat(30, 8, passenger.SeatNumber, "1", 0, "C", false, 0, "")
			if booking.Schedule != nil {
				pdf.CellFormat(40, 8, formatCurrency(booking.Schedule.Price), "1", 1, "R", false, 0, "")
			} else {
				pdf.CellFormat(40, 8, "-", "1", 1, "R", false, 0, "")
			}
		}
		pdf.Ln(5)
	}

	// Payment Summary
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(0, 8, "PAYMENT SUMMARY", "", 1, "L", false, 0, "")
	pdf.Ln(8)
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(120, 6, "Total Amount:", "", 0, "L", false, 0, "")
	pdf.CellFormat(40, 6, formatCurrency(booking.TotalAmount), "", 1, "R", false, 0, "")
	pdf.Ln(5)

	// Payment Status
	if booking.Status == "success" {
		pdf.SetFont("Arial", "B", 10)
		pdf.CellFormat(120, 6, "Payment Status:", "", 0, "L", false, 0, "")
		pdf.SetTextColor(0, 150, 0) // Green
		pdf.CellFormat(40, 6, "PAID", "", 1, "R", false, 0, "")
		pdf.SetTextColor(0, 0, 0) // Reset to black
	} else {
		pdf.SetFont("Arial", "B", 10)
		pdf.CellFormat(120, 6, "Payment Status:", "", 0, "L", false, 0, "")
		pdf.SetTextColor(200, 0, 0) // Red
		pdf.CellFormat(40, 6, "PENDING", "", 1, "R", false, 0, "")
		pdf.SetTextColor(0, 0, 0) // Reset to black
	}

	// Footer
	pdf.Ln(20)
	pdf.SetFont("Arial", "I", 8)
	pdf.CellFormat(0, 6, "Thank you for choosing Malaka Shuttle!", "", 1, "C", false, 0, "")
	pdf.CellFormat(0, 6, fmt.Sprintf("Generated on %s", time.Now().Format("02 January 2006, 15:04 WIB")), "", 1, "C", false, 0, "")

	// Save PDF
	return pdf.OutputFileAndClose(outputPath)
}

// StatusColor represents RGB color
type StatusColor struct {
	R, G, B int
}

// getStatusColor returns color based on booking status
func getStatusColor(status string) StatusColor {
	switch status {
	case "success":
		return StatusColor{0, 150, 0} // Green
	case "pending":
		return StatusColor{255, 165, 0} // Orange
	case "waiting_verification":
		return StatusColor{0, 0, 255} // Blue
	case "rejected":
		return StatusColor{200, 0, 0} // Red
	case "expired":
		return StatusColor{128, 128, 128} // Gray
	default:
		return StatusColor{0, 0, 0} // Black
	}
}

// formatCurrency formats number to Indonesian Rupiah currency
func formatCurrency(amount float64) string {
	// Convert to string without decimal
	str := strconv.FormatFloat(amount, 'f', 0, 64)

	// Add thousand separators
	formatted := addThousandSeparator(str)

	return "Rp " + formatted
}

// addThousandSeparator adds thousand separators to a number string
func addThousandSeparator(str string) string {
	n := len(str)
	if n <= 3 {
		return str
	}

	var result strings.Builder
	for i, char := range str {
		if i > 0 && (n-i)%3 == 0 {
			result.WriteString(".")
		}
		result.WriteRune(char)
	}

	return result.String()
}

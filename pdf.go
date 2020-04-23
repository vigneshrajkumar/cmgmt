package main

import (
	"fmt"
	"net/http"

	"github.com/jung-kurt/gofpdf"
)

func handleGETPDFMember(w http.ResponseWriter, r *http.Request) {

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Hello, world")
	err := pdf.OutputFileAndClose("hello.pdf")
	if err != nil {
		fmt.Println(err)
	}

	// pdf.byt

	w.Header().Set("Content-Type", "application/pdf ")
	// w.Write(pdf.Bytes())
}

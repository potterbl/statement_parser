package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/otiai10/gosseract"
)

func ParsePDFFileWithOCR(file *multipart.FileHeader, languages []string) ([]string, error) {
	tmpFile, err := os.CreateTemp("", "*.pdf")
	if err != nil {
		return []string{}, fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	fileOpened, err := file.Open()
	if err != nil {
		return []string{}, fmt.Errorf("failed to open file: %w", err)
	}
	defer fileOpened.Close()

	if _, err := io.Copy(tmpFile, fileOpened); err != nil {
		return []string{}, fmt.Errorf("failed to copy to temp file: %w", err)
	}
	tmpFile.Close()

	// Create image from PDF
	outputPrefix := strings.TrimSuffix(tmpFile.Name(), ".pdf")
	imagePaths, err := ConvertPDFToImage(tmpFile.Name(), outputPrefix)
	if err != nil {
		return []string{}, fmt.Errorf("failed to convert PDF to image: %w", err)
	}

	client := gosseract.NewClient()
	defer client.Close()

	client.SetLanguage(languages...)

	var (
		chunks    = make([]string, len(imagePaths))
		allErrors []error
		wg        sync.WaitGroup
		mu        sync.Mutex
	)

	for _, imagePath := range imagePaths {
		wg.Add(1)
		go func() {
			defer os.Remove(imagePath)
			defer wg.Done()

			client.SetImage(imagePath)
			t, err := client.Text()

			mu.Lock()
			defer mu.Unlock()

			if err != nil {
				allErrors = append(allErrors, fmt.Errorf("OCR failed for image %s: %w", imagePath, err))
				return
			}

			chunks = append(chunks, t)
		}()
	}

	wg.Wait()

	if len(allErrors) > 0 {
		return chunks, fmt.Errorf("OCR errors occurred: %v", allErrors)
	}
	return chunks, nil
}

func ConvertPDFToImage(pdfPath string, outputPrefix string) ([]string, error) {
	cmd := exec.Command("pdftoppm", "-png", pdfPath, outputPrefix)
	err := cmd.Run()
	if err != nil {
		return []string{}, fmt.Errorf("error converting PDF to image: %w", err)
	}

	files, err := filepath.Glob(outputPrefix + "-*.png")
	if err != nil || len(files) == 0 {
		return []string{}, fmt.Errorf("no PNG files created at prefix: %s", outputPrefix)
	}

	return files, nil
}

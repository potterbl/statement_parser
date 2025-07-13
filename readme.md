# Bank Statement Parser for Ukrainian Banks

A powerful Go-based tool for parsing bank statements from Ukrainian banks using advanced OCR (Optical Character Recognition) technology. This parser extracts transaction data from PDF bank statements and converts them into a structured JSON format, making financial data analysis seamless and efficient.

## 📜 Table of Contents

- [Features](#-features)
- [Supported Banks](#-supported-banks)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Usage Examples](#-usage-examples)
- [API Reference](#-api-reference)
- [Project Structure](#-project-structure)
- [Contributing](#contributing)
- [License](#-license)
- [Contact](#-contact)
- [Acknowledgments](#-acknowledgments)


## 🌟 Features

- **PDF to Text Conversion**: Converts PDF bank statements to text using OCR
- **Multi-language Support**: Supports multiple languages for OCR processing
- **Modular Architecture**: Extensible design for adding support for different banks
- **Concurrent Processing**: Efficient processing of multi-page PDF documents
- **Structured Output**: Returns transaction data in a standardized JSON format

## 🏦 Supported Banks

- **Monobank** - Full statement parsing support
- **Privat24** - Full statement parsing support
- *More banks coming soon!*

## Prerequisites

Before running the application, ensure you have the following dependencies installed:

### System Dependencies

#### Linux (Ubuntu/Debian)

1. **Update package list**:
   ```bash
   sudo apt update
   ```

2. **Install pdftoppm** (part of poppler-utils):
   ```bash
   sudo apt install poppler-utils
   ```

3. **Install Tesseract OCR**:
   ```bash
   sudo apt install tesseract-ocr
   ```

4. **Install language packages for Tesseract**:
   ```bash
   # Ukrainian language support
   sudo apt install tesseract-ocr-ukr

   # English language support (usually included by default)
   sudo apt install tesseract-ocr-eng

   # Russian language support
   sudo apt install tesseract-ocr-rus
   ```

5. **Verify installation**:
   ```bash
   # Check pdftoppm
   pdftoppm -h

   # Check tesseract
   tesseract --version

   # List available languages
   tesseract --list-langs
   ```

#### Linux (CentOS/RHEL/Fedora)

1. **For CentOS/RHEL** (enable EPEL repository first):
   ```bash
   # CentOS 7
   sudo yum install epel-release
   sudo yum install poppler-utils tesseract tesseract-langpack-ukr tesseract-langpack-eng tesseract-langpack-rus

   # CentOS 8/RHEL 8
   sudo dnf install epel-release
   sudo dnf install poppler-utils tesseract tesseract-langpack-ukr tesseract-langpack-eng tesseract-langpack-rus
   ```

2. **For Fedora**:
   ```bash
   sudo dnf install poppler-utils tesseract tesseract-langpack-ukr tesseract-langpack-eng tesseract-langpack-rus
   ```

#### macOS

1. **Install Homebrew** (if not already installed):
   ```bash
   /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
   ```

2. **Install dependencies using Homebrew**:
   ```bash
   # Install poppler (includes pdftoppm)
   brew install poppler

   # Install tesseract
   brew install tesseract
   ```

3. **Install language packages**:
   ```bash
   # Install additional language packs
   brew install tesseract-lang
   ```

   Or install specific languages:
   ```bash
   # Note: Language packs are typically included with tesseract-lang
   # You can also download specific .traineddata files manually

   # Check available languages after installation
   tesseract --list-langs
   ```

4. **Manual language installation** (if needed):
   ```bash
   # Create tessdata directory if it doesn't exist
   mkdir -p /usr/local/share/tessdata

   # Download Ukrainian language data
   curl -L https://github.com/tesseract-ocr/tessdata/raw/main/ukr.traineddata -o /usr/local/share/tessdata/ukr.traineddata

   # Download Russian language data
   curl -L https://github.com/tesseract-ocr/tessdata/raw/main/rus.traineddata -o /usr/local/share/tessdata/rus.traineddata
   ```

5. **Verify installation**:
   ```bash
   # Check pdftoppm
   pdftoppm -h

   # Check tesseract
   tesseract --version

   # List available languages
   tesseract --list-langs
   ```

#### Windows

1. **Install pdftoppm**:
   - Download Poppler for Windows from: https://blog.alivate.com.au/poppler-windows/
   - Extract to a folder (e.g., `C:\poppler`)
   - Add `C:\poppler\Library\bin` to your system PATH

2. **Install Tesseract OCR**:
   - Download from: https://github.com/UB-Mannheim/tesseract/wiki
   - Run the installer and follow the setup wizard
   - Make sure to select additional language packs during installation
   - Add Tesseract installation directory to your system PATH

3. **Verify installation**:
   ```cmd
   pdftoppm -h
   tesseract --version
   tesseract --list-langs
   ```

### Go Dependencies

The project uses Go modules for dependency management. The main dependency is:

- `github.com/otiai10/gosseract` - Go wrapper for Tesseract OCR

## Installation

### As a Go Library

Add the parser to your Go project:

```bash
go get github.com/potterbl/statement_parser@latest
```

Then import in your Go code:

```go
import (
    "github.com/potterbl/statement_parser"
)
```

## 📖 Usage Examples

### Basic PDF Parsing

```go
package main

import (
    "log"
    "fmt"
    "github.com/potterbl/statement_parser"
    "github.com/potterbl/statement_parser/pkg/types"
)

func main() {
    // Create a Monobank parser
    parser := statement_parser.NewStatementParser(types.BankNameMono)

    // Parse PDF with Ukrainian and English language support
    transactions, err := parser.ParsePDF(fileHeader, []string{"ukr", "eng"})
    if err != nil {
        log.Fatal(err)
    }

    // Process and display transactions
    for _, transaction := range transactions {
        fmt.Printf("Date: %s, Amount: %.2f %s, Description: %s\n",
            transaction.Date, transaction.Amount, transaction.Currency, transaction.Name)
    }
}
```

## 🔍 API Reference

### Transaction Structure

```go
type Transaction struct {
    Type     string  `json:"type"`     // Transaction type
    Amount   float64 `json:"amount"`   // Transaction amount
    Currency string  `json:"currency"` // Currency code
    Name     string  `json:"name"`     // Transaction description
    Date     string  `json:"date"`     // Transaction date
}
```

### Bank Parser Interface

```go
type BankParser interface {
    ParsePDF(file multipart.FileHeader, languages []string) ([]Transaction, error)
    GetBankName() BankName
}
```

## 🗂 Project Structure

```
statement_parser/
├── LICENSE            # Project License
├── Makefile           # Makefile for easy project management
├── changelog.md       # Project changelog
├── go.mod             # Go module definition
├── go.sum             # Go module checksums
├── main.go            # Main application entrypoint
├── pkg/               # Core library packages
│   ├── consts/        # Constants and configuration
│   ├── types/         # Type definitions
│   └── utils/         # Utility functions
└── readme.md          # Project documentation
```

## How It Works

1. **PDF Upload**: The system accepts PDF bank statements as input
2. **PDF to Image Conversion**: Each PDF page is converted to PNG images using `pdftoppm`
3. **OCR Processing**: Images are processed using Tesseract OCR with specified languages
4. **Text Parsing**: Extracted text is parsed to identify transaction patterns
5. **Data Extraction**: Transaction details (date, amount, description, etc.) are extracted
6. **Structured Output**: Data is returned as structured Transaction objects

## Adding Support for New Banks

To add support for a new bank:

1. **Define Bank Constant**:
   ```go
   // In pkg/types/types.go
   const (
       BankNameMono    BankName = "Mono"
       BankNameNewBank BankName = "NewBank" // Add your bank
   )
   ```

2. **Create Parser Implementation**:
   ```go
   // Create pkg/utils/statement_parser/newbank_parser.go
   type NewBankParser struct{}
   
   func (p *NewBankParser) ParsePDF(file multipart.FileHeader, languages []string) ([]types.Transaction, error) {
       // Implement bank-specific parsing logic
   }
   
   func (p *NewBankParser) GetBankName() types.BankName {
       return types.BankNameNewBank
   }
   ```

3. **Update Parser Factory**:
   ```go
   // In main.go
   func NewStatementParser(bankName types.BankName) types.BankParser {
       switch bankName {
       case types.BankNameMono:
           return &MonoBankParser{}
       case types.BankNameNewBank:
           return &NewBankParser{}
       default:
           return nil
       }
   }
   ```

## Error Handling

The parser includes comprehensive error handling for:
- Invalid PDF files
- OCR processing failures
- Temporary file creation/cleanup issues
- Parsing errors for malformed statements

## Performance Considerations

- **Concurrent Processing**: Multiple PDF pages are processed concurrently for better performance
- **Temporary File Management**: Automatic cleanup of temporary files to prevent disk space issues
- **Memory Optimization**: Efficient handling of large PDF files

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📄 License

Distributed under the MIT License. See `LICENSE` for more information.

## 📧 Contact

Project Maintainer: potterbl
Project Link: https://github.com/potterbl/statement_parser

## 🙏 Acknowledgments

- [Tesseract OCR](https://github.com/tesseract-ocr/tesseract)
- [Poppler](https://poppler.freedesktop.org/)
- Go Community
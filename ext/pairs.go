package ext

import (
	"fmt"
	"strings"
)

type Category string

type Provider string

func (p Provider) Topic(prefix string) string {
	var result string
	val := strings.ToLower(
		strings.Replace(string(p), "CONVERSION_PROVIDER_", "", 1),
	)
	if prefix != "" {
		result = fmt.Sprintf("%s.%s", prefix, val)
	} else {
		result = val
	}
	return result
}

const (
	Writer  Category = "CONVERSION_PAIR_CATEGORY_WRITER"
	Calc    Category = "CONVERSION_PAIR_CATEGORY_CALC"
	Impress Category = "CONVERSION_PAIR_CATEGORY_IMPRESS"
	Draw    Category = "CONVERSION_PAIR_CATEGORY_DRAW"
)

const (
	LibreOffice Provider = "CONVERSION_PROVIDER_LIBRE_OFFICE"
)

var (
	Pairs = NewConversionPairs()
)

type extInfo struct {
	Ext      string
	Category Category
	Pairs    map[string]Provider
}

type ConversionPairs struct {
	Items map[string]extInfo
}

func (p *ConversionPairs) register(from string, to []string, category Category, provider Provider) {
	opts, ok := p.Items[from]
	if !ok {
		opts = extInfo{from, category, make(map[string]Provider)}
	}
	for _, ext := range to {
		opts.Pairs[ext] = provider
	}
	p.Items[from] = opts
}

func (p *ConversionPairs) GetProvider(from, to string) *Provider {
	items, ok := p.Items[from]
	if !ok {
		return nil
	}

	provider, ok := items.Pairs[to]
	if !ok {
		return nil
	}
	return &provider
}

func NewConversionPairs() *ConversionPairs {
	result := ConversionPairs{make(map[string]extInfo)}
	// Writer
	result.register("odt", []string{"doc", "docx", "rtf", "html", "txt", "fodt", "pdf"}, Writer, LibreOffice)
	result.register("doc", []string{"odt", "docx", "rtf", "html", "txt", "fodt", "pdf"}, Writer, LibreOffice)
	result.register("docx", []string{"odt", "doc", "rtf", "html", "txt", "fodt", "pdf"}, Writer, LibreOffice)
	result.register("rtf", []string{"odt", "doc", "docx", "html", "txt", "fodt", "pdf"}, Writer, LibreOffice)
	result.register("txt", []string{"odt", "doc", "docx", "rtf", "html", "fodt", "pdf"}, Writer, LibreOffice)
	result.register("html", []string{"odt", "doc", "docx", "rtf", "txt", "fodt", "pdf"}, Writer, LibreOffice)
	result.register("fodt", []string{"odt", "doc", "docx", "rtf", "html", "txt", "pdf"}, Writer, LibreOffice)

	// Calc
	result.register("ods", []string{"xls", "xlsx", "csv", "html", "pdf"}, Calc, LibreOffice)
	result.register("xls", []string{"ods", "xlsx", "csv", "html", "pdf"}, Calc, LibreOffice)
	result.register("xlsx", []string{"ods", "xls", "csv", "html", "pdf"}, Calc, LibreOffice)
	result.register("csv", []string{"ods", "xls", "xlsx", "html", "pdf"}, Calc, LibreOffice)
	result.register("html", []string{"ods", "xls", "xlsx", "csv", "pdf"}, Calc, LibreOffice)

	// Impress
	result.register("odp", []string{"ppt", "pptx", "pdf"}, Impress, LibreOffice)
	result.register("ppt", []string{"odp", "pptx", "pdf"}, Impress, LibreOffice)
	result.register("pptx", []string{"odp", "ppt", "pdf"}, Impress, LibreOffice)

	// Draw
	result.register("odg", []string{"cdr", "pdf"}, Draw, LibreOffice)
	result.register("cdr", []string{"odg", "pdf"}, Draw, LibreOffice)

	return &result
}

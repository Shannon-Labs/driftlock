package parsers

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// FileType represents the type of file we're parsing
type FileType string

const (
	FileTypeJSON     FileType = "json"
	FileTypeNDJSON   FileType = "ndjson"
	FileTypeCSV      FileType = "csv"
	FileTypeLog      FileType = "log"
	FileTypeCustom   FileType = "custom"
	FileTypeUnknown  FileType = "unknown"
)

// FileInfo represents information about a parsed file
type FileInfo struct {
	Type        FileType     `json:"type"`
	Size        int64        `json:"size"`
	LineCount   int          `json:"line_count"`
	Headers     []string     `json:"headers,omitempty"`
	Encoding    string       `json:"encoding"`
	ProcessedAt time.Time    `json:"processed_at"`
	Metadata    interface{}  `json:"metadata,omitempty"`
}

// ParseEvent represents a single parsed event from any file type
type ParseEvent struct {
	Index    int             `json:"index"`
	Data     json.RawMessage `json:"data"`
	Raw      string          `json:"raw,omitempty"`
	Metadata interface{}     `json:"metadata,omitempty"`
}

// Parser interface for different file types
type Parser interface {
	DetectFileType(reader io.Reader) (FileType, error)
	Parse(reader io.Reader, events chan<- ParseEvent) error
	GetFileInfo() FileInfo
}

// DetectFileType automatically detects the file type based on content and extension
func DetectFileType(filename string, reader io.Reader) (FileType, error) {
	// First try to detect by file extension
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".json":
		return FileTypeJSON, nil
	case ".ndjson", ".jsonl":
		return FileTypeNDJSON, nil
	case ".csv":
		return FileTypeCSV, nil
	case ".log", ".txt":
		return FileTypeLog, nil
	}

	// If extension doesn't give us a clear answer, detect by content
	buf := make([]byte, 1024)
	n, err := reader.Read(buf)
	if err != nil && err != io.EOF {
		return FileTypeUnknown, fmt.Errorf("failed to read file header: %w", err)
	}

	content := string(buf[:n])

	// Reset reader for future parsing
	if seeker, ok := reader.(io.Seeker); ok {
		seeker.Seek(0, io.SeekStart)
	}

	// Detect by content patterns
	if strings.HasPrefix(strings.TrimSpace(content), "[") {
		return FileTypeJSON, nil
	}

	// Check if first line looks like JSON
	firstLine := strings.Split(content, "\n")[0]
	if strings.HasPrefix(strings.TrimSpace(firstLine), "{") {
		return FileTypeNDJSON, nil
	}

	// Check for CSV patterns
	if strings.Contains(firstLine, ",") && !strings.Contains(firstLine, "{") {
		return FileTypeCSV, nil
	}

	// Default to log format
	return FileTypeLog, nil
}

// JSONParser handles JSON array files
type JSONParser struct {
	fileInfo FileInfo
}

func NewJSONParser() *JSONParser {
	return &JSONParser{}
}

func (p *JSONParser) DetectFileType(reader io.Reader) (FileType, error) {
	return FileTypeJSON, nil
}

func (p *JSONParser) Parse(reader io.Reader, events chan<- ParseEvent) error {
	decoder := json.NewDecoder(reader)

	// Read opening bracket
	token, err := decoder.Token()
	if err != nil {
		return fmt.Errorf("failed to read JSON start: %w", err)
	}

	if token != json.Delim('[') {
		return fmt.Errorf("expected JSON array start '['")
	}

	index := 0
	for decoder.More() {
		var rawData json.RawMessage
		if err := decoder.Decode(&rawData); err != nil {
			return fmt.Errorf("failed to decode JSON object at index %d: %w", index, err)
		}

		events <- ParseEvent{
			Index: index,
			Data:  rawData,
		}
		index++
	}

	// Read closing bracket
	token, err = decoder.Token()
	if err != nil {
		return fmt.Errorf("failed to read JSON end: %w", err)
	}

	if token != json.Delim(']') {
		return fmt.Errorf("expected JSON array end ']'")
	}

	p.fileInfo = FileInfo{
		Type:        FileTypeJSON,
		LineCount:   index,
		Encoding:    "utf-8",
		ProcessedAt: time.Now(),
	}

	return nil
}

func (p *JSONParser) GetFileInfo() FileInfo {
	return p.fileInfo
}

// NDJSONParser handles Newline Delimited JSON files
type NDJSONParser struct {
	fileInfo FileInfo
}

func NewNDJSONParser() *NDJSONParser {
	return &NDJSONParser{}
}

func (p *NDJSONParser) DetectFileType(reader io.Reader) (FileType, error) {
	return FileTypeNDJSON, nil
}

func (p *NDJSONParser) Parse(reader io.Reader, events chan<- ParseEvent) error {
	scanner := bufio.NewScanner(reader)
	index := 0

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue // Skip empty lines
		}

		// Validate JSON format
		if !json.Valid([]byte(line)) {
			return fmt.Errorf("invalid JSON at line %d: %s", index+1, line)
		}

		events <- ParseEvent{
			Index: index,
			Data:  json.RawMessage(line),
			Raw:   line,
		}
		index++
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	p.fileInfo = FileInfo{
		Type:        FileTypeNDJSON,
		LineCount:   index,
		Encoding:    "utf-8",
		ProcessedAt: time.Now(),
	}

	return nil
}

func (p *NDJSONParser) GetFileInfo() FileInfo {
	return p.fileInfo
}

// CSVParser handles CSV files
type CSVParser struct {
	fileInfo FileInfo
	hasHeaders bool
}

func NewCSVParser() *CSVParser {
	return &CSVParser{}
}

func (p *CSVParser) DetectFileType(reader io.Reader) (FileType, error) {
	return FileTypeCSV, nil
}

func (p *CSVParser) Parse(reader io.Reader, events chan<- ParseEvent) error {
	csvReader := csv.NewReader(reader)

	// Read header row
	headers, err := csvReader.Read()
	if err != nil {
		return fmt.Errorf("failed to read CSV headers: %w", err)
	}

	p.hasHeaders = true
	p.fileInfo.Headers = headers

	index := 0
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read CSV row %d: %w", index+1, err)
		}

		// Convert CSV row to JSON object
		jsonObj := make(map[string]interface{})
		for i, header := range headers {
			var value interface{} = ""
			if i < len(record) {
				value = record[i]
			}
			jsonObj[header] = value
		}

		jsonData, err := json.Marshal(jsonObj)
		if err != nil {
			return fmt.Errorf("failed to convert CSV row %d to JSON: %w", index+1, err)
		}

		events <- ParseEvent{
			Index: index,
			Data:  jsonData,
			Raw:   strings.Join(record, ","),
		}
		index++
	}

	p.fileInfo = FileInfo{
		Type:        FileTypeCSV,
		LineCount:   index,
		Headers:     headers,
		Encoding:    "utf-8",
		ProcessedAt: time.Now(),
	}

	return nil
}

func (p *CSVParser) GetFileInfo() FileInfo {
	return p.fileInfo
}

// LogParser handles log files with various formats
type LogParser struct {
	fileInfo FileInfo
	pattern  string
}

func NewLogParser() *LogParser {
	return &LogParser{}
}

func (p *LogParser) DetectFileType(reader io.Reader) (FileType, error) {
	return FileTypeLog, nil
}

func (p *LogParser) Parse(reader io.Reader, events chan<- ParseEvent) error {
	scanner := bufio.NewScanner(reader)
	index := 0

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue // Skip empty lines
		}

		// Try to parse common log formats
		jsonObj := p.parseLogLine(line)

		jsonData, err := json.Marshal(jsonObj)
		if err != nil {
			// If JSON marshaling fails, create a simple message object
			fallbackObj := map[string]interface{}{
				"message": line,
				"index":   index,
			}
			jsonData, _ = json.Marshal(fallbackObj)
		}

		events <- ParseEvent{
			Index: index,
			Data:  jsonData,
			Raw:   line,
		}
		index++
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading log file: %w", err)
	}

	p.fileInfo = FileInfo{
		Type:        FileTypeLog,
		LineCount:   index,
		Encoding:    "utf-8",
		ProcessedAt: time.Now(),
		Metadata: map[string]interface{}{
			"pattern": p.pattern,
		},
	}

	return nil
}

func (p *LogParser) parseLogLine(line string) map[string]interface{} {
	result := make(map[string]interface{})

	// Common log format patterns
	patterns := []struct {
		name   string
		parser func(string) map[string]interface{}
	}{
		{
			name: "apache_common",
			parser: func(s string) map[string]interface{} {
				// Simplified Apache common log format parsing
				parts := strings.Fields(s)
				if len(parts) >= 7 {
					return map[string]interface{}{
						"ip":        parts[0],
						"timestamp": parts[3] + " " + parts[4],
						"method":    parts[5],
						"path":      parts[6],
						"protocol":  parts[7],
						"raw":       s,
					}
				}
				return nil
			},
		},
		{
			name: "json_log",
			parser: func(s string) map[string]interface{} {
				var jsonObj map[string]interface{}
				if err := json.Unmarshal([]byte(s), &jsonObj); err == nil {
					jsonObj["raw"] = s
					return jsonObj
				}
				return nil
			},
		},
		{
			name: "timestamp_level_message",
			parser: func(s string) map[string]interface{} {
				// Try to extract timestamp, log level, and message
				if strings.Contains(s, " ") {
					parts := strings.SplitN(s, " ", 3)
					if len(parts) >= 3 {
						return map[string]interface{}{
							"timestamp": parts[0],
							"level":     parts[1],
							"message":   parts[2],
							"raw":       s,
						}
					}
				}
				return nil
			},
		},
	}

	for _, pattern := range patterns {
		if parsed := pattern.parser(line); parsed != nil {
			p.pattern = pattern.name
			return parsed
		}
	}

	// Default parsing - just treat as message
	return map[string]interface{}{
		"message": line,
		"raw":     line,
	}
}

func (p *LogParser) GetFileInfo() FileInfo {
	return p.fileInfo
}

// GetParser returns the appropriate parser for the given file type
func GetParser(fileType FileType) Parser {
	switch fileType {
	case FileTypeJSON:
		return NewJSONParser()
	case FileTypeNDJSON:
		return NewNDJSONParser()
	case FileTypeCSV:
		return NewCSVParser()
	case FileTypeLog:
		return NewLogParser()
	default:
		return NewNDJSONParser() // Default to NDJSON parser
	}
}

// ParseFile is a convenience function that detects file type and parses it
func ParseFile(filename string, reader io.Reader, events chan<- ParseEvent) (FileInfo, error) {
	fileType, err := DetectFileType(filename, reader)
	if err != nil {
		return FileInfo{}, err
	}

	parser := GetParser(fileType)
	if err := parser.Parse(reader, events); err != nil {
		return FileInfo{}, err
	}

	fileInfo := parser.GetFileInfo()
	fileInfo.Type = fileType

	return fileInfo, nil
}

// ValidateFileSize checks if file size is within allowed limits
func ValidateFileSize(size int64, maxSize int64) error {
	if size > maxSize {
		return fmt.Errorf("file size %d bytes exceeds maximum allowed size %d bytes", size, maxSize)
	}
	if size == 0 {
		return fmt.Errorf("file is empty")
	}
	return nil
}

// SanitizeFilename removes potentially dangerous characters from filenames
func SanitizeFilename(filename string) string {
	// Remove path separators and other potentially dangerous characters
	dangerous := []string{"/", "\\", "..", "\x00"}
	sanitized := filename

	for _, char := range dangerous {
		sanitized = strings.ReplaceAll(sanitized, char, "_")
	}

	// Limit filename length
	if len(sanitized) > 255 {
		sanitized = sanitized[:255]
	}

	return sanitized
}

// IsAllowedExtension checks if file extension is allowed
func IsAllowedExtension(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	allowedExts := []string{".json", ".ndjson", ".jsonl", ".csv", ".log", ".txt"}

	for _, allowed := range allowedExts {
		if ext == allowed {
			return true
		}
	}
	return false
}

// CountLines safely counts lines in a file with a maximum limit
func CountLines(reader io.Reader, maxLines int) (int, error) {
	scanner := bufio.NewScanner(reader)
	count := 0

	for scanner.Scan() && count < maxLines {
		count++
	}

	if err := scanner.Err(); err != nil {
		return 0, err
	}

	return count, nil
}
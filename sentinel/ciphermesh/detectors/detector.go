package detectors

import (
	"context"
	"regexp"
	"time"
)

// DetectionResult represents a detected sensitive data item
type DetectionResult struct {
	// ID is a unique identifier for this detection
	ID string `json:"id"`

	// Type is the type of data detected (e.g., "pii", "phi", "pci")
	Type string `json:"type"`

	// Subtype is the specific subtype (e.g., "ssn", "credit_card")
	Subtype string `json:"subtype"`

	// Confidence is the confidence score (0.0 to 1.0)
	Confidence float64 `json:"confidence"`

	// Start is the start position of the detected data
	Start int `json:"start"`

	// End is the end position of the detected data
	End int `json:"end"`

	// Text is the detected text
	Text string `json:"text"`

	// Context is the surrounding context
	Context string `json:"context"`

	// DetectedAt is when the detection was made
	DetectedAt time.Time `json:"detected_at"`
}

// Detector is the interface for all data detectors
type Detector interface {
	// Detect identifies sensitive data in the provided text
	Detect(ctx context.Context, text string) ([]DetectionResult, error)

	// GetType returns the type of data this detector identifies
	GetType() string

	// GetName returns the name of this detector
	GetName() string
}

// RegexDetector implements detection using regular expressions
type RegexDetector struct {
	name          string
	dataType      string
	subtype       string
	pattern       *regexp.Regexp
	confidence    float64
	contextWindow int
}

// NewRegexDetector creates a new regex-based detector
func NewRegexDetector(name, dataType, subtype, pattern string, confidence float64, contextWindow int) (*RegexDetector, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	return &RegexDetector{
		name:          name,
		dataType:      dataType,
		subtype:       subtype,
		pattern:       re,
		confidence:    confidence,
		contextWindow: contextWindow,
	}, nil
}

// Detect identifies sensitive data using regex patterns
func (rd *RegexDetector) Detect(ctx context.Context, text string) ([]DetectionResult, error) {
	var results []DetectionResult

	matches := rd.pattern.FindAllStringSubmatchIndex(text, -1)
	for _, match := range matches {
		if len(match) < 2 {
			continue
		}

		start, end := match[0], match[1]
		detectedText := text[start:end]

		// Extract context
		contextStart := start - rd.contextWindow
		if contextStart < 0 {
			contextStart = 0
		}

		contextEnd := end + rd.contextWindow
		if contextEnd > len(text) {
			contextEnd = len(text)
		}

		contextText := text[contextStart:contextEnd]

		result := DetectionResult{
			ID:         generateID(),
			Type:       rd.dataType,
			Subtype:    rd.subtype,
			Confidence: rd.confidence,
			Start:      start,
			End:        end,
			Text:       detectedText,
			Context:    contextText,
			DetectedAt: time.Now(),
		}

		results = append(results, result)
	}

	return results, nil
}

// GetType returns the type of data this detector identifies
func (rd *RegexDetector) GetType() string {
	return rd.dataType
}

// GetName returns the name of this detector
func (rd *RegexDetector) GetName() string {
	return rd.name
}

// generateID creates a unique identifier for detection results
func generateID() string {
	// In a real implementation, this would generate a proper UUID
	// For now, we'll use a simplified approach
	return time.Now().Format("20060102150405.000")
}

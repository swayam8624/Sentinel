package detector

import (
	"context"
	"fmt"
	"math"
)

// ViolationDetector detects security violations in prompts and responses
type ViolationDetector struct {
	signatureStore SignatureStore
	ruleEngine     RuleEngine
	embeddingModel EmbeddingModel
	thresholds     DetectionThresholds
}

// DetectionThresholds defines thresholds for different types of detections
type DetectionThresholds struct {
	ViolationSimilarity float64 `json:"violation_similarity"`
	ReflectConfidence   float64 `json:"reflect_confidence"`
}

// DetectionResult represents the result of a violation detection
type DetectionResult struct {
	Score          float64            `json:"score"`
	Confidence     float64            `json:"confidence"`
	ViolationType  string             `json:"violation_type"`
	Details        map[string]float64 `json:"details"`
	Recommendation string             `json:"recommendation"`
}

// SignatureStore interface for storing and retrieving attack signatures
type SignatureStore interface {
	// Search finds signatures similar to the input text
	Search(ctx context.Context, text string) ([]SignatureMatch, error)

	// AddSignature adds a new signature to the store
	AddSignature(ctx context.Context, signature Signature) error
}

// Signature represents an attack signature
type Signature struct {
	ID          string    `json:"id"`
	Embedding   []float64 `json:"embedding"`
	Label       string    `json:"label"`
	Description string    `json:"description"`
	Active      bool      `json:"active"`
}

// SignatureMatch represents a match against a signature
type SignatureMatch struct {
	SignatureID string  `json:"signature_id"`
	Score       float64 `json:"score"`
	Label       string  `json:"label"`
}

// RuleEngine interface for rule-based detection
type RuleEngine interface {
	// Evaluate evaluates rules against input text
	Evaluate(ctx context.Context, text string) ([]RuleMatch, error)
}

// RuleMatch represents a match against a rule
type RuleMatch struct {
	RuleID      string  `json:"rule_id"`
	Score       float64 `json:"score"`
	Description string  `json:"description"`
}

// EmbeddingModel interface for generating text embeddings
type EmbeddingModel interface {
	// Generate generates an embedding for the input text
	Generate(ctx context.Context, text string) ([]float64, error)
}

// NewViolationDetector creates a new violation detector
func NewViolationDetector(
	signatureStore SignatureStore,
	ruleEngine RuleEngine,
	embeddingModel EmbeddingModel,
	thresholds DetectionThresholds) *ViolationDetector {

	return &ViolationDetector{
		signatureStore: signatureStore,
		ruleEngine:     ruleEngine,
		embeddingModel: embeddingModel,
		thresholds:     thresholds,
	}
}

// Detect detects violations in the provided text
func (vd *ViolationDetector) Detect(ctx context.Context, text string) (*DetectionResult, error) {
	// Initialize result
	result := &DetectionResult{
		Details: make(map[string]float64),
	}

	// Get embedding for the text
	embedding, err := vd.embeddingModel.Generate(ctx, text)
	if err != nil {
		return nil, fmt.Errorf("failed to generate embedding: %w", err)
	}
	_ = embedding // Explicitly mark as used

	// Check against signature store
	signatureMatches, err := vd.signatureStore.Search(ctx, text)
	if err != nil {
		return nil, fmt.Errorf("failed to search signatures: %w", err)
	}

	// Evaluate rules
	ruleMatches, err := vd.ruleEngine.Evaluate(ctx, text)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate rules: %w", err)
	}

	// Calculate signature-based score
	signatureScore := vd.calculateSignatureScore(signatureMatches)
	result.Details["signature_score"] = signatureScore

	// Calculate rule-based score
	ruleScore := vd.calculateRuleScore(ruleMatches)
	result.Details["rule_score"] = ruleScore

	// Calculate overall score (weighted average)
	overallScore := (signatureScore*0.6 + ruleScore*0.4)
	result.Score = overallScore

	// Determine confidence based on the strength of matches
	confidence := vd.calculateConfidence(signatureMatches, ruleMatches)
	result.Confidence = confidence

	// Determine violation type and recommendation
	violationType, recommendation := vd.determineViolationType(signatureMatches, ruleMatches)
	result.ViolationType = violationType
	result.Recommendation = recommendation

	return result, nil
}

// calculateSignatureScore calculates a score based on signature matches
func (vd *ViolationDetector) calculateSignatureScore(matches []SignatureMatch) float64 {
	if len(matches) == 0 {
		return 0.0
	}

	// Find the highest scoring match
	maxScore := 0.0
	for _, match := range matches {
		if match.Score > maxScore {
			maxScore = match.Score
		}
	}

	return maxScore
}

// calculateRuleScore calculates a score based on rule matches
func (vd *ViolationDetector) calculateRuleScore(matches []RuleMatch) float64 {
	if len(matches) == 0 {
		return 0.0
	}

	// Calculate weighted average of rule scores
	totalScore := 0.0
	totalWeight := 0.0

	for _, match := range matches {
		// Weight rules based on their score
		weight := match.Score
		totalScore += match.Score * weight
		totalWeight += weight
	}

	if totalWeight == 0 {
		return 0.0
	}

	return totalScore / totalWeight
}

// calculateConfidence calculates the confidence level of the detection
func (vd *ViolationDetector) calculateConfidence(signatureMatches []SignatureMatch, ruleMatches []RuleMatch) float64 {
	// Start with base confidence
	confidence := 0.5

	// Increase confidence based on strong signature matches
	for _, match := range signatureMatches {
		if match.Score > vd.thresholds.ViolationSimilarity {
			confidence += match.Score * 0.3
		}
	}

	// Increase confidence based on rule matches
	for _, match := range ruleMatches {
		if match.Score > 0.7 {
			confidence += match.Score * 0.2
		}
	}

	// Cap confidence at 1.0
	if confidence > 1.0 {
		confidence = 1.0
	}

	return confidence
}

// determineViolationType determines the type of violation and recommendation
func (vd *ViolationDetector) determineViolationType(signatureMatches []SignatureMatch, ruleMatches []RuleMatch) (string, string) {
	// Check for strong signature matches first
	for _, match := range signatureMatches {
		if match.Score > vd.thresholds.ViolationSimilarity {
			switch match.Label {
			case "jailbreak":
				return "jailbreak_attempt", "Block or reframe the prompt"
			case "injection":
				return "injection_attack", "Block and sanitize the input"
			case "exfiltration":
				return "data_exfiltration", "Encrypt output and block tools"
			default:
				return "known_attack", "Block the request"
			}
		}
	}

	// Check for rule violations
	for _, match := range ruleMatches {
		if match.Score > 0.8 {
			return "policy_violation", "Reframe or block based on policy"
		}
	}

	// If scores are moderate, suggest reflection
	if len(signatureMatches) > 0 || len(ruleMatches) > 0 {
		return "suspicious_content", "Reflect and potentially reframe"
	}

	// No violations detected
	return "none", "Allow the request"
}

// IsViolation determines if the detection result indicates a violation
func (vd *ViolationDetector) IsViolation(result *DetectionResult) bool {
	return result.Score > vd.thresholds.ViolationSimilarity
}

// ShouldReflect determines if the content should be reflected upon
func (vd *ViolationDetector) ShouldReflect(result *DetectionResult) bool {
	return result.Score > vd.thresholds.ReflectConfidence
}

// cosineSimilarity calculates the cosine similarity between two vectors
func cosineSimilarity(a, b []float64) float64 {
	if len(a) != len(b) {
		return 0.0
	}

	dotProduct := 0.0
	magnitudeA := 0.0
	magnitudeB := 0.0

	for i := 0; i < len(a); i++ {
		dotProduct += a[i] * b[i]
		magnitudeA += a[i] * a[i]
		magnitudeB += b[i] * b[i]
	}

	if magnitudeA == 0 || magnitudeB == 0 {
		return 0.0
	}

	return dotProduct / (math.Sqrt(magnitudeA) * math.Sqrt(magnitudeB))
}

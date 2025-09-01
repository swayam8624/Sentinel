package core
package core

import (
	"context"
	"fmt"
	"io"
	"time"
)

// CutoverManager handles mid-stream cutovers when violations are detected
type CutoverManager struct {
	detector        ViolationDetector
	maxCutoverTime  time.Duration
	bufferSize      int
}

// ViolationDetector interface for detecting violations
type ViolationDetector interface {
	// DetectStream detects violations in a stream of data
	DetectStream(ctx context.Context, reader io.Reader) (*DetectionResult, error)
}

// DetectionResult represents the result of a violation detection
type DetectionResult struct {
	Score         float64            `json:"score"`
	Confidence    float64            `json:"confidence"`
	ViolationType string             `json:"violation_type"`
	Details       map[string]float64 `json:"details"`
}

// NewCutoverManager creates a new cutover manager
func NewCutoverManager(detector ViolationDetector, maxCutoverTime time.Duration, bufferSize int) *CutoverManager {
	return &CutoverManager{
		detector:        detector,
		maxCutoverTime:  maxCutoverTime,
		bufferSize:      bufferSize,
	}
}

// CutoverResult represents the result of a cutover operation
type CutoverResult struct {
	CutoverPerformed bool          `json:"cutover_performed"`
	ViolationDetected bool        `json:"violation_detected"`
	DetectionResult   *DetectionResult `json:"detection_result"`
	CutoverTime       time.Duration    `json:"cutover_time"`
	BytesProcessed    int64            `json:"bytes_processed"`
}

// StreamWithCutover processes a stream with mid-stream cutover capability
func (cm *CutoverManager) StreamWithCutover(
	ctx context.Context,
	inputReader io.Reader,
	outputWriter io.Writer,
	violationHandler func(*DetectionResult) error) (*CutoverResult, error) {
	
	result := &CutoverResult{
		CutoverPerformed:  false,
		ViolationDetected: false,
	}
	
	// Create a context with timeout for cutover
	cutoverCtx, cancel := context.WithTimeout(ctx, cm.maxCutoverTime)
	defer cancel()
	
	// Buffer for streaming data
	buffer := make([]byte, cm.bufferSize)
	
	// Start time for performance tracking
	startTime := time.Now()
	
	// Process stream in chunks
	for {
		// Check if context is cancelled
		select {
		case <-cutoverCtx.Done():
			return result, cutoverCtx.Err()
		default:
		}
		
		// Read chunk
		n, err := inputReader.Read(buffer)
		if err != nil && err != io.EOF {
			return result, fmt.Errorf("failed to read from stream: %w", err)
		}
		
		// Update bytes processed
		result.BytesProcessed += int64(n)
		
		// If we have data, check for violations
		if n > 0 {
			// Create a reader for the chunk
			chunkReader := &ChunkReader{data: buffer[:n]}
			
			// Check for violations in the chunk
			detectionResult, detectErr := cm.detector.DetectStream(cutoverCtx, chunkReader)
			if detectErr != nil {
				// Log error but continue processing
				continue
			}
			
			// If violation detected, handle it
			if detectionResult != nil && detectionResult.Score > 0.7 { // Threshold for action
				result.ViolationDetected = true
				result.DetectionResult = detectionResult
				
				// Call violation handler
				handlerErr := violationHandler(detectionResult)
				if handlerErr != nil {
					return result, fmt.Errorf("violation handler failed: %w", handlerErr)
				}
				
				// Perform cutover
				result.CutoverPerformed = true
				result.CutoverTime = time.Since(startTime)
				
				// Stop processing
				break
			}
			
			// Write chunk to output
			_, writeErr := outputWriter.Write(buffer[:n])
			if writeErr != nil {
				return result, fmt.Errorf("failed to write to output: %w", writeErr)
			}
		}
		
		// If we reached EOF, we're done
		if err == io.EOF {
			break
		}
	}
	
	// If no cutover was performed, set the time
	if !result.CutoverPerformed {
		result.CutoverTime = time.Since(startTime)
	}
	
	return result, nil
}

// ChunkReader is a simple reader for a chunk of data
type ChunkReader struct {
	data []byte
	pos  int
}

// Read implements the io.Reader interface
func (cr *ChunkReader) Read(p []byte) (n int, err error) {
	if cr.pos >= len(cr.data) {
		return 0, io.EOF
	}
	
	n = copy(p, cr.data[cr.pos:])
	cr.pos += n
	
	if cr.pos >= len(cr.data) {
		err = io.EOF
	}
	
	return n, err
}

// IsCutoverPerformed determines if a cutover was performed
func (cm *CutoverManager) IsCutoverPerformed(result *CutoverResult) bool {
	return result.CutoverPerformed
}

// GetCutoverTime returns the time it took to perform the cutover
func (cm *CutoverManager) GetCutoverTime(result *CutoverResult) time.Duration {
	return result.CutoverTime
}

// GetBytesProcessed returns the number of bytes processed
func (cm *CutoverManager) GetBytesProcessed(result *CutoverResult) int64 {
	return result.BytesProcessed
}

// IsViolationDetected determines if a violation was detected
func (cm *CutoverManager) IsViolationDetected(result *CutoverResult) bool {
	return result.ViolationDetected
}

// GetDetectionResult returns the detection result
func (cm *CutoverManager) GetDetectionResult(result *CutoverResult) *DetectionResult {
	return result.DetectionResult
}

// WithinTimeLimit determines if the cutover was performed within the time limit
func (cm *CutoverManager) WithinTimeLimit(result *CutoverResult) bool {
	return result.CutoverTime <= cm.maxCutoverTime
}
package streaming

import (
	"context"
	"fmt"
	"io"
)

// StreamProcessor handles streaming redaction of data
type StreamProcessor struct {
	buffer        []byte
	bufferSize    int
	chunkSize     int
	processorFunc func([]byte) ([]byte, error)
}

// NewStreamProcessor creates a new stream processor
func NewStreamProcessor(chunkSize int, processorFunc func([]byte) ([]byte, error)) *StreamProcessor {
	return &StreamProcessor{
		buffer:        make([]byte, 0, chunkSize*2),
		bufferSize:    chunkSize * 2,
		chunkSize:     chunkSize,
		processorFunc: processorFunc,
	}
}

// Process processes a stream of data
func (sp *StreamProcessor) Process(ctx context.Context, reader io.Reader, writer io.Writer) error {
	// Create a buffer for reading chunks
	chunk := make([]byte, sp.chunkSize)
	
	for {
		// Check if context is cancelled
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		
		// Read a chunk
		n, err := reader.Read(chunk)
		if err != nil && err != io.EOF {
			return fmt.Errorf("failed to read from stream: %w", err)
		}
		
		// If we have data, process it
		if n > 0 {
			// Append to buffer
			sp.buffer = append(sp.buffer, chunk[:n]...)
			
			// Process complete chunks
			for len(sp.buffer) >= sp.chunkSize {
				// Take a chunk to process
				processChunk := sp.buffer[:sp.chunkSize]
				sp.buffer = sp.buffer[sp.chunkSize:]
				
				// Process the chunk
				processedChunk, processErr := sp.processorFunc(processChunk)
				if processErr != nil {
					return fmt.Errorf("failed to process chunk: %w", processErr)
				}
				
				// Write processed chunk
				_, writeErr := writer.Write(processedChunk)
				if writeErr != nil {
					return fmt.Errorf("failed to write to stream: %w", writeErr)
				}
			}
		}
		
		// If we reached EOF, process any remaining data
		if err == io.EOF {
			// Process any remaining data in buffer
			if len(sp.buffer) > 0 {
				processedChunk, processErr := sp.processorFunc(sp.buffer)
				if processErr != nil {
					return fmt.Errorf("failed to process final chunk: %w", processErr)
				}
				
				_, writeErr := writer.Write(processedChunk)
				if writeErr != nil {
					return fmt.Errorf("failed to write final chunk: %w", writeErr)
				}
			}
			
			// Done
			return nil
		}
	}
}

// StreamingRedactor handles streaming redaction with context awareness
type StreamingRedactor struct {
	detectorFunc func([]byte) ([]Detection, error)
	redactorFunc func([]byte, []Detection) ([]byte, error)
	chunkSize    int
}

// Detection represents a detected sensitive item in streaming data
type Detection struct {
	Start     int    `json:"start"`
	End       int    `json:"end"`
	Type      string `json:"type"`
	Subtype   string `json:"subtype"`
	Confidence float64 `json:"confidence"`
}

// NewStreamingRedactor creates a new streaming redactor
func NewStreamingRedactor(
	detectorFunc func([]byte) ([]Detection, error),
	redactorFunc func([]byte, []Detection) ([]byte, error),
	chunkSize int) *StreamingRedactor {
	
	return &StreamingRedactor{
		detectorFunc: detectorFunc,
		redactorFunc: redactorFunc,
		chunkSize:    chunkSize,
	}
}

// RedactStream redacts sensitive data from a streaming source
func (sr *StreamingRedactor) RedactStream(ctx context.Context, reader io.Reader, writer io.Writer) error {
	// For streaming redaction, we need to be careful about boundary conditions
	// where sensitive data might span chunk boundaries
	
	// Create a sliding window buffer to handle boundary conditions
	window := make([]byte, 0, sr.chunkSize*3)
	overlap := sr.chunkSize // Amount of overlap between chunks
	
	for {
		// Check if context is cancelled
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		
		// Read a chunk
		chunk := make([]byte, sr.chunkSize)
		n, err := reader.Read(chunk)
		if err != nil && err != io.EOF {
			return fmt.Errorf("failed to read from stream: %w", err)
		}
		
		// If we have data, add it to the window
		if n > 0 {
			window = append(window, chunk[:n]...)
		}
		
		// Process the window if it's large enough or we've reached EOF
		if len(window) >= sr.chunkSize || (err == io.EOF && len(window) > 0) {
			// Determine how much to process
			processSize := len(window)
			if err != io.EOF && len(window) > sr.chunkSize {
				// Leave overlap for next iteration
				processSize = sr.chunkSize
			}
			
			// Extract the portion to process
			toProcess := window[:processSize]
			
			// Detect sensitive data
			detections, detectErr := sr.detectorFunc(toProcess)
			if detectErr != nil {
				return fmt.Errorf("failed to detect sensitive data: %w", detectErr)
			}
			
			// Redact sensitive data
			redacted, redactErr := sr.redactorFunc(toProcess, detections)
			if redactErr != nil {
				return fmt.Errorf("failed to redact sensitive data: %w", redactErr)
			}
			
			// Write redacted data
			_, writeErr := writer.Write(redacted)
			if writeErr != nil {
				return fmt.Errorf("failed to write redacted data: %w", writeErr)
			}
			
			// Update window - keep overlap if not at EOF
			if err != io.EOF && len(window) > processSize {
				// Keep overlap portion
				window = append(window[:0], window[processSize-overlap:]...)
			} else {
				// Clear window
				window = window[:0]
			}
		}
		
		// If we reached EOF, we're done
		if err == io.EOF {
			return nil
		}
	}
}

// StreamingProcessorFunc is a function that processes a chunk of data
type StreamingProcessorFunc func([]byte) ([]byte, error)

// StreamingDetectorFunc is a function that detects sensitive data in a chunk
type StreamingDetectorFunc func([]byte) ([]Detection, error)

// StreamingRedactorFunc is a function that redacts detections from a chunk
type StreamingRedactorFunc func([]byte, []Detection) ([]byte, error)package streaming

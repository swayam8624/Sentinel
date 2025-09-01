package detectors

import (
	"context"
	"fmt"
	"sort"
	"sync"
)

// DetectorManager coordinates multiple detectors
type DetectorManager struct {
	detectors []Detector
	mutex     sync.RWMutex
}

// NewDetectorManager creates a new detector manager
func NewDetectorManager() *DetectorManager {
	return &DetectorManager{
		detectors: make([]Detector, 0),
	}
}

// AddDetector adds a detector to the manager
func (dm *DetectorManager) AddDetector(detector Detector) {
	dm.mutex.Lock()
	defer dm.mutex.Unlock()

	dm.detectors = append(dm.detectors, detector)
}

// RemoveDetector removes a detector from the manager by name
func (dm *DetectorManager) RemoveDetector(name string) error {
	dm.mutex.Lock()
	defer dm.mutex.Unlock()

	for i, detector := range dm.detectors {
		if detector.GetName() == name {
			// Remove the detector by slicing
			dm.detectors = append(dm.detectors[:i], dm.detectors[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("detector %s not found", name)
}

// GetDetectors returns all detectors
func (dm *DetectorManager) GetDetectors() []Detector {
	dm.mutex.RLock()
	defer dm.mutex.RUnlock()

	// Create a copy to avoid external modification
	detectors := make([]Detector, len(dm.detectors))
	copy(detectors, dm.detectors)

	return detectors
}

// Detect runs all detectors on the provided text
func (dm *DetectorManager) Detect(ctx context.Context, text string) ([]DetectionResult, error) {
	dm.mutex.RLock()
	defer dm.mutex.RUnlock()

	var allResults []DetectionResult

	// Run each detector concurrently
	type detectorResult struct {
		results []DetectionResult
		err     error
	}

	resultChan := make(chan detectorResult, len(dm.detectors))

	// Launch goroutines for each detector
	for _, detector := range dm.detectors {
		go func(d Detector) {
			results, err := d.Detect(ctx, text)
			resultChan <- detectorResult{results: results, err: err}
		}(detector)
	}

	// Collect results
	for i := 0; i < len(dm.detectors); i++ {
		result := <-resultChan
		if result.err != nil {
			// Log error but continue with other detectors
			// In a real implementation, we'd use a proper logger
			continue
		}

		allResults = append(allResults, result.results...)
	}

	// Sort results by position for consistent ordering
	sortDetectionResults(allResults)

	return allResults, nil
}

// sortDetectionResults sorts detection results by their start position
func sortDetectionResults(results []DetectionResult) {
	// Use standard library sort
	sort.Slice(results, func(i, j int) bool {
		return results[i].Start < results[j].Start
	})
}

// GetDetectorsByType returns detectors filtered by type
func (dm *DetectorManager) GetDetectorsByType(dataType string) []Detector {
	dm.mutex.RLock()
	defer dm.mutex.RUnlock()

	var filtered []Detector
	for _, detector := range dm.detectors {
		if detector.GetType() == dataType {
			filtered = append(filtered, detector)
		}
	}

	return filtered
}

// ClearDetectors removes all detectors
func (dm *DetectorManager) ClearDetectors() {
	dm.mutex.Lock()
	defer dm.mutex.Unlock()

	dm.detectors = make([]Detector, 0)
}

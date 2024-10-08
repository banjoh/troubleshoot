package preflight

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	analyze "github.com/replicatedhq/troubleshoot/pkg/analyze"
	troubleshootv1beta2 "github.com/replicatedhq/troubleshoot/pkg/apis/troubleshoot/v1beta2"
	"k8s.io/klog/v2"
)

// Analyze runs the analyze phase of preflight checks
func (c ClusterCollectResult) Analyze() []*analyze.AnalyzeResult {
	return doAnalyze(c.Context, c.AllCollectedData, c.Spec.Spec.Analyzers, nil, "")
}

// Analyze runs the analyze phase of host preflight checks
func (c HostCollectResult) Analyze() []*analyze.AnalyzeResult {
	return doAnalyze(c.Context, c.AllCollectedData, nil, c.Spec.Spec.Analyzers, "")
}

// Analyze runs the analyze phase of host preflight checks.
//
// Runs the analysis for each node and aggregates the results.
func (c RemoteCollectResult) Analyze() []*analyze.AnalyzeResult {
	var results []*analyze.AnalyzeResult
	for nodeName, nodeResult := range c.AllCollectedData {
		var strResult = make(map[string]string)
		if err := json.Unmarshal(nodeResult, &strResult); err != nil {
			analyzeResult := &analyze.AnalyzeResult{
				IsFail:  true,
				Title:   "Remote Result Parser Failed",
				Message: err.Error(),
			}
			results = append(results, analyzeResult)
			continue
		}

		var byteResult = make(map[string][]byte)
		for k, v := range strResult {
			byteResult[k] = []byte(v)

		}
		results = append(results, doAnalyze(c.Context, byteResult, nil, c.Spec.Spec.Analyzers, nodeName)...)
	}
	return results
}

func doAnalyze(
	ctx context.Context,
	allCollectedData map[string][]byte,
	analyzers []*troubleshootv1beta2.Analyze,
	hostAnalyzers []*troubleshootv1beta2.HostAnalyze,
	nodeName string,
) []*analyze.AnalyzeResult {
	getCollectedFileContents := func(fileName string) ([]byte, error) {
		contents, ok := allCollectedData[fileName]
		if !ok {
			return nil, fmt.Errorf("file %s was not collected", fileName)
		}

		return contents, nil
	}
	getChildCollectedFileContents := func(prefix string, excludeFiles []string) (map[string][]byte, error) {
		matching := make(map[string][]byte)
		for k, v := range allCollectedData {
			if strings.HasPrefix(k, prefix) {
				matching[k] = v
			}
		}

		for k, v := range allCollectedData {
			if ok, _ := filepath.Match(prefix, k); ok {
				matching[k] = v
			}
		}

		if len(excludeFiles) > 0 {
			for k := range matching {
				for _, ex := range excludeFiles {
					if ok, _ := filepath.Match(ex, k); ok {
						delete(matching, k)
					}
				}
			}
		}

		if len(matching) == 0 {
			return nil, fmt.Errorf("File not found: %s", prefix)
		}
		return matching, nil
	}

	analyzeResults := []*analyze.AnalyzeResult{}
	for _, analyzer := range analyzers {
		analyzeResult, err := analyze.Analyze(ctx, analyzer, getCollectedFileContents, getChildCollectedFileContents)
		if err != nil {
			strict, strictErr := HasStrictAnalyzer(analyzer)
			if strictErr != nil {
				klog.Errorf("failed to determine if analyzer %v is strict: %s", analyzer, strictErr)
			}

			title := "Analyzer Failed"
			analyzerInst := analyze.GetAnalyzer(analyzer)
			if analyzerInst != nil {
				title = analyzerInst.Title()
			}

			analyzeResult = []*analyze.AnalyzeResult{
				{
					Strict:  strict,
					IsFail:  true,
					Title:   title,
					Message: err.Error(),
				},
			}
		}

		if analyzeResult != nil {
			analyzeResults = append(analyzeResults, analyzeResult...)
		}
	}

	for _, hostAnalyzer := range hostAnalyzers {
		analyzeResult := analyze.HostAnalyze(ctx, hostAnalyzer, getCollectedFileContents, getChildCollectedFileContents)
		analyzeResults = append(analyzeResults, analyzeResult...)
	}

	// Add the nodename to the result title if provided.
	if nodeName != "" {
		for _, result := range analyzeResults {
			result.Title = fmt.Sprintf("%s (%s)", result.Title, nodeName)
		}
	}
	return analyzeResults
}

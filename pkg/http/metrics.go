package http

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/mwantia/metrics-merger/pkg/common"
)

func HandleMetrics(cache *common.MetricsCache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metrics := cache.GetAllMetrics()
		fmt.Fprint(w, metrics)
	}
}

func FetchEndpointBody(endpoint common.EndpointConfig, label string) (string, error) {
	resp, err := http.Get(endpoint.Address)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return AddLabelToMetrics(string(body), label, endpoint.Name), nil
}

func AddLabelToMetrics(metrics, name, value string) string {
	label := fmt.Sprintf(`%s="%s"`, name, value)

	// Regular expression to identify lines with and without existing labels
	re := regexp.MustCompile(`(^[a-zA-Z_:][a-zA-Z0-9_:]*)(\{[^}]*\})?(.*)`)

	var result []string
	for _, line := range strings.Split(metrics, "\n") {
		if line == "" || strings.HasPrefix(line, "#") {
			// Preserve empty lines and comment lines as-is
			result = append(result, line)
			continue
		}

		matches := re.FindStringSubmatch(line)
		if len(matches) == 4 {
			metricName := matches[1]
			existingLabels := matches[2]
			metricValue := matches[3]

			if existingLabels == "" {
				// No existing labels, add new label set
				result = append(result, fmt.Sprintf(`%s{%s}%s`, metricName, label, metricValue))
			} else {
				// Existing labels found, insert new label into existing label set
				updatedLabels := strings.TrimSuffix(existingLabels, "}") + "," + label + "}"
				result = append(result, fmt.Sprintf(`%s%s%s`, metricName, updatedLabels, metricValue))
			}
		} else {
			// If the line doesn't match the metric format, add it unchanged
			result = append(result, line)
		}
	}

	return strings.Join(result, "\n")
}

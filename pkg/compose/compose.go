package compose

import (
	"encoding/json"
	"sort"
	"strings"

	"github.com/sirupsen/logrus"
)

// Docker Compose labels.
const (
	// ComposeDependsOnLabel lists container names this container depends on from Docker Compose, comma-separated.
	ComposeDependsOnLabel = "com.docker.compose.depends_on"
	// ComposeProjectLabel specifies the project name of the container in Docker Compose.
	ComposeProjectLabel = "com.docker.compose.project"
	// ComposeServiceLabel specifies the service name of the container in Docker Compose.
	ComposeServiceLabel = "com.docker.compose.service"
	// ComposeContainerNumber specifies the container number of the container in Docker Compose.
	ComposeContainerNumber = "com.docker.compose.container-number"
)

// ServiceDependency represents a service dependency from Docker Compose depends_on.
type ServiceDependency struct {
	// ServiceName is the name of the dependent service.
	ServiceName string
	// RestartExplicitlyDisabled indicates whether restart is explicitly set to false.
	// When true, the dependent container should NOT be restarted when its dependency restarts.
	RestartExplicitlyDisabled bool
}

// GetServiceNames returns the list of service names from the given dependencies.
// This is a convenience function for backward compatibility.
func GetServiceNames(dependencies []ServiceDependency) []string {
	result := make([]string, len(dependencies))
	for i, dep := range dependencies {
		result[i] = dep.ServiceName
	}

	return result
}

// ParseDependsOnLabel parses the Docker Compose depends_on label value.
//
// It handles both JSON format (Docker Compose v2+) and comma-separated string format.
// Returns a slice of ServiceDependency containing service names and restart configuration.
//
// Supported formats:
//
//	Short form: depends_on: [db, redis]
//	Long form with condition: depends_on: { db: { condition: service_started } }
//	Long form with restart: depends_on: { db: { condition: service_started, restart: true } }
//
// Parameters:
//   - labelValue: The raw label value from com.docker.compose.depends_on.
//
// Returns:
//   - []ServiceDependency: List of service dependencies with their restart configuration.
func ParseDependsOnLabel(labelValue string) []ServiceDependency {
	if labelValue == "" {
		return nil
	}

	clog := logrus.WithField("label_value", labelValue)
	clog.Debug("Parsing compose depends-on label")

	// Try to parse as JSON first (Docker Compose v2+ format)
	if strings.HasPrefix(strings.TrimSpace(labelValue), "{") {
		var dependsOn map[string]any

		err := json.Unmarshal([]byte(labelValue), &dependsOn)
		if err != nil {
			clog.WithError(err).Debug("Failed to parse as JSON, falling back to string parsing")
		} else {
			dependencies := parseLongForm(dependsOn)
			// Sort for consistent ordering
			sort.Slice(dependencies, func(i, j int) bool {
				return dependencies[i].ServiceName < dependencies[j].ServiceName
			})
			clog.WithField("parsed_dependencies", dependencies).
				Debug("Parsed JSON format compose depends-on label")

			return dependencies
		}
	}

	// Fall back to string parsing (legacy format)
	deps := strings.Split(labelValue, ",")
	dependencies := make([]ServiceDependency, 0, len(deps))

	// Parse comma-separated list of service:condition:required
	for _, dep := range deps {
		dep = strings.TrimSpace(dep)
		if dep == "" {
			continue
		}

		clog.WithField("parsing_dep", dep).Debug("Parsing individual dependency")
		// Parse colon-separated format: service:condition:required
		parts := strings.Split(dep, ":")

		serviceName := strings.TrimSpace(parts[0])
		if serviceName != "" {
			dependencies = append(dependencies, ServiceDependency{
				ServiceName:               serviceName,
				RestartExplicitlyDisabled: false,
			})
		}
	}

	clog.WithField("parsed_dependencies", dependencies).
		Debug("Completed parsing string format compose depends-on label")

	return dependencies
}

// parseLongForm parses the long-form JSON format (Docker Compose v2+).
// This format supports condition and restart properties.
func parseLongForm(dependsOn map[string]any) []ServiceDependency {
	result := make([]ServiceDependency, 0, len(dependsOn))

	for serviceName, config := range dependsOn {
		if serviceName == "" {
			continue
		}

		// Default: restart is enabled (not explicitly disabled)
		restartDisabled := false

		// Try to parse the config as a map to extract restart property
		if configMap, ok := config.(map[string]any); ok {
			if restartVal, exists := configMap["restart"]; exists {
				// If restart is explicitly set to false, mark it as disabled
				if restartBool, isBool := restartVal.(bool); isBool {
					restartDisabled = !restartBool // restart: false means disabled
				}
			}
		}

		result = append(result, ServiceDependency{
			ServiceName:               serviceName,
			RestartExplicitlyDisabled: restartDisabled,
		})
	}

	return result
}

// GetProjectName extracts the project name from Docker Compose labels.
//
// If the com.docker.compose.project label is present, returns its value.
// Otherwise, returns an empty string.
//
// Parameters:
//   - labels: Map of container labels.
//
// Returns:
//   - string: Project name if present, empty string otherwise.
func GetProjectName(labels map[string]string) string {
	if labels == nil {
		return ""
	}

	projectName, ok := labels[ComposeProjectLabel]
	if !ok {
		return ""
	}

	logrus.WithFields(logrus.Fields{
		"label": ComposeProjectLabel,
		"value": projectName,
	}).Debug("Retrieved compose project name")

	return projectName
}

// GetServiceName extracts the service name from Docker Compose labels.
//
// If the com.docker.compose.service label is present, returns its value.
// Otherwise, returns an empty string.
//
// Parameters:
//   - labels: Map of container labels.
//
// Returns:
//   - string: Service name if present, empty string otherwise.
func GetServiceName(labels map[string]string) string {
	if labels == nil {
		return ""
	}

	serviceName, ok := labels[ComposeServiceLabel]
	if !ok {
		return ""
	}

	logrus.WithFields(logrus.Fields{
		"label": ComposeServiceLabel,
		"value": serviceName,
	}).Debug("Retrieved compose service name")

	return serviceName
}

// GetContainerNumber extracts the container number from the Docker Compose labels.
//
// If the ComposeContainerNumber label is present, returns its value.
// Otherwise, returns an empty string.
//
// Parameters:
//   - labels: Map of container labels.
//
// Returns:
//   - string: Container replica number if present, empty string otherwise.
func GetContainerNumber(labels map[string]string) string {
	if labels == nil {
		return ""
	}

	containerNumber, ok := labels[ComposeContainerNumber]
	if !ok {
		return ""
	}

	logrus.WithFields(logrus.Fields{
		"label": ComposeContainerNumber,
		"value": containerNumber,
	}).Debug("Retrieved container replica number")

	return containerNumber
}

// Package compose provides functionality for handling Docker Compose-specific logic,
// including parsing depends_on labels and extracting service names for dependency management.
//
// Key components:
//   - ServiceDependency: Represents a service dependency with name and restart configuration.
//   - ParseDependsOnLabel: Parses the Docker Compose depends_on label value and returns
//     ServiceDependency objects containing service name and restart configuration.
//   - GetServiceNames: Extracts service names from ServiceDependency slice for backward compatibility.
//   - GetServiceName: Extracts the service name from Docker Compose labels
//     using the com.docker.compose.service label.
//
// Usage example:
//
//	labels := map[string]string{
//		"com.docker.compose.service": "myservice",
//	}
//	dependencies := compose.ParseDependsOnLabel("postgres:service_started:required,redis")
//	// Get service names for backward compatibility
//	serviceNames := compose.GetServiceNames(dependencies)
//	// Or access individual dependency info
//	for _, dep := range dependencies {
//		fmt.Printf("Service: %s, Restart disabled: %v\n", dep.ServiceName, dep.RestartExplicitlyDisabled)
//	}
//	serviceName := compose.GetServiceName(labels)
package compose

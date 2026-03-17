package compose

import (
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Compose", func() {
	ginkgo.DescribeTable(
		"ParseDependsOnLabel",
		func(input string, expected []ServiceDependency) {
			result := ParseDependsOnLabel(input)
			gomega.Expect(result).To(gomega.Equal(expected))
		},
		ginkgo.Entry(
			"returns nil for empty label",
			"",
			nil,
		),
		ginkgo.Entry(
			"parses single service",
			"postgres",
			[]ServiceDependency{{ServiceName: "postgres", RestartExplicitlyDisabled: false}},
		),
		ginkgo.Entry(
			"parses multiple services",
			"postgres,redis",
			[]ServiceDependency{
				{ServiceName: "postgres", RestartExplicitlyDisabled: false},
				{ServiceName: "redis", RestartExplicitlyDisabled: false},
			},
		),
		ginkgo.Entry(
			"trims whitespace",
			" postgres , redis ",
			[]ServiceDependency{
				{ServiceName: "postgres", RestartExplicitlyDisabled: false},
				{ServiceName: "redis", RestartExplicitlyDisabled: false},
			},
		),
		ginkgo.Entry(
			"parses colon-separated format",
			"postgres:service_started:required,redis:service_healthy",
			[]ServiceDependency{
				{ServiceName: "postgres", RestartExplicitlyDisabled: false},
				{ServiceName: "redis", RestartExplicitlyDisabled: false},
			},
		),
		ginkgo.Entry(
			"ignores empty parts",
			"postgres,,redis",
			[]ServiceDependency{
				{ServiceName: "postgres", RestartExplicitlyDisabled: false},
				{ServiceName: "redis", RestartExplicitlyDisabled: false},
			},
		),
		ginkgo.Entry(
			"parses JSON format with condition only",
			`{"database":{"condition":"service_started"}}`,
			[]ServiceDependency{{ServiceName: "database", RestartExplicitlyDisabled: false}},
		),
		ginkgo.Entry(
			"parses JSON format with multiple services",
			`{"database":{"condition":"service_started"},"cache":{"condition":"service_healthy"}}`,
			[]ServiceDependency{
				{ServiceName: "cache", RestartExplicitlyDisabled: false},
				{ServiceName: "database", RestartExplicitlyDisabled: false},
			},
		),
		ginkgo.Entry(
			"parses JSON format with restart true",
			`{"database":{"condition":"service_started","restart":true}}`,
			[]ServiceDependency{{ServiceName: "database", RestartExplicitlyDisabled: false}},
		),
		ginkgo.Entry(
			"parses JSON format with restart false",
			`{"database":{"condition":"service_started","restart":false}}`,
			[]ServiceDependency{{ServiceName: "database", RestartExplicitlyDisabled: true}},
		),
		ginkgo.Entry(
			"parses JSON format with multiple services and mixed restart",
			`{"database":{"condition":"service_started","restart":false},"cache":{"condition":"service_healthy","restart":true}}`,
			[]ServiceDependency{
				{ServiceName: "cache", RestartExplicitlyDisabled: false},
				{ServiceName: "database", RestartExplicitlyDisabled: true},
			},
		),
		ginkgo.Entry(
			"parses JSON format with some services missing restart property",
			`{"database":{"condition":"service_started","restart":false},"cache":{"condition":"service_healthy"}}`,
			[]ServiceDependency{
				{ServiceName: "cache", RestartExplicitlyDisabled: false},
				{ServiceName: "database", RestartExplicitlyDisabled: true},
			},
		),
	)

	ginkgo.DescribeTable(
		"GetServiceNames",
		func(dependencies []ServiceDependency, expected []string) {
			result := GetServiceNames(dependencies)
			gomega.Expect(result).To(gomega.Equal(expected))
		},
		ginkgo.Entry("returns empty slice for nil", nil, []string{}),
		ginkgo.Entry(
			"returns service names",
			[]ServiceDependency{
				{ServiceName: "postgres", RestartExplicitlyDisabled: false},
				{ServiceName: "redis", RestartExplicitlyDisabled: true},
			},
			[]string{"postgres", "redis"},
		),
	)

	ginkgo.DescribeTable(
		"GetServiceName",
		func(labels map[string]string, expected string) {
			result := GetServiceName(labels)
			gomega.Expect(result).To(gomega.Equal(expected))
		},
		ginkgo.Entry("returns empty string for nil labels", nil, ""),
		ginkgo.Entry("returns empty string for empty labels", map[string]string{}, ""),
		ginkgo.Entry(
			"returns empty string when label not present",
			map[string]string{"other": "value"},
			"",
		),
		ginkgo.Entry(
			"returns service name when label present",
			map[string]string{ComposeServiceLabel: "web"},
			"web",
		),
	)

	ginkgo.DescribeTable(
		"GetProjectName",
		func(labels map[string]string, expected string) {
			result := GetProjectName(labels)
			gomega.Expect(result).To(gomega.Equal(expected))
		},
		ginkgo.Entry("returns empty string for nil labels", nil, ""),
		ginkgo.Entry("returns empty string for empty labels", map[string]string{}, ""),
		ginkgo.Entry(
			"returns empty string when label not present",
			map[string]string{"other": "value"},
			"",
		),
		ginkgo.Entry(
			"returns project name when label present",
			map[string]string{ComposeProjectLabel: "myproject"},
			"myproject",
		),
	)

	ginkgo.DescribeTable(
		"GetContainerNumber",
		func(labels map[string]string, expected string) {
			result := GetContainerNumber(labels)
			gomega.Expect(result).To(gomega.Equal(expected))
		},
		ginkgo.Entry("returns empty string for nil labels", nil, ""),
		ginkgo.Entry("returns empty string for empty labels", map[string]string{}, ""),
		ginkgo.Entry(
			"returns empty string when label not present",
			map[string]string{"other": "value"},
			"",
		),
		ginkgo.Entry(
			"returns container number when label present",
			map[string]string{ComposeContainerNumber: "1"},
			"1",
		),
	)
})

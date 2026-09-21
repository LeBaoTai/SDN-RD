package router

import "strings"

type ConnectionState struct {
	Connection string
	Telemetry  string
	Config     string
}

func (c *ConnectionState) CheckState() string {
	var sb strings.Builder
	sb.WriteString("\nConnection: ")
	sb.WriteString(c.Connection)
	sb.WriteString("\nTelemetry: ")
	sb.WriteString(c.Telemetry)
	sb.WriteString("\nConfig: ")
	sb.WriteString(c.Config)
	return sb.String()
}

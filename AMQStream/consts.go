package AMQStream

// Configuration
const (
	BootstrapServers                 = "BootstrapServers"
	GroupId                          = "GroupId"
	SchemaRegistry                   = "SchemaRegistry"
	SessionTimeoutMs                 = "SessionTimeoutMs"
	SecurityProtocol                 = "SecurityProtocol"
	AutoOffsetReset                  = "AutoOffsetReset"
	SslCertificateLocation           = "SslCertificateLocation"
	MillisecondsTimeout              = "MillisecondsTimeout"
	ConsumerDebug                    = "ConsumerDebug"
	MaxRetry                         = "MaxRetry"
	AutoRegisterSchemas              = "AutoRegisterSchemas"
	MessageMaxBytes                  = "MessageMaxBytes"
	ApplicationName                  = "ApplicationName"
	SchemaUrl                        = "SchemaUrl"
	PartitionAssignmentStrategy      = "PartitionAssignmentStrategy"
	EnableSslCertificateVerification = "EnableSslCertificateVerification"
)

// Headers
const (
	RetryCount    = "RetryCount"
	Reason        = "Reason_%v"
	Deadline      = "Deadline"
	Remitente     = "remitente"
	CrossDeadline = "cross-deadline"
)

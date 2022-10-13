package AMQStream

// Configuration
const (
	BootstrapServers                 = "BootstrapServer"
	GroupId                          = "GroupId"
	SchemaRegistry                   = "SchemaUrl"
	SessionTimeoutMs                 = "SessionTimeoutMs"
	SecurityProtocol                 = "Protocol"
	AutoOffsetReset                  = "AutoOffsetReset"
	SslCertificateLocation           = "SslCertificateLocation"
	MillisecondsTimeout              = "MillisecondsTimeout"
	ConsumerDebug                    = "ConsumerDebug"
	MaxRetry                         = "MaxRetry"
	AutoRegisterSchemas              = "AutoRegisterSchemas"
	MessageMaxBytes                  = "MessageMaxBytes"
	ApplicationName                  = "ApplicationName"
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

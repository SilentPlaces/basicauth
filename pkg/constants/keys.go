package constants

// MySQL configuration keys in Consul
const (
	MySQLHostKey               = "config/mysql/connection/host"
	MySQLPortKey               = "config/mysql/connection/port"
	MySQLUserKey               = "config/mysql/connection/user"
	MySQLPasswordKey           = "config/mysql/connection/password"
	MySQLDBKey                 = "config/mysql/connection/db"
	MySQLMaxLifetimeSecondsKey = "config/mysql/connection/maxLifeTime"
	MySQLIdleConnectionsKey    = "config/mysql/connection/idleConnections"
	MySQLMaxOpenConnectionsKey = "config/mysql/connection/maxOpenConnections"
)

// Redis configuration keys in Consul
const (
	RedisHostKey     = "config/redis/connection/host"
	RedisPortKey     = "config/redis/connection/port"
	RedisPasswordKey = "config/redis/connection/password"
)

// Mail Server config keus
const (
	SMTPHostKey     = "config/mail/connection/host"
	SMTPPortKey     = "config/mail/connection/port"
	SMTPUsernameKey = "config/mail/connection/username"
	SMTPPasswordKey = "config/mail/connection/password"
)

// General config keys of application
const (
	GeneralDomainKey                                = "config/general/domain"
	GeneralHTTPListenerPortKey                      = "config/general/httpListenerPort"
	GeneralRegisterMailVerificationTimeInSecondsKey = "config/general/register/mailVerificationTimeInSeconds"
	GeneralRegisterHostVerificationMailAddressKey   = "config/general/register/hostVerificationMailAddress"
	GeneralRegisterVerificationMailTextKey          = "config/general/register/verificationMailText"
	GeneralMaxVerificationMailCountInDay            = "config/general/register/maxVerificationMailInCountInDay"
)

// Registration password config keys
const (
	KeyRegistrationPasswordMinLength      = "/config/registration/password/minLength"
	KeyRegistrationPasswordRequireUpper   = "/config/registration/password/requireUpper"
	KeyRegistrationPasswordRequireLower   = "/config/registration/password/requireLower"
	KeyRegistrationPasswordRequireNumber  = "/config/registration/password/requireNumber"
	KeyRegistrationPasswordRequireSpecial = "/config/registration/password/requireSpecial"
)

// Environment variable keys
const (
	EnvKeyConsulAddress  = "CONSUL_ADDRESS"
	EnvKeyConsulScheme   = "CONSUL_SCHEME"
	EnvKeyVaultAddr      = "VAULT_ADDR"
	EnvKeyMountPath      = "VAULT_MOUNT_PATH"
	EnvKeySecretPath     = "VAULT_SECRET_PATH"
	EnvKeyVaultToken     = "VAULT_TOKEN"
	EnvKeyAppEnvironment = "APP_ENV"
)

// Security configs keys fetched from vault
const (
	VaultJWTSecretKey        = "jwtSecret"
	VaultJWTRefreshSecretKey = "jwtRefreshSecret"
)

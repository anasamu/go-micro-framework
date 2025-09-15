package core

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	// Use the new go-micro-libs library
	microservices "github.com/anasamu/go-micro-libs"
	"github.com/anasamu/go-micro-libs/api"
	"github.com/anasamu/go-micro-libs/auth"
	"github.com/anasamu/go-micro-libs/cache"
	"github.com/anasamu/go-micro-libs/circuitbreaker"
	"github.com/anasamu/go-micro-libs/communication"
	"github.com/anasamu/go-micro-libs/database"
	"github.com/anasamu/go-micro-libs/database/migrations"
	"github.com/anasamu/go-micro-libs/discovery"
	"github.com/anasamu/go-micro-libs/email"
	"github.com/anasamu/go-micro-libs/event"
	"github.com/anasamu/go-micro-libs/failover"
	"github.com/anasamu/go-micro-libs/filegen"
	"github.com/anasamu/go-micro-libs/logging"
	"github.com/anasamu/go-micro-libs/messaging"
	"github.com/anasamu/go-micro-libs/middleware"
	"github.com/anasamu/go-micro-libs/monitoring"
	"github.com/anasamu/go-micro-libs/payment"
	"github.com/anasamu/go-micro-libs/ratelimit"
	"github.com/anasamu/go-micro-libs/scheduling"
	"github.com/anasamu/go-micro-libs/storage"
)

// Use types from go-micro-libs
type (
	APIManager            = microservices.APIManager
	ConfigManager         = microservices.ConfigManager
	LoggingManager        = microservices.LoggingManager
	MonitoringManager     = microservices.MonitoringManager
	DatabaseManager       = microservices.DatabaseManager
	MigrationManager      = migrations.MigrationManager
	AuthManager           = microservices.AuthManager
	MiddlewareManager     = microservices.MiddlewareManager
	CommunicationManager  = microservices.CommunicationManager
	AIManager             = microservices.AIManager
	StorageManager        = microservices.StorageManager
	MessagingManager      = microservices.MessagingManager
	SchedulingManager     = microservices.SchedulingManager
	BackupManager         = microservices.BackupManager
	ChaosManager          = microservices.ChaosManager
	FailoverManager       = microservices.FailoverManager
	EventManager          = microservices.EventManager
	DiscoveryManager      = microservices.DiscoveryManager
	CacheManager          = microservices.CacheManager
	RateLimitManager      = microservices.RateLimitManager
	CircuitBreakerManager = microservices.CircuitBreakerManager
	FileGenManager        = microservices.FileGenManager
	PaymentManager        = microservices.PaymentManager
	EmailManager          = microservices.EmailManager
)

// Bootstrap manages the initialization and lifecycle of all microservices components
type Bootstrap struct {
	// Core managers using existing libraries
	configManager        *ConfigManager
	loggingManager       *LoggingManager
	monitoringManager    *MonitoringManager
	databaseManager      *DatabaseManager
	migrationManager     *MigrationManager
	authManager          *AuthManager
	middlewareManager    *MiddlewareManager
	communicationManager *CommunicationManager

	// Optional managers using existing libraries
	apiManager            *APIManager
	aiManager             *AIManager
	storageManager        *StorageManager
	messagingManager      *MessagingManager
	schedulingManager     *SchedulingManager
	backupManager         *BackupManager
	chaosManager          *ChaosManager
	failoverManager       *FailoverManager
	eventManager          *EventManager
	discoveryManager      *DiscoveryManager
	cacheManager          *CacheManager
	rateLimitManager      *RateLimitManager
	circuitBreakerManager *CircuitBreakerManager
	filegenManager        *FileGenManager
	paymentManager        *PaymentManager
	emailManager          *EmailManager

	// Framework configuration
	config *FrameworkConfig
	logger *logrus.Logger
	mu     sync.RWMutex
}

// FrameworkConfig holds framework configuration
type FrameworkConfig struct {
	Service    ServiceConfig    `yaml:"service"`
	Server     ServerConfig     `yaml:"server"`
	Database   *DatabaseConfig  `yaml:"database,omitempty"`
	Auth       *AuthConfig      `yaml:"auth,omitempty"`
	Messaging  *MessagingConfig `yaml:"messaging,omitempty"`
	Monitoring MonitoringConfig `yaml:"monitoring"`
	Optional   OptionalConfig   `yaml:"optional"`
}

// ServiceConfig holds service configuration
type ServiceConfig struct {
	Name        string `yaml:"name"`
	Version     string `yaml:"version"`
	Port        int    `yaml:"port"`
	Environment string `yaml:"environment"`
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Host         string        `yaml:"host"`
	Port         int           `yaml:"port"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	IdleTimeout  time.Duration `yaml:"idle_timeout"`
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Providers map[string]interface{} `yaml:"providers"`
}

// AuthConfig holds authentication configuration
type AuthConfig struct {
	Providers map[string]interface{} `yaml:"providers"`
}

// MessagingConfig holds messaging configuration
type MessagingConfig struct {
	Providers map[string]interface{} `yaml:"providers"`
}

// MonitoringConfig holds monitoring configuration
type MonitoringConfig struct {
	Providers map[string]interface{} `yaml:"providers"`
}

// OptionalConfig holds optional features configuration
type OptionalConfig struct {
	API            map[string]interface{} `yaml:"api,omitempty"`
	AI             map[string]interface{} `yaml:"ai,omitempty"`
	Storage        map[string]interface{} `yaml:"storage,omitempty"`
	Scheduling     map[string]interface{} `yaml:"scheduling,omitempty"`
	Backup         map[string]interface{} `yaml:"backup,omitempty"`
	Chaos          map[string]interface{} `yaml:"chaos,omitempty"`
	Failover       map[string]interface{} `yaml:"failover,omitempty"`
	Event          map[string]interface{} `yaml:"event,omitempty"`
	Discovery      map[string]interface{} `yaml:"discovery,omitempty"`
	Cache          map[string]interface{} `yaml:"cache,omitempty"`
	RateLimit      map[string]interface{} `yaml:"ratelimit,omitempty"`
	CircuitBreaker map[string]interface{} `yaml:"circuitbreaker,omitempty"`
	FileGen        map[string]interface{} `yaml:"filegen,omitempty"`
	Payment        map[string]interface{} `yaml:"payment,omitempty"`
	Email          map[string]interface{} `yaml:"email,omitempty"`
}

// Constructor functions using go-micro-libs
func NewAPIManager(config interface{}, logger *logrus.Logger) *APIManager {
	if cfg, ok := config.(*api.ManagerConfig); ok {
		return microservices.NewAPIManager(cfg, logger)
	}
	return microservices.NewAPIManager(nil, logger)
}
func NewConfigManager() *ConfigManager { return microservices.NewConfigManager() }
func NewLoggingManager(config interface{}, logger *logrus.Logger) *LoggingManager {
	if cfg, ok := config.(*logging.ManagerConfig); ok {
		return microservices.NewLoggingManager(cfg, logger)
	}
	return microservices.NewLoggingManager(nil, logger)
}
func NewMonitoringManager(config interface{}, logger *logrus.Logger) *MonitoringManager {
	if cfg, ok := config.(*monitoring.ManagerConfig); ok {
		return microservices.NewMonitoringManager(cfg, logger)
	}
	return microservices.NewMonitoringManager(nil, logger)
}
func NewDatabaseManager(config interface{}, logger *logrus.Logger) *DatabaseManager {
	if cfg, ok := config.(*database.ManagerConfig); ok {
		return microservices.NewDatabaseManager(cfg, logger)
	}
	return microservices.NewDatabaseManager(nil, logger)
}
func NewAuthManager(config interface{}, logger *logrus.Logger) *AuthManager {
	if cfg, ok := config.(*auth.ManagerConfig); ok {
		return microservices.NewAuthManager(cfg, logger)
	}
	return microservices.NewAuthManager(nil, logger)
}
func NewMiddlewareManager(config interface{}, logger *logrus.Logger) *MiddlewareManager {
	if cfg, ok := config.(*middleware.ManagerConfig); ok {
		return microservices.NewMiddlewareManager(cfg, logger)
	}
	return microservices.NewMiddlewareManager(nil, logger)
}
func NewCommunicationManager(config interface{}, logger *logrus.Logger) *CommunicationManager {
	if cfg, ok := config.(*communication.ManagerConfig); ok {
		return microservices.NewCommunicationManager(cfg, logger)
	}
	return microservices.NewCommunicationManager(nil, logger)
}
func NewAIManager() *AIManager { return microservices.NewAIManager() }
func NewStorageManager(config interface{}, logger *logrus.Logger) *StorageManager {
	if cfg, ok := config.(*storage.ManagerConfig); ok {
		return microservices.NewStorageManager(cfg, logger)
	}
	return microservices.NewStorageManager(nil, logger)
}
func NewMessagingManager(config interface{}, logger *logrus.Logger) *MessagingManager {
	if cfg, ok := config.(*messaging.ManagerConfig); ok {
		return microservices.NewMessagingManager(cfg, logger)
	}
	return microservices.NewMessagingManager(nil, logger)
}
func NewSchedulingManager(config interface{}, logger *logrus.Logger) *SchedulingManager {
	if cfg, ok := config.(*scheduling.ManagerConfig); ok {
		return microservices.NewSchedulingManager(cfg, logger)
	}
	return microservices.NewSchedulingManager(nil, logger)
}
func NewBackupManager() *BackupManager { return microservices.NewBackupManager() }
func NewChaosManager() *ChaosManager   { return microservices.NewChaosManager() }
func NewFailoverManager(config interface{}, logger *logrus.Logger) *FailoverManager {
	if cfg, ok := config.(*failover.ManagerConfig); ok {
		return microservices.NewFailoverManager(cfg, logger)
	}
	return microservices.NewFailoverManager(nil, logger)
}
func NewEventManager(config interface{}, logger *logrus.Logger) *EventManager {
	if cfg, ok := config.(*event.ManagerConfig); ok {
		return microservices.NewEventManager(cfg, logger)
	}
	return microservices.NewEventManager(nil, logger)
}
func NewDiscoveryManager(config interface{}, logger *logrus.Logger) *DiscoveryManager {
	if cfg, ok := config.(*discovery.ManagerConfig); ok {
		return microservices.NewDiscoveryManager(cfg, logger)
	}
	return microservices.NewDiscoveryManager(nil, logger)
}
func NewCacheManager(config interface{}, logger *logrus.Logger) *CacheManager {
	if cfg, ok := config.(*cache.ManagerConfig); ok {
		return microservices.NewCacheManager(cfg, logger)
	}
	return microservices.NewCacheManager(nil, logger)
}
func NewRateLimitManager(config interface{}, logger *logrus.Logger) *RateLimitManager {
	if cfg, ok := config.(*ratelimit.ManagerConfig); ok {
		return microservices.NewRateLimitManager(cfg, logger)
	}
	return microservices.NewRateLimitManager(nil, logger)
}
func NewCircuitBreakerManager(config interface{}, logger *logrus.Logger) *CircuitBreakerManager {
	if cfg, ok := config.(*circuitbreaker.ManagerConfig); ok {
		return microservices.NewCircuitBreakerManager(cfg, logger)
	}
	return microservices.NewCircuitBreakerManager(nil, logger)
}
func NewFilegenManager(config interface{}) (*FileGenManager, error) {
	if cfg, ok := config.(*filegen.ManagerConfig); ok {
		return microservices.NewFileGenManager(cfg)
	}
	return microservices.NewFileGenManager(nil)
}
func NewPaymentManager(config interface{}, logger *logrus.Logger) *PaymentManager {
	if cfg, ok := config.(*payment.ManagerConfig); ok {
		return microservices.NewPaymentManager(cfg, logger)
	}
	return microservices.NewPaymentManager(nil, logger)
}
func NewEmailManager(config interface{}, logger *logrus.Logger) *EmailManager {
	if cfg, ok := config.(*email.ManagerConfig); ok {
		return microservices.NewEmailManager(cfg, logger)
	}
	return microservices.NewEmailManager(nil, logger)
}

// Config functions using go-micro-libs
func DefaultAPIManagerConfig() interface{} { return api.DefaultManagerConfig() }
func DefaultManagerConfig() interface{}    { return map[string]interface{}{} }
func DefaultMonitoringManagerConfig() interface{} {
	return monitoring.DefaultManagerConfig()
}
func DefaultDatabaseManagerConfig() interface{}      { return database.DefaultManagerConfig() }
func DefaultAuthManagerConfig() interface{}          { return auth.DefaultManagerConfig() }
func DefaultMiddlewareManagerConfig() interface{}    { return middleware.DefaultManagerConfig() }
func DefaultCommunicationManagerConfig() interface{} { return communication.DefaultManagerConfig() }
func DefaultStorageManagerConfig() interface{}       { return storage.DefaultManagerConfig() }
func DefaultMessagingManagerConfig() interface{} {
	return messaging.DefaultManagerConfig()
}
func DefaultSchedulingManagerConfig() interface{}     { return map[string]interface{}{} }
func DefaultFailoverManagerConfig() interface{}       { return map[string]interface{}{} }
func DefaultEventManagerConfig() interface{}          { return map[string]interface{}{} }
func DefaultDiscoveryManagerConfig() interface{}      { return map[string]interface{}{} }
func DefaultCacheManagerConfig() interface{}          { return map[string]interface{}{} }
func DefaultRateLimitManagerConfig() interface{}      { return map[string]interface{}{} }
func DefaultCircuitBreakerManagerConfig() interface{} { return map[string]interface{}{} }
func DefaultFilegenManagerConfig() interface{}        { return map[string]interface{}{} }
func DefaultPaymentManagerConfig() interface{}        { return payment.DefaultManagerConfig() }
func DefaultEmailManagerConfig() interface{}          { return email.DefaultManagerConfig() }

// NewBootstrap creates a new bootstrap instance
func NewBootstrap(config *FrameworkConfig, logger *logrus.Logger) *Bootstrap {
	if logger == nil {
		logger = logrus.New()
	}

	return &Bootstrap{
		config: config,
		logger: logger,
	}
}

// Initialize initializes all configured components
func (b *Bootstrap) Initialize(ctx context.Context) error {
	b.logger.Info("Initializing microservices framework...")

	// Initialize core components
	if err := b.initializeCoreComponents(ctx); err != nil {
		return fmt.Errorf("failed to initialize core components: %w", err)
	}

	// Initialize optional components
	if err := b.initializeOptionalComponents(ctx); err != nil {
		return fmt.Errorf("failed to initialize optional components: %w", err)
	}

	b.logger.Info("Microservices framework initialized successfully")
	return nil
}

// initializeCoreComponents initializes core components
func (b *Bootstrap) initializeCoreComponents(ctx context.Context) error {
	// Initialize configuration manager
	b.configManager = NewConfigManager()
	b.logger.Info("Configuration manager initialized")

	// Initialize logging manager
	b.loggingManager = NewLoggingManager(
		DefaultManagerConfig(),
		b.logger,
	)
	b.logger.Info("Logging manager initialized")

	// Initialize monitoring manager
	b.monitoringManager = NewMonitoringManager(
		DefaultMonitoringManagerConfig(),
		b.logger,
	)
	b.logger.Info("Monitoring manager initialized")

	// Initialize database manager if configured
	if b.config.Database != nil {
		b.databaseManager = NewDatabaseManager(
			DefaultDatabaseManagerConfig(),
			b.logger,
		)
		b.logger.Info("Database manager initialized")
		
		// Initialize migration manager for database
		b.migrationManager = migrations.NewMigrationManager(
			nil, // Will be set when database provider is available
			b.logger,
		)
		b.logger.Info("Migration manager initialized")
	}

	// Initialize auth manager if configured
	if b.config.Auth != nil {
		b.authManager = NewAuthManager(
			DefaultAuthManagerConfig(),
			b.logger,
		)
		b.logger.Info("Auth manager initialized")
	}

	// Initialize middleware manager
	b.middlewareManager = NewMiddlewareManager(
		DefaultMiddlewareManagerConfig(),
		b.logger,
	)
	b.logger.Info("Middleware manager initialized")

	// Initialize communication manager
	b.communicationManager = NewCommunicationManager(
		DefaultCommunicationManagerConfig(),
		b.logger,
	)
	b.logger.Info("Communication manager initialized")

	return nil
}

// initializeOptionalComponents initializes optional components
func (b *Bootstrap) initializeOptionalComponents(ctx context.Context) error {
	// Initialize API manager if configured
	if b.config.Optional.API != nil {
		b.apiManager = NewAPIManager(
			DefaultAPIManagerConfig(),
			b.logger,
		)
		b.logger.Info("API manager initialized")
	}

	// Initialize AI manager if configured
	if b.config.Optional.AI != nil {
		b.aiManager = NewAIManager()
		b.logger.Info("AI manager initialized")
	}

	// Initialize storage manager if configured
	if b.config.Optional.Storage != nil {
		b.storageManager = NewStorageManager(
			DefaultStorageManagerConfig(),
			b.logger,
		)
		b.logger.Info("Storage manager initialized")
	}

	// Initialize messaging manager if configured
	if b.config.Messaging != nil {
		b.messagingManager = NewMessagingManager(
			DefaultMessagingManagerConfig(),
			b.logger,
		)
		b.logger.Info("Messaging manager initialized")
	}

	// Initialize scheduling manager if configured
	if b.config.Optional.Scheduling != nil {
		b.schedulingManager = NewSchedulingManager(
			DefaultSchedulingManagerConfig(),
			b.logger,
		)
		b.logger.Info("Scheduling manager initialized")
	}

	// Initialize backup manager if configured
	if b.config.Optional.Backup != nil {
		b.backupManager = NewBackupManager()
		b.logger.Info("Backup manager initialized")
	}

	// Initialize chaos manager if configured
	if b.config.Optional.Chaos != nil {
		b.chaosManager = NewChaosManager()
		b.logger.Info("Chaos manager initialized")
	}

	// Initialize failover manager if configured
	if b.config.Optional.Failover != nil {
		b.failoverManager = NewFailoverManager(
			DefaultFailoverManagerConfig(),
			b.logger,
		)
		b.logger.Info("Failover manager initialized")
	}

	// Initialize event manager if configured
	if b.config.Optional.Event != nil {
		b.eventManager = NewEventManager(
			DefaultEventManagerConfig(),
			b.logger,
		)
		b.logger.Info("Event manager initialized")
	}

	// Initialize discovery manager if configured
	if b.config.Optional.Discovery != nil {
		b.discoveryManager = NewDiscoveryManager(
			DefaultDiscoveryManagerConfig(),
			b.logger,
		)
		b.logger.Info("Discovery manager initialized")
	}

	// Initialize cache manager if configured
	if b.config.Optional.Cache != nil {
		b.cacheManager = NewCacheManager(
			DefaultCacheManagerConfig(),
			b.logger,
		)
		b.logger.Info("Cache manager initialized")
	}

	// Initialize rate limit manager if configured
	if b.config.Optional.RateLimit != nil {
		b.rateLimitManager = NewRateLimitManager(
			DefaultRateLimitManagerConfig(),
			b.logger,
		)
		b.logger.Info("Rate limit manager initialized")
	}

	// Initialize circuit breaker manager if configured
	if b.config.Optional.CircuitBreaker != nil {
		b.circuitBreakerManager = NewCircuitBreakerManager(
			DefaultCircuitBreakerManagerConfig(),
			b.logger,
		)
		b.logger.Info("Circuit breaker manager initialized")
	}

	// Initialize file generation manager if configured
	if b.config.Optional.FileGen != nil {
		var err error
		b.filegenManager, err = NewFilegenManager(DefaultFilegenManagerConfig())
		if err != nil {
			return fmt.Errorf("failed to initialize file generation manager: %w", err)
		}
		b.logger.Info("File generation manager initialized")
	}

	// Initialize payment manager if configured
	if b.config.Optional.Payment != nil {
		b.paymentManager = NewPaymentManager(
			DefaultPaymentManagerConfig(),
			b.logger,
		)
		b.logger.Info("Payment manager initialized")
	}

	// Initialize email manager if configured
	if b.config.Optional.Email != nil {
		b.emailManager = NewEmailManager(
			DefaultEmailManagerConfig(),
			b.logger,
		)
		b.logger.Info("Email manager initialized")
	}

	return nil
}

// Start starts all initialized components
func (b *Bootstrap) Start(ctx context.Context) error {
	b.logger.Info("Starting microservices framework...")

	// Start core components
	if err := b.startCoreComponents(ctx); err != nil {
		return fmt.Errorf("failed to start core components: %w", err)
	}

	// Start optional components
	if err := b.startOptionalComponents(ctx); err != nil {
		return fmt.Errorf("failed to start optional components: %w", err)
	}

	b.logger.Info("Microservices framework started successfully")
	return nil
}

// startCoreComponents starts core components
func (b *Bootstrap) startCoreComponents(ctx context.Context) error {
	// Start monitoring
	if b.monitoringManager != nil {
		// Connect to default monitoring provider
		if err := b.monitoringManager.Connect(ctx, "prometheus"); err != nil {
			return fmt.Errorf("failed to start monitoring: %w", err)
		}
	}

	// Start database connections
	if b.databaseManager != nil {
		// Connect to configured databases
		for provider := range b.config.Database.Providers {
			if err := b.databaseManager.Connect(ctx, provider); err != nil {
				return fmt.Errorf("failed to connect to database %s: %w", provider, err)
			}
		}
		
		// Run database migrations if migration manager is available
		if b.migrationManager != nil {
			if err := b.runDatabaseMigrations(ctx); err != nil {
				return fmt.Errorf("failed to run database migrations: %w", err)
			}
		}
	}

	// Start communication server
	if b.communicationManager != nil {
		if err := b.communicationManager.Start(ctx, "http", map[string]interface{}{}); err != nil {
			return fmt.Errorf("failed to start communication: %w", err)
		}
	}

	return nil
}

// startOptionalComponents starts optional components
func (b *Bootstrap) startOptionalComponents(ctx context.Context) error {
	// API manager is ready to use (no explicit start needed)

	// Email manager is ready to use (no explicit start needed)

	// Start messaging
	if b.messagingManager != nil {
		for provider := range b.config.Messaging.Providers {
			if err := b.messagingManager.Connect(ctx, provider); err != nil {
				return fmt.Errorf("failed to connect to messaging %s: %w", provider, err)
			}
		}
	}

	// Discovery manager is ready to use (no explicit start needed)

	return nil
}

// Stop stops all components gracefully
func (b *Bootstrap) Stop(ctx context.Context) error {
	b.logger.Info("Stopping microservices framework...")

	// Stop components in reverse order
	if b.communicationManager != nil {
		b.communicationManager.Stop(ctx, "http")
	}

	if b.apiManager != nil {
		b.apiManager.Close()
	}

	if b.emailManager != nil {
		b.emailManager.Close()
	}

	if b.messagingManager != nil {
		b.messagingManager.Close()
	}

	if b.databaseManager != nil {
		b.databaseManager.Close()
	}

	if b.monitoringManager != nil {
		b.monitoringManager.Close()
	}

	b.logger.Info("Microservices framework stopped successfully")
	return nil
}

// HealthCheck performs health check on all components
func (b *Bootstrap) HealthCheck(ctx context.Context) map[string]interface{} {
	health := make(map[string]interface{})

	// Check core components
	if b.monitoringManager != nil {
		health["monitoring"] = b.monitoringManager.HealthCheckAll(ctx)
	}

	if b.databaseManager != nil {
		health["database"] = b.databaseManager.HealthCheck(ctx)
	}

	if b.authManager != nil {
		health["auth"] = b.authManager.HealthCheck(ctx)
	}

	if b.apiManager != nil {
		health["api"] = b.apiManager.HealthCheck(ctx)
	}

	if b.emailManager != nil {
		health["email"] = b.emailManager.HealthCheck(ctx)
	}

	if b.messagingManager != nil {
		health["messaging"] = b.messagingManager.HealthCheck(ctx)
	}

	if b.discoveryManager != nil {
		health["discovery"] = "ready"
	}

	if b.cacheManager != nil {
		health["cache"] = "ready"
	}

	return health
}

// GetManager returns a specific manager by name
func (b *Bootstrap) GetManager(name string) interface{} {
	b.mu.RLock()
	defer b.mu.RUnlock()

	switch name {
	case "config":
		return b.configManager
	case "logging":
		return b.loggingManager
	case "monitoring":
		return b.monitoringManager
	case "database":
		return b.databaseManager
	case "auth":
		return b.authManager
	case "middleware":
		return b.middlewareManager
	case "communication":
		return b.communicationManager
	case "api":
		return b.apiManager
	case "ai":
		return b.aiManager
	case "storage":
		return b.storageManager
	case "messaging":
		return b.messagingManager
	case "scheduling":
		return b.schedulingManager
	case "backup":
		return b.backupManager
	case "chaos":
		return b.chaosManager
	case "failover":
		return b.failoverManager
	case "event":
		return b.eventManager
	case "discovery":
		return b.discoveryManager
	case "cache":
		return b.cacheManager
	case "ratelimit":
		return b.rateLimitManager
	case "circuitbreaker":
		return b.circuitBreakerManager
	case "filegen":
		return b.filegenManager
	case "payment":
		return b.paymentManager
	case "email":
		return b.emailManager
	case "migration":
		return b.migrationManager
	default:
		return nil
	}
}

// runDatabaseMigrations runs database migrations using the migration manager
func (b *Bootstrap) runDatabaseMigrations(ctx context.Context) error {
	// Get default database provider
	provider, err := b.databaseManager.GetDefaultProvider()
	if err != nil {
		return fmt.Errorf("failed to get default database provider: %w", err)
	}
	
	// Set the provider in migration manager
	b.migrationManager = migrations.NewMigrationManager(provider, b.logger)
	
	// Initialize migration table
	if err := b.migrationManager.Initialize(ctx); err != nil {
		return fmt.Errorf("failed to initialize migration table: %w", err)
	}
	
	// Create CLI manager for migrations
	cliManager := migrations.NewCLIManager(provider, "./migrations", b.logger)
	
	// Apply pending migrations
	if err := cliManager.Up(ctx); err != nil {
		return fmt.Errorf("failed to apply database migrations: %w", err)
	}
	
	b.logger.Info("Database migrations applied successfully")
	return nil
}

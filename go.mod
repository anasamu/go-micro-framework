module github.com/anasamu/go-micro-framework

go 1.21

require (

	// All microservices-library-go dependencies
	github.com/anasamu/microservices-library-go/ai v0.1.0
	github.com/anasamu/microservices-library-go/auth v0.1.0
	github.com/anasamu/microservices-library-go/backup v0.1.0
	github.com/anasamu/microservices-library-go/cache v0.1.0
	github.com/anasamu/microservices-library-go/chaos v0.1.0
	github.com/anasamu/microservices-library-go/circuitbreaker v0.1.0
	github.com/anasamu/microservices-library-go/communication v0.1.0
	github.com/anasamu/microservices-library-go/config v0.1.0
	github.com/anasamu/microservices-library-go/database v0.1.0
	github.com/anasamu/microservices-library-go/discovery v0.1.0
	github.com/anasamu/microservices-library-go/event v0.1.0
	github.com/anasamu/microservices-library-go/failover v0.1.0
	github.com/anasamu/microservices-library-go/filegen v0.1.0
	github.com/anasamu/microservices-library-go/logging v0.1.0
	github.com/anasamu/microservices-library-go/messaging v0.1.0
	github.com/anasamu/microservices-library-go/middleware v0.1.0
	github.com/anasamu/microservices-library-go/monitoring v0.1.0
	github.com/anasamu/microservices-library-go/payment v0.1.0
	github.com/anasamu/microservices-library-go/ratelimit v0.1.0
	github.com/anasamu/microservices-library-go/scheduling v0.1.0
	github.com/anasamu/microservices-library-go/storage v0.1.0
	github.com/sirupsen/logrus v1.9.3
	// Core framework dependencies
	github.com/spf13/cobra v1.7.0

	// Testing
	github.com/stretchr/testify v1.8.4 // indirect
)

require (
	github.com/anasamu/microservices-library-go/ai/providers/anthropic v0.0.0-00010101000000-000000000000 // indirect
	github.com/anasamu/microservices-library-go/ai/providers/deepseek v0.0.0-00010101000000-000000000000 // indirect
	github.com/anasamu/microservices-library-go/ai/providers/google v0.0.0-00010101000000-000000000000 // indirect
	github.com/anasamu/microservices-library-go/ai/providers/openai v0.0.0-00010101000000-000000000000 // indirect
	github.com/anasamu/microservices-library-go/ai/providers/xai v0.0.0-00010101000000-000000000000 // indirect
	github.com/anasamu/microservices-library-go/ai/types v0.0.0 // indirect
	github.com/anasamu/microservices-library-go/auth/types v0.0.0-00010101000000-000000000000 // indirect
	github.com/anasamu/microservices-library-go/cache/types v0.0.0-00010101000000-000000000000 // indirect
	github.com/anasamu/microservices-library-go/chaos/types v0.0.0 // indirect
	github.com/anasamu/microservices-library-go/circuitbreaker/types v0.0.0-00010101000000-000000000000 // indirect
	github.com/anasamu/microservices-library-go/config/types v0.0.0-00010101000000-000000000000 // indirect
	github.com/anasamu/microservices-library-go/discovery/types v0.0.0-00010101000000-000000000000 // indirect
	github.com/anasamu/microservices-library-go/event/types v0.0.0-20250910142242-8bec92b8b0f4 // indirect
	github.com/anasamu/microservices-library-go/failover/types v0.0.0 // indirect
	github.com/anasamu/microservices-library-go/filegen/providers/csv v0.0.0-00010101000000-000000000000 // indirect
	github.com/anasamu/microservices-library-go/filegen/providers/custom v0.0.0-00010101000000-000000000000 // indirect
	github.com/anasamu/microservices-library-go/filegen/providers/docx v0.0.0-00010101000000-000000000000 // indirect
	github.com/anasamu/microservices-library-go/filegen/providers/excel v0.0.0-00010101000000-000000000000 // indirect
	github.com/anasamu/microservices-library-go/filegen/providers/pdf v0.0.0-00010101000000-000000000000 // indirect
	github.com/anasamu/microservices-library-go/filegen/types v0.0.0 // indirect
	github.com/anasamu/microservices-library-go/logging/types v0.0.0-00010101000000-000000000000 // indirect
	github.com/anasamu/microservices-library-go/middleware/types v0.0.0-00010101000000-000000000000 // indirect
	github.com/anasamu/microservices-library-go/monitoring/types v0.0.0-20250910142242-8bec92b8b0f4 // indirect
	github.com/anasamu/microservices-library-go/ratelimit/types v0.0.0-00010101000000-000000000000 // indirect
	github.com/anasamu/microservices-library-go/scheduling/types v0.0.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826 // indirect
	github.com/richardlehane/mscfb v1.0.4 // indirect
	github.com/richardlehane/msoleps v1.0.3 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/unidoc/unioffice v1.26.0 // indirect
	github.com/xuri/efp v0.0.0-20231025114914-d1ff6096ae53 // indirect
	github.com/xuri/excelize/v2 v2.8.0 // indirect
	github.com/xuri/nfp v0.0.0-20230919160717-d98342af3f05 // indirect
	golang.org/x/crypto v0.18.0 // indirect
	golang.org/x/net v0.20.0 // indirect
	golang.org/x/sys v0.16.0 // indirect
	golang.org/x/text v0.14.0 // indirect
)

// Replace with local microservices-library-go packages
replace github.com/anasamu/microservices-library-go/ai => ../microservices-library-go/ai

replace github.com/anasamu/microservices-library-go/ai/providers/anthropic => ../microservices-library-go/ai/providers/anthropic

replace github.com/anasamu/microservices-library-go/ai/providers/deepseek => ../microservices-library-go/ai/providers/deepseek

replace github.com/anasamu/microservices-library-go/ai/providers/google => ../microservices-library-go/ai/providers/google

replace github.com/anasamu/microservices-library-go/ai/providers/openai => ../microservices-library-go/ai/providers/openai

replace github.com/anasamu/microservices-library-go/ai/providers/xai => ../microservices-library-go/ai/providers/xai

replace github.com/anasamu/microservices-library-go/ai/types => ../microservices-library-go/ai/types

replace github.com/anasamu/microservices-library-go/auth => ../microservices-library-go/auth

replace github.com/anasamu/microservices-library-go/auth/providers/authentication/jwt => ../microservices-library-go/auth/providers/authentication/jwt

replace github.com/anasamu/microservices-library-go/auth/providers/authentication/oauth => ../microservices-library-go/auth/providers/authentication/oauth

replace github.com/anasamu/microservices-library-go/auth/providers/authentication/twofa => ../microservices-library-go/auth/providers/authentication/twofa

replace github.com/anasamu/microservices-library-go/auth/providers/authorization/abac => ../microservices-library-go/auth/providers/authorization/abac

replace github.com/anasamu/microservices-library-go/auth/providers/authorization/acl => ../microservices-library-go/auth/providers/authorization/acl

replace github.com/anasamu/microservices-library-go/auth/providers/authorization/rbac => ../microservices-library-go/auth/providers/authorization/rbac

replace github.com/anasamu/microservices-library-go/auth/types => ../microservices-library-go/auth/types

replace github.com/anasamu/microservices-library-go/backup => ../microservices-library-go/backup

replace github.com/anasamu/microservices-library-go/backup/providers/gcs => ../microservices-library-go/backup/providers/gcs

replace github.com/anasamu/microservices-library-go/backup/providers/local => ../microservices-library-go/backup/providers/local

replace github.com/anasamu/microservices-library-go/backup/providers/s3 => ../microservices-library-go/backup/providers/s3

replace github.com/anasamu/microservices-library-go/cache => ../microservices-library-go/cache

replace github.com/anasamu/microservices-library-go/cache/providers/memcache => ../microservices-library-go/cache/providers/memcache

replace github.com/anasamu/microservices-library-go/cache/providers/memory => ../microservices-library-go/cache/providers/memory

replace github.com/anasamu/microservices-library-go/cache/providers/redis => ../microservices-library-go/cache/providers/redis

replace github.com/anasamu/microservices-library-go/cache/types => ../microservices-library-go/cache/types

replace github.com/anasamu/microservices-library-go/chaos => ../microservices-library-go/chaos

replace github.com/anasamu/microservices-library-go/chaos/providers/http => ../microservices-library-go/chaos/providers/http

replace github.com/anasamu/microservices-library-go/chaos/providers/kubernetes => ../microservices-library-go/chaos/providers/kubernetes

replace github.com/anasamu/microservices-library-go/chaos/providers/messaging => ../microservices-library-go/chaos/providers/messaging

replace github.com/anasamu/microservices-library-go/chaos/types => ../microservices-library-go/chaos/types

replace github.com/anasamu/microservices-library-go/circuitbreaker => ../microservices-library-go/circuitbreaker

replace github.com/anasamu/microservices-library-go/circuitbreaker/providers/custom => ../microservices-library-go/circuitbreaker/providers/custom

replace github.com/anasamu/microservices-library-go/circuitbreaker/providers/gobreaker => ../microservices-library-go/circuitbreaker/providers/gobreaker

replace github.com/anasamu/microservices-library-go/circuitbreaker/types => ../microservices-library-go/circuitbreaker/types

replace github.com/anasamu/microservices-library-go/communication => ../microservices-library-go/communication

replace github.com/anasamu/microservices-library-go/communication/providers/graphql => ../microservices-library-go/communication/providers/graphql

replace github.com/anasamu/microservices-library-go/communication/providers/grpc => ../microservices-library-go/communication/providers/grpc

replace github.com/anasamu/microservices-library-go/communication/providers/http => ../microservices-library-go/communication/providers/http

replace github.com/anasamu/microservices-library-go/communication/providers/quic => ../microservices-library-go/communication/providers/quic

replace github.com/anasamu/microservices-library-go/communication/providers/sse => ../microservices-library-go/communication/providers/sse

replace github.com/anasamu/microservices-library-go/communication/providers/websocket => ../microservices-library-go/communication/providers/websocket

replace github.com/anasamu/microservices-library-go/config => ../microservices-library-go/config

replace github.com/anasamu/microservices-library-go/config/providers/consul => ../microservices-library-go/config/providers/consul

replace github.com/anasamu/microservices-library-go/config/providers/env => ../microservices-library-go/config/providers/env

replace github.com/anasamu/microservices-library-go/config/providers/file => ../microservices-library-go/config/providers/file

replace github.com/anasamu/microservices-library-go/config/providers/vault => ../microservices-library-go/config/providers/vault

replace github.com/anasamu/microservices-library-go/config/types => ../microservices-library-go/config/types

replace github.com/anasamu/microservices-library-go/database => ../microservices-library-go/database

replace github.com/anasamu/microservices-library-go/database/cmd/migrate => ../microservices-library-go/database/cmd/migrate

replace github.com/anasamu/microservices-library-go/database/migrations => ../microservices-library-go/database/migrations

replace github.com/anasamu/microservices-library-go/database/providers/cassandra => ../microservices-library-go/database/providers/cassandra

replace github.com/anasamu/microservices-library-go/database/providers/cockroachdb => ../microservices-library-go/database/providers/cockroachdb

replace github.com/anasamu/microservices-library-go/database/providers/elasticsearch => ../microservices-library-go/database/providers/elasticsearch

replace github.com/anasamu/microservices-library-go/database/providers/influxdb => ../microservices-library-go/database/providers/influxdb

replace github.com/anasamu/microservices-library-go/database/providers/mariadb => ../microservices-library-go/database/providers/mariadb

replace github.com/anasamu/microservices-library-go/database/providers/mongodb => ../microservices-library-go/database/providers/mongodb

replace github.com/anasamu/microservices-library-go/database/providers/mysql => ../microservices-library-go/database/providers/mysql

replace github.com/anasamu/microservices-library-go/database/providers/postgresql => ../microservices-library-go/database/providers/postgresql

replace github.com/anasamu/microservices-library-go/database/providers/redis => ../microservices-library-go/database/providers/redis

replace github.com/anasamu/microservices-library-go/database/providers/sqlite => ../microservices-library-go/database/providers/sqlite

replace github.com/anasamu/microservices-library-go/discovery => ../microservices-library-go/discovery

replace github.com/anasamu/microservices-library-go/discovery/providers/consul => ../microservices-library-go/discovery/providers/consul

replace github.com/anasamu/microservices-library-go/discovery/providers/etcd => ../microservices-library-go/discovery/providers/etcd

replace github.com/anasamu/microservices-library-go/discovery/providers/kubernetes => ../microservices-library-go/discovery/providers/kubernetes

replace github.com/anasamu/microservices-library-go/discovery/providers/static => ../microservices-library-go/discovery/providers/static

replace github.com/anasamu/microservices-library-go/discovery/types => ../microservices-library-go/discovery/types

replace github.com/anasamu/microservices-library-go/event => ../microservices-library-go/event

replace github.com/anasamu/microservices-library-go/event/providers/kafka => ../microservices-library-go/event/providers/kafka

replace github.com/anasamu/microservices-library-go/event/providers/nats => ../microservices-library-go/event/providers/nats

replace github.com/anasamu/microservices-library-go/event/providers/postgresql => ../microservices-library-go/event/providers/postgresql

replace github.com/anasamu/microservices-library-go/event/types => ../microservices-library-go/event/types

replace github.com/anasamu/microservices-library-go/failover => ../microservices-library-go/failover

replace github.com/anasamu/microservices-library-go/failover/providers/consul => ../microservices-library-go/failover/providers/consul

replace github.com/anasamu/microservices-library-go/failover/providers/kubernetes => ../microservices-library-go/failover/providers/kubernetes

replace github.com/anasamu/microservices-library-go/failover/types => ../microservices-library-go/failover/types

replace github.com/anasamu/microservices-library-go/filegen => ../microservices-library-go/filegen

replace github.com/anasamu/microservices-library-go/filegen/providers/csv => ../microservices-library-go/filegen/providers/csv

replace github.com/anasamu/microservices-library-go/filegen/providers/custom => ../microservices-library-go/filegen/providers/custom

replace github.com/anasamu/microservices-library-go/filegen/providers/docx => ../microservices-library-go/filegen/providers/docx

replace github.com/anasamu/microservices-library-go/filegen/providers/excel => ../microservices-library-go/filegen/providers/excel

replace github.com/anasamu/microservices-library-go/filegen/providers/pdf => ../microservices-library-go/filegen/providers/pdf

replace github.com/anasamu/microservices-library-go/filegen/types => ../microservices-library-go/filegen/types

replace github.com/anasamu/microservices-library-go/logging => ../microservices-library-go/logging

replace github.com/anasamu/microservices-library-go/logging/providers/console => ../microservices-library-go/logging/providers/console

replace github.com/anasamu/microservices-library-go/logging/providers/elasticsearch => ../microservices-library-go/logging/providers/elasticsearch

replace github.com/anasamu/microservices-library-go/logging/providers/file => ../microservices-library-go/logging/providers/file

replace github.com/anasamu/microservices-library-go/logging/types => ../microservices-library-go/logging/types

replace github.com/anasamu/microservices-library-go/messaging => ../microservices-library-go/messaging

replace github.com/anasamu/microservices-library-go/messaging/providers/kafka => ../microservices-library-go/messaging/providers/kafka

replace github.com/anasamu/microservices-library-go/messaging/providers/nats => ../microservices-library-go/messaging/providers/nats

replace github.com/anasamu/microservices-library-go/messaging/providers/rabbitmq => ../microservices-library-go/messaging/providers/rabbitmq

replace github.com/anasamu/microservices-library-go/messaging/providers/sqs => ../microservices-library-go/messaging/providers/sqs

replace github.com/anasamu/microservices-library-go/middleware => ../microservices-library-go/middleware

replace github.com/anasamu/microservices-library-go/middleware/providers/auth => ../microservices-library-go/middleware/providers/auth

replace github.com/anasamu/microservices-library-go/middleware/providers/cache => ../microservices-library-go/middleware/providers/cache

replace github.com/anasamu/microservices-library-go/middleware/providers/chaos => ../microservices-library-go/middleware/providers/chaos

replace github.com/anasamu/microservices-library-go/middleware/providers/circuitbreaker => ../microservices-library-go/middleware/providers/circuitbreaker

replace github.com/anasamu/microservices-library-go/middleware/providers/communication => ../microservices-library-go/middleware/providers/communication

replace github.com/anasamu/microservices-library-go/middleware/providers/failover => ../microservices-library-go/middleware/providers/failover

replace github.com/anasamu/microservices-library-go/middleware/providers/logging => ../microservices-library-go/middleware/providers/logging

replace github.com/anasamu/microservices-library-go/middleware/providers/messaging => ../microservices-library-go/middleware/providers/messaging

replace github.com/anasamu/microservices-library-go/middleware/providers/monitoring => ../microservices-library-go/middleware/providers/monitoring

replace github.com/anasamu/microservices-library-go/middleware/providers/ratelimit => ../microservices-library-go/middleware/providers/ratelimit

replace github.com/anasamu/microservices-library-go/middleware/providers/storage => ../microservices-library-go/middleware/providers/storage

replace github.com/anasamu/microservices-library-go/middleware/types => ../microservices-library-go/middleware/types

replace github.com/anasamu/microservices-library-go/monitoring => ../microservices-library-go/monitoring

replace github.com/anasamu/microservices-library-go/monitoring/providers/elasticsearch => ../microservices-library-go/monitoring/providers/elasticsearch

replace github.com/anasamu/microservices-library-go/monitoring/providers/jaeger => ../microservices-library-go/monitoring/providers/jaeger

replace github.com/anasamu/microservices-library-go/monitoring/providers/prometheus => ../microservices-library-go/monitoring/providers/prometheus

replace github.com/anasamu/microservices-library-go/monitoring/types => ../microservices-library-go/monitoring/types

replace github.com/anasamu/microservices-library-go/payment => ../microservices-library-go/payment

replace github.com/anasamu/microservices-library-go/payment/providers/midtrans => ../microservices-library-go/payment/providers/midtrans

replace github.com/anasamu/microservices-library-go/payment/providers/paypal => ../microservices-library-go/payment/providers/paypal

replace github.com/anasamu/microservices-library-go/payment/providers/stripe => ../microservices-library-go/payment/providers/stripe

replace github.com/anasamu/microservices-library-go/payment/providers/xendit => ../microservices-library-go/payment/providers/xendit

replace github.com/anasamu/microservices-library-go/ratelimit => ../microservices-library-go/ratelimit

replace github.com/anasamu/microservices-library-go/ratelimit/providers/inmemory => ../microservices-library-go/ratelimit/providers/inmemory

replace github.com/anasamu/microservices-library-go/ratelimit/providers/redis => ../microservices-library-go/ratelimit/providers/redis

replace github.com/anasamu/microservices-library-go/ratelimit/types => ../microservices-library-go/ratelimit/types

replace github.com/anasamu/microservices-library-go/scheduling => ../microservices-library-go/scheduling

replace github.com/anasamu/microservices-library-go/scheduling/providers/cron => ../microservices-library-go/scheduling/providers/cron

replace github.com/anasamu/microservices-library-go/scheduling/providers/redis => ../microservices-library-go/scheduling/providers/redis

replace github.com/anasamu/microservices-library-go/scheduling/types => ../microservices-library-go/scheduling/types

replace github.com/anasamu/microservices-library-go/storage => ../microservices-library-go/storage

replace github.com/anasamu/microservices-library-go/storage/providers/azure => ../microservices-library-go/storage/providers/azure

replace github.com/anasamu/microservices-library-go/storage/providers/gcs => ../microservices-library-go/storage/providers/gcs

replace github.com/anasamu/microservices-library-go/storage/providers/minio => ../microservices-library-go/storage/providers/minio

replace github.com/anasamu/microservices-library-go/storage/providers/s3 => ../microservices-library-go/storage/providers/s3

replace github.com/anasamu/microservices-library-go/storage/types => ../microservices-library-go/storage/types

replace github.com/anasamu/microservices-library-go/utils => ../microservices-library-go/utils

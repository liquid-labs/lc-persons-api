module github.com/Liquid-Labs/lc-persons-api

require (
	github.com/Liquid-Labs/catalyst-core-api v1.0.0-prototype.15
	github.com/Liquid-Labs/go-api v1.0.0-protottype.0
	github.com/Liquid-Labs/go-rest v1.0.0-prototype.4
	github.com/Liquid-Labs/lc-authentication-api v0.0.0-20190817161517-b440787415e4
	github.com/Liquid-Labs/lc-entities-model v1.0.0-alpha.0
	github.com/Liquid-Labs/lc-locations-model v1.0.0-alpha.1
	github.com/Liquid-Labs/lc-persons-model v1.0.0-alpha.1
	github.com/Liquid-Labs/lc-rdb-service v1.0.0-alpha.1
	github.com/Liquid-Labs/strkit v0.0.0-20190818184832-9e3e35dcfc9c
	github.com/Liquid-Labs/terror v1.0.0-alpha.1
	github.com/golang/mock v1.3.1
	github.com/gorilla/mux v1.7.0
	github.com/stretchr/testify v1.4.0
)

replace github.com/Liquid-Labs/catalyst-core-api => ../catalyst-core-api

replace github.com/Liquid-Labs/lc-authentication-api => ../lc-authentication-api

replace github.com/Liquid-Labs/lc-entities-model => ../lc-entities-model

replace github.com/Liquid-Labs/lc-locations-model => ../lc-locations-model

replace github.com/Liquid-Labs/lc-persons-model => ../lc-persons-model

replace github.com/Liquid-Labs/lc-rdb-service => ../lc-rdb-service

replace github.com/Liquid-Labs/lc-users-model => ../lc-users-model

replace github.com/Liquid-Labs/go-api => ../go-api

replace github.com/Liquid-Labs/go-rest => ../go-rest

replace github.com/Liquid-Labs/terror => ../terror

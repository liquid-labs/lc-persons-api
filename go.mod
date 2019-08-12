module github.com/Liquid-Labs/lc-persons-api

require (
	github.com/Liquid-Labs/catalyst-core-api v1.0.0-prototype.15
	github.com/Liquid-Labs/catalyst-firewrap v1.0.0-prototype.0
	github.com/Liquid-Labs/catalyst-persons-api v1.0.0-prototype.8
	github.com/Liquid-Labs/go-api v1.0.0-protottype.0
	github.com/Liquid-Labs/go-nullable-mysql v1.0.2
	github.com/Liquid-Labs/go-rest v1.0.0-prototype.4
	github.com/gorilla/mux v1.7.0
	github.com/stretchr/testify v1.3.0
)

replace github.com/Liquid-Labs/catalyst-core-api => ../catalyst-core-api

replace github.com/Liquid-Labs/catalyst-firewrap => ../catalyst-firewrap

replace github.com/Liquid-Labs/go-api => ../go-api

replace github.com/Liquid-Labs/go-rest => ../go-rest

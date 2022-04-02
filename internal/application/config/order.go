package config

type order struct {
	Quantity       int64  `env:"ORDER_QUANTITY" envDefault:"25"`
	SourceCurrency string `env:"ORDER_SOURCE_CURRENCY" envDefault:"BNB"`
	TargetCurrency string `env:"ORDER_TARGET_CURRENCY" envDefault:"USDT"`
	Price          string `env:"ORDER_PRICE" envDefault:"42"`
}

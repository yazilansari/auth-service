package config

type Tenant struct {
	TenantCode  string
	CountryCode string
}

var TenantMap = map[string]Tenant{
	"ae.ahmedalmaghribi.com": {
		TenantCode:  "AE",
		CountryCode: "AE",
	},
	"ksa.ahmedalmaghribi.com": {
		TenantCode:  "SA",
		CountryCode: "SA",
	},
	"qa.ahmedalmaghribi.com": {
		TenantCode:  "QA",
		CountryCode: "QA",
	},
	"kw.ahmedalmaghribi.com": {
		TenantCode:  "KW",
		CountryCode: "KW",
	},
	"bh.ahmedalmaghribi.com": {
		TenantCode:  "BH",
		CountryCode: "BH",
	},
	"om.ahmedalmaghribi.com": {
		TenantCode:  "OM",
		CountryCode: "OM",
	},
}

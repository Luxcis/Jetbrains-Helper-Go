package helper

type ProductCache struct {
	Name        string `json:"name"`
	ProductCode string `json:"productCode"`
	IconClass   string `json:"iconClass"`
}

type PluginCache struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	ProductCode  string `json:"productCode"`
	PricingModel string `json:"pricingModel"`
	Icon         string `json:"icon"`
}

type PluginInfo struct {
	ID           int          `json:"id"`
	PurchaseInfo PurchaseInfo `json:"purchaseInfo"`
}

type PurchaseInfo struct {
	ProductCode string `json:"productCode"`
}

type PluginList struct {
	Plugins []Plugin `json:"plugins"`
	Total   int      `json:"total"`
}

type Plugin struct {
	Id           int        `json:"id"`
	Name         string     `json:"name"`
	Preview      string     `json:"preview"`
	Downloads    int        `json:"downloads"`
	PricingModel string     `json:"pricingModel"`
	Organization string     `json:"organization"`
	Icon         string     `json:"icon"`
	PreviewImage string     `json:"previewImage"`
	Rating       float64    `json:"rating"`
	VendorInfo   VendorInfo `json:"vendorInfo"`
}

type VendorInfo struct {
	Name       string `json:"name"`
	IsVerified bool   `json:"isVerified"`
}

type LicensePart struct {
	LicenseID    string    `json:"licenseId"`
	LicenseeName string    `json:"licenseeName"`
	AssigneeName string    `json:"assigneeName"`
	Products     []Product `json:"products"`
	Metadata     string    `json:"metadata"`
}

type Product struct {
	Code         string `json:"code"`
	FallbackDate string `json:"fallbackDate"`
	PaidUpTo     string `json:"paidUpTo"`
}

type GenerateLicenseReqBody struct {
	LicenseName  string `json:"licenseName"`
	AssigneeName string `json:"assigneeName"`
	ExpiryDate   string `json:"expiryDate"`
	ProductCode  string `json:"productCode"`
}

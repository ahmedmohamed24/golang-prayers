package resources

type AzanResource struct {
	Message  string `json:"message"`
	Date     string `json:"date"`
	Fajr     string `json:"Fajr"`
	Sunrise  string `json:"sunrise"`
	Zuhr     string `json:"zuhr"`
	Asr      string `json:"asr"`
	Maghrib  string `json:"maghrib"`
	Isha     string `json:"isha"`
	City     string `json:"city"`
	TimeZone string `json:"timezone"`
}

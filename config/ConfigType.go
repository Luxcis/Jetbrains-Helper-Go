package config

var Conf Config

type Config struct {
	Mode    string  `required:"true" default:"release"`
	License License `required:"true"`
}
type License struct {
	LicenseName  string `required:"true" default:"Azurlane"`
	AssigneeName string `required:"true" default:"Yamato"`
	ExpiryDate   string `required:"true" default:"2030-12-31"`
}

package configs

import _ "embed"

//go:embed  config.yaml
var Config []byte

//go:embed  reserved-usernames.json
var ReservedUsernames []byte

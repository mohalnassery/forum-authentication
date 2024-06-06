module forum

go 1.22

require (
	github.com/gofrs/uuid v4.4.0+incompatible
	github.com/google/go-github/v52 v52.0.0
	github.com/mattn/go-sqlite3 v1.14.22
	golang.org/x/crypto v0.21.0
	golang.org/x/oauth2 v0.20.0
	golang.org/x/time v0.5.0
)

require (
	cloud.google.com/go/compute/metadata v0.3.0 // indirect
	github.com/ProtonMail/go-crypto v0.0.0-20230217124315-7d5c6f04bbb8 // indirect
	github.com/cloudflare/circl v1.1.0 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	golang.org/x/sys v0.18.0 // indirect
)

// Removed unnecessary indirect dependencies
// require github.com/google/go-cmp v0.6.0 // indirect
// require cloud.google.com/go/compute/metadata v0.3.0 // indirect
// require github.com/ProtonMail/go-crypto v0.0.0-20230217124315-7d5c6f04bbb8 // indirect
// require github.com/cloudflare/circl v1.1.0 // indirect
// require github.com/google/go-github v17.0.0+incompatible // indirect
// require github.com/google/go-github/v62 v62.0.1-0.20240528231835-6257442bcdfa // indirect
// require github.com/google/go-querystring v1.1.0 // indirect
// require github.com/huandu/facebook v2.3.1+incompatible // indirect

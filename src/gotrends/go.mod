module gotrends

go 1.13

require (
	github.com/deckarep/golang-set v1.7.1 // indirect
	github.com/mingrammer/commonregex v1.0.1 // indirect
	gonum.org/v1/gonum v0.0.0-20191013112747-a2e99c9265e9 // indirect
	gopkg.in/jdkato/prose.v2 v2.0.0-20190814032740-822d591a158c
	gopkg.in/neurosnap/sentences.v1 v1.0.6 // indirect
	redditclient v0.0.0-00010101000000-000000000000
)

replace redditclient => ../redditclient

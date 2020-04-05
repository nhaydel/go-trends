module github.com/nhaydel/go-trends

go 1.14

require (
	github.com/deckarep/golang-set v1.7.1 // indirect
	github.com/mingrammer/commonregex v1.0.1 // indirect
	gonum.org/v1/gonum v0.7.0 // indirect
	gopkg.in/jdkato/prose.v2 v2.0.0-20190814032740-822d591a158c
	gopkg.in/neurosnap/sentences.v1 v1.0.6 // indirect
)

replace "github.com/nhaydel/go-trends/internal/trendsmap" v0.0.0 => ./internal/trendsmap

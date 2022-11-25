module github.com/feduxorg/letter-generator

go 1.12

require (
	github.com/cpuguy83/go-md2man/v2 v2.0.2 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/libgit2/git2go/v34 v34.0.0
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.9.0
	github.com/stretchr/testify v1.8.1
	github.com/urfave/cli v1.22.10
	golang.org/x/crypto v0.3.0 // indirect
	golang.org/x/net v0.2.0
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v2 v2.4.0
)

replace github.com/libgit2/git2go/v34 => ./ext_deps/git2go

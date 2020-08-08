module fenix

go 1.14

require (
	github.com/clagraff/argparse v1.0.1 // indirect
	github.com/coreos/etcd v3.3.22+incompatible // indirect
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/go-systemd v0.0.0-20191104093116-d3cd4ed1dbcf // indirect
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f // indirect
	github.com/galdor/go-cmdline v1.1.1 // indirect
	github.com/gogo/protobuf v1.3.1 // indirect
	github.com/google/uuid v1.1.1
	github.com/julienschmidt/httprouter v1.3.0
	github.com/mattn/go-runewidth v0.0.9 // indirect
	github.com/mitchellh/mapstructure v1.3.2 // indirect
	github.com/nsf/termbox-go v0.0.0-20200418040025-38ba6e5628f1 // indirect
	go.etcd.io/etcd v3.3.22+incompatible
	go.uber.org/zap v1.15.0 // indirect
	golang.org/x/crypto v0.0.0-20200709230013-948cd5f35899
	google.golang.org/grpc v1.30.0 // indirect
)

replace (
	github.com/coreos/etcd => github.com/ozonru/etcd v3.3.20-grpc1.27-origmodule+incompatible
	google.golang.org/grpc => google.golang.org/grpc v1.27.0
)

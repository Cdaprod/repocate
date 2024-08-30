module github.com/cdaprod/repocate

go 1.19

require (
    // github.com/cdaprod/registry-service v0.1.0 
	replace github.com/Cdaprod/registry-service => ../registry-service
)

require (
	github.com/cheggaaa/pb/v3 v3.1.5
	github.com/docker/docker v23.0.3+incompatible
	github.com/fatih/color v1.17.0
	github.com/spf13/cobra v1.2.1
)

require (
	github.com/Microsoft/go-winio v0.6.2 // indirect
	github.com/VividCortex/ewma v1.2.0 // indirect
	github.com/docker/distribution v2.7.1+incompatible // indirect
	github.com/docker/go-connections v0.5.0 // indirect
	github.com/docker/go-units v0.5.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.15 // indirect
	github.com/moby/term v0.5.0 // indirect
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.1.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/sys v0.22.0 // indirect
	golang.org/x/time v0.0.0-20220210224613-90d013bbcef8 // indirect
	gotest.tools/v3 v3.5.1 // indirect
)

.PHONY: all test clean zip mac docker

### バージョンの定義
VERSION     := "v2.0.0"
COMMIT      := $(shell git rev-parse --short HEAD)
WD          := $(shell pwd)
### コマンドの定義
GO          = go
GO_BUILD    = $(GO) build
GO_TEST     = $(GO) test -v
GO_LDFLAGS  = -ldflags="-s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT)"
ZIP          = zip

### ターゲットパラメータ
DIST = dist
SRC = ./main.go ./sdrpower.go ./syslog.go ./monitor.go
TARGETS_MAC     =  $(DIST)/twSdrPower.darwin
TARGETS_LINUX     =  $(DIST)/twSdrPower $(DIST)/twSdrPower.arm $(DIST)/twSdrPower.arm64
GO_PKGROOT  = ./...

### PHONY ターゲットのビルドルール
all: $(TARGETS_LINUX) $(TARGETS_MAC)
test:
	env GOOS=$(GOOS) $(GO_TEST) $(GO_PKGROOT)
clean:
	rm -rf $(TARGETS_MAC) $(TARGETS_LINUX) $(DIST)/*.zip
mac: $(TARGETS_MAC)
zip: $(TARGETS_MAC) $(TARGETS_LINUX)
	cd dist && $(ZIP) twSdrPower_mac.zip twSdrPower.darwin
	cd dist && $(ZIP) twSdrPower_linux.zip twSdrPower twSdrPower.arm*

### 実行ファイルのビルドルール
$(DIST)/twSdrPower.darwin: $(SRC)
	env CGO_ENABLED=1 $(GO_BUILD) $(GO_LDFLAGS) -o $@
$(DIST)/twSdrPower.arm64: $(SRC)
	docker run --rm -v "$(WD)":/twSdrPower -w /twSdrPower golang:1.25 /twSdrPower/mkarm64.sh $(DIST) $(VERSION) $(COMMIT)
$(DIST)/twSdrPower.arm: $(SRC)
	docker run --rm -v "$(WD)":/twSdrPower -w /twSdrPower golang:1.25 /twSdrPower/mkarm.sh $(DIST) $(VERSION) $(COMMIT)
$(DIST)/twSdrPower: $(SRC)
	docker run --rm -v "$(WD)":/twSdrPower -w /twSdrPower golang:1.25 /twSdrPower/mklinux.sh $(DIST) $(VERSION) $(COMMIT)

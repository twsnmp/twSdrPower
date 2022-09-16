.PHONY: all test clean zip mac docker

### バージョンの定義
VERSION     := "v1.0.0"
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
TARGETS     =  $(DIST)/twSdrPower.app $(DIST)/twSdrPower $(DIST)/twSdrPower.arm64
GO_PKGROOT  = ./...

### PHONY ターゲットのビルドルール
all: $(TARGETS)
test:
	env GOOS=$(GOOS) $(GO_TEST) $(GO_PKGROOT)
clean:
	rm -rf $(TARGETS) $(DIST)/*.zip
mac: $(DIST)/twSdrPower.app
zip: $(TARGETS)
	cd dist && $(ZIP) twSdrPower_mac.zip twSdrPower.app
	cd dist && $(ZIP) twSdrPower_linux_amd64.zip twSdrPower
	cd dist && $(ZIP) twSdrPower_linux_arm64.zip twSdrPower.arm64

### 実行ファイルのビルドルール
$(DIST)/twSdrPower.app: $(SRC)
	env GO111MODULE=on GOOS=darwin GOARCH=amd64 $(GO_BUILD) $(GO_LDFLAGS) -o $@
$(DIST)/twSdrPower.arm64: $(SRC)
	docker run --rm -v "$(WD)":/twSdrPower -w /twSdrPower golang:1.17 /twSdrPower/mkarm64.sh $(DIST) $(VERSION) $(COMMIT)
$(DIST)/twSdrPower: $(SRC)
	docker run --rm -v "$(WD)":/twSdrPower -w /twSdrPower golang:1.17 /twSdrPower/mklinux.sh $(DIST) $(VERSION) $(COMMIT)

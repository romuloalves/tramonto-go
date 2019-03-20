GO ?= go

test:
	$(GO) test -cover ./...
.PHONY: test

android:
	@mkdir -p dist/

	@rm -f ./dist/Tramonto-sources.jar ./dist/Tramonto.aar

	@gomobile bind -o ./dist/Tramonto.aar -target=android gitlab.com/romuloalves/go-ipfs-react-native
.PHONY: android

ios:
	@mkdir -p dist/

	@rm -rf ./dist/Tramonto.framework

	@gomobile bind -o ./dist/Tramonto.framework -target=ios gitlab.com/romuloalves/go-ipfs-react-native
.PHONY:ios

todo:
	@rg TODO:
.PHONY: todo
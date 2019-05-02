GO ?= go

test:
	$(GO) test -cover ./...
.PHONY: test

android:
	@mkdir -p dist/

	@rm -f ./dist/Tramonto-sources.jar ./dist/Tramonto.aar

	GIN_MODE=release @gomobile bind -o ./dist/Tramonto.aar -target=android gitlab.com/tramonto-one/go-tramonto/tramonto
.PHONY: android

ios:
	@mkdir -p dist/

	@rm -rf ./dist/Tramonto.framework

	GIN_MODE=release @gomobile bind -o ./dist/Tramonto.framework -target=ios gitlab.com/tramonto-one/go-tramonto/tramonto
.PHONY:ios

todo:
	@rg TODO:
.PHONY: todo
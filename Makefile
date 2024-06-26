.EXPORT_ALL_VARIABLES:
	GIN_MODE=release

test:
	@go test -cover ./...
.PHONY: test

android:
	@mkdir -p dist/

	@rm -f ./dist/Tramonto-sources.jar ./dist/Tramonto.aar

	@gomobile bind -o ./dist/Tramonto.aar -target=android gitlab.com/tramonto-one/go-tramonto/tramonto
.PHONY: android

ios:
	@mkdir -p dist/

	@rm -rf ./dist/Tramonto.framework

	@gomobile bind -o ./dist/Tramonto.framework -target=ios gitlab.com/tramonto-one/go-tramonto/tramonto
.PHONY:ios

todo:
	@rg TODO:
.PHONY: todo

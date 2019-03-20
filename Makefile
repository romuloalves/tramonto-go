test:
	@go test -cover ./...
.PHONY: test

android:
	@gomobile bind -o Tramonto.aar -target=android gitlab.com/romuloalves/go-ipfs-react-native
.PHONY: android

ios:
	@gomobile bind -o Tramonto.framework -target=ios gitlab.com/romuloalves/go-ipfs-react-native
.PHONY:ios
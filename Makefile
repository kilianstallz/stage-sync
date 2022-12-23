
commit:
	@npx git-cz

release:
	@sh ./scripts/release.sh

test:
	@go test ./...

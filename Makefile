configure-pre-commit:
	@pip install pre-commit
	@pre-commit install &&  pre-commit run --all-files

lint:
	@gofumpt -l -w .
check-lint:
	@golangci-lint run ./... -v --timeout 3m --config .golangci.yml

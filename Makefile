.PHONY: test integrate

test:
	go test -mod=vendor -tags="unit" .

integration:
	go test -mod=vendor -tags="integration" .

# makefile
protoc:
	@echo "Генерация GRPC кода..."
	cd internal/delivery/grpc && \
	protoc -I /home/bloomguy/GolandProjects/Ozon-Internship/urlshortener/internal/delivery/grpc --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. /home/bloomguy/GolandProjects/Ozon-Internship/urlshortener/internal/delivery/grpc/urlshortener.proto

test:
	go test -v ./internal_test/...

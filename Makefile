# Название пакета 
PACKAGE="main"
# Генерация кода 
generate: 
	protoc -Iapi --go_opt=module=$(PACKAGE) --go_out=. \
	--go-grpc_opt=module=$(PACKAGE) --go-grpc_out=. \
	api/*.proto
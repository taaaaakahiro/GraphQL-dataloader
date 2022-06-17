dck=docker-compose

conf:
	${dck} config
up:
	${dck} up --build -d
down:
	${dck} down -v
res:
	${dck} restart
sh-mysql:
	${dck} exec mysql8 sh
test:
	@echo "Test Start"
	go test ./pkg/... -v -count=1


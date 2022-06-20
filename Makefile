dck=docker-compose

conf:
	${dck} config

up:
	${dck} up --build -d

down:
	${dck} down -v

re:
	${dck} down -v
	${dck} up --build -d

sh-mysql:
	${dck} exec mysql8 sh

test:
	@echo "★★★ Test Start ★★★"
	go test ./pkg/... -count=1


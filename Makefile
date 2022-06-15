dk=docker-compose

conf:
	${dk} config
up:
	${dk} up --build -d
down:
	${dk} down -v
res:
	${dk} restart
sh-mysql:
	${dk} exec mysql8 sh 

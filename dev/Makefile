export UID := $(shell id -u)

up: ${HOME}/.local/share/serviceweaver
	docker compose up -d

${HOME}/.local/share/serviceweaver:
	mkdir -p ${HOME}/.local/share/serviceweaver
	chown -R ${USER} ${HOME}/.local/share/serviceweaver

down:
	docker compose down

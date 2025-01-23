
# make release version=v0.1.0
release:
	docker build -f Dockerfile_logistic_package_api -t arslanovdi/logistic_package_api:$(version) .
	docker build -f Dockerfile_events -t arslanovdi/events:$(version) .
	docker build -f Dockerfile_retranslator -t arslanovdi/retranslator:$(version) .
	docker build -f Dockerfile_tgbot -t arslanovdi/tgbot:$(version) .

# make pull version=v0.1.0
pull:
	docker push arslanovdi/logistic_package_api:$(version)
	docker push arslanovdi/events:$(version)
	docker push arslanovdi/retranslator:$(version)
	docker push arslanovdi/tgbot:$(version)

configs:
	kubectl create configmap events --from-file=events/config.yml
	kubectl create configmap logistic-package-api --from-file=logistic-package-api/config.yml
	kubectl create configmap tgbot --from-file=telegram_bot/config.yml
	kubectl create configmap postgres --from-file=logistic-package-api/scripts/init-database.sh

secrets:
	kubectl create secret generic tgbot --from-literal=TOKEN=7012140868:AAHAkiK606qFalhnX7Cm3d8aDRTIw5m3NWw
	kubectl create secret generic logistic-package-api --from-literal=PASSWORD=P@$$w0rd
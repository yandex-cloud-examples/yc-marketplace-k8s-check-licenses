build:
	docker build -t sample-check-license .

run:
	docker run -p 8080:8080 sample-check-license

invoke:
	curl localhost:8080

run-kube:
	kubectl apply -f kube-sample-check-license-pod.yaml

remove-kube:
	kubectl delete -f kube-sample-check-license-pod.yaml

tag-to-registry:
	docker tag sample-check-license cr.yandex/<your registry id>/sample-check-license:1.0.0

push-to-registry:
	docker push cr.yandex/<your registry id>/sample-check-license:1.0.0

install: build tag-to-registry push-to-registry run-kube

all:
	echo Nothing

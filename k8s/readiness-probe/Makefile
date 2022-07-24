KUBE = kubectl -n sandbox
VARIANT = recover

default:
	@echo ok

image:
	docker build -t local/k8s-sandbox/readiness-probe:dev .

apply:
	$(KUBE) apply -k variants/$(VARIANT)

delete:
	$(KUBE) delete -k variants/$(VARIANT)

watch:
	$(KUBE) get po -w

describe:
	$(KUBE) get po -w

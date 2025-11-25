NAMESPACE=social-network

build-base:
	docker build -t social-network/go-base -f docker/go/base.Dockerfile .

apply-namespace:
	kubectl apply -f k8s/namespace.yaml

apply-db:
	kubectl apply -f k8s/users-db/pvc.yaml
	kubectl apply -f k8s/users-db/deployment.yaml
	kubectl apply -f k8s/users-db/service.yaml

build-users:
	docker build -t social-network/users:dev -f services/users/Dockerfile .

push-users:
	# colima uses local docker -- no push needed
	echo "Local image loaded."

deploy-users:
	kubectl apply -f k8s/users/deployment.yaml
	kubectl apply -f k8s/users/service.yaml

run-migrations:
	kubectl delete job users-migrate -n $(NAMESPACE) --ignore-not-found=true
	kubectl apply -f k8s/users-db/migration-job.yaml

logs-users:
	kubectl logs -l app=users -n $(NAMESPACE) -f

logs-db:
	kubectl logs -l app=users-db -n $(NAMESPACE) -f

all: apply-namespace apply-db build-users deploy-users run-migrations

reset:
	kubectl delete namespace $(NAMESPACE) --ignore-not-found=true
	kubectl create namespace $(NAMESPACE)

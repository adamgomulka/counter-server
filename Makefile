GCP_PROJECT_ID := nyt-adam-gomulka-interview
GCP_COMPUTE_REGION := us-east4
GCP_KEY_FILE := nyt-adam-gomulka-interview-6b516e0cfe83.json
GKE_CLUSTER_NAME := nyt-interview-server
PATH := $(PATH):$(PWD)/tools:$(PWD)/tools/google-cloud-sdk/bin
KUBECONFIG := $(PWD)/.kubeconfig

SHELL := /bin/bash

VERSION := 1.1.0
BUILD := `git rev-parse HEAD`

.DEFAULT_GOAL := all

all: get_utils auth build create config_kubernetes deploy
all_nodl: auth build create config_kubernetes deploy

get_utils:
	wget https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/google-cloud-sdk-241.0.0-linux-x86_64.tar.gz && tar -xvf google-cloud-sdk-241.0.0-linux-x86_64.tar.gz -C tools/ && rm google-cloud-sdk-241.0.0-linux-x86_64.tar.gz
	wget https://releases.hashicorp.com/terraform/0.11.13/terraform_0.11.13_linux_amd64.zip && unzip terraform_0.11.13_linux_amd64.zip -d tools/ && rm terraform_0.11.13_linux_amd64.zip
	wget -O tools/kubectl https://storage.googleapis.com/kubernetes-release/release/v1.14.0/bin/linux/amd64/kubectl && chmod +x tools/kubectl
	wget https://storage.googleapis.com/kubernetes-helm/helm-v2.13.1-linux-amd64.tar.gz && tar -xvf helm-v2.13.1-linux-amd64.tar.gz -C tools/ linux-amd64/helm linux-amd64/tiller --strip=1 && rm helm-v2.13.1-linux-amd64.tar.gz 
auth:
	gcloud auth activate-service-account --key-file $(GCP_KEY_FILE)
	gcloud auth configure-docker

build:
	docker build -t gcr.io/$(GCP_PROJECT_ID)/counter-server:$(VERSION) -t gcr.io/$(GCP_PROJECT_ID)/counter-server:$(BUILD) .
	docker push gcr.io/$(GCP_PROJECT_ID)/counter-server:$(VERSION)
	docker push gcr.io/$(GCP_PROJECT_ID)/counter-server:$(BUILD)

create:
	terraform init tf/
	terraform apply -auto-approve tf/ 

config_kubernetes:
	gcloud container clusters get-credentials $(GKE_CLUSTER_NAME) --region $(GCP_COMPUTE_REGION)
	kubectl --namespace kube-system create sa tiller
	kubectl create clusterrolebinding tiller --clusterrole cluster-admin --serviceaccount=kube-system:tiller
	helm init --service-account tiller
	sleep 30

deploy:
	helm install --name nyt-interview-server --namespace nyt-interview-server chart/counter-server

clean:
	rm -rf tools/*
	rm terraform.*
	rm .kubeconfig

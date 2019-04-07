GCP_PROJECT_ID := nyt-adam-gomulka-interview
GCP_COMPUTE_REGION := us-east4
GCP_KEY_FILE := nyt-adam-gomulka-interview-6b516e0cfe83.json
GKE_CLUSTER_NAME := nyt-interview-server
PATH := $(PATH):$(PWD)/tools:$(PWD)/tools/google-cloud-sdk/bin
KUBECONFIG := $(PWD)/.kubeconfig

SHELL := /bin/bash

VERSION := 1.0.0
BUILD := `git rev-parse HEAD`

all: get_gcloud_sdk auth build create_gke config_kubernetes deploy

get_gcloud_sdk:
	wget https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/google-cloud-sdk-241.0.0-linux-x86_64.tar.gz && tar -xvf google-cloud-sdk-241.0.0-linux-x86_64.tar.gz -C tools/

auth:
	gcloud auth activate-service-account --key-file $(GCP_KEY_FILE)
	gcloud auth configure-docker

build:
	docker build -t gcr.io/$(GCP_PROJECT_ID)/nyt-interview-server:$(VERSION) -t gcr.io/$(GCP_PROJECT_ID)/nyt-interview-server:$(BUILD) .
	docker push gcr.io/$(GCP_PROJECT_ID)/nyt-interview-server:$(VERSION)
	docker push gcr.io/$(GCP_PROJECT_ID)/nyt-interview-server:$(BUILD)

create_gke:
	terraform init tf/
	terraform apply -auto-approve tf/ 

config_kubernetes:
	gcloud container clusters get-credentials $(GKE_CLUSTER_NAME) --region $(GCP_COMPUTE_REGION)
	kubectl --namespace kube-system create sa tiller
	kubectl create clusterrolebinding tiller --clusterrole cluster-admin --serviceaccount=kube-system:tiller
	helm init --service-account tiller
	sleep 30

deploy:
	helm install --name nyt-interview-server --namespace nyt-interview-server chart/nyt-interview-server

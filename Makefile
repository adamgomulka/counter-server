GCP_PROJECT_ID := nyt-adam-gomulka-interview
GCP_COMPUTE_ZONE := us-east4-a
GCP_KEY_FILE := nyt-adam-gomulka-interview-6b516e0cfe83.json
GKE_CLUSTER_NAME := nyt-server

SHELL := /bin/bash

VERSION := 1.0.0
BUILD := `git rev-parse HEAD`

all: create deploy

auth:
	gcloud auth activate-service-account --key-file $(GCP_KEY_FILE)
	gcloud auth configure-docker

build:
	docker build -t gcr.io/$(GCP_PROJECT_ID)/nyt-server:$(VERSION) -t gcr.io/$(GCP_PROJECT_ID)/nyt-server:$(BUILD) .
	docker push gcr.io/$(GCP_PROJECT_ID)/nyt-server:$(VERSION)
	docker push gcr.io/$(GCP_PROJECT_ID)/nyt-server:$(BUILD)

create: 
	gcloud auth activate-service-account --key-file $(GCP_KEY_FILE)
	gcloud config set project $(GCP_PROJECT_ID)
	gcloud config set compute/zone $(GCP_COMPUTE_ZONE)
	gcloud components update
	gcloud container clusters create $(GKE_CLUSER_NAME)

deploy:
	kubectl run nyt-server --image=gcr.io/$(GCP_PROJECT_ID)/nyt-server:$(VERSION) --port 8080
	kubectl expose deployment nyt-server --type=LoadBalancer --port 80 --target-port 8080

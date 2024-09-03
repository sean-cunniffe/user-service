#!/bin/bash
def createNamespace() {
    echo "Creating namespace testing"
    kubectl create namespace testing
}

def createDockerPullSecret() {   
    echo "Creating secret regcred"
    kubectl create secret generic regcred --from-file=~/.docker/config.json --type=kubernetes.io/dockerconfigjson --namespace testing
}

def minikubeStart() {
    echo "Starting minikube"
    minikube start --driver=docker
}

def shareTestImage() {
    echo "Sharing test image"
    docker build -f docker/Dockerfile -t seancunniffe/user-service:testing .
    minikube image load seancunniffe/user-service:testing
}

def installWithTestImage() {
    echo "Installing user-service with test image"
    helm install user-service charts/user-service --set image.tag=testing --namespace testing
}
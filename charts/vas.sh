helm repo add jetstack https://charts.jetstack.io
helm repo update

kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.15.3/cert-manager.crds.yaml

echo 'installing integration chart'

cd integration/
helm uninstall integration
helm install integration .
cd -
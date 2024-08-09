helm repo add jetstack https://charts.jetstack.io
helm repo update

helm uninstall cert-manager -n cert-manager

helm install cert-manager jetstack/cert-manager \
  --namespace cert-manager \
  --create-namespace \
  --version v1.12.2 \
  --set installCRDs=true

cd integration/
helm uninstall integration
helm install integration .
cd -
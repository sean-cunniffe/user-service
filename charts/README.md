# README

This README provides instructions for running `vas.sh` to install Cert Manager and the Integration Chart.

## Prerequisites
Before proceeding, ensure that you have the following prerequisites installed:

- Kubernetes cluster
- Helm

## Installation Steps
Follow these steps to install Cert Manager and the Integration Chart:


1. Change directory to the charts folder:

    ```shell
    cd charts
    ```

2. Run `vas.sh` script:

    ```shell
    ./vas.sh
    ```

3. Wait for the installation to complete.

4. Verify the installation:

    ```shell
    kubectl get pods
    ```

    Ensure that all the required pods are running.

Congratulations! You have successfully installed Cert Manager and the Integration Chart.

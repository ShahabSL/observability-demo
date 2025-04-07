# Go + Docker + Kubernetes Learning Demo

This project is a simple demonstration created primarily for learning the basics of:

*   Go (Golang) application development (specifically exposing metrics)
*   Dockerizing an application using multi-stage builds
*   Deploying applications to Kubernetes using Minikube
*   Using Helm to install complex applications (Prometheus & Grafana)
*   Setting up basic observability with Prometheus and Grafana in Kubernetes

**Disclaimer:** This is purely a learning exercise and not intended for production use.

## Technologies Used

*   Go (Golang)
*   Docker / Docker Desktop
*   Minikube
*   kubectl
*   Helm
*   Prometheus Operator (via kube-prometheus-stack Helm chart)
*   Grafana (via kube-prometheus-stack Helm chart)

## Prerequisites

*   Go toolchain installed ([https://go.dev/doc/install](https://go.dev/doc/install))
*   Docker Desktop installed and running ([https://www.docker.com/products/docker-desktop/](https://www.docker.com/products/docker-desktop/))
*   Minikube installed ([https://minikube.sigs.k8s.io/docs/start/](https://minikube.sigs.k8s.io/docs/start/))
*   kubectl installed (usually comes with Docker Desktop or can be installed separately)
*   Helm installed ([https://helm.sh/docs/intro/install/](https://helm.sh/docs/intro/install/))

## Setup and Running

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/ShahabSL/observability-demo.git
    cd observability-demo/app
    ```

2.  **Start Minikube:**
    ```bash
    minikube start
    ```

3.  **Build the Go application image into Minikube's Docker environment:**
    (This makes the image available locally within the cluster, necessary for `imagePullPolicy: Never`)
    ```bash
    minikube image build -t observability-demo:latest .
    ```

4.  **Deploy the Go application to Kubernetes:**
    ```bash
    kubectl apply -f deployment.yaml
    ```
    (This creates the Deployment and ClusterIP Service for the app in the `default` namespace)

5.  **Set up the Monitoring Stack (Prometheus & Grafana):**
    *   Create a dedicated namespace:
        ```bash
        kubectl create namespace monitoring
        ```
    *   Add the Prometheus community Helm repository:
        ```bash
        helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
        helm repo update # Optional, but good practice
        ```
    *   Install the kube-prometheus-stack chart:
        ```bash
        helm install prometheus prometheus-community/kube-prometheus-stack --namespace monitoring
        ```
        (Wait for all pods in the `monitoring` namespace to be running: `kubectl get pods -n monitoring`)

6.  **Configure Prometheus to scrape the Go application:**
    *   Ensure the Service definition in `deployment.yaml` has the port named `http` (required by `service-monitor.yaml`). Example:
        ```yaml
        # In deployment.yaml -> Service spec -> ports:
        - name: http # <-- Name required by ServiceMonitor
          port: 80
          targetPort: 8080
        ```
    *   Apply the ServiceMonitor custom resource:
        ```bash
        kubectl apply -f service-monitor.yaml
        ```
        (This tells the Prometheus Operator to find the `observability-demo` service via its label and scrape metrics from the `/metrics` endpoint on the port named `http`).

## Accessing Services

*   **Go Application:**
    Since it's a ClusterIP service, use `port-forward` or `minikube service`:
    ```bash
    # Option 1: Port Forward
    kubectl port-forward service/observability-demo 8080:80
    # Access via http://localhost:8080 in your browser

    # Option 2: Minikube Service (will open in browser or give URL)
    minikube service observability-demo
    ```

*   **Prometheus UI:**
    ```bash
    kubectl port-forward -n monitoring service/prometheus-kube-prometheus-prometheus 9090:9090
    # Access via http://localhost:9090
    ```
    You can check the `/targets` endpoint in the UI to see if it's successfully scraping your `observability-demo` service.

*   **Grafana UI:**
    *   Get the admin password:
        ```powershell
        # Windows PowerShell:
        kubectl --namespace monitoring get secrets prometheus-grafana -o jsonpath='{.data.admin-password}' | ForEach-Object { [System.Text.Encoding]::UTF8.GetString([System.Convert]::FromBase64String($_)) }

        # Linux/macOS/WSL:
        # kubectl --namespace monitoring get secrets prometheus-grafana -o jsonpath="{.data.admin-password}" | base64 --decode ; echo
        ```
    *   Port forward to Grafana:
        ```bash
        kubectl port-forward -n monitoring service/prometheus-grafana 3000:3000
        ```
    *   Access via `http://localhost:3000`
    *   Login with username `admin` and the password obtained above. The Prometheus data source should be pre-configured.

## Key Files

*   `main.go`: Simple Go HTTP server exposing Prometheus metrics at `/metrics`.
*   `Dockerfile`: Multi-stage Docker build for creating a minimal Go application image.
*   `deployment.yaml`: Kubernetes manifests for deploying the Go app (Deployment) and exposing it internally (Service).
*   `service-monitor.yaml`: Prometheus Operator custom resource to enable automatic scraping of the Go application's metrics.
*   `docker-compose.yml` / `prometheus.yml`: Used for initial Docker-only setup (less relevant for the Kubernetes part).

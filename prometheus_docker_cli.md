1. Run prometheus
    ```bash
    docker run \
        --name prometheus \
        -p 9090:9090 --rm \
        --network=host \
        -d -v $(pwd)/prometheus.yml:/etc/prometheus/prometheus.yml \
        prom/prometheus
    ```

2. Run Grafana
    ```bash
    docker run -d --name=grafana --net=host grafana/grafana:5.0.0
    ```
3. Run Scraper for Postgres [Reference](https://github.com/prometheus-community/postgres_exporter)
    - Replace IP address with postgress IP or launch inside kubernetes only with dns resolution

4. Create ConfigMap before kube deployment
    ```bash
    kubectl create configmap prom-config --from-file=./prometheus.yml
    ```
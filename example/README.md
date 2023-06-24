# Prometheus monitoring for shmipc-go

## Usage Example

### Server

See [server](./shmipc_server)

### Client

See [client](./shmipc_client)

## HOW-TO-RUN

1. install docker and start docker
2. change `{{ INET_IP }}` to local ip in prometheus.yml
3. run Prometheus and Grafana  
   `docker-compose up`
4. run shmipc_client and shmipc_server   
   `sh run_shmipc_client_server.sh`
5. visit `http://localhost:3000`, the account password is `admin` by default
6. configure Prometheus data sources
    1. `Configuration`
    2. `Data Source`
    3. `Add data source`
    4. Select `Prometheus` and fill the URL with `http://prometheus:9090`
    5. click `Save & Test` after configuration to test if it works
7. add dashboard `Create` -> `dashboard`, add monitoring metrics such as throughput and pct99 according to your needs,
   for example:

    - shmipc-client all in used share memory in bytes

   `all_in_used_share_memory_in_bytes{job="shmipc-client"}`

    - shmipc-server active stream count

   `active_stream_count{job="shmipc-server"}`

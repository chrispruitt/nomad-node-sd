# nomad-node-sd

Simple service to update a prometheus service discovery file for nomad server and client nodes.

## Docker

Run the below command and view the results of nomad-nodes.yaml in the current directory.

```bash
docker run -it \
	-e NOMAD_API_URL="http://127.0.0.1:4646" \
	-e OUTPUT_FILE_PATH="/tmp/nomad-nodes.yaml" \
	-v "$(pwd):/tmp" \
	chrispruitt/nomad-node-sd
```

## Environment Variables

| ENV VAR          | Description                                                                           | Default            |
| --------         | -------                                                                               | --------           |
| NOMAD_API_URL    | The base URL of nomad - "http://localhost:4646"                                       |                    |
| OUTPUT_FILE_PATH | The full path of file_sd_config                                                       | "nomad-nodes.yaml" |
| REFRESH_INTERVAL | The refresh interval in seconds                                                       | 300                |
| NOMAD_NODE_PORT  | The port on which each node is running nomad. This will be used to build the targets. | 4646               |


## Putting it all together

Example prometheus scrape configs (assuming you are writing running the service above to write to `/etc/prometheus/sd-configs/nomad-nodes.yaml`)

```yaml
scrape_configs:
  - job_name: 'nomad'
    metrics_path: /v1/metrics
    params:
      format: ['prometheus']
    file_sd_configs:
      - files:
        - '/etc/prometheus/sd-configs/nomad-nodes.yaml'
        refresh_interval: 5m
```

# blockspider

An open source blockchain crawler and kafka producer.

### Requires

- go 1.20.x or greater
- Apache Kafka 3.6.x or greater
- a geth-like rpc endpoint

### Get the source

```shell
git clone https://github.com/iquidus/blockspider blockspider
```

### build

```shell
cd blockspider && make spiderd
```

### Configure

```shell
cp ./config.json.example ./config.json
```

_Make required changes in config.json_

```json
{
  "chainId": 8, // chainId of target network
  "crawler": {
    // crawler settings
    "start": 0, // start block
    "interval": "10000ms", // polling interval. e.g 0.5 * target block time
    "routines": 1, // go routines
    "cache": 128, // number of blocks to keep in local cache. Must be larger than reorgs.
    "kafka": {
      "events": [
        {
          "broker": "localhost:9092",
          "topic": "events",
          "addresses": [],
          "topics": []
        }
      ],
      "blocks": {
        "broker": "localhost:9092",
        "topic": "blocks"
      }
    }
  },
  "rpc": {
    "type": "http",
    "endpoint": "http://127.0.0.1:8588"
  },
  "state": {
    "path": "~/.blockspider/ubiq-mainnet.json"
  }
}
```

### Run

```shell
./build/bin spiderd -c config.json
```

### Kafka

[Download](https://www.apache.org/dyn/closer.cgi?path=/kafka/3.6.0/kafka_2.13-3.6.0.tgz) the latest Kafka release and extract it

```shell
tar -xzf kafka_2.13-3.6.0.tgz
cd kafka_2.13-3.6.0
```

Run the following commands in order to start all services in the correct order

```shell
# Start the ZooKeeper service
bin/zookeeper-server-start.sh config/zookeeper.properties
```

Open another terminal session and run

```shell
# Start the Kafka broker service
bin/kafka-server-start.sh config/server.properties
```

Once all services have successfully launched, you will have a basic Kafka environment running and ready to use.

Create topics

```shell
bin/kafka-topics.sh --create --topic blocks --bootstrap-server localhost:9092
bin/kafka-topics.sh --create --topic events --bootstrap-server localhost:9092
```

## Install
```
brew install go  # install golang
cd $HOME
git clone https://github.com/luozhaoyu/GoHttpsHashSample.git
docker-compose up
./test.py  # python lib requests is required
```

## TODOs
* How would your implementation scale if this were a high throughput service, and how could you improve that?
  * Is read or write throughput high? What's the consistency requirement?
    * To scale write: the only way is to shard its service by messages
    * To scale read: we can consider put move slaves for read operations
      * if it needs strong consistency, then we can apply quorums on this replicaset and tune its R W quorum accordingly
      * otherwise, we can set several servers to serve the read operations for the same shards given replication lag is acceptable; we can also setup a bunch of memcache in front of it
  * Does the service need persistent in case crashed? We can add memory snapshot for this
* How would you deploy this to the cloud, using the provider of your choosing? What would the architecture look like? What tools would you use?
  * Terraform + AWS autoscaling
  * There are 2 ways for sharding: consistent hash and normal manual sharding. Let's say we are using consistent hash
    * all traffic will go to a Zookeeper/etcd like service maintains the consistent hash ring and will route traffic to the designated shard group
    * assume we need strong consistency, then each shard group will use quorum based strategy that, e.g., say quorum is 3, we set the quorum-read as 2 and quorum-write as 2. (this quorum could tune according to real traffic pattern)
* How would you monitor this service? What metrics would you collect? How would you act on these metrics?
  * Graphite/Datadog/Prometheus/Uber-M3
  * Metrics are: num_requests, latency, errors, CPU/Memory/Network IO
    * more granularity per API
    * more granualrity by datacenter/host
  * It depends on our SLA. Only alarm when SLA violated such as latency, error regression. The others should tune to be warning level
    * We could add post-alert hook to automatically scale more hosts if num_requests increased

## Keys
generate private key
`openssl genrsa -out server.key 2048`

generate public key
`openssl req -new -x509 -sha256 -key server.key -out localhost.crt -days 3650`

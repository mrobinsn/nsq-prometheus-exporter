# nsq-prometheus-exporter

Forked from https://github.com/caozhipan/nsq-prometheus-exporter

## Build
```bash
go build -o nsq-prometheus-exporter main.go
```

## Docker

- Build & Push
```bash
docker build -t nsq-prometheus-exporter .
docker tag nsq-prometheus-exporter mrobinsn/nsq-prometheus-exporter:$VERSION
docker push mrobinsn/nsq-prometheus-exporter:$VERSION
```

- Run
```bash
docker run -p 9527:9527 mrobinsn/nsq-prometheus-exporter:$VERSION -nsq.lookupd.address=192.168.31.1:4161,192.168.31.2:4161
```


# 1. 示例应用
先部署一个使用 SpringBoot 和 MongoDB 开发的示例应用。首先部署一个 MongoDB 应用
```
$ kubectl apply -f mongo.yml
service/mongo created
statefulset.apps/mongo created
$ kubectl get pods -n elastic -l app=mongo             
NAME      READY   STATUS    RESTARTS   AGE
mongo-0   1/1     Running   0          34m
```
接下来部署 SpringBoot 的 API 应用，这里我们通过 NodePort 类型的 Service 服务来暴露该服务
```
$ kubectl apply -f spring-boot-simple.yaml 
service/spring-boot-simple created
deployment.apps/spring-boot-simple created
$ kubectl get pods -n elastic -l app=spring-boot-simple
NAME                                  READY   STATUS    RESTARTS   AGE
spring-boot-simple-64795494bf-hqpcj   1/1     Running   0          24m
$ kubectl get svc -n elastic -l app=spring-boot-simple
NAME                 TYPE       CLUSTER-IP      EXTERNAL-IP   PORT(S)          AGE
spring-boot-simple   NodePort   10.109.55.134   <none>        8080:31847/TCP   84s
```
应用部署完成后，我们就可以通过地址 http://:31847 访问应用
```
$ curl -X GET  http://*******:31847/
Greetings from Spring Boot!
```
发送一个 POST 请求：
```
$ curl -X POST http://*******:31847/message -d 'hello world'
{"id":"5ef55c130d53190001bf74d2","message":"hello+world=","postedAt":"2020-06-26T02:23:15.860+0000"}
```

# 2. ElasticSearch集群
要建立一个 Elastic 技术的监控栈，当然首先我们需要部署 ElasticSearch，它是用来存储所有的指标、日志和追踪的数据库，这里我们通过3个不同角色的可扩展的节点组成一个集群。

## 2.1. 安装 ElasticSearch 主节点
设置集群的第一个节点为 Master 主节点，来负责控制整个集群。首先创建一个 ConfigMap 对象，用来描述集群的一些配置信息，以方便将 ElasticSearch 的主节点配置到集群中并开启安全认证功能。

然后创建一个 Service 对象，在 Master 节点下，我们只需要通过用于集群通信的 9300 端口进行通信。

最后使用一个 Deployment 对象来定义 Master 节点应用。
```
$ kubectl apply  -f elasticsearch-master.configmap.yaml \
                 -f elasticsearch-master.service.yaml \
                 -f elasticsearch-master.deployment.yaml

configmap/elasticsearch-master-config created
service/elasticsearch-master created
deployment.apps/elasticsearch-master created
$ kubectl get pods -n elastic -l app=elasticsearch
NAME                                    READY   STATUS    RESTARTS   AGE
elasticsearch-master-6f666cbbd-r9vtx    1/1     Running   0          111m
```

## 2.2. 安装 ElasticSearch 数据节点
现在我们需要安装的是集群的数据节点，它主要来负责集群的数据托管和执行查询。和 master 节点一样，我们使用一个 ConfigMap 对象来配置我们的数据节点
``` 
elasticsearch-data.configmap.yaml
```
可以看到之前的 master 配置非常类似，不过需要注意的是属性 node.data=true。

同样只需要通过 9300 端口和其他节点进行通信：
```
elasticsearch-data.service.yaml
```
最后创建一个 StatefulSet 的控制器，因为可能会有多个数据节点，每一个节点的数据不是一样的，需要单独存储，所以也使用了一个 volumeClaimTemplates 来分别创建存储卷
```
# elasticsearch-data.statefulset.yaml
```

```
$ kubectl apply -f elasticsearch-data.configmap.yaml \
                -f elasticsearch-data.service.yaml \
                -f elasticsearch-data.statefulset.yaml

configmap/elasticsearch-data-config created
service/elasticsearch-data created
statefulset.apps/elasticsearch-data created

$ kubectl get pods -n elastic -l app=elasticsearch
NAME                                    READY   STATUS    RESTARTS   AGE
elasticsearch-data-0                    1/1     Running   0          90m
elasticsearch-master-6f666cbbd-r9vtx    1/1     Running   0          111m
```
## 2.3. 安装 ElasticSearch 客户端节点
同样使用一个 ConfigMap 对象来配置该节点：

客户端节点需要暴露两个端口，9300端口用于与集群的其他节点进行通信，9200端口用于 HTTP API。

使用一个 Deployment 对象来描述客户端节点：

```
$ kubectl apply  -f elasticsearch-client.configmap.yaml \
                 -f elasticsearch-client.service.yaml \
                 -f elasticsearch-client.deployment.yaml

configmap/elasticsearch-client-config created
service/elasticsearch-client created
deployment.apps/elasticsearch-client created

$ kubectl get pods -n elastic -l app=elasticsearch
NAME                                    READY   STATUS    RESTARTS   AGE
elasticsearch-client-788bffcc98-hh2s8   1/1     Running   0          83m
elasticsearch-data-0                    1/1     Running   0          91m
elasticsearch-master-6f666cbbd-r9vtx    1/1     Running   0          112m
```

可以通过如下所示的命令来查看集群的状态变化：
```
$ kubectl logs -f -n elastic \
  $(kubectl get pods -n elastic | grep elasticsearch-master | sed -n 1p | awk '{print $1}') \
  | grep "Cluster health status changed from"

{"type": "server", "timestamp": "2020-06-26T03:31:21,353Z", "level": "INFO", "component": "o.e.c.r.a.AllocationService", "cluster.name": "elasticsearch", "node.name": "elasticsearch-master", "message": "Cluster health status changed from [RED] to [GREEN] (reason: [shards started [[.monitoring-es-7-2020.06.26][0]]]).", "cluster.uuid": "SS_nyhNiTDSCE6gG7z-J4w", "node.id": "BdVScO9oQByBHR5rfw-KDA"  }
```

## 2.4. 生成密码
我们启用了 xpack 安全模块来保护我们的集群，所以我们需要一个初始化的密码。我们可以执行如下所示的命令，在客户端节点容器内运行 bin/elasticsearch-setup-passwords 命令来生成默认的用户名和密码：
```
$ kubectl exec $(kubectl get pods -n elastic | grep elasticsearch-client | sed -n 1p | awk '{print $1}') \
    -n elastic \
    -- bin/elasticsearch-setup-passwords auto -b

Changed password for user apm_system
PASSWORD apm_system = 3Lhx61s6woNLvoL5Bb7t

Changed password for user kibana_system
PASSWORD kibana_system = NpZv9Cvhq4roFCMzpja3

Changed password for user kibana
PASSWORD kibana = NpZv9Cvhq4roFCMzpja3

Changed password for user logstash_system
PASSWORD logstash_system = nNnGnwxu08xxbsiRGk2C

Changed password for user beats_system
PASSWORD beats_system = fen759y5qxyeJmqj6UPp

Changed password for user remote_monitoring_user
PASSWORD remote_monitoring_user = mCP77zjCATGmbcTFFgOX

Changed password for user elastic
PASSWORD elastic = wmxhvsJFeti2dSjbQEAH
```

注意需要将 elastic 用户名和密码也添加到 Kubernetes 的 Secret 对象中：
```
$ kubectl create secret generic elasticsearch-pw-elastic \
    -n elastic \
    --from-literal password=wmxhvsJFeti2dSjbQEAH
secret/elasticsearch-pw-elastic created
```

# 3. Kibana
ElasticSearch 集群安装完成后，接着我们可以来部署 Kibana，这是 ElasticSearch 的数据可视化工具，它提供了管理 ElasticSearch 集群和可视化数据的各种功能。

```
$ kubectl apply  -f kibana.configmap.yaml \
                 -f kibana.service.yaml \
                 -f kibana.deployment.yaml

configmap/kibana-config created
service/kibana created
deployment.apps/kibana created

$ kubectl logs -f -n elastic $(kubectl get pods -n elastic | grep kibana | sed -n 1p | awk '{print $1}') \
     | grep "Status changed from yellow to green"

{"type":"log","@timestamp":"2020-06-26T04:20:38Z","tags":["status","plugin:elasticsearch@7.8.0","info"],"pid":6,"state":"green","message":"Status changed from yellow to green - Ready","prevState":"yellow","prevMsg":"Waiting for Elasticsearch"}
```
当状态变成 green 后，我们就可以通过 NodePort 端口 30474 去浏览器中访问 Kibana 服务了
```
$ kubectl get svc kibana -n elastic   
NAME     TYPE       CLUSTER-IP      EXTERNAL-IP   PORT(S)          AGE
kibana   NodePort   10.101.121.31   <none>        5601:30474/TCP   8m18s
```
使用上面我们创建的 Secret 对象的 elastic 用户和生成的密码即可登录,到这里我们就安装成功了 ElasticSearch 与 Kibana，它们将为我们来存储和可视化我们的应用数据（监控指标、日志和追踪）服务。

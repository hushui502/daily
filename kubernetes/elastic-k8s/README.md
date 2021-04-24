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

# 4. kube-state-metrics
首先，我们需要安装 kube-state-metrics，这个组件是一个监听 Kubernetes API 的服务，可以暴露每个资源对象状态的相关指标数据。
```
$ git clone https://github.com/kubernetes/kube-state-metrics.git
$ cd kube-state-metrics
# 执行安装命令
$ kubectl apply -f examples/standard/  
clusterrolebinding.rbac.authorization.k8s.io/kube-state-metrics configured
clusterrole.rbac.authorization.k8s.io/kube-state-metrics configured
deployment.apps/kube-state-metrics configured
serviceaccount/kube-state-metrics configured
service/kube-state-metrics configured
$ kubectl get pods -n kube-system -l app.kubernetes.io/name=kube-state-metrics
NAME                                  READY   STATUS    RESTARTS   AGE
kube-state-metrics-6d7449fc78-mgf4f   1/1     Running   0          88s
```
# 5. Metricbeat
由于我们需要监控所有的节点，所以我们需要使用一个 DaemonSet 控制器来安装 Metricbeat。

首先，使用一个 ConfigMap 来配置 Metricbeat，然后通过 Volume 将该对象挂载到容器中的 /etc/metricbeat.yaml 中去。配置文件中包含了 ElasticSearch 的地址、用户名和密码，以及 Kibana 配置，我们要启用的模块与抓取频率等信息。

ElasticSearch 的 indice 生命周期表示一组规则，可以根据 indice 的大小或者时长应用到你的 indice 上。比如可以每天或者每次超过 1GB 大小的时候对 indice 进行轮转，我们也可以根据规则配置不同的阶段。由于监控会产生大量的数据，很有可能一天就超过几十G的数据，所以为了防止大量的数据存储，我们可以利用 indice 的生命周期来配置数据保留，这个在 Prometheus 中也有类似的操作。

我们配置成每天或每次超过5GB的时候就对 indice 进行轮转，并删除所有超过10天的 indice 文件，我们这里只保留10天监控数据完全足够了。
```
metricbeat.indice-lifecycle.configmap.yml
```

接下来就可以来编写 Metricbeat 的 DaemonSet 资源对象清单
```
metricbeat.daemonset.yml
```

需要注意的将上面的两个 ConfigMap 挂载到容器中去，由于需要 Metricbeat 获取宿主机的相关信息，所以我们这里也挂载了一些宿主机的文件到容器中去，比如 proc 目录，cgroup 目录以及 dockersock 文件。

由于 Metricbeat 需要去获取 Kubernetes 集群的资源对象信息，所以同样需要对应的 RBAC 权限声明，由于是全局作用域的，所以这里我们使用 ClusterRole 进行声明：
```
metricbeat.permissions.yml
```

```
$ kubectl apply  -f metricbeat.settings.configmap.yml \
                 -f metricbeat.indice-lifecycle.configmap.yml \
                 -f metricbeat.daemonset.yml \
                 -f metricbeat.permissions.yml

configmap/metricbeat-config configured
configmap/metricbeat-indice-lifecycle configured
daemonset.extensions/metricbeat created
clusterrolebinding.rbac.authorization.k8s.io/metricbeat created
clusterrole.rbac.authorization.k8s.io/metricbeat created
serviceaccount/metricbeat created

$ kubectl get pods -n elastic -l app=metricbeat   
NAME               READY   STATUS    RESTARTS   AGE
metricbeat-2gstq   1/1     Running   0          18m
metricbeat-99rdb   1/1     Running   0          18m
metricbeat-9bb27   1/1     Running   0          18m
metricbeat-cgbrg   1/1     Running   0          18m
metricbeat-l2csd   1/1     Running   0          18m
metricbeat-lsrgv   1/1     Running   0          18m
```

# 6. Filebeat
## 安装配置 Filebeat
安装配置 Filebeat来收集 Kubernetes 集群中的日志数据，然后发送到 ElasticSearch 去中，Filebeat 是一个轻量级的日志采集代理，还可以配置特定的模块来解析和可视化应用（比如数据库、Nginx 等）的日志格式。

和 Metricbeat 类似，Filebeat 也需要一个配置文件来设置和 ElasticSearch 的链接信息、和 Kibana 的连接已经日志采集和解析的方式。
我们配置采集 /var/log/containers/ 下面的所有日志数据，并且使用 inCluster 的模式访问 Kubernetes 的 APIServer，获取日志数据的 Meta 信息，将日志直接发送到 Elasticsearch。
```
filebeat.settings.configmap.yml
```

此外还通过 policy_file 定义了 indice 的回收策略：
```
filebeat.indice-lifecycle.configmap.yml
```

同样为了采集每个节点上的日志数据，我们这里使用一个 DaemonSet 控制器，使用上面的配置来采集节点的日志。
```
filebeat.daemonset.yml
```
我们这里使用的是 Kubeadm 搭建的集群，默认 Master 节点是有污点的，所以如果还想采集 Master 节点的日志，还必须加上对应的容忍，我这里不采集就没有添加容忍了。

此外由于需要获取日志在 Kubernetes 集群中的 Meta 信息，比如 Pod 名称、所在的命名空间等，所以 Filebeat 需要访问 APIServer，自然就需要对应的 RBAC 权限了，所以还需要进行权限声明：
```
filebeat.permission.yml
```

```
$ kubectl apply  -f filebeat.settings.configmap.yml \
                 -f filebeat.indice-lifecycle.configmap.yml \
                 -f filebeat.daemonset.yml \
                 -f filebeat.permissions.yml 

configmap/filebeat-config created
configmap/filebeat-indice-lifecycle created
daemonset.apps/filebeat created
clusterrolebinding.rbac.authorization.k8s.io/filebeat created
clusterrole.rbac.authorization.k8s.io/filebeat created
serviceaccount/filebeat created
```

当所有的 Filebeat 和 Logstash 的 Pod 都变成 Running 状态后，证明部署完成。现在我们就可以进入到 Kibana 页面中去查看日志了。左侧菜单 Observability → Logs

如果集群中要采集的日志数据量太大，直接将数据发送给 ElasticSearch，对 ES 压力比较大，这种情况一般可以加一个类似于 Kafka 这样的中间件来缓冲下，或者通过 Logstash 来收集 Filebeat 的日志。

这里我们就完成了使用 Filebeat 采集 Kubernetes 集群的日志，在下篇文章中，我们继续学习如何使用 Elastic APM 来追踪 Kubernetes 集群应用。

# 7. Elastic APM
Elastic APM 是 Elastic Stack 上用于应用性能监控的工具，它允许我们通过收集传入请求、数据库查询、缓存调用等方式来实时监控应用性能。这可以让我们更加轻松快速定位性能问题。

Elastic APM 是兼容 OpenTracing 的，所以我们可以使用大量现有的库来跟踪应用程序性能。

比如我们可以在一个分布式环境（微服务架构）中跟踪一个请求，并轻松找到可能潜在的性能瓶颈。

Elastic APM 通过一个名为 APM-Server 的组件提供服务，用于收集agent数据, 并向 ElasticSearch 发送追踪数据，再通过Kibana查看

## 安装 APM-Server
首先我们需要在 Kubernetes 集群上安装 APM-Server 来收集 agent 的追踪数据，并转发给 ElasticSearch，这里同样我们使用一个 ConfigMap 来配置：
```
apm.configmap.yml
```
APM-Server 需要暴露 8200 端口来让 agent 转发他们的追踪数据，新建一个对应的 Service 对象即可：
```
apm.service.yml
```
然后使用一个 Deployment 资源对象管理即可：
```
apm.deployment.yml
```
直接部署上面的几个资源对象：
```
$ kubectl apply  -f apm.deployment.yml \
                 -f apm.service.yml \
                 -f apm.deployment.yml

configmap/apm-server-config created
service/apm-server created
deployment.extensions/apm-server created

$ kubectl get pods -n elastic -l app=apm-server
NAME                          READY   STATUS    RESTARTS   AGE
apm-server-667bfc5cff-zj8nq   1/1     Running   0          12m
```

接下来我们可以在之前部署的 Spring-Boot 应用上安装一个 agent 应用。
## 配置Java Agent
接下来我们在示例应用程序 spring-boot-simple 上配置一个 Elastic APM Java agent。

首先我们需要把 elastic-apm-agent-1.8.0.jar 这个 jar 包程序内置到应用容器中去，在构建镜像的 Dockerfile 文件中添加一行如下所示的命令直接下载该 JAR 包即可：
```
Dockerfile
```
然后需要在示例应用中添加上如下依赖关系，这样我们就可以集成 open-tracing 的依赖库或者使用 Elastic APM API 手动检测。
```
<dependency>
    <groupId>co.elastic.apm</groupId>
    <artifactId>apm-agent-api</artifactId>
    <version>${elastic-apm.version}</version>
</dependency>
<dependency>
    <groupId>co.elastic.apm</groupId>
    <artifactId>apm-opentracing</artifactId>
    <version>${elastic-apm.version}</version>
</dependency>
<dependency>
    <groupId>io.opentracing.contrib</groupId>
    <artifactId>opentracing-spring-cloud-mongo-starter</artifactId>
    <version>${opentracing-spring-cloud.version}</version>
</dependency>
```
修改之前的java Deployment 部署的 Spring-Boot 应用，需要开启 Java agent 并且要连接到 APM-Server。

然后重新部署上面的示例应用：
```
$ kubectl apply -f spring-boot-simple.yml

$ kubectl get pods -n elastic -l app=spring-boot-simple
NAME                                 READY   STATUS    RESTARTS   AGE
spring-boot-simple-fb5564885-tf68d   1/1     Running   0          5m11s

$ kubectl get svc -n elastic -l app=spring-boot-simple
NAME                 TYPE       CLUSTER-IP      EXTERNAL-IP   PORT(S)          AGE
spring-boot-simple   NodePort   10.109.55.134   <none>        8080:31847/TCP   9d
```

### get messages
```
$ curl -X GET http://********:31847/message
```

### get messages (慢请求)
```
$ curl -X GET http://******:31847/message?sleep=3000
```
### get messages (error)
```
curl -X GET http://******:31847/message?error=true
```

现在我们去到 Kibana 页面中路由到 APM 页面，我们应该就可以看到 spring-boot-simple 应用的数据了。

# 小结
到这里我们就完成了使用 Elastic Stack 进行 Kubernetes 环境的全栈监控，通过监控指标、日志、性能追踪来了解我们的应用各方面运行情况，加快我们排查和解决各种问题。
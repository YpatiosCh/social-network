Got it! You’re looking to deploy a **Kafka cluster in Kubernetes (K8s)** with custom configurations like:

```
num.partitions=10
max.poll.records=50
fetch.max.bytes=5028800
```

I’ll give you a **step-by-step guide** to do this in a production-friendly way.

---

## **1. Choose a Kafka Deployment Method in Kubernetes**

There are a few options:

1. **Strimzi Operator (Recommended)**

   * A Kubernetes operator for running Apache Kafka clusters. Handles deployment, scaling, configuration, and rolling upgrades.
2. **Bitnami Kafka Helm Chart**

   * Simple Helm-based deployment.
3. **Manual StatefulSet**

   * You create the StatefulSet, ConfigMaps, and Services yourself. Gives full control but more work.

✅ I recommend **Strimzi** for production-grade Kafka clusters.

---

## **2. Install Strimzi Kafka Operator**

```bash
kubectl create namespace kafka
kubectl apply -f https://strimzi.io/install/latest?namespace=kafka -n kafka
```

This deploys the operator in the `kafka` namespace.

---

## **3. Define a Kafka Cluster Custom Resource (CR)**

Create a YAML file (e.g., `kafka-cluster.yaml`) like this:

```yaml
apiVersion: kafka.strimzi.io/v1beta2
kind: Kafka
metadata:
  name: my-cluster
  namespace: kafka
spec:
  kafka:
    replicas: 3
    listeners:
      - name: plain
        port: 9092
        type: internal
        tls: false
    config:
      num.partitions: 10
      offsets.topic.replication.factor: 3
      transaction.state.log.replication.factor: 3
      transaction.state.log.min.isr: 2
      log.message.format.version: "3.4"
    storage:
      type: persistent-claim
      size: 20Gi
      class: standard
  zookeeper:
    replicas: 3
    storage:
      type: persistent-claim
      size: 10Gi
      class: standard
  entityOperator:
    topicOperator: {}
    userOperator: {}
```

**Notes on your configs**:

* `num.partitions=10` → Added under `kafka.config.num.partitions`.
* Consumer configs like `max.poll.records` and `fetch.max.bytes` are **client-side configs**, not broker-side. You’ll need to configure them when creating Kafka consumers in your application.

---

## **4. Apply the Cluster YAML**

```bash
kubectl apply -f kafka-cluster.yaml -n kafka
```

Check pods:

```bash
kubectl get pods -n kafka
```

You should see 3 Kafka brokers + 3 ZooKeeper pods.

---

## **5. Deploy Kafka Clients with Configs**

For a Kafka consumer:

```properties
bootstrap.servers=my-cluster-kafka-bootstrap.kafka:9092
group.id=my-group
max.poll.records=50
fetch.max.bytes=5028800
```

> These are **application-side consumer configs**, not part of the Kafka broker. So your Spring Boot or Python Kafka consumer would use these in its config.

---

## **6. Optional: Expose Kafka Outside K8s**

For external access, you can add:

```yaml
listeners:
  - name: external
    port: 9094
    type: nodeport
    tls: false
```

or use **LoadBalancer** if your cloud supports it.

---

### ✅ Summary

1. Install Strimzi operator.
2. Create Kafka cluster CR with broker configs (num.partitions, etc.).
3. Deploy consumer apps with client configs (max.poll.records, fetch.max.bytes).
4. Expose Kafka if needed.

---

I can also give you a **ready-to-use Helm/Strimzi setup with all your configs included** that you can deploy with one command.

Do you want me to do that?

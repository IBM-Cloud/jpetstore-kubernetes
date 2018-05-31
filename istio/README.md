# In-Depth telemetry using Istio

## Before you continue:
Follow the steps till [Build and Push the container images](https://github.com/IBM-Cloud/ModernizeDemo#build-and-push-the-container-images) on the main README.md file.

Once the container images are pushed, continue with the below steps for In-Depth Telemetry - obtain uniform metrics, logs, traces across different services using Istio.

You will be using add-ons like Zipkin, Promethus, Grafana, Servicegraph & Weavescope.

*Note: Some configurations and features of the Istio platform are still under development and are subject to change based on user feedback. Allow a few months for stabilization before you use Istio in production.*


## Setup istio

[Istio](https://www.ibm.com/cloud/info/istio) is an open platform to connect, secure, and manage a network of microservices, also known as a service mesh, on cloud platforms such as Kubernetes in IBM® Cloud Kubernetes Service. With Istio, you can manage network traffic, load balance across microservices, enforce access policies, verify service identity, and more.

Install Istio in your cluster.

1. Get the latest version by using curl:

   ```
   curl -L https://git.io/getLatestIstio | sh -
   ```

2. Add the `istioctl` client to your PATH. For example, run the following command on a MacOS or Linux system:

   ```
   export PATH=$PWD/istio-0.7.1/bin:$PATH
   ```

3. Change the directory to the Istio file location.

   ```
   cd <FILEPATH>/istio-0.7.1
   ```

4. Install Istio on the Kubernetes cluster. Istio is deployed in the Kubernetes namespace `istio-system`.

   ```
   kubectl apply -f install/kubernetes/istio.yaml
   ```

   **Note**: If you need to enable mutual TLS authentication between sidecars, you can install the `istio-auth` file instead: `kubectl apply -f install/kubernetes/istio-auth.yaml`

5. Ensure that the Kubernetes services `istio-pilot`, `istio-mixer`, and `istio-ingress` are fully deployed before you continue.

   ```
   kubectl get svc -n istio-system
   ```

   ```
   NAME            TYPE           CLUSTER-IP       EXTERNAL-IP      PORT(S)                                                            AGE
   istio-ingress   LoadBalancer   172.21.xxx.xxx   169.xx.xxx.xxx   80:31176/TCP,443:30288/TCP                                         2m
   istio-mixer     ClusterIP      172.21.xxx.xxx     <none>           9091/TCP,15004/TCP,9093/TCP,9094/TCP,9102/TCP,9125/UDP,42422/TCP   2m
   istio-pilot     ClusterIP      172.21.xxx.xxx    <none>           15003/TCP,443/TCP                                                  2m
   ```

6. Ensure the corresponding pods `istio-pilot-*`, `istio-mixer-*`, `istio-ingress-*`, and `istio-ca-*`are also fully deployed before you continue.

   ```
   kubectl get pods -n istio-system
   ```

   ```
   istio-ca-3657790228-j21b9           1/1       Running   0          5m
   istio-ingress-1842462111-j3vcs      1/1       Running   0          5m
   istio-pilot-2275554717-93c43        1/1       Running   0          5m
   istio-mixer-2104784889-20rm8        2/2       Running   0          5m
   ```

Good work! You successfully installed Istio into your cluster. Next, deploy the JPetStore sample app into your cluster.

## Deploy

There are two different ways to deploy the three micro-services to a Kubernetes cluster:

- Updating yaml files with the right values and then running  `kubectl create` and `kubectl inject` (recommended with istio 0.7.1)
- Using [Helm](https://helm.sh/) to provide values for templated charts (Coming soon)

### Option 1: Deploy using YAML files(recommended)

When you deploy JPetStore, Envoy sidecar proxies are injected as containers into your app microservices' pods before the microservice pods are deployed. Istio uses an extended version of the Envoy proxy to mediate all inbound and outbound traffic for all microservices in the service mesh. For more about Envoy, see the [Istio documentation ![External link icon](https://console.bluemix.net/docs/api/content/icons/launch-glyph.svg?lang=en)](https://istio.io/docs/concepts/what-is-istio/overview.html#envoy).

For this option, you need to update the YAML files to point to your registry namespace and Kubernetes cluster Ingress subdomain:

1. Edit **jpetstore/jpetstore.yaml** and **jpetstore/jpetstore-watson.yaml** and replace all instances of:

  - `<CLUSTER DOMAIN>` with your Ingress Subdomain (`ibmcloud cs cluster-get CLUSTER_NAME`)
  - `<REGISTRY NAMESPACE>` with your Image registry URL. For example:`registry.ng.bluemix.net/mynamespace`

  When you deploy JPetStore, Envoy sidecar proxies are injected as containers into your app microservices' pods before the microservice pods are deployed. Istio uses an extended version of the Envoy proxy to mediate all inbound and outbound traffic for all microservices in the service mesh. For more about Envoy, see the [Istio documentation ![External link icon](https://console.bluemix.net/docs/api/content/icons/launch-glyph.svg?lang=en)](https://istio.io/docs/concepts/what-is-istio/overview.html#envoy).

2. Deploy the JPetstore app and database. The `kube-inject` command adds Envoy to the `jpetstore.yaml` file and uses this updated file to deploy the app. When the JPetstore app and database microservices deploy, the Envoy sidecar is also deployed in each microservice pod.

   ````
   kubectl create -f <(istioctl kube-inject --debug -f jpetstore/jpetstore.yaml)
   ````

3. This creates the MMSSearch microservice with Envoy sidecar

   ```
   kubectl create -f <(istioctl kube-inject --debug -f jpetstore/jpetstore-watson.yaml)
   ```

### Option 2: Deploy with Helm (coming soon)

## You're Done!

You are now ready to use the UI to shop for a pet or query the store by texting a picture of what you're looking at:

1. Access the java jpetstore application web UI for JPetstore at `http://jpetstore.<Ingress Subdomain>/`

   ![](/readme_images/petstore.png)

2. Access the mmssearch app and start uploading images from `pet-images` directory.

   ![](/readme_images/webchat.png)


### Load Generation for demo purposes

In a demo situation you might want to generate load for your application (it will help illustrate the various features in the dashboard). This can be done through the loadtest package:

```shell
# Use npm to install loadtest
npm install -g loadtest

# Generate increasing load (make sure to replace <myclustername> with the name of your cluster)
loadtest http://jpetstore.<yourclustername>.us-south.containers.mybluemix.net/
```

**Note:** Rerun the loadtest before every step to see the metrics and logging in realtime.

## In-depth telemetry from the service mesh through dashboards

With the application responding to traffic the graphs will start highlighting what's happening under the covers.

### Distributed tracing with Zipkin

Navigate to the folder where you have initially installed **Istio** and run the below command to install **Zipkin** addon

```
kubectl apply -f install/kubernetes/addons/zipkin.yaml
```

Setup access to the Zipkin dashboard URL using port-forwarding:

```
kubectl port-forward -n istio-system $(kubectl get pod -n istio-system -l app=zipkin -o jsonpath='{.items[0].metadata.name}') 9411:9411 &
```

Then open your browser at [http://localhost:9411](http://localhost:9411/)

![](images/Zipkin.png)

### Logs & Metrics collection and monitoring with Promethus

Navigate to the folder where you have initially installed **Istio** and run the below command to install **Promethus** addon

```
kubectl apply -f install/kubernetes/addons/prometheus.yaml
```

Under `istio` folder, a YAML file is provided to hold configuration for the new metric and log stream that Istio will generate and collect automatically. Navigate to `istio` folder and push the new configuration by running the below command

```
istioctl create -f istio-monitoring.yaml
```

In Kubernetes environments, execute the following command:

```
kubectl -n istio-system port-forward $(kubectl -n istio-system get pod -l app=prometheus -o jsonpath='{.items[0].metadata.name}') 9090:9090 &   
```

Visit <http://localhost:9090/graph> in your web browser and look for metrics starting with `istio`

![](images/Promethus.png)

### Visualizing Metrics with Grafana

Remember to install **Promethus** addon before following the steps below

1. To view Istio metrics in a graphical dashboard install the Grafana add-on.

   Point to the folder where you have install istio and In Kubernetes environments, execute the following command:

   ```sh
   kubectl apply -f install/kubernetes/addons/grafana.yaml
   ```

2. Verify that the service is running in your cluster.

   In Kubernetes environments, execute the following command:

   ```sh
   kubectl -n istio-system get svc grafana
   ```

   The output will be similar to:

   ```
   NAME      CLUSTER-IP      EXTERNAL-IP   PORT(S)    AGE
   grafana   10.59.247.103   <none>        3000/TCP   2m
   ```

3. Open the Istio Dashboard via the Grafana UI.

   In Kubernetes environments, execute the following command:

   ```sh
   kubectl -n istio-system port-forward $(kubectl -n istio-system get pod -l app=grafana -o jsonpath='{.items[0].metadata.name}') 3000:3000 &
   ```

   Visit <http://localhost:3000/dashboard/db/istio-dashboard> in your web browser.



   ![Istio Dashboard](images/grafana_istio.png)



   ![Mixer Dashboard](images/grafana_istio_mixer.png)





### Generating a service graph

Follow the steps mentioned in the [link on istio documentation](https://istio.io/docs/tasks/telemetry/servicegraph.html) to generate a servicegraph. Your servicegraph should look similar to the image below

![](images/servicegraph.png)

## Visualise Cluster using Weave Scope
While Service Graph displays a high-level overview of how systems are connected, a tool called Weave Scope provides a powerful visualisation and debugging tool for the entire cluster.

Using Scope it's possible to see what processes are running within each pod and which pods are communicating with each other. This allows users to understand how Istio and their application is behaving.

Scope is deployed onto a Kubernetes cluster with the command

```sh
kubectl apply -f "https://cloud.weave.works/k8s/scope.yaml?k8s-version=$(kubectl version | base64 | tr -d '\n')"
```

**Open Scope in Your Browser**

```sh
kubectl port-forward -n weave "$(kubectl get -n weave pod --selector=weave-scope-component=app -o jsonpath='{.items..metadata.name}')" 4040
```

The URL is: http://localhost:4040.

![](images/weavescope.png)

## Clean up

Pointing to the folder where you have install istio.

```shell
kubectl delete -f install/kubernetes/istio.yaml
```

## Related Content

- [More information on istio.io](https://istio.io/docs/guides/telemetry.html)

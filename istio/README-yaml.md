# Obtain and visualise uniform metrics, logs, traces across different services using Istio.

*Work in Progress*

You will be using add-ons like Jaeger, Promethus, Grafana, & Weavescope to collect, query and visualize metrics, logs & traces.

## Pre-req

Start by following the instructions in the parent [README](../README.md) to deploy the secrets and applications using Kubernetes-YAML. In this guide you will *uninstall* the applications, install Istio into the cluster and then redeploy the applications.

Run the below commands

```sh
kubectl delete -f jpetstore.yaml
kubectl delete -f jpetstore-watson.yaml
```

## Setup istio

[Istio](https://www.ibm.com/cloud/info/istio) is an open platform to connect, secure, and manage a network of microservices, also known as a service mesh, on cloud platforms such as Kubernetes in IBM® Cloud Kubernetes Service. With Istio, you can manage network traffic, load balance across microservices, enforce access policies, verify service identity, and more.

Install Istio in your cluster.

1. From your home directory, run this command to download the latest version by using curl:

   ```
   # from ~
   curl -L https://git.io/getLatestIstio | sh -
   ```

2. Change the directory to the Istio file location.

   ```
   cd istio-0.8.0
   ```

3. Add the `istioctl` client to your PATH. For example, run the following command on a MacOS or Linux system:

   ```
   export PATH=$PWD/bin:$PATH
   ```
   Run `istioctl version` to confirm successful setup. Next, deploy the JPetStore sample app into your cluster.


## Deploy

When you deploy JPetStore, Envoy sidecar proxies are injected as containers into your app microservices' pods before the microservice pods are deployed. Istio uses an extended version of the Envoy proxy to mediate all inbound and outbound traffic for all microservices in the service mesh. For more about Envoy, see the [Istio documentation](https://istio.io/docs/concepts/what-is-istio/overview.html#envoy).

For this option, you need to update the YAML files to point to your registry namespace and Kubernetes cluster Ingress subdomain:

1. a) Install Istio without enabling [mutual TLS authentication](https://istio.io/docs/concepts/security/mutual-tls/) between sidecars. Choose this option for clusters with existing applications, applications where services with an Istio sidecar need to be able to communicate with other non-Istio Kubernetes services, and applications that use [liveness and readiness probes](https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-probes/), headless services, or StatefulSets.

  ```sh
  $ kubectl apply -f install/kubernetes/istio-demo.yaml
  ```

  OR

  b) Install Istio and enforce mutual TLS authentication between sidecars by default. Use this option **only on a fresh kubernetes cluster** where newly deployed workloads are guaranteed to have Istio sidecars installed.

  ```sh
  $ kubectl apply -f install/kubernetes/istio-demo-auth.yaml
  ```

2. The Istio-Sidecar-injector will automatically inject Envoy containers into your application pods assuming running in namespaces labeled with `istio-injection=enabled`

   ```sh
   kubectl label namespace <namespace> istio-injection=enabled
   ```

3. Deploy the JPetstore app and database. When the JPetstore app and database microservices deploy, the Envoy sidecar is also deployed in each microservice pod.

   ````sh
   kubectl create -f jpetstore/jpetstore.yaml
   ````

4. This creates the MMSSearch microservice with Envoy sidecar

   ```sh
   kubectl create -f jpetstore/jpetstore-watson.yaml
   ```

5. By default, Istio-enabled services are unable to access URLs outside of the cluster because iptables is used in the pod to transparently redirect all outbound traffic to the sidecar proxy, which only handles intra-cluster destinations.

   Create an `ServiceEntry` to allow access to an external HTTPS service:

   ```sh
   kubectl create -f istio/egressgateway.yaml
   ```
   Notice that we also create a corresponding `DestinationRule` to initiate TLS for connections to the HTTPS service. Callers must access this service using HTTP on port 443 and Istio will upgrade the connection to HTTPS.

## You're Done!

You are now ready to use the UI to shop for a pet or query the store by texting a picture of what you're looking at:

1. Access the java jpetstore application web UI for JPetstore at `http://jpetstore.<Ingress Subdomain>/`

   ![](../readme_images/petstore.png)
2. Access the mmssearch app and start uploading images from `pet-images` directory.

   ![](../readme_images/webchat.png)


### Load Generation for demo purposes

In a demo situation you might want to generate load for your application (it will help illustrate the various features in the dashboard). This can be done through the loadtest package:

```shell
# Use npm to install loadtest
npm install -g loadtest

# Generate increasing load (make sure to replace <Ingress Subdomain> with the name of your ingress subdomain)
loadtest http://jpetstore.<Ingress Subdomain>/
```

**Note:** Rerun the loadtest before every step to see the metrics and logging in realtime.

## In-depth telemetry from the service mesh through dashboards

With the application responding to traffic, the graphs will start highlighting what's happening under the covers.

### Logs & Metrics collection and monitoring with Promethus

Under `istio` folder of JPetstore app, a YAML file is provided to hold configuration for the new metric and log stream that Istio will generate and collect automatically. On your terminal or command prompt, navigate to `istio` folder and push the new configuration by running the below command

```sh
istioctl create -f istio-monitoring.yaml
```

In Kubernetes environments, execute the following command:

```sh
kubectl -n istio-system port-forward $(kubectl -n istio-system get pod -l app=prometheus -o jsonpath='{.items[0].metadata.name}') 9090:9090 &   
```

Visit <http://localhost:9090/graph> in your web browser and look for metrics starting with `istio`

![](images/prometheus.png)

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

### Distributed tracing with Jaeger

In `istio/ingressgateway.yaml`,Replace the `<Ingress Subdomain>` with your subdomain and run the below 

```sh
kubectl create -f istio/ingressgateway.yaml
```

Setup access to the Jaeger dashboard URL using port-forwarding:

```sh
kubectl port-forward -n istio-system $(kubectl get pod -n istio-system -l app=jaeger -o jsonpath='{.items[0].metadata.name}') 16686:16686 &
```

Then open your browser at [http://localhost:16686](http://localhost:16686/) -> Select a trace and click **Find Traces**. If you click on the top (most recent) trace, you should see the details corresponding to your latest refresh. 

![](images/jaeger.png)

## Clean up

If you installed Istio with `istio-demo.yaml`. Point to the folder where you have installed istio and run the below command:

```shell
kubectl delete -f install/kubernetes/istio-demo.yaml
```

## Related Content

- [More information on istio.io](https://istio.io/docs/guides/telemetry.html)

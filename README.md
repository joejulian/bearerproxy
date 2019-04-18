# bearerproxy
A proxy service to add a Bearer token to the headers

## Usage

`bearerproxy` is configured with environment variables:

* `PORT` is the TCP port to listen on for connections
* `TARGET_URL` is the url for which to proxy
* `BEARER_TOKEN` is the bearer token to add to the request

The intended use is as a kubernetes sidecar container:

```yaml
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  labels:
    k8s-app: kubernetes-dashboard
  name: kubernetes-dashboard
  namespace: kube-system
spec:
  replicas: 1
  selector:
    matchLabels:
      k8s-app: kubernetes-dashboard
  template:
    metadata:
      labels:
        k8s-app: kubernetes-dashboard
    spec:
      containers:
      - name: kubernetes-dashboard
        image: gcr.io/google_containers/kubernetes-dashboard-amd64:v1.10.1
        ports:
          - name: insecureport
            containerPort: 9090
            protocol: TCP
        args:
          - --enable-insecure-login=true
      - name: bearerproxy
        image: quay.io/joejulian/bearerproxy:latest
        ports:
          - containerPort: 9091
            protocol: TCP
        env:
          - name: "PORT"
            value: "9091"
          - name: "TARGET_URL"
            value: "http://kubernetes-dashboard:9090"
          - name: "BEARER_TOKEN"
            valueFrom:
              secretKeyRef:
                name: dashboard-token-zyx98
                key: token
      serviceAccountName: kubernetes-dashboard
---
apiVersion: v1
kind: Service
metadata:
  name: kubernetes-dashboard
  namespace: kube-system
  labels:
    k8s-app: kubernetes-dashboard
spec:
  selector:
    k8s-app: kubernetes-dashboard
  ports:
  - port: 9090
    targetPort: 9090
```

You would then put a service and ingress controller in front of 9091 where the ingress controller is responsible for adequately providing authentication of the user.


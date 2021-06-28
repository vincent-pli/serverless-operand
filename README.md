<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [serverless-operand](#serverless-operand)
- [Tutorial](#tutorial)
- [Initial Operand as 0](#initial-operand-as-0)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# serverless-operand
This repo is a simple example of how to use Knative to manage your operand, and enable the operand as serverless model. This can help to reduce the footprint of your operator when it is running.

The operator will create a nginx server when related cr was detected. The nginx server will be wrappered as a [knative service](https://github.com/knative/serving/blob/main/pkg/apis/serving/v1/service_types.go)

# Tutorial
- Clone the peoject

```console
git clone git@github.com:vincent-pli/serverless-operand.git
```

- Deploy the controller

```console
make deploy IMG="quay.io/pengli1707/serverless-app:v0.0.1"
```

- Check if deploy succeed
```console
[root@symtest11 serverless-operand]# kubectl -n serverless-operand-system get po
NAME                                                     READY   STATUS    RESTARTS   AGE
serverless-operand-controller-manager-66cf8b4fd4-wpmbp   2/2     Running   0          14s
```

- Create CR  
```console
kubectl create -f config/samples/ibm.dev_v1alpha1_expressapp.yaml
```

- Check if `KSVC` appeared  
```console
[root@symtest11 serverless-operand]# kubectl get ksvc
NAME                     URL                                                 LATESTCREATED                  LATESTREADY                    READY   REASON
expressapp-sample-ksvc   http://expressapp-sample-ksvc.default.example.com   expressapp-sample-ksvc-00001   expressapp-sample-ksvc-00001   True
```

- Check `pod` which belong to `KSVC` appeared
```console
[root@symtest11 serverless-operand]# kubectl get po
NAME                                                       READY   STATUS    RESTARTS   AGE
expressapp-sample-ksvc-00001-deployment-576946dfdb-t2h6p   2/2     Running   0          8s
```

- Active the `pod`
The `pod` will be terminated(scale to 0, it's configurable) after 60s, the we could active it by send a request: 

```console
curl -H "Host: expressapp-sample-ksvc.default.example.com" http://9.28.237.203:31319`    
```

The `9.28.237.203` is the public ip of your node  
The `31319` is the host port for 80 in svc: `istio-system/istio-ingressgateway` 

- Check again the pod  
```console
[root@symtest11 rbac]# kubectl get po
NAME                                                       READY   STATUS    RESTARTS   AGE
expressapp-sample-ksvc-00001-deployment-576946dfdb-2lwpm   2/2     Running   0          31s
```

# Initial Operand as 0

Make initial-scale as 0 (after create `KSVC` the pod will not startup until request comming).

```console
kubectl -n knative-serving edit cm config-autoscaler
```

then add below setting to the `data`

```
  allow-zero-initial-scale: "true"
  initial-scale: "0"
```

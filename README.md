# serverless-operand
The operator will create a nginx server when related cr was detected. The nginx server will wrapper as a [knative service](https://github.com/knative/serving/blob/main/pkg/apis/serving/v1/service_types.go)

# Take a try
- Clone the peoject:  
`git clone git@github.com:vincent-pli/serverless-operand.git`

- Deploy the controller  
`make deploy IMG="quay.io/pengli1707/serverless-app:v0.0.1"`

- Check if deploy succedd
```
[root@symtest11 serverless-operand]# kubectl -n serverless-operand-system get po
NAME                                                     READY   STATUS    RESTARTS   AGE
serverless-operand-controller-manager-66cf8b4fd4-wpmbp   2/2     Running   0          14s
```

- Create CR  
`kubectl create -f config/samples/ibm.dev_v1alpha1_expressapp.yaml`

- Check if `KSVC` appeared  
```
[root@symtest11 serverless-operand]# kubectl get ksvc
NAME                     URL                                                 LATESTCREATED                  LATESTREADY                    READY   REASON
expressapp-sample-ksvc   http://expressapp-sample-ksvc.default.example.com   expressapp-sample-ksvc-00001   expressapp-sample-ksvc-00001   True
```

- Check `pod` which belong to `KSVC` appeared
```
[root@symtest11 serverless-operand]# kubectl get po
NAME                                                       READY   STATUS    RESTARTS   AGE
expressapp-sample-ksvc-00001-deployment-576946dfdb-t2h6p   2/2     Running   0          8s
```

- Active the `pod`
About 60s the `pod` will terminated(scale to 0, it's configurable), the we could active it by send a request:   
`curl -H "Host: http://expressapp-sample-ksvc.default.example.com" http://9.28.237.203:31319`    
The `9.28.237.203` is the public ip of your node  
The `31319` is the host port for 80 in svc: `istio-system/istio-ingressgateway` 

- Check again the pod  
```
[root@symtest11 rbac]# kubectl get po
NAME                                                       READY   STATUS    RESTARTS   AGE
expressapp-sample-ksvc-00001-deployment-576946dfdb-2lwpm   2/2     Running   0          31s
```

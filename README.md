# Kubernetes-Namespace-Permisison-Bot

This is a kubernetes operator to automatically provide certain people with permission upon namespace creation.


## Installation

In order to install the operator, first retrieve the example [deployment yaml](https://github.com/MrRiptide/Kubernetes-Namespace-Permission-Bot/blob/main/deploy.yaml) and save it as deploy.yaml
Once you have the deployment yaml, apply it into your kubernetes cluster using 

```bash
kubectl apply -f deploy.yaml
```

## Configuration

In order to configure this operator to your liking, you must edit the configmap. The easiest way of doing this is through editing a `config.json` file and then using that to provide the configmap with configuration

Once you've edited `config.json` to your liking, you can update the configmap with 

```bash
kubectl create configmap permission-bot-config --namespace=namespace-permission-bot-system --from-file=config.json
```

Then you must reboot the operator to update it's version of the config with

```bash
kubectl rollout restart deployment namespace-permission-bot-controller-manager --namespace=namespace-permission-bot-system
```

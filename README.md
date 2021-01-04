# Kubernetes-Namespace-Permisison-Bot

This is a kubernetes operator to automatically provide certain people with permission upon namespace creation.


## Installation

In order to install the operator, first retrieve the example [deployment yaml](https://github.com/MrRiptide/Kubernetes-Namespace-Permission-Bot/blob/main/deploy.yaml) and save it as deploy.yaml
Once you have the deployment yaml, apply it into your kubernetes cluster using 

```bash
kubectl apply -f deploy.yaml
```

*Note: If you are planning on using roles other than `admin` `edit` and `view`, make sure to create a ClusterRoleBinding that gives the operator's service account the role so that it can give those permissions.*

## Configuration

In order to configure this operator to your liking, you must edit the configmap. The easiest way of doing this is through editing a `config.json` file and then using that to provide the configmap with configuration

### Components of the config file

- `mode`: Sets the mode of the operator. `authoritative` makes it update the rolebindings on each event call, while `passive` will only set the rolebindings on namespace creation(Default: `authoritative`)
- `delimiter`: Defines the delimiter used in namespaces for use in detecting the prefix. Can be set to `""` to not require a specific delimiter(Default: `"-"`)
- `label_key`: Defines the global key to be used when checking labels to determine the group(Default: `"group"`)
- `groups`: Slice of groups to be used when creating rolebindings
  - `prefix`: Defines the prefix that will cause this group's roles to be applied to a namespace
  - `label`: Defines the label value that will cause this group's roles to be applied to a namespace
  - `label_key`: Defines the label to be used when checking labels for this group only. Will override the global setting
  - `roles`: Map of the roles and subjects to be created into a rolebinding. The key is the ClusterRole that will be used and the value is a slice of subjects to apply the RoleBinding to.
    - `type`: The type of the subject, can be `User`, `ServiceAccount`, or `Group`
    - `name`: The name of the subject


### Applying changes

Once you've edited `config.json` to your liking, you can update the configmap with 

```bash
kubectl create configmap permission-bot-config --namespace=namespace-permission-bot-system --from-file=config.json
```

Then you must reboot the operator to update it's version of the config with

```bash
kubectl rollout restart deployment namespace-permission-bot-controller-manager --namespace=namespace-permission-bot-system
```

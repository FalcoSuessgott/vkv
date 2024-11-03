# Kubernetes

`vkv` comes in container images, which enable you to run scheduled snapshots in a kubernetes cluster.

The idea is to schedule a cronjob which snapshots a vault server and writes the snapshot files to a persistent volume.

Here is a minimum working `k3s` using `local-storage` example:

## create the volume directories

```bash
# on a k3s node
mkdir -p /data/volume/pv1
chmod 777 /data/volume/pv1 # for testing
```

## create a pv

```yaml
apiVersion: v1
kind: PersistentVolume
metadata:
  name: local-pv
spec:
  capacity:
    storage: 5Gi
  accessModes:
  - ReadWriteOnce
  persistentVolumeReclaimPolicy: Retain
  storageClassName: local-storage
  local:
    path: /data/volumes/pv1
  nodeAffinity:
    required:
      nodeSelectorTerms:
      - matchExpressions:
        - key: kubernetes.io/hostname
          operator: In
          values:
          - worker-node # change
```

## create a pvc
```yaml
kind: PersistentVolumeClaim
apiVersion: v1
metadata:
  name: pvc
spec:
  accessModes:
  - ReadWriteOnce
  storageClassName: local-storage
  resources:
    requests:
      storage: 5Gi
```

## create a cronjob
```yaml
apiVersion: batch/v1
kind: CronJob
metadata:
  name: vkv
spec:
  schedule: "* * * * *" # runs every minute
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - name: vkv
            image: falcosuessgott/vkv:latest # stick to a version later
            imagePullPolicy: IfNotPresent
            command: ["/bin/sh", "-c"]
            args:
              - /vkv snapshot save -d /mnt/vkv-export-$(date '+%Y%m%d%H%M%S')
            env:
              - name: VAULT_SKIP_VERIFY
                value: "true"
              - name: VAULT_ADDR
                value: https://vault-server:8200 # change to Vault API address
              - name: VAULT_TOKEN
                value: hvs.xxxx # change to your token
            volumeMounts:
              - name: local-persistent-storage
                mountPath: /mnt
          restartPolicy: OnFailure
          volumes:
            - name: local-persistent-storage
              persistentVolumeClaim:
                claimName: pvc
```

## verify snapshots
if everything went correct, you should see the following:

```bash
ls -l /data/volumes/pv1/
total 0
drwxr-xr-x. 2 root root 108  5. Jan 09:50 vkv-export-20230105095000
drwxr-xr-x. 2 root root 108  5. Jan 09:51 vkv-export-20230105095100
```

## some last thoughts
Obviously this approach is just for development purposes. In order to make it production ready, you should consider changing some things, such as:

* inject the environments from a ConfigMap
* inject the token from a Secret
* Or obtain the token using Vaults kubernetes auth engine and the [Vault Agent injector](https://developer.hashicorp.com/vault/docs/platform/k8s/injector)
* run the cronjob daily
* update the permission of the volumes
* backup the pv

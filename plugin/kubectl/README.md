# Using kfilt as a kubectl Plugin

kfilt can be used as a kubectl plugin, allowing you to run it as `kubectl kfilt` instead of as a standalone binary.

## What are kubectl Plugins?

kubectl plugins are standalone executable files that extend kubectl with additional functionality. When you have an executable named `kubectl-foo` in your PATH, you can invoke it as `kubectl foo`.

## Installation

### Option 1: Install via Krew (Recommended)

[Krew](https://krew.sigs.k8s.io/) is the plugin manager for kubectl. Once you have [Krew installed](https://krew.sigs.k8s.io/docs/user-guide/setup/install/), you can install kfilt with:

```bash
kubectl krew install kfilt
```

### Option 2: Manual Installation

1. Download the appropriate `kubectl-kfilt` archive for your platform from the [releases page](https://github.com/ryane/kfilt/releases)

2. Extract the archive:
   ```bash
   # For Linux/macOS
   tar -xzf kubectl-kfilt_linux_amd64.tar.gz

   # For Windows (use your preferred extraction tool)
   unzip kubectl-kfilt_windows_amd64.zip
   ```

3. Move the `kubectl-kfilt` binary to a directory in your PATH:
   ```bash
   # For Linux/macOS
   sudo mv kubectl-kfilt /usr/local/bin/

   # Make it executable (Linux/macOS only)
   sudo chmod +x /usr/local/bin/kubectl-kfilt

   # For Windows, move kubectl-kfilt.exe to a directory in your PATH
   ```

4. Verify the installation:
   ```bash
   kubectl kfilt --version
   ```

## Usage

Once installed, you can use kfilt as a kubectl plugin. All the same flags and options work, but you invoke it with `kubectl kfilt` instead of just `kfilt`.

### Examples

#### Filter resources from kubectl output

Get all resources named "nginx-ingress-controller" regardless of kind:
```bash
kubectl get all -A -oyaml | kubectl kfilt -i name=nginx-ingress-controller
```

Filter ConfigMaps from a specific namespace:
```bash
kubectl get all -n kube-system -oyaml | kubectl kfilt -k ConfigMap
```

#### Use with Helm

Only output rendered Service resources from a Helm chart:
```bash
helm template my-chart | kubectl kfilt -i kind=service
```

Exclude all Secrets before applying a Chart to a cluster:
```bash
helm template my-chart | kubectl kfilt -x kind=secret | kubectl apply -f -
```

#### Use with Kustomize

Only output the ConfigMaps in a Kustomize base:
```bash
kustomize build ./base | kubectl kfilt -i kind=ConfigMap
```

Output only resources named "the-deployment":
```bash
kustomize build ./base | kubectl kfilt -i name=the-deployment
```

#### Filter and apply resources

Apply only Deployment resources from a manifest:
```bash
kubectl kfilt -f manifests.yaml -k Deployment | kubectl apply -f -
```

Delete all resources except ConfigMaps and Secrets:
```bash
kubectl kfilt -f resources.yaml -x kind=ConfigMap -x kind=Secret | kubectl delete -f -
```

## Advanced Usage

All kfilt functionality is available when using it as a kubectl plugin. See the main [README](../README.md) for comprehensive documentation on:

- Including/excluding resources by kind, name, namespace, labels, etc.
- Using wildcards in filters
- Label selectors
- Working with files and URLs

## Verifying Plugin Installation

To see all installed kubectl plugins:

```bash
kubectl krew list
```

or

```bash
kubectl plugin list
```

This should show `kubectl-kfilt` in the list of available plugins.

## Updating the Plugin

### With Krew
```bash
kubectl krew upgrade kfilt
```

### Manual Update
Download the latest release and repeat the manual installation steps.

## Uninstalling

### With Krew
```bash
kubectl krew uninstall kfilt
```

### Manual Uninstall
```bash
# Remove the binary from your PATH
sudo rm /usr/local/bin/kubectl-kfilt
```

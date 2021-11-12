**This example combines three `goal` features:**
1. Environmental executions: `goal run apply --on dev`
2. Built-in `kubectl_context` assertion upon execution
   to prevent accidental runs on wrong environment.
3. Built-in `approval` assertion to ask user to config execution.

**List available kubectl contexts with:**

```
$ kubectl config get-contexts
```

**Usage:**

```shell
$ goal run pods --on dev
$ goal run pods --on stage
```

```shell
$ goal run apply --on dev
$ goal run apply --on stage
```

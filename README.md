## Go Aliases

![GitHub release (latest by date)](https://img.shields.io/github/v/release/aaabramov/goal) [![Tests](https://github.com/aaabramov/goal/actions/workflows/test.yml/badge.svg?branch=master)](https://github.com/aaabramov/goal/actions/workflows/test.yml) [![Coverage Status](https://coveralls.io/repos/github/aaabramov/goal/badge.svg?branch=feature/coverage)](https://coveralls.io/github/aaabramov/goal?branch=feature/coverage) ![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/aaabramov/goal)

Allows you to create local aliases withing directory/repository with proper assertions upon executions.

**Motivation:**

- Simplify executing scoped repetitive commands
- Avoid executing commands on wrong environment (e.g. _kubectl_, _terraform_, _helm_, _etc._)
- Automatically generate OpsDoc from available goals. No need to read through whole README file to start operating on your infrastructure.

## Install

Install via `brew`:

```shell
# Will be simplified
brew tap aaabramov/goal https://github.com/aaabramov/goal
brew install aaabramov/goal/goal
```

## Usage

Run `goal init` in directory where aliases will be used. This will generate example `goal.yaml` file. Use it as a reference to define your own aliases. 

```shell
$ goal init
⌛ Generating default goal.yaml file
✅ Generated default goal.yaml file. Try running `goal` to see available goals.
```

### List goals

Simply type `goal` to see list of available goals and their dependencies:

```shell
$ goal
Available goals:
+---------------------+-------------+-----------------------------------------------------------------+-----------------------------+--------------------------------------------------+
|        GOAL         | ENVIRONMENT |                               CLI                               |         DESCRIPTION         |                    ASSERTIONS                    |
+---------------------+-------------+-----------------------------------------------------------------+-----------------------------+--------------------------------------------------+
| gcloud-ssh          | dev         | gcloud compute ssh dev-vm --zone=us-central1-c                  | SSH to dev                  | 1. gcloud.project == "dev-project"               |
+                     +-------------+-----------------------------------------------------------------+-----------------------------+--------------------------------------------------+
|                     | stage       | gcloud compute ssh stage-vm --zone=us-central1-c                | SSH to stage                | 1. gcloud.project == "stage-project"             |
+---------------------+-------------+-----------------------------------------------------------------+-----------------------------+--------------------------------------------------+
| helm-upgrade        | dev         | helm upgrade release-name -f values.yaml -f values/dev.yaml .   | helm upgrade on dev         | 1. kubectl.context == "gke_project_region_dev"   |
+                     +-------------+-----------------------------------------------------------------+-----------------------------+--------------------------------------------------+
|                     | stage       | helm upgrade release-name -f values.yaml -f values/stage.yaml . | helm upgrade on stage       | 1. kubectl.context == "gke_project_region_stage" |
+---------------------+-------------+-----------------------------------------------------------------+-----------------------------+--------------------------------------------------+
| k8s-apply           | dev         | kubectl apply -f deployment.yaml                                | kubectl apply on dev        | 1. kubectl.context == "gke_project_region_dev"   |
|                     |             |                                                                 |                             | 2. Manual approval                               |
+                     +-------------+                                                                 +-----------------------------+--------------------------------------------------+
|                     | stage       |                                                                 | kubectl apply on stage      | 1. kubectl.context == "gke_project_region_stage" |
|                     |             |                                                                 |                             | 2. Manual approval                               |
+---------------------+-------------+-----------------------------------------------------------------+-----------------------------+--------------------------------------------------+
| terraform-apply     | dev         | terraform apply -var-file vars/dev.tfvars                       | Terraform apply on dev      | 1. terraform.workspace == "dev"                  |
+                     +-------------+-----------------------------------------------------------------+-----------------------------+--------------------------------------------------+
|                     | stage       | terraform apply -var-file vars/stage.tfvars                     | Terraform apply on stage    | 1. terraform.workspace == "stage"                |
+---------------------+-------------+-----------------------------------------------------------------+-----------------------------+--------------------------------------------------+
| terraform-workspace |             | terraform workspace show                                        | Current terraform workspace |                                                  |
+---------------------+-------------+-----------------------------------------------------------------+-----------------------------+--------------------------------------------------+
| test                |             | go test -v ./...                                                | Run go tests                |                                                  |
+---------------------+-------------+-----------------------------------------------------------------+-----------------------------+--------------------------------------------------+
```

### Define simple local aliases

```yaml
pods:
  desc: Get nginx pods
  cmd: kubectl
  args:
    - get
    - pods
    - -l
    - app=nginx
svc:
  desc: Get nginx services
  cmd: kubectl
  args:
    - get
    - svc
    - -l
    - app=nginx
```

### Define goal with assertions

```yaml
# This example demonstrates how to use custom assertions upon executions.

my-assertion:
  desc: Ultimate Question of Life?
  cmd: echo
  args:
    - -n
    - $((40 + 2))
  
my-goal:
  desc: The Answer to the Ultimate Question of Life
  assert:
    - desc: If answer is 42..
      ref: my-assertion # references another goal
      expect: '42'
      fix: # CLI on how to fix
    - approve: yes # ask user to config execution  
  cmd: echo
  args:
    - The Answer to the Ultimate Question of Life, the Universe, and Everything is 42
```

### Built-in assertions

| Tool      | Example                                  |
|-----------|------------------------------------------|
| approval  | [examples/kubectl](examples/kubectl)     |
| kubectl   | [examples/kubectl](examples/kubectl)     |
| helm      | [examples/helm](examples/helm)           |
| terraform | [examples/terraform](examples/terraform) |
| gcloud    | [examples/gcloud](examples/gcloud)       |

## goal vs Makefile
_TODO_

## Project plan

- [X] Pipe STDIN for "yes/no" inputs, etc.
- [X] Add `assert.fix`. Display when assertion failed, e.g. `terraform workspace select dev`
- [X] Add "environment" management to avoid tf-plan-dev, tf-plan-stage, tf-plan-prod, etc. E.g. `goal tf-apply --on dev`
  & `goal.env: dev` matches
- [X] Support `-f my-goal.yaml`
- [X] Validate empty goal cmd
- [X] Validate empty assertion ref
- [X] Add `goal init` which simply generated example `goal.yaml`
- [X] Add predefined assertions:
    - [X] `k8s_cluster`
    - [X] `terraform_workspace`
    - [X] `gcloud_project`
- [X] `Check if current kubectl context is "gke_project_region_stage"` -> `kubectl.context == "gke_project_region_stage"`
- [ ] Assertions
  - [X] ref output
  - [X] support multiple assertions
  - [ ] recursive assertions?
  - [ ] raw CLI output -- bad pattern?
- [ ] Simpler `brew tap aaabramov/goal`
- [ ] Manual approvals for proceeding like `assert.approval`
- [ ] Add "depends on" other task like switch to dev?
    - [ ] Recursive dependencies
- [ ] Global aliases in `$HOME` directory?
- [ ] Self-autocompletion via [https://github.com/posener/complete](complete) library
- [ ] Support both goal.yaml & goal.yml
- [ ] Generate simple markdown file from `goal.yaml` (ops-doc)
- [ ] `goal add GOAL_NAME` -- check if already exists
- [ ] rework `Fatal` with `err`
- [ ] suggest `fix?` when precondition failed with `yes/no` prompt
- [ ] shared description from `goal.name` if there is no specific for env goal
- [ ] add to readme about `source <(goal completion zsh)`
- [ ] `did you forget "--on env"` when command name is found but env is required
- [ ] highlight commands & errors using https://github.com/fatih/color

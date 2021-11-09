## Go Aliases

![GitHub release (latest by date)](https://img.shields.io/github/v/release/aaabramov/goal) [![Tests](https://github.com/aaabramov/goal/actions/workflows/test.yml/badge.svg?branch=master)](https://github.com/aaabramov/goal/actions/workflows/test.yml) ![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/aaabramov/goal)

Allows you to create local aliases withing directory/repository with proper assertions upon executions.

**The idea behind is to:**

- simplify executing scoped repetitive commands 
- avoid executing commands on wrong environment (e.g. _kubectl_, _terraform_, _helm_, _etc._)

## Install

Install via `brew`:

```shell
# Will be simplified
brew tap aaabramov/goal https://github.com/aaabramov/goal
brew install aaabramov/goal/goal
```

## Usage

Create `goal.yaml` file in directory where aliases will be used:

```yaml
workspace:
  desc: Current terraform workspace
  cmd: terraform
  args:
    - workspace
    - show

tf-plan:
  envs:
    dev:
      desc: Terraform plan on dev
      assert:
        desc: Check if on dev workspace
        ref: workspace # References goal above
        expect: dev    # Checks whether trimmed output from 'ref' goal is equal to "dev"
      cmd: terraform
      args:
        - plan
        - -var-file
        - vars/dev.tfvars
    stage:
      desc: Terraform plan on stage
      assert:
        desc: Check if on stage workspace
        ref: workspace # References goal above
        expect: stage  # Checks whether trimmed output from 'ref' goal is equal to "stage"
      cmd: terraform
      args:
        - plan
        - -var-file
        - vars/stage.tfvars

tf-apply:
  envs:
    dev:
      desc: Terraform apply on dev
      assert:
        desc: Check if on dev workspace
        ref: workspace # References goal above
        expect: dev    # Checks whether trimmed output from 'ref' goal is equal to "dev"
      cmd: terraform
      args:
        - apply
        - -var-file
        - vars/dev.tfvars
    stage:
      desc: Terraform apply on stage
      assert:
        desc: Check if on stage workspace
        ref: workspace # References goal above
        expect: stage  # Checks whether trimmed output from 'ref' goal is equal to "stage"
      cmd: terraform
      args:
        - apply
        - -var-file
        - vars/stage.tfvars
```

Simply type `goal` to see list of available goals and their dependencies:

```shell
$ goal
Available goals:
+-----------+-------------+--------------------------------+-----------------------------+--------------------------------+
|   GOAL    | ENVIRONMENT |              CLI               |         DESCRIPTION         |           ASSERTIONS           |
+-----------+-------------+--------------------------------+-----------------------------+--------------------------------+
| workspace |             | terraform workspace show       | Current terraform workspace |                                |
+-----------+-------------+--------------------------------+-----------------------------+--------------------------------+
| tf-plan   | dev         | terraform plan -var-file       | Terraform plan on dev       | [workspace] Check if on dev    |
|           |             | vars/dev.tfvars                |                             | workspace                      |
+-----------+-------------+--------------------------------+-----------------------------+--------------------------------+
| tf-plan   | stage       | terraform plan -var-file       | Terraform plan on stage     | [workspace] Check if on stage  |
|           |             | vars/stage.tfvars              |                             | workspace                      |
+-----------+-------------+--------------------------------+-----------------------------+--------------------------------+
| tf-apply  | dev         | terraform apply -var-file      | Terraform apply on dev      | [workspace] Check if on dev    |
|           |             | vars/dev.tfvars                |                             | workspace                      |
+-----------+-------------+--------------------------------+-----------------------------+--------------------------------+
| tf-apply  | stage       | terraform apply -var-file      | Terraform apply on stage    | [workspace] Check if on stage  |
|           |             | vars/stage.tfvars              |                             | workspace                      |
+-----------+-------------+--------------------------------+-----------------------------+--------------------------------+
```

Let's see if _goal_ would allow us to apply terraform configuration on wrong environment:

```shell
$ terraform workspace show
dev
$ goal tf-apply --on stage
⚙️ Exec tf-apply-stage
⌛ Check precondition: Check if on stage workspace
❗ Precondition failed: workspace
	Output:   "dev"
	Expected: "stage"
	CLI: terraform workspace show
```

## Idea behind

1. Local alias management  
   To avoid typing repeatable commands
2. AssD - Aliases as a Documentation :D  
   No need to read through whole README file to start operating on you infrastructure

## goal vs Makefile

## Project plan

- [X] Pipe STDIN for "yes/no" inputs, etc.
- [X] Add `assert.fix`. Display when assertion failed, e.g. `terraform workspace select dev`
- [X] Add "environment" management to avoid tf-plan-dev, tf-plan-stage, tf-plan-prod, etc. E.g. `goal tf-apply --on dev` & `goal.env: dev` matches
- [X] Support `-f my-goal.yaml`
- [X] Validate empty goal cmd
- [X] Validate empty assertion ref
- [X] Add `goal init` which simply generated example `goal.yaml`
- [ ] Simpler `brew tap aaabramov/goal`
- [ ] Manual approvals for proceeding like `assert.approval`
- [ ] Add "depends on" other task like switch to dev?
  - [ ] Recursive dependencies
- [ ] Assertions
    - [X] ref output
    - [ ] recursive assertions?
    - [ ] raw CLI output -- bad pattern?
- [ ] Global aliases in `$HOME` directory?
- [ ] Self-autocompletion via [https://github.com/posener/complete](complete) library
- [ ] Support both goal.yaml & goal.yml
- [ ] Generate simple markdown file from `goal.yaml` (ops-doc)
- [ ] Add predefined assertions like `k8s_cluster`, `terraform_workspace`, `etc.`
- [ ] `goal add GOAL_NAME` -- check if already exists

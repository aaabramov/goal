## Go Aliases

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

tf-apply-dev:
  desc: Terraform apply on dev
  assert:
    desc: Check if on dev workspace
    ref: workspace # References goal above
    equals: dev    # Checks whether trimmed output from 'ref' goal is equal to "dev"
  cmd: terraform
  args:
    - apply
    - -var-file
    - vars/dev.tfvars

tf-apply-stage:
  desc: Terraform apply on stage
  assert:
    desc: Check if on stage workspace
    ref: workspace # References goal above
    equals: stage  # Checks whether trimmed output from 'ref' goal is equal to "stage"
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
+----------------+--------------------------------+-----------------------------+--------------------------------+
|      GOAL      |              CLI               |         DESCRIPTION         |           ASSERTIONS           |
+----------------+--------------------------------+-----------------------------+--------------------------------+
| tf-apply-dev   | terraform apply -var-file      | Terraform apply on dev      | [workspace] Check if on dev    |
|                | vars/dev.tfvars                |                             | workspace                      |
| tf-apply-stage | terraform apply -var-file      | Terraform apply on stage    | [workspace] Check if on stage  |
|                | vars/stage.tfvars              |                             | workspace                      |
| workspace      | terraform workspace show       | Current terraform workspace |                                |
+----------------+--------------------------------+-----------------------------+--------------------------------+
```

Let's see if _goal_ would allow us to apply terraform configuration on wrong environment:

```shell
$ terraform workspace show
dev
$ goal tf-apply-stage
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

- [ ] Pipe STDIN for "yes/no" inputs, etc.
- [ ] Simpler `brew tap aaabramov/goal`
- [ ] Add manual approve step
- [ ] Add "environment" management to avoid tf-plan-dev, tf-plan-stage, tf-plan-prod, etc.
- [ ] Add "depends on" other task like switch to dev?
- [ ] Recursive dependencies
- [ ] Manual approvals for proceeding
- [ ] Assertions
    - [ ] ref output
    - [ ] recursive assertions
    - [ ] raw CLI output -- bad pattern?
- [ ] Global aliases in `$HOME` directory?
- [ ] Self-autocompletion via [https://github.com/posener/complete](complete) library
- [ ] Generate ops-doc from commands
- [ ] Support both goal.yaml & goal.yml
- [ ] Support `-f my-goal.yaml`
- [ ] Add `goal init` which simply generated example `goal.yaml`
- [ ] Generate simply markdown file from `goal.yaml`

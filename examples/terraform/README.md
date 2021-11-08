# Test it out

```shell
$ terraform workspace new dev
$ terraform workspace new stage
$ goal
Available goals:
+-----------+-------------+--------------------------------+-----------------------------+--------------------------------+
|   GOAL    | ENVIRONMENT |              CLI               |         DESCRIPTION         |           ASSERTIONS           |
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
| workspace |             | terraform workspace show       | Current terraform workspace |                                |
+-----------+-------------+--------------------------------+-----------------------------+--------------------------------+
$ terraform workspace show
stage
$ goal tf-plan --on dev
Running on dev
⚙️  Exec tf-plan
⌛ Check precondition: Check if on dev workspace
❗ Precondition failed: workspace
	Output:   "stage"
	Expected: "dev"
	CLI: terraform workspace show
```

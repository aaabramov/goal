# Test it out

```shell
$ terraform workspace new dev
Created and switched to workspace "dev"!

$ terraform workspace new stage
Created and switched to workspace "stage"!

$ goal
Available goals:
+-------+-------------+---------------------------------------------+--------------------------+--------------------------------------------------+
| GOAL  | ENVIRONMENT |                     CLI                     |       DESCRIPTION        |                    ASSERTIONS                    |
+-------+-------------+---------------------------------------------+--------------------------+--------------------------------------------------+
| apply | dev         | terraform apply -var-file vars/dev.tfvars   | Terraform apply on dev   | Check if selected terraform workspace is "dev"   |
+       +-------------+---------------------------------------------+--------------------------+--------------------------------------------------+
|       | stage       | terraform apply -var-file vars/stage.tfvars | Terraform apply on stage | Check if selected terraform workspace is "stage" |
+-------+-------------+---------------------------------------------+--------------------------+--------------------------------------------------+
| plan  | dev         | terraform plan -var-file vars/dev.tfvars    | Terraform plan on dev    | Check if selected terraform workspace is "dev"   |
+       +-------------+---------------------------------------------+--------------------------+--------------------------------------------------+
|       | stage       | terraform plan -var-file vars/stage.tfvars  | Terraform plan on stage  | Check if selected terraform workspace is "stage" |
+-------+-------------+---------------------------------------------+--------------------------+--------------------------------------------------+

$ terraform workspace show
stage

# Let's see if goal would allow us to apply terraform configuration on wrong environment:
$ goal tf-plan --on dev
Running on dev
⚙️  Exec tf-plan
⌛ Check precondition: Check if on dev workspace
❗ Precondition failed: workspace
	Output:   "stage"
	Expected: "dev"
	CLI: terraform workspace show
```

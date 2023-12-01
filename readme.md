**ghostcloud** is an innovative web hosting blockchain powered by Cosmos SDK and the CometBFT consensus engine.

[![build](https://img.shields.io/circleci/build/github/liftedinit/ghostcloud-cosmos/main)](https://app.circleci.com/pipelines/github/liftedinit/ghostcloud-cosmos)
[![coverage](https://img.shields.io/codecov/c/github/liftedinit/ghostcloud-cosmos)](https://app.codecov.io/gh/liftedinit/ghostcloud-cosmos)

<!-- TOC -->
  * [Requirements](#requirements-)
  * [Get started](#get-started-)
    * [Configure](#configure)
  * [How to use](#how-to-use)
    * [Deploying a new instance](#deploying-a-new-instance)
    * [Update an existing deployment](#update-an-existing-deployment)
    * [Remove an existing deployment](#remove-an-existing-deployment)
    * [List all deployments](#list-all-deployments)
  * [Developers](#developers)
<!-- TOC -->

## Requirements 

Before setting up ghostcloud, ensure your system meets the following prerequisites:

- Ignite CLI, version 0.27.2
- Go programming language, version 1.21 or higher

## Get started 

To jumpstart your **ghostcloud** blockchain development, run the following command:

```
ignite chain serve
```

The `serve` command performs a series of actions: it installs necessary dependencies, compiles your blockchain's source code, initializes the default configuration files, and finally, launches your blockchain in a local development environment.


### Configure

To tailor your development blockchain, modify the settings in the `config.yml` file according to your requirements.

## How to use

This section describes how to interact with your **ghostcloud** blockchain using the command-line interface (CLI). 

### Deploying a new instance

To create a new deployment on your blockchain, use the command format below:

```shell
ghostcloudd tx ghostcloud create [NAME] [PAYLOAD] --from [KEY] --gas auto --yes
```

where
- `[NAME]` is your chosen name for the deployment.
- `[PAYLOAD]` is the path to either a directory or a zip file containing your deployment's content. 
- `[KEY]` is the name of the key used for signing the transaction.

Optional flags:
- `--description "[DESCRIPTION]"` - A brief description of the deployment (optional).
- `--domain [DOMAIN]` - The domain that will be associated with this deployment (optional).

Important considerations: 
- The `[PAYLOAD]` must have an `index.html` file located at the root. 
- The size of the `[PAYLOAD]` is limited to a maximum of 5MB.

Example usage:
```shell
ghostcloudd tx ghostcloud create foobar ~/foobar.zip --from alice --gas auto --yes
```

In the example above, a new deployment named foobar is created using the contents of `~/foobar.zip`, signed with the key alice. 
The transaction is set to calculate the gas automatically and is confirmed without further prompts.

### Update an existing deployment

```shell
ghostcloudd tx ghostcloud update [NAME] [DESCRIPTION] [DOMAIN] [flags] --from [KEY] --gas auto --yes
```

where
- `[NAME]` is the name of the deployment to update.
- `[DESCRIPTION]` is the new description for the deployment.
- `[DOMAIN]` is the new domain associated with the deployment.
- `[KEY]` is the name of the key to use for signing the transaction.

Available flags:
  - `--website-payload [PATH]` - (Optional) Provide the path to the new website payload, which can be a directory or a zip file.

Important considerations:
- If using `--website-payload`, ensure that the payload contains an `index.html` file at its root and that the total size does not exceed 5MB.
- `[NAME]`, `[DESCRIPTION]`, and `[DOMAIN]` must be provided together. If you do not wish to update a particular field, use the existing value.

Example usage:
```shell
ghostcloudd tx ghostcloud update myapp "Updated version" "newexample.com" --website-payload ~/newapp.zip \
  --from alice --gas auto --yes 
```

In this example, the `myapp` deployment is updated with new contents from `~/newapp.zip`, a new description, and a new domain, all signed with the key alice. 
The `--gas auto` flag allows the transaction to automatically calculate the gas needed, and `--yes` confirms the transaction without additional prompts.

### Remove an existing deployment

```shell
ghostcloudd tx ghostcloud remove [NAME] --from [KEY] --gas auto --yes
```

where
- `[NAME]` is the name of the deployment to remove.
- `[KEY]` is the name of the key to use for signing the transaction.

Example usage:
```shell
ghostcloudd tx ghostcloud remove myapp --from alice --gas auto --yes
```

In this example, the `myapp` deployment is removed from the blockchain, signed with the key alice.
The `--gas auto` flag allows the transaction to automatically calculate the gas needed, and `--yes` confirms the transaction without additional prompts.

### List all deployments

```shell
ghostcloudd q ghostcloud list
```

Available flags:
- `--filter-by string` - Apply a filter to the list of deployments. The filter can be `creator`.
- `--filter-value string` - The value to use for the filter.
- `--filter-operator string` - The operator to use for the filter. The operator can be `equal`.

Example usage:
```shell
ghostcloudd q ghostcloud list --filter-by creator \
  --filter-value gc13q9lpjm0zwse4msn5l6anwzznphxgnzxssf86x \
  --filter-operator equal
```

In this example, the command will return the list of deployments created by the address `gc13q9lpjm0zwse4msn5l6anwzznphxgnzxssf86x`. 


## Developers

Use the provided `Makefile` to execute common operations:

```
$ make help
help                           Display this help screen
lint                           Run linter (golangci-lint)
format                         Run formatter (goimports)
coverage                       Run coverage report
test                           Run tests
```
# Hephy Terraform Provider #

This repository contains a terraform provider for interacting with Hephy clusters.

The terraform provider implements 1 resources currently:

* `hephy_initial_admin` Manage the initial admin of a hephy cluster
* `hephy_user` Manage a user on a hephy cluster
* `hephy_key` Manage an SSH key for a user on a hephy cluster

Check the `docs` folder for detailed documentation of the Terraform Provider.

As the hephy controller API does not support `update` for most objects, such operations have been omitted from the resources implemented, so any update of those resources will cause a delete/create.

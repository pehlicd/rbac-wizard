<div align="center" style="padding-top: 20px">
    <img src="/assets/rancher-rbac-wizard.jpg?raw=true" width="220" style="background-color: blue;">
</div>

# Rancher RBAC Wizard

![go version](https://img.shields.io/github/go-mod/go-version/alegrey91/rancher-rbac-wizard)
![release](https://img.shields.io/github/v/release/alegrey91/rancher-rbac-wizard?filter=v*)
![license](https://img.shields.io/github/license/alegrey91/rancher-rbac-wizard)
[![go report](https://goreportcard.com/badge/github.com/alegrey91/rancher-rbac-wizard)](https://goreportcard.com/report/github.com/alegrey91/rancher-rbac-wizard)

Rancher RBAC Wizard is a tool that helps you visualize and analyze the RBAC configurations of your Kubernetes cluster. It provides a graphical representation of the Rancher Kubernetes RBAC objects (such as `GlobalRoles`, `Projects`, `Clusters`, `RoleTemplates`, etc).

📌 Note: This project is a fork of [pehlicd/rbac-wizard](https://github.com/pehlicd/rbac-wizard). Please refer to the original repository for the initial implementation and history.

<div align="center">


| Demo                                       |
|--------------------------------------------|
| <img src="/assets/rbac-wizard-demo.gif" /> |

</div>

## How to install

```bash
go install github.com/alegrey91/rancher-rbac-wizard@latest
```

## How to use

Using Rancher RBAC Wizard is super simple. Just run the following command:

```bash
rancher-rbac-wizard serve
```

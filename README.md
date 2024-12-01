<div align="center" style="padding-top: 20px">
    <img src="/ui/public/rbac-wizard-logo-embedded.svg?raw=true" width="120" style="background-color: blue; border-radius: 50%;">
</div>

# RBAC Wizard

![go version](https://img.shields.io/github/go-mod/go-version/pehlicd/rbac-wizard)
![release](https://img.shields.io/github/v/release/pehlicd/rbac-wizard?filter=v*)
![helm release](https://img.shields.io/github/v/release/pehlicd/rbac-wizard?filter=rbac-wizard*&logo=helm)
![license](https://img.shields.io/github/license/pehlicd/rbac-wizard)
[![go report](https://goreportcard.com/badge/github.com/pehlicd/rbac-wizard)](https://goreportcard.com/report/github.com/pehlicd/rbac-wizard)

RBAC Wizard is a tool that helps you visualize and analyze the RBAC configurations of your Kubernetes cluster. It provides a graphical representation of the Kubernetes RBAC objects.

<div align="center">


| Demo                                       |
|--------------------------------------------|
| <img src="/assets/rbac-wizard-demo.gif" /> |

</div>

## How to install

### Helm

Since rbac-wizard is capable of getting kubernetes clientset from the cluster ease free, you can also install it on your cluster using Helm with 3 simple steps!

```bash
# to add the Helm repository
helm repo add rbac-wizard https://rbac-wizard.pehli.dev
# to pull the latest Helm charts
helm pull rbac-wizard/rbac-wizard
# to install the Helm charts with the default values
helm install rbac-wizard rbac-wizard/rbac-wizard --namespace rbac-wizard --create-namespace
```

### Homebrew

```bash
brew tap pehlicd/rbac-wizard https://github.com/pehlicd/rbac-wizard
brew install rbac-wizard
```

### Using go install

```bash
go install github.com/pehlicd/rbac-wizard@latest
```

## How to use

Using RBAC Wizard is super simple. Just run the following command:

```bash
rbac-wizard serve
```

## How to contribute

If you'd like to contribute to RBAC Wizard, feel free to submit pull requests or open issues on the [GitHub repository](https://github.com/pehlicd/rbac-wizard). Your feedback and contributions are highly appreciated!

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

Developed by [Furkan Pehlivan](https://github.com/pehlicd) - [Project Repository](https://github.com/pehlicd/rbac-wizard)
<div align="center" style="padding-top: 5px">
    <img src="/ui/public/rbac-wizard-logo-embedded.svg?raw=true" width="120">
</div>

# RBAC Wizard

RBAC Wizard is a tool that helps you visualize and analyze the RBAC configurations of your Kubernetes cluster. It provides a graphical representation of the Kubernetes RBAC objects.

<div align="center">


| Demo                                       |
|--------------------------------------------|
| <img src="/assets/rbac-wizard-demo.gif" /> |

</div>

## How to install

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
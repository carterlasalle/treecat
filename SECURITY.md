# Security Policy

Security fixes are applied to the latest release and the `main` branch.

Report vulnerabilities privately through GitHub's **Security → Report a vulnerability** flow for this repository. Include the affected version, operating system, reproduction steps, impact, and any suggested remediation. Avoid attaching private source trees or generated context files that may contain credentials.

The project aims to acknowledge reports within 72 hours and provide an initial severity assessment within seven days.

## Safe usage

`treecat` can collect the contents of an entire directory. Review its output before sharing it with another person or service. Existing `.gitignore` rules are honored by default, but ignore files are not a substitute for secret scanning or a deliberate extension/path allowlist.

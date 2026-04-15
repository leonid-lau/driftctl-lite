# driftctl-lite

A lightweight CLI tool to detect infrastructure drift between Terraform state and live cloud resources.

---

## Installation

```bash
go install github.com/yourusername/driftctl-lite@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/driftctl-lite.git
cd driftctl-lite
go build -o driftctl-lite .
```

---

## Usage

Point `driftctl-lite` at your Terraform state file and let it compare against your live cloud environment:

```bash
# Scan using a local state file
driftctl-lite scan --state terraform.tfstate --provider aws

# Scan using a remote S3 backend
driftctl-lite scan --state s3://my-bucket/terraform.tfstate --provider aws --region us-east-1
```

### Example Output

```
[✔] aws_s3_bucket.my-bucket        — in sync
[✗] aws_security_group.web         — DRIFT DETECTED (ingress rules modified)
[✗] aws_iam_role.lambda_exec       — DRIFT DETECTED (policy detached)

Summary: 3 resources checked | 2 drifted | 1 in sync
```

---

## Supported Providers

- AWS (initial support)
- GCP *(coming soon)*
- Azure *(coming soon)*

---

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.

```bash
go test ./...
```

---

## License

This project is licensed under the [MIT License](LICENSE).
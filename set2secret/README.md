## Overview
The `set2secret` application is designed to generate Kubernetes Secret manifests from a collection of environment variables referred to as "sets". The tool allows for the extraction of these sets, comparison to referenced files (for validation), and generation of a Kubernetes Secret manifest, which can then be applied to a Kubernetes cluster.

## Features
1. Environment Set Extraction: Extract environment variables based on provided set names and create a combined variable set.
2. Reference File Comparison: Compare the extracted variables to those defined in provided reference files to ensure all required placeholders are present and identify redundant variables.
3. Kubernetes Secret Generation: Generate a Kubernetes Secret manifest using the final variable set.

## Usage
To execute the `set2secret` application, use the following command structure:

```shell
set2secret -s [variable_set_name ...] -n [k8s_secret_name] [-r path_to_reference_file ...]
```

Where
* **-s, --set**: Name(s) of the variable set(s) to extract from the environment. Multiple sets can be provided.
* **-n, --secret-name**: The desired name of the Kubernetes Secret.
* **-r, --ref**: (Optional) Paths to reference files containing placeholders to be validated against the extracted variable set.

### Example:
```shell
set2secret -s SET1 SET2 -n my-k8s-secret -r ./config.yaml
```
This command will:

1. Extract variables from the environment based on the sets `SET1` and `SET2`.
2. Compare the extracted variables to placeholders within `config.yaml`.
3. Generate a Kubernetes Secret manifest named `my-k8s-secret`.

## Output
1. **Variables Extraction**: Lists extracted variables and values.
2. **Reference Comparison**: Reports any missing variables not found in the reference files or redundant variables that aren't needed.
3. **Kubernetes Secret Generation**: Outputs the location where the generated Kubernetes Secret manifest is saved.

## Considerations
* The application expects variable sets to be defined in YAML format within environment variables.
* By default, the application identifies placeholders in the format `{{PLACEHOLDER_NAME}}` within the reference files.
* All variables within the generated Kubernetes Secret manifest are encoded using Base64 as required by Kubernetes.

## Logging
The application provides detailed logging for better traceability:

* **INFO**: General information about the ongoing process.
* **WARN**: Warnings about missing or redundant variables.
* **FATAL**: Critical errors that cause the application to exit.

## Dependencies
The `set2secret` application leverages several utilities and external libraries:

* **cicd_envsubst/utils/env_var**: Operations related to environment variables.
* **cicd_envsubst/utils/file**: File operations.
* **cicd_envsubst/utils/path**: Path-related utilities.
* **cicd_envsubst/utils/placeholder**: Operations related to placeholders.
* **github.com/jessevdk/go-flags**: CLI argument parsing.
* **github.com/thoas/go-funk**: Array operations and helpers.
* **gopkg.in/yaml.v3**: YAML serialization and deserialization.

## Contribution & Support
For bug reports, feature requests, and contributions, please open an issue in the application's repository. Feedback and contributions are always welcome!
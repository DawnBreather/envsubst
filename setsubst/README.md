## Overview
The `setsubst` application focuses on substituting environment variable placeholders within target files or directories based on the provided "sets". Given a set of environment variables (stored in YAML format), the application reads these variables and replaces placeholders within target files with the appropriate values from the sets.

## Features
1. **Environment Set Extraction**: Extracts environment variables based on provided set names.
2. **Placeholder Substitution**: Substitutes placeholders within target files or directories using values from the extracted sets.

## Usage
To execute the `setsubst` application, use the following command structure:

```shell
setsubst -s [variable_set_name ...]
```
Where:
* **-s, --set**: Specifies the name(s) of the variable set(s) to extract from the environment.
<p></p>
After the set name(s), provide the paths to the target files or directories as positional arguments.

### Example:
```shell
setsubst -s SET1 SET2 ./config.yaml ./templates/
```
This command will:
1. Extract variables from the environment based on the sets `SET1` and `SET2`.
2. Substitute placeholders within `config.yaml` and all files within the `templates/` directory.

## Output
1. **Variables Extraction**: Lists extracted variables and their values.
2. **Placeholder Substitution**: Indicates which files or directories underwent the substitution process.

## Considerations
* The application expects variable sets to be defined in YAML format within environment variables.
* Placeholders in the format `{{PLACEHOLDER_NAME}}` within target files will be substituted with the appropriate values from the sets.
* The application is capable of processing both individual files and entire directories. When a directory is specified, all files within the directory will be processed recursively.

## Logging
The application provides detailed logging for better traceability:
* INFO: General information about the ongoing process.
* FATAL: Critical errors that cause the application to exit.

## Dependencies
The `setsubst` application leverages several utilities and external libraries:
* **cicd_envsubst/utils/env_var**: Provides operations related to environment variables.
* **cicd_envsubst/utils/file**: Includes utilities for file operations like reading, writing, and searching.
* **cicd_envsubst/utils/path**: Offers path-related utilities.
* **github.com/goccy/go-yaml**: Used for YAML serialization and deserialization.
* **github.com/jessevdk/go-flags**: Parses CLI arguments.

## Contribution & Support
For bug reports, feature requests, and contributions, please open an issue in the application's repository. Feedback and contributions are always welcome!
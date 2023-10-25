## Overview
### envsubst
The `envsubst` application is designed to replace placeholders in files with their respective environment variables. Its purpose can be particularly useful in CI/CD pipelines where environment-specific values need to be injected into configuration files or other types of resources.
#### Usage
The application provides various command-line options to customize its behavior:

* **-p, --prefix**: Placeholder prefix. Default is "{{".
* **-s, --suffix**: Placeholder suffix. Default is "}}".
* **-m, --regex-mask**: Placeholder regex mask. This defines the pattern that the content of the placeholder should match. Default is "[A-Z_0-9]*".
* **-e, --env-file**: Source of environment variables. If specified, the application will attempt to read and set environment variables from this file.

##### Example:
```shell
cicd_envsubst -p "{{" -s "}}" -m "[A-Z_0-9]*" -e "path/to/env/file" path/to/file_or_directory
```
This command will process `path/to/file_or_directory`, replacing placeholders that match the default format with corresponding environment variables. It will also read environment variables from `path/to/env/file` and use them for replacements.
### envmake
The `envmake` application is designed to compile an environment (.env) file based on the placeholders found within specified files and directories. The placeholders are identified and then either replaced with their actual values from the existing environment variables, or retained as placeholders in the compiled environment file.
#### Usage
To execute the `envmake` application, use the following command structure:

```shell
envmake [destination_env_file] [path_to_file_or_directory ...]
```
Where
* **destination_env_file**: Path to the file where the compiled environment values should be saved.
* **path_to_file_or_directory**: Paths to files or directories that should be processed for placeholders. Multiple paths can be provided.

##### Example:
```shell
envmake ./output.env ./config.yaml ./settings/
```
This command will process the `config.yaml` file and all files within the `settings/` directory, identifying and extracting placeholders. The resulting .env file will be saved as `output.env`.
### set2secret
The `set2secret` application is designed to generate Kubernetes Secret manifests from a collection of environment variables referred to as "sets". The tool allows for the extraction of these sets, comparison to referenced files (for validation), and generation of a Kubernetes Secret manifest, which can then be applied to a Kubernetes cluster.
#### Usage
To execute the `set2secret` application, use the following command structure:

```shell
set2secret -s [variable_set_name ...] -n [k8s_secret_name] [-r path_to_reference_file ...]
```

Where
* **-s, --set**: Name(s) of the variable set(s) to extract from the environment. Multiple sets can be provided.
* **-n, --secret-name**: The desired name of the Kubernetes Secret.
* **-r, --ref**: (Optional) Paths to reference files containing placeholders to be validated against the extracted variable set.

##### Example:
```shell
set2secret -s SET1 SET2 -n my-k8s-secret -r ./config.yaml
```
This command will:

1. Extract variables from the environment based on the sets `SET1` and `SET2`.
2. Compare the extracted variables to placeholders within `config.yaml`.
3. Generate a Kubernetes Secret manifest named `my-k8s-secret`.

### setsubst
The `setsubst` application focuses on substituting environment variable placeholders within target files or directories based on the provided "sets". Given a set of environment variables (stored in YAML format), the application reads these variables and replaces placeholders within target files with the appropriate values from the sets.
#### Usage
To execute the `setsubst` application, use the following command structure:

```shell
setsubst -s [variable_set_name ...]
```
Where:
* **-s, --set**: Specifies the name(s) of the variable set(s) to extract from the environment.
<p></p>
After the set name(s), provide the paths to the target files or directories as positional arguments.

##### Example:
```shell
export SET1="
VARIABLE1: THIS is a value for VARIABLE1
VARIABLE2: |
  This is a value
  for VARIABLE2"

export SET2="
VARIABLE3: THIS is a value for VARIABLE3
VARIABLE4: |
  This is a value
  for VARIABLE4"

setsubst -s SET1 SET2 ./config.yaml ./templates/
```
This command will:
1. Extract variables from the environment based on the sets `SET1` and `SET2`.
2. Substitute placeholders within `config.yaml` and all files within the `templates/` directory.
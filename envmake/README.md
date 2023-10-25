## Overview
The `envmake` application is designed to compile an environment (.env) file based on the placeholders found within specified files and directories. The placeholders are identified and then either replaced with their actual values from the existing environment variables, or retained as placeholders in the compiled environment file.

## Features
1. Dynamically identifies and extracts placeholders from specified files and directories.
2. The application checks if the extracted placeholder matches an existing environment variable. If so, it will use the value of that environment variable. Otherwise, the placeholder remains unchanged in the generated .env file.
3. Supports the processing of individual files as well as entire directories.
4. Provides logging for better traceability.

## Usage
To execute the `envmake` application, use the following command structure:

```shell
envmake [destination_env_file] [path_to_file_or_directory ...]
```
Where
* **destination_env_file**: Path to the file where the compiled environment values should be saved.
* **path_to_file_or_directory**: Paths to files or directories that should be processed for placeholders. Multiple paths can be provided.

### Example:
```shell
envmake ./output.env ./config.yaml ./settings/
```
This command will process the `config.yaml` file and all files within the `settings/` directory, identifying and extracting placeholders. The resulting .env file will be saved as `output.env`.

## Placeholders
By default, the application identifies placeholders in the format `{{PLACEHOLDER_NAME}}`. These placeholders are extracted from the specified files and directories.

## Output
The generated `.env` file will contain lines in the format:

```
PLACEHOLDER_NAME=value
ANOTHER_PLACEHOLDER={{ANOTHER_PLACEHOLDER}}
```
If an environment variable matching `PLACEHOLDER_NAME` exists, its value will be used. Otherwise, the placeholder remains unchanged, as shown with `ANOTHER_PLACEHOLDER`.

Logging:
The application provides log outputs for better traceability:

* **INFO**: Lists the destination file and the files/directories being processed.
* **PROCESSING**: Indicates the start of the placeholder identification process.
* **DONE**: Indicates the end of the process, signifying that the destination .env file has been generated.

## Dependencies:
The `envmake` application relies on a few utilities provided within the `cicd_envsubst` package and some external libraries:

* **cicd_envsubst/utils/env_var**: For operations related to environment variables.
* **cicd_envsubst/utils/file**: For file operations.
* **cicd_envsubst/utils/path**: For path-related utilities.
* **cicd_envsubst/utils/placeholder**: For operations related to placeholders.
* **github.com/thoas/go-funk**: For certain array operations like uniqueness.

## Contribution & Support:
For bug reports, feature requests, and contributions, please open an issue in the application's repository. Feedback and contributions are always welcome!
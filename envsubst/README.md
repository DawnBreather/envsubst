## Overview
The `envsubst` application is designed to replace placeholders in files with their respective environment variables. Its purpose can be particularly useful in CI/CD pipelines where environment-specific values need to be injected into configuration files or other types of resources.

## Features
1. Customizable placeholder prefix and suffix: By default, the application identifies placeholders in the format `{{PLACEHOLDER_NAME}}`, but this can be customized.
2. Ability to provide a custom regex mask for the placeholder content.
3. Option to load environment variables from an external file.
4. Processes both individual files and entire directories, replacing placeholders recursively in all files found in a directory.
5. Provides logging to trace the progress of the replacements.

## Usage
The application provides various command-line options to customize its behavior:

* **-p, --prefix**: Placeholder prefix. Default is "{{".
* **-s, --suffix**: Placeholder suffix. Default is "}}".
* **-m, --regex-mask**: Placeholder regex mask. This defines the pattern that the content of the placeholder should match. Default is "[A-Z_0-9]*".
* **-e, --env-file**: Source of environment variables. If specified, the application will attempt to read and set environment variables from this file.

### Example:

```shell
cicd_envsubst -p "{{" -s "}}" -m "[A-Z_0-9]*" -e "path/to/env/file" path/to/file_or_directory
```
This command will process `path/to/file_or_directory`, replacing placeholders that match the default format with corresponding environment variables. It will also read environment variables from `path/to/env/file` and use them for replacements.

### Environment File Format
If you choose to use an external file to define environment variables (-e option), the expected format is:
```shell
VAR_NAME=value
ANOTHER_VAR=another_value
# This is a comment and will be ignored
```
### Logging
The application provides log outputs to trace its progress:

* **INFO**: Lists the files and directories being processed.
* **PROCESSING**: Indicates the start of the placeholder replacement process.
* **DONE**: Indicates the end of the placeholder replacement process.

## Dependencies
The application makes use of several external libraries:
* **github.com/DawnBreather/go-commons/file**: For file operations.
* **github.com/DawnBreather/go-commons/logger**: For logging.
* **github.com/DawnBreather/go-commons/path**: For path-related utilities.
* **github.com/jessevdk/go-flags**: For command-line argument parsing.

## Contribution & Support
For bug reports, features requests, and contributions, please open an issue in the application's repository. Your feedback is valuable and much appreciated!
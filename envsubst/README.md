## Overview
The `envsubst` application is designed to replace placeholders in files with their respective environment variables. Its purpose can be particularly useful in CI/CD pipelines where environment-specific values need to be injected into configuration files or other types of resources.

## Features
1. Customizable placeholder prefix and suffix: By default, the application identifies placeholders in the format `{{PLACEHOLDER_NAME}}`, but this can be customized.
2. Ability to provide a custom regex mask for the placeholder content.
3. Option to load environment variables from an external file.
4. Processes both individual files and entire directories, replacing placeholders recursively in all files found in a directory.
5. Provides logging to trace the progress of the replacements.

#### Usage
The application provides various command-line options to customize its behavior:

* **-p, --prefix**: Placeholder prefix. Default is "{{".
* **-x, --suffix**: Placeholder suffix. Default is "}}".
* **-m, --regex-mask**: Placeholder regex mask. This defines the pattern that the content of the placeholder should match. Default is "[A-Z_0-9]*".
* **-e, --env-file**: Source of environment variables. If specified, the application will attempt to read and set environment variables from this file.
* **-s, --set**: Specifies the name(s) of the variable set(s) to extract from the environment

```shell
envsubst -s [variable_set_name ...] -p [prefix] -x [suffix] -m [regexp_mask] [path_to_file_or_directory ...]
```

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

envsubst -s SET1 -s SET2 -p "{{" -x "}}" -m "[A-Z_0-9]*" -e "path/to/env/file" path/to/file_or_directory another/to/file_or_directory
```
This command will:
1. Extract variables from the env file `path/to/env/file`.
2. Extract variables from the environment based on the sets `SET1` and `SET2`.
3. Extract all environment variables from the environment.
4. Merge all the extracted variables (lowest priority to highest priority: env file, sets, environment).
5. Substitute placeholders within `path/to/file_or_directory` recursively.

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
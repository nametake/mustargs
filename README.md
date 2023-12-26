# mustargs

mustargs checks Go functions that do not adhere to the specified argument rules.

## Installation

```console
go install github.com/nametake/mustargs/cmd/mustargs@latest
```

## Usage

```console
go vet -vettool=`which mustargs` -mustargs.config=$(pwd)/config.yaml .
```

## Rule Definition

Rules are configured in a YAML file.

The overall structure of the rules that can be configured is shown below:

```yaml
---
# Multiple rules can be configured. Each rule is independent.
rules:
  # These are the target arguments. Multiple configurations are possible within each rule.
  # An error will occur if not all argument rules are followed.
  - args:
      # 'type' is the name of the type that should be included in the function's arguments.
      # This is a required field.
      - type: Context
        # 'index' is the position of the argument. 0 represents the first argument of the function.
        # If not specified, it will not result in an error if included anywhere.
        index: 0 
        pkg: context # 'pkg' is the package name of the argument.
      - type: TenantID
        index: 1
        pkg: github.com/nametake/mustargs/domain
        # 'pkg_name' is used when giving a name to the imported package.
        # If not specified, it will follow the default import rules.
        pkg_name: dm
        # 'ptr' specifies whether the argument type is a pointer.
        ptr: true
        is_array: true
    # These are the patterns for the functions targeted by the rule.
    # Multiple patterns can be specified as a list.
    # Patterns support regular expressions.
    # Each pattern is an AND condition, and the list of patterns is an OR condition.
    #
    # file_patterns: # File patterns.
    #   - usecase/tenant_*.go
    #   - mysql/tenant_*.go
    # ignore_file_patterns: # Patterns to exclude files.
    #   - .*_gen.go
    # func_patterns: # Function name patterns.
    #   - Get.*
    #   - Update.*
    # ignore_func_patterns: # Patterns to exclude function names.
    #   - ^New.*
    # recv_patterns: # Receiver patterns.
    #   - ^Tenant.*Usecase$
    #   - ^Tenant.*DB$
    # ignore_recv_patterns: # Patterns to exclude receivers.
    #   - ^TenantDBDebug$
```

If you want to examine detailed behavior, please check the test data in the `testdata/src` directory.

## Unsupported

The following argument types are not supported:

- `map`
- `func`
- `chan`
- `interface`
- `ellipsis`

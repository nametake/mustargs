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

The overall structure of the rules that can be configured is shown below. Parameters without description are all optional.

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
        # A negative value can be specified to indicate arguments from the end.
        index: 0
        pkg: context # 'pkg' is the package name of the argument.
      - type: TenantID
        index: 1
        pkg: github.com/nametake/mustargs/domain
        # 'pkg_name' is used when giving a name to the imported package.
        # If not specified, it will follow the default import rules.
        pkg_name: dm
        # 'is_ptr' specifies whether the argument type is a pointer.
        is_ptr: true
        # 'is_array' specifies whether the argument type is a array or slice.
        is_array: true
    # These are the patterns for the functions targeted by the rule.
    # Multiple patterns can be specified as a list.
    # Patterns support regular expressions.
    # Each pattern is an AND condition, and the list of patterns is an OR condition.
    file_patterns: # File patterns.
      - usecase/tenant_*.go
      - mysql/tenant_*.go
    ignore_file_patterns: # Patterns to ignore files.
      - .*_gen.go
    func_patterns: # Function name patterns.
      - Get.*
      - Update.*
    ignore_func_patterns: # Patterns to ignore function names.
      - ^New.*
    recv_patterns: # Receiver patterns.
      - ^Tenant.*Usecase$
      - ^Tenant.*DB$
    ignore_recv_patterns: # Patterns to ignore receivers.
      - ^TenantDBDebug$
```

## Example

```yaml
---
rules:
  - args:
      - type: Context
        pkg: context
        index: 0
      - type: TenantID
        index: 1
    recv_patterns:
      - ^Usecase$
  - args:
      - type: Context
        pkg: context
        index: 0
      - type: Tx
        pkg: database/sql
        index: 1
        is_ptr: true
    recv_patterns:
      - ^DB$
  - args:
      - type: int
        index: -1
      - type: int
        index: -2
    recv_patterns:
      - ^DB$
    func_patterns:
      - ^GetMultiple.*
```

```go
package example

import (
	"context"
	"database/sql"
)

type TenantID string

type Usecase struct{}

func (u *Usecase) GetUser(ctx context.Context, tenantID TenantID, userID string) {
}

func (u *Usecase) GetPost(ctx context.Context, userID string) { // ERROR
}

type DB struct{}

func (db *DB) GetUser(ctx context.Context, tx *sql.Tx, tenantID TenantID, userID string) {
}

func (db *DB) GetPost(ctx context.Context, tenantID TenantID, postID string) { // ERROR
}

func (db *DB) GetMultipleUsers(ctx context.Context, tx *sql.Tx, tenantID TenantID, limit, offset int) {
}

func (db *DB) GetMultiplePosts(ctx context.Context, tx *sql.Tx, tenantID TenantID) { // ERROR
}
```

If you want to examine detailed behavior, please check the test data in the `testdata/src` directory.

## Unsupported

The following argument types are not supported:

- `map`
- `func`
- `chan`
- `interface`
- `ellipsis`

# Package [cloudeng.io/linux/keyrings](https://pkg.go.dev/cloudeng.io/linux/keyrings?tab=doc)

```go
import cloudeng.io/linux/keyrings
```


## Types
### Type Option
```go
type Option func(o *options)
```
Option represents an option for configuring a keyrings.T

### Functions

```go
func WithKeyring(keyring keyctl.Keyring) Option
```
WithKeyring specifies the keyring to use for storing and retrieving secrets.




### Type T
```go
type T struct {
	// contains filtered or unexported fields
}
```
T provides access to linux kernel keyrings.

### Functions

```go
func New(opts ...Option) *T
```
New creates a new keyrings.T. If WithKeyring is not specified then the
session keyring is used.



### Methods

```go
func (s *T) Delete(ctx context.Context, name string) error
```
Delete removes a secret from the kernel keyring.


```go
func (s *T) ReadFile(name string) ([]byte, error)
```


```go
func (s *T) ReadFileCtx(ctx context.Context, name string) ([]byte, error)
```
ReadFileCtx reads a secret from the kernel keyring.


```go
func (s *T) WriteFile(name string, data []byte, mode fs.FileMode) error
```


```go
func (s *T) WriteFileCtx(ctx context.Context, name string, data []byte, _ fs.FileMode) error
```
WriteFileCtx writes a secret to the kernel keyring. The fs.FileMode
parameter is ignored.








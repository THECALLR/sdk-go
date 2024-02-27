sdk-go
======

SDK in Go for the CALLR API.

Works with Go 1.22+, using standard packages only.

## Quick start

```go
import callr "github.com/THECALLR/sdk-go/v2"

func main() {
    // Api Key Auth (use the customer portal to generate keys)
    api := callr.NewWithAPIKeyAuth("key")
    
    result, err := api.Call(context.Background(), "method", params...)
```

## Usage

### Login-As

If you have sub accounts and want to manage them with your master account, you can use the Login-As feature.

```go
import callr "github.com/THECALLR/sdk-go/v2"

func main() {
    // Api Key Auth (use the customer portal to generate keys)
    api := callr.NewWithAPIKeyAuth("key") // master account key

    if err := api.SetLoginAsSubAccountRef("<subAccountRef>"); err != nil {
        log.Fatalf("[error] cannot login as: %s\n", err)
    }

    // all following calls will be done as the sub account
    result, err := api.Call(context.Background(), "method", params...)
```

### Return and error handling

The SDK returns a json.RawMessage, which can be unmarshalled to the expected type, or an error.
The error may be of type *JSONRPCError if the error comes from the API, of type *HTTPError if the error comes from the HTTP transport, or a native error otherwise.

```go
	result, err := api.Call(ctx, "sms.send", "SMS", os.Args[1], "Hello, world", nil)

	// error management
	if err != nil {
		switch e := err.(type) {
		case *callr.JSONRPCError:
			slog.Error("JSON-RPC Error",
				"code", e.Code,
				"message", e.Message,
				"data", e.Data)
		case *callr.HTTPError:
			slog.Error("HTTP Error",
				"code", e.Code,
				"message", e.Message)
		default:
			slog.Error("Other error", "error", err)
		}
		os.Exit(1)
	}

	// the sms.send JSON-RPC method returns a string
	var hash string

	if err := json.Unmarshal(result, &hash); err != nil {
		slog.Error("Error unmarshalling result", "error", err)
		os.Exit(1)
	}

	slog.Info("SMS sent", "hash", hash)
```

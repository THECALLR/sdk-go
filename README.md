sdk-go
======

SDK in Go for the CALLR API.

Works with Go 1.11+ (this is a go module), using standard packages only.

## Quick start

```go
import callr "github.com/THECALLR/sdk-go"

func main() {
    // use Api Key Auth (recommended) - use the customer portal to generate keys
    api := callr.NewWithAPIKeyAuth("key")
    
    result, err := api.Call("method", params...)
```

## Usage

### Broadcast

#### Broadcast messages to a target

```go
target := map[string]interface{}{
    "number": "+33123456789",
    "timeout": 30,
}
messages := []interface{}{
    131,
    132,
    "TTS|TTS_EN-GB_SERENA|Hello world! how are you ? I hope you enjoy this call. good bye."
}

options := map[string]interface{}{
    "cdr_field": "userData",
    "cli": "BLOCKED",
    "loop": 2,
}

result, err := api.Call("calls.broadcast_1", target, messages, options)
```

##### Without options

```go
target := map[string]interface{}{
    "number": "+33123456789",
    "timeout": 30,
}

messages := []interface{}{
    131,
    132,
    134,
}

result, err := api.Call("calls.broadcast_1", target, messages, nil)
```

*Method*
* [calls.broadcast_1](https://www.callr.com/docs/api/services/calls/#calls.broadcast_1)

*Objects*
* [Target](https://www.callr.com/docs/objects/#Target)
* [Calls.Broadcast1.Options](https://www.callr.com/docs/objects/#Calls.Broadcast1.Options)

********************************************************************************

### REALTIME

#### Create a REALTIME app with a callback URL

```go
options := map[string]interface{}{
    "url": "http://yourdomain.com/realtime_callback_url",
}

result, err := api.Call("apps.create", "REALTIME10", "Your app name", options)
```

*Method*
* [apps.create](https://www.callr.com/docs/api/services/apps/#apps.create)

*Objects*
* [REALTIME10](https://www.callr.com/docs/objects/#REALTIME10)
* [App](https://www.callr.com/docs/objects/#App)

#### Start a REALTIME outbound call

```go
target := map[string]interface{}{
    "number": "+33132456789",
    "timeout": 30,
}

callOptions := map[string]interface{}{
    "cdr_field": "42",
    "cli": "BLOCKED",
}

result, err := api.Call("calls.realtime", "appHash", target, callOptions)
```

*Method*
* [calls.realtime](https://www.callr.com/docs/api/services/calls/#calls.realtime)

*Objects*
* [Target](https://www.callr.com/docs/objects/#Target)
* [REALTIME10.Call.Options](https://www.callr.com/docs/objects/#REALTIME10.Call.Options)

#### Inbound Calls - Assign a phone number to a REALTIME app

```go
result, err := api.Call("apps.assign_did", "appHash", "DID ID")
```

*Method*
* [apps.assign_did](https://www.callr.com/docs/api/services/apps/#apps.assign_did)

*Objects*
* [App](https://www.callr.com/docs/objects/#App)
* [DID](https://www.callr.com/docs/objects/#DID)

********************************************************************************

### Sending SMS

#### Without options

```go
result, err := api.Call("sms.send", "SMS", "+33123456789", "Hello, world", nil)
```

*Method*
* [sms.send](https://www.callr.com/docs/api/services/sms/#sms.send)

#### Personalized sender

> Your sender must have been authorized and respect the [sms_sender](https://www.callr.com/docs/formats/#sms_sender) format

```go
result, err := api.Call("sms.send", "Your Brand", "+33123456789", "Hello world!", nil)
```

*Method*
* [sms.send](https://www.callr.com/docs/api/services/sms/#sms.send)

#### If you want to receive replies, do not set a sender - we will automatically use a shortcode

```go
result, err := api.Call("sms.send", "", "+33123456789", "Hello world!", nil)
```

*Method*
- [sms.send](https://www.callr.com/docs/api/services/sms/#sms.send)

#### Force GSM encoding

```go
optionSMS := map[string]interface{}{
    "force_encoding": "GSM",
}

result, err := api.Call("sms.send", "", "+33123456789", "Hello world!", optionSMS)
```

*Method*
* [sms.send](https://www.callr.com/docs/api/services/sms/#sms.send)

*Objects*
* [SMS.Options](https://www.callr.com/docs/objects/#SMS.Options)

#### Long SMS (availability depends on carrier)

```go
var text bytes.Buffer

text.WriteString("Some super mega ultra long text to test message longer than 160 characters ")
text.WriteString("Some super mega ultra long text to test message longer than 160 characters ")
text.WriteString("Some super mega ultra long text to test message longer than 160 characters")

result, err := api.Call("sms.send", "SMS", "+33123456789", text.String(), nil)
```

*Method*
* [sms.send](https://www.callr.com/docs/api/services/sms/#sms.send)

#### Specify your SMS nature (alerting or marketing)

```go
optionSMS := map[string]interface{}{
    "nature": "ALERTING",
}

result, err := api.Call("sms.send", "SMS", "+33123456789", "Hello world!", optionSMS)
```

*Method*
* [sms.send](https://www.callr.com/docs/api/services/sms/#sms.send)

*Objects*
* [SMS.Options](https://www.callr.com/docs/objects/#SMS.Options)

#### Custom data

```go
optionSMS := map[string]interface{}{
    "user_data": "42",
}

result, err := api.Call("sms.send", "SMS", "+33123456789", "Hello world!", optionSMS)
```

*Method*
* [sms.send](https://www.callr.com/docs/api/services/sms/#sms.send)

*Objects*
* [SMS.Options](https://www.callr.com/docs/objects/#SMS.Options)


#### Delivery Notification - set webhook endpoint to receive notifications

```go
optionSMS := map[string]interface{}{
       	"webhook": map[string]interface{}{ 
			"endpoint":"http://yourdomain.com/webhook_endpoint", 
		},
    }


result, err := api.Call("sms.send", "SMS", "+33123456789", "Hello world!", optionSMS)
```

*Method*
* [sms.send](https://www.callr.com/docs/api/services/sms/#sms.send)

*Objects*
* [SMS.Options](https://www.callr.com/docs/objects/#SMS.Options)


### Inbound SMS - set webhook endpoint to receive inbound messages (MO) and replies

> **Do not set a sender if you want to receive replies** - we will automatically use a shortcode.

```go
optionSMS := map[string]interface{}{
       	"webhook": map[string]interface{}{ 
			"endpoint":"http://yourdomain.com/webhook_endpoint", 
		},
    }

result, err := api.Call("sms.send", "", "+33123456789", "Hello world!", optionSMS)
```

*Method*
* [sms.send](https://www.callr.com/docs/api/services/sms/#sms.send)

*Objects*
* [SMS.Options](https://www.callr.com/docs/objects/#SMS.Options)


### Get an SMS
```go
result, err := api.Call("sms.get", "SMSHASH")
```

*Method*
* [sms.get](https://www.callr.com/docs/api/services/sms/#sms.get)

*Objects*
* [SMS](https://www.callr.com/docs/objects/#SMS)

********************************************************************************

### DIDs

#### List available countries with DID availability
```go
result, err := api.Call("did/areacode.countries")
```

*Method*
* [did/areacode.countries](https://www.callr.com/docs/api/services/did/areacode/#did/areacode.countries)

*Objects*
* [DID.Country](https://www.callr.com/docs/objects/#DID.Country)

#### Get area codes available for a specific country and DID type

```go
result, err := api.Call("did/areacode.get_list", "US", nil)
```

*Method*
* [did/areacode.get_list](https://www.callr.com/docs/api/services/did/areacode/#did/areacode.get_list)

*Objects*
* [DID.AreaCode](https://www.callr.com/docs/objects/#DID.AreaCode)

#### Get DID types available for a specific country
```go
result, err := api.Call("did/areacode.types", "US")
```

*Method*
* [did/areacode.types](https://www.callr.com/docs/api/services/did/areacode/#did/areacode.types)

*Objects*
* [DID.Type](https://www.callr.com/docs/objects/#DID.Type)

#### Buy a DID (after a reserve)

```go
result, err := api.Call("did/store.buy_order", "OrderToken")
```

*Method*
* [did/store.buy_order](https://www.callr.com/docs/api/services/did/store/#did/store.buy_order)

*Objects*
* [DID.Store.BuyStatus](https://www.callr.com/docs/objects/#DID.Store.BuyStatus)

#### Cancel your order (after a reserve)

```go
result, err := api.Call("did/store.cancel_order", "OrderToken")
```

*Method*
* [did/store.cancel_order](https://www.callr.com/docs/api/services/did/store/#did/store.cancel_order)

#### Cancel a DID subscription

```go
result, err := api.Call("did/store.cancel_subscription", "DID ID")
```

*Method*
* [did/store.cancel_subscription](https://www.callr.com/docs/api/services/did/store/#did/store.cancel_subscription)

#### View your store quota status

```go
result, err := api.Call("did/store.get_quota_status")
```

*Method*
* [did/store.get_quota_status](https://www.callr.com/docs/api/services/did/store/#did/store.get_quota_status)

*Objects*
* [DID.Store.QuotaStatus](https://www.callr.com/docs/objects/#DID.Store.QuotaStatus)

#### Get a quote without reserving a DID

```go
result, err := api.Call("did/store.get_quote", 0, "GOLD", 1)
```

*Method*
* [did/store.get_quote](https://www.callr.com/docs/api/services/did/store/#did/store.get_quote)

*Objects/
* [DID.Store.Quote](https://www.callr.com/docs/objects/#DID.Store.Quote)

#### Reserve a DID

```go
result, err := api.Call("did/store.reserve", 0, "GOLD", 1, "RANDOM")
```

*Method*
* [did/store.reserve](https://www.callr.com/docs/api/services/did/store/#did/store.reserve)

*Objects*
* [DID.Store.Reservation](https://www.callr.com/docs/objects/#DID.Store.Reservation)

#### View your order

```go
result, err := api.Call("did/store.view_order", "OrderToken")
```

*Method*
* [did/store.buy_order](https://www.callr.com/docs/api/services/did/store/#did/store.view_order)

*Objects*
* [DID.Store.Reservation](https://www.callr.com/docs/objects/#DID.Store.Reservation)

********************************************************************************

### Media

#### List your medias

```go
result, err := api.Call("media/library.get_list", nil)
```

*Method*
* [media/library.get_list](https://www.callr.com/docs/api/services/media/library/#media/library.get_list)

#### Create an empty media

```go
result, err := api.Call("media/library.create", "name")
```

*Method*
* [media/library.create](https://www.callr.com/docs/api/services/media/library/#media/library.create)

#### Upload a media

```go
mediaID := 0

result, err := api.Call("media/library.set_content_from_file", mediaID, "imported temporary file name")
```

*Method*
* [media/library.set_content_from_file](https://www.callr.com/docs/api/services/media/library/#media/library.set_content_from_file)

#### Use Text-to-Speech

```go
mediaID := 0

result, err := api.Call("media/tts.set_content", mediaID, "Hello world!", "TTS-EN-GB_SERENA", nil)
```

*Method*
* [media/tts.set_content](https://www.callr.com/docs/api/services/media/tts/#media/tts.set_content)

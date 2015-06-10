sdk-go
======

SDK in Go for CALLR API

## Quick start
Install the Go way

    go get github.com/THECALLR/sdk-go

Or get sources from Github

## Initialize your code

```go
import "github.com/THECALLR/sdk-go"
```

## Usage
### Sending SMS

#### Without options

```go
result, err := callr.Call("sms.send", "CALLR", "+33123456789", "Hello, world", nil)
```

*Method*
* [sms.send](http://thecallr.com/docs/api/services/sms/#sms.send)

#### Personalized sender

> Your sender must have been authorized and respect the [sms_sender](http://thecallr.com/docs/formats/#sms_sender) format

```go
result, err := callr.Call("sms.send", "Your Brand", "+33123456789", "Hello world!", nil)
```

*Method*
* [sms.send](http://thecallr.com/docs/api/services/sms/#sms.send)

#### If you want to receive replies, do not set a sender - we will automatically use a shortcode

```go
result, err := callr.Call("sms.send", "", "+33123456789", "Hello world!", nil)
```

*Method*
- [sms.send](http://thecallr.com/docs/api/services/sms/#sms.send)

#### Force GSM encoding

```go
optionSMS := map[string]interface{}{
    "force_encoding": "GSM",
}

result, err := callr.Call("sms.send", "", "+33123456789", "Hello world!", optionSMS)
```

*Method*
* [sms.send](http://thecallr.com/docs/api/services/sms/#sms.send)

*Objects*
* [SMS.Options](http://thecallr.com/docs/objects/#SMS.Options)

#### Long SMS (availability depends on carrier)

```go
var text bytes.Buffer

text.WriteString("Some super mega ultra long text to test message longer than 160 characters ")
text.WriteString("Some super mega ultra long text to test message longer than 160 characters ")
text.WriteString("Some super mega ultra long text to test message longer than 160 characters")

result, err := callr.Call("sms.send", "CALLR", "+33123456789", text.String(), nil)
```

*Method*
* [sms.send](http://thecallr.com/docs/api/services/sms/#sms.send)

#### Specify your SMS nature (alerting or marketing)

```go
optionSMS := map[string]interface{}{
    "nature": "ALERTING",
}

result, err := callr.Call("sms.send", "CALLR", "+33123456789", "Hello world!", optionSMS)
```

*Method*
* [sms.send](http://thecallr.com/docs/api/services/sms/#sms.send)

*Objects*
* [SMS.Options](http://thecallr.com/docs/objects/#SMS.Options)

#### Custom data

```go
optionSMS := map[string]interface{}{
    "user_data": "42",
}

result, err := callr.Call("sms.send", "CALLR", "+33123456789", "Hello world!", optionSMS)
```

*Method*
* [sms.send](http://thecallr.com/docs/api/services/sms/#sms.send)

*Objects*
* [SMS.Options](http://thecallr.com/docs/objects/#SMS.Options)


#### Delivery Notification - set URL to receive notifications

```go
optionSMS := map[string]interface{}{
    "push_dlr_enabled": "42",
    "push_dlr_url": "http://yourdomain.com/push_delivery_path",
    // "push_dlr_url_auth": "login:password", // needed if you use Basic HTTP Authentication
}

result, err := callr.Call("sms.send", "CALLR", "+33123456789", "Hello world!", optionSMS)
```

*Method*
* [sms.send](http://thecallr.com/docs/api/services/sms/#sms.send)

*Objects*
* [SMS.Options](http://thecallr.com/docs/objects/#SMS.Options)


### Inbound SMS - set URL to receive inbound messages (MO) and replies

> **Do not set a sender if you want to receive replies** - we will automatically use a shortcode.

```go
optionSMS := map[string]interface{}{
    "push_mo_enabled": true,
    "push_mo_url": "http://yourdomain.com/mo_delivery_path",
    // "push_mo_url_auth": "login:password" // needed if you use Basic HTTP Authentication
}

result, err := callr.Call("sms.send", "", "+33123456789", "Hello world!", optionSMS)
```

*Method*
* [sms.send](http://thecallr.com/docs/api/services/sms/#sms.send)

*Objects*
* [SMS.Options](http://thecallr.com/docs/objects/#SMS.Options)


### Get an SMS
```go
result, err := callr.Call("sms.get", "SMSHASH")
```

*Method*
* [sms.get](http://thecallr.com/docs/api/services/sms/#sms.get)

*Objects*
* [SMS](http://thecallr.com/docs/objects/#SMS)

### SMS Global Settings

#### Get settings
```go
result, err := callr.Call("sms.get_settings")
```

*Method*
* [sms.get_settings](http://thecallr.com/docs/api/services/sms/#sms.get_settings)

*Objects*
* [SMS.settings](http://thecallr.com/docs/objects/#SMS.Settings)


#### Set settings

> Add options that you want to change in the object

```go
settings := map[string]interface{}{
    "push_dlr_enabled": true,
    "push_dlr_url": "http://yourdomain.com/push_delivery_path",
    "push_mo_enabled": true,
    "push_mo_url": "http://yourdomain.com/mo_delivery_path"
}

result, err := callr.Call("sms.set_settings", settings)
```

> Returns the updated settings.

*Method*
* [sms.set_settings](http://thecallr.com/docs/api/services/sms/#sms.set_settings)

*Objects*
* [SMS.settings](http://thecallr.com/docs/objects/#SMS.Settings)

********************************************************************************

### REALTIME

#### Create a REALTIME app with a callback URL

```go
options := map[string]interface{}{
    "url": "http://yourdomain.com/realtime_callback_url",
}

result, err := callr.Call("apps.create", "REALTIME10", "Your app name", options)
```

*Method*
* [apps.create](http://thecallr.com/docs/api/services/apps/#apps.create)

*Objects*
* [REALTIME10](http://thecallr.com/docs/objects/#REALTIME10)
* [App](http://thecallr.com/docs/objects/#App)

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

result, err := callr.Call("dialr/call.realtime", "appHash", target, callOptions)
```

*Method*
* [dialr/call.realtime](http://thecallr.com/docs/api/services/dialr/call/#dialr/call.realtime)

*Objects*
* [Target](http://thecallr.com/docs/objects/#Target)
* [REALTIME10.Call.Options](http://thecallr.com/docs/objects/#REALTIME10.Call.Options)

#### Inbound Calls - Assign a phone number to a REALTIME app

```go
result, err := callr.call("apps.assign_did", "appHash", "DID ID")
```

*Method*
* [apps.assign_did](http://thecallr.com/docs/api/services/apps/#apps.assign_did)

*Objects*
* [App](http://thecallr.com/docs/objects/#App)
* [DID](http://thecallr.com/docs/objects/#DID)

********************************************************************************

### DIDs

#### List available countries with DID availability
```go
result, err := callr.Call("did/areacode.countries")
```

*Method*
* [did/areacode.countries](http://thecallr.com/docs/api/services/did/areacode/#did/areacode.countries)

*Objects*
* [DID.Country](http://thecallr.com/docs/objects/#DID.Country)

#### Get area codes available for a specific country and DID type

```go
result, err := callr.Call("did/areacode.get_list", "US", nil)
```

*Method*
* [did/areacode.get_list](http://thecallr.com/docs/api/services/did/areacode/#did/areacode.get_list)

*Objects*
* [DID.AreaCode](http://thecallr.com/docs/objects/#DID.AreaCode)

#### Get DID types available for a specific country
```go
result, err := callr.Call("did/areacode.types", "US")
```

*Method*
* [did/areacode.types](http://thecallr.com/docs/api/services/did/areacode/#did/areacode.types)

*Objects*
* [DID.Type](http://thecallr.com/docs/objects/#DID.Type)

#### Buy a DID (after a reserve)

```go
result, err := callr.Call("did/store.buy_order", "OrderToken")
```

*Method*
* [did/store.buy_order](http://thecallr.com/docs/api/services/did/store/#did/store.buy_order)

*Objects*
* [DID.Store.BuyStatus](http://thecallr.com/docs/objects/#DID.Store.BuyStatus)

#### Cancel your order (after a reserve)

```go
result, err := callr.Call("did/store.cancel_order", "OrderToken")
```

*Method*
* [did/store.cancel_order](http://thecallr.com/docs/api/services/did/store/#did/store.cancel_order)

#### Cancel a DID subscription

```go
result, err := callr.Call("did/store.cancel_subscription", "DID ID")
```

*Method*
* [did/store.cancel_subscription](http://thecallr.com/docs/api/services/did/store/#did/store.cancel_subscription)

#### View your store quota status

```go
result, err := callr.Call("did/store.get_quota_status")
```

*Method*
* [did/store.get_quota_status](http://thecallr.com/docs/api/services/did/store/#did/store.get_quota_status)

*Objects*
* [DID.Store.QuotaStatus](http://thecallr.com/docs/objects/#DID.Store.QuotaStatus)

#### Get a quote without reserving a DID

```go
result, err := callr.Call("did/store.get_quote", 0, "GOLD", 1)
```

*Method*
* [did/store.get_quote](http://thecallr.com/docs/api/services/did/store/#did/store.get_quote)

*Objects/
* [DID.Store.Quote](http://thecallr.com/docs/objects/#DID.Store.Quote)

#### Reserve a DID

```go
result, err := callr.Call("did/store.reserve", 0, "GOLD", 1, "RANDOM")
```

*Method*
* [did/store.reserve](http://thecallr.com/docs/api/services/did/store/#did/store.reserve)

*Objects*
* [DID.Store.Reservation](http://thecallr.com/docs/objects/#DID.Store.Reservation)

#### View your order

```go
result, err := callr.Call("did/store.view_order", "OrderToken")
```

*Method*
* [did/store.buy_order](http://thecallr.com/docs/api/services/did/store/#did/store.view_order)

*Objects*
* [DID.Store.Reservation](http://thecallr.com/docs/objects/#DID.Store.Reservation)

********************************************************************************

### Conferencing

#### Create a conference room

```go
params := map[string]interface{}{
    "open": true,
}
access := []interface{}{}

result, err := callr.Call("conference/10.create_room", "room name", params, access)
```

*Method*
* [conference/10.create_room](http://thecallr.com/docs/api/services/conference/10/#conference/10.create_room)

*Objects*
* [CONFERENCE10](http://thecallr.com/docs/objects/#CONFERENCE10)
* [CONFERENCE10.Room.Access](http://thecallr.com/docs/objects/#CONFERENCE10.Room.Access)

#### Assign a DID to a room

```go
result, err := callr.Call("conference/10.assign_did", "Room ID", "DID ID")
```

*Method*
* [conference/10.assign_did](http://thecallr.com/docs/api/services/conference/10/#conference/10.assign_did)

#### Create a PIN protected conference room

```go
params := map[string]interface{}{
    "open": true,
}
access := []interface{}{
    map[string]interface{}{
        "pin": "1234",
        "level": "GUEST",
    },
    map[string]interface{}{
        "pin": "4321",
        "level": "ADMIN",
        "phone_number": "+33123456789",
    },
}

result, err := callr.Call("conference/10.create_room", "room name", params, access)
```

*Method*
* [conference/10.create_room](http://thecallr.com/docs/api/services/conference/10/#conference/10.create_room)

*Objects*
* [CONFERENCE10](http://thecallr.com/docs/objects/#CONFERENCE10)
* [CONFERENCE10.Room.Access](http://thecallr.com/docs/objects/#CONFERENCE10.Room.Access)

#### Call a room access

```go
result, err := callr.Call("conference/10.call_room_access", "Room Access ID", "BLOCKED", true)
```

*Method*
* [conference/10.call_room_access](http://thecallr.com/docs/api/services/conference/10/#conference/10.call_room_access)

********************************************************************************

### Media

#### List your medias

```go
result, err := callr.Call("media/library.get_list", nil)
```

*Method*
* [media/library.get_list](http://thecallr.com/docs/api/services/media/library/#media/library.get_list)

#### Create an empty media

```go
result, err := callr.Call("media/library.create", "name")
```

*Method*
* [media/library.create](http://thecallr.com/docs/api/services/media/library/#media/library.create)

#### Upload a media

```go
media_id := 0

result, err := callr.Call("media/library.set_content", media_id, "text", "base64_audio_data")
```

*Method*
* [media/library.set_content](http://thecallr.com/docs/api/services/media/library/#media/library.set_content)

#### Use Text-to-Speech

```go
media_id := 0

result, err := callr.Call("media/tts.set_content", media_id, "Hello world!", "TTS-EN-GB_SERENA", nil)
```

*Method*
* [media/tts.set_content](http://thecallr.com/docs/api/services/media/tts/#media/tts.set_content)

********************************************************************************

### CDR

#### Get inbound or outbound CDRs
```go
from := "YYYY-MM-DD HH:MM:SS"
to := "YYYY-MM-DD HH:MM:SS"

result, err := callr.Call("cdr.get", "OUT", from, to, nil, nil)
```

*Method*
* [cdr.get](http://thecallr.com/docs/api/services/cdr/#cdr.get)

*Objects*
* [CDR.In](http://thecallr.com/docs/objects/#CDR.In)
* [CDR.Out](http://thecallr.com/docs/objects/#CDR.Out)

********************************************************************************

### SENDR

#### Broadcast messages to a target (BETA)

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

result, err := callr.Call("sendr/simple.broadcast_1", target, messages, options)
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

result, err := callr.Call("sendr/simple.broadcast_1", target, messages, nil)
```

*Method*
* [sendr/simple.broadcast_1](http://thecallr.com/docs/api/services/sendr/simple/#sendr/simple.broadcast_1)

*Objects*
* [Target](http://thecallr.com/docs/objects/#Target)
* [SENDR.Simple.Broadcast1.Options](http://thecallr.com/docs/objects/#SENDR.Simple.Broadcast1.Options)

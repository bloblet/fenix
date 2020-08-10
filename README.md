# Fenix v6

[![Codeac](https://static.codeac.io/badges/2-281254941.svg "Codeac.io")](https://app.codeac.io/github/bloblet/fenix)
![Tests](https://github.com/bloblet/fenix/workflows/Tests/badge.svg)

# Fenix 6.0.1 API Documentation

**This document is a work in progress and is subject to change!**

All endpoints require either a `Authorization` header with the [`Basic`](https://developer.mozilla.org/en-US/docs/Web/HTTP/Authentication#Basic_authentication_scheme) authentication scheme or a `X-Token` header containing an account token.  

(Note `/create` only accepts the `Basic` method)  

Fenix has two APIs, a [WebSockets](https://developer.mozilla.org/en-US/docs/Glossary/WebSockets) API, and a [REST](https://developer.mozilla.org/en-US/docs/Glossary/REST) HTTP API.  The following information is shared by both APIs.

# Response format
All responses will be a JSON object containing a `s` boolean.  If fenix encounters errors processing your request, this will be false and an error code will be added to the response.  Otherwise, the response data (if any) will be in the `d` key.

```json
{
    "s": true,
    "d": {
        "id": "abcdef-ghijk-lmnop"
    }
}
```

```json
{
    "s": false,
    "e": "ERR_USEREXISTS",
    "m": "That user already exists!"
}
```

### Globally raised errors

The following errors may be raised on any API endpoint:
 - `ERR_INVALIDLOGIN` when either the `Basic` auth is invalid or `token` is invalid.
 - `ERR_INVALIDREQUEST` when required parameters are missing or are an invalid
    type.
 - `ERR_INTERNALERROR` when something breaks in the server.

# HTTP API endpoints

## POST `/v6.0.1/create`

Creates a new Fenix user.

Parameters:
- `username`:  The user's desired username

Errors raised:
- `ERR_USERNAMETOOLONG` when `username` is above 32 characters (`400 Bad Request`)
- `ERR_USEREXISTS` when the email already belongs to a user (`403 Forbidden`)

Returns:

[User](#User)

# Dataclasses

# User
Parameters:
- `ID`: User's ID (String)
- `Token`: User's access token (String)
- `Email`: User's associated email (String)
- `Username`: User's current username (String)
- `Discriminator`: User's current discriminator (String)
- `Servers`: 
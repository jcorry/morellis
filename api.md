---
title: Morellis API Docs
---
A search engine for ice cream parlor flavors.

## Authentication
All REST API endpoints are authenticated using JWT. User management is presently beyond the scope
of this API and will be performed manually by the API developers. To add or delete a user contact
John <[jcorry@gmail.com](mailto:jcorry@gmail.com)>

### `POST /authentication`
Exchange credentials for a JWT that will be used to authenticate subsequent API requests.
#### Request body
```$xslt
{
    "email": "scooper@morellis.com",
    "password": "MyP@55w0rd"
}
```
#### Response
```$xslt
{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
}
```

The token value will be included in the `authorization` request header as a `Bearer ...` token.

## Flavors
### `GET /flavor`
Gets a list of the available flavors. 

### `POST /flavor`
Adds a new flavor to the catalog.

`ingredients` is an array of strings, each being an ingredient contributing to this ice cream's unique
flavor profile. Ice cream obviously includes 'cream' and 'sugar' so these can be ommitted. Use words
that indicate the flavor: 'chcoloate', 'macadamia nut', 'peanut butter', 'strawberry', etc. The more complete
this list is, the better we're able to offer potential matches to customers seeking flavors matching their
preferences.
#### Request body
```$xslt
{
    "name": "The name of the flavor",
    "description": "A textual description of the flavor",
    "ingredients": [],
    "created": "2017-09-14 00:00:32"
}
```

### `DELETE /flavor/{flavorID}`
Removes a flavor from the store flavor portfolio.
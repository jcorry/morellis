#Morellis API Docs

A notifications API for informing customers when their favorite ice cream flavor is available at Morellis Gourmet Ice Cream, the best ice cream shop in Atlanta!

## Authentication
All REST API endpoints are authenticated using JWT. User management is presently beyond the scope
of this API and will be performed manually by the API developers. To add or delete a user contact
John <[jcorry@gmail.com](mailto:jcorry@gmail.com)>

### `POST /user/authenticate` (unimplemented)
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

The token value will be included in the `authorization` request header as a `Bearer ...` token in requests requiring 
authentication. The token is used to identify the user by their internal User ID and contains no personally identifying user data.

## Flavors
### `GET /flavor`
Gets a list of the available flavors.

#### Request Params
- **count** (Integer: `25`) Describes the number of records that will be returned in the `items` property of the response.
- **offset** (Integer: `0`) Describes the start position of the records that will be returned in the `items` propery of the response.
- **sortBy** (String: `name`)(unimplemented) Name of the field the results will be sorted by, default `name`

#### Response body
- `items` contains the data. Each item contains:
    - id (Integer)
    - name (String)
    - description (String)
    - ingredients (Array)
    - created (Datetime)
    
```$xslt
{
  "items": [
    {
      "id": 31,
      "name": "Banana Cream Pie",
      "description": "Our signature, intensely flavored banana ice cream includes a fluffy marshmallow swirl and crunchy vanilla wafers. DELICIOUS!",
      "ingredients": [
        {
          "id": 24,
          "name": "Banana"
        }
      ],
      "created": "2019-03-03T05:29:37Z"
    },
    {
      "id": 13,
      "name": "Blueberry Corncake",
      "description": "The perfect combination of taste and texture! Our delicious Blueberry Corncake ice cream includes bits of house made corncake with a sweet and tart, wild blueberry swirl.",
      "ingredients": [
        {
          "id": 18,
          "name": "Blueberry"
        },
        {
          "id": 19,
          "name": "Corncake"
        }
      ],
      "created": "2019-03-03T05:29:36Z"
    },
    ...
  ],
  "meta": {
    "count": 26,
    "sortBy": "created",
    "start": 0,
    "totalRecords": 28
  }
}
```

### `POST /flavor`
Adds a new flavor to the catalog.

`ingredients` is an array of strings, each being an ingredient contributing to this ice cream's unique
flavor profile. Ice cream obviously includes 'cream' and 'sugar' so these can be omitted. Use words
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
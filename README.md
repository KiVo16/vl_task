# VL REST

VL Rest was created to satify requirements of interview task provided by VoiceLab. It's based on Gorilla Mux and Gorm library for database manipulations. Only SQLite is supported.


## Endpoints
| Path        | Method         | Query Params | Requires Body | Mandatory Body keys |   Description |
| ------------- |:----------:| :----: | :-----:| :----: |  ----:|
| /users      | `GET` |`limit (string)` `offset (string)` | `false` | - |  Returns users. Both parameters are optional. |
| /users      | `POST` | - | `true` | `name (string) ` | Creates user with specified name |
| /users/{`user_id`}/{`record_id`}      | `POST` |  - | `false` | - |  Assigns existing record to existing user  |
| /records      | `GET` | `user_name (string)` `type (string)` | `false` | - | Counts records.  |
| /records      | `POST` |  - | `true` | `name (string)` `type (string)` | Creates new record with type and name. |

## Response Format
Response always contains 2 keys: 
- `status` - response status
- `response` - actual response from a given endpoint

Example response for: `GET` `/users?limit=1`
```json
{
    "status": 0,
    "response": [
        {
            "id": 0,
            "name": "Michał"
        }
    ]
}
```

## Flags
| Flag        | Default value           | Description  |
| ------------- |:-------------:| -----:|
| `-db`      | `test.db` | Path to SQLite database file. New database will be created if provided path points to non-existing file. |
| `-load-sample-data`      | `false`      |   Loads sample data |

## Flag usage example

If base.db doesn't exists - program will create new database and fill it with sample data.
```sh
./server -load-sample-data -db=./base.db
```
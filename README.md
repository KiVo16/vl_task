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
| Response        | Code           | 
| ------------- |:-------------:| 
| `ResponseStatusOK`      | `0` | 

## Error Format
Mandatory keys:
- `http_code` - HTTP code
- `error_code` - internal api error code

Optional keys:
- `message` - human readable message
- `refers_to` - points to part where error occured. For example if request's body requires `name` as `string` and got `int` then value of this field will be `name`, because error refers to `name` key.
- `detailed_error` - detailed error derived from goland `error` type. For example error threw by database.

```json
{
    "http_code": 400,
    "error_code": 1,
    "refers_to": "name",
    "message": "Expected string got int"
}
```


| Error        | Code           | 
| ------------- |:-------------:| 
| `ErrValueNotFound`      | `0` | 
| `ErrValueInvalidType`   | `1` |
| `ErrBodyMissing`      | `2` | 
| `ErrBodyRead`      | `3` | 
| `ErrJsonInvalid`      | `4` |
| `ErrInsertData`      | `5` | 
| `ErrGetData`      | `6` | 
| `ErrForeignKey`      | `7` | 

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

##### Progress bar is printed 100% correctly on Windows Terminal and any Linux terminal

```sh
Generating sample data...  38% [====>          ]  [10s:17s]
```
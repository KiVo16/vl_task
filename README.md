# VL REST

VL Rest was created to satify requirements of interview task provided by VoiceLab. It's based on Gorilla Mux and Gorm library for database manipulations. Only SQLite is supported.

## Endpoints
| Path        | Method         | Query Params | Required Body | Mandatory Body keys | Description |
| ------------- |:-------------:| :----: | :-----:| :----: |   ----:|
| /users      | `GET` | `limit (int)` `offset (int)` | `false` | - |Returns users. Both parameters are optional. |
| /users      | `POST` | - | `true` | `name (string)` |Creates user with specified name |
| /users/{user_id}/{record_id}      | `POST` |  | `false` | - |Creates user with specified name |
| /records      | `GET` | `user_name (string)` `type (string)` | `false` | - | Counts records. Counts records of users with specific name if `user_name` is provided and filter it to specific type if `type` params is present |
| /records      | `POST` | - | `true` | `name` `type` | Creates new record with type and name. |

## Flags
| Flag        | Default value           | Description  |
| ------------- |:-------------:| -----:|
| `-db`      | `test.db` | Path to SQLite database file. New database will be created if provided path points to non-existing file. |
| `-load-sample-data`      | `false`      |   Loads sample data |

## Examples

If base.db doesn't exists - program will create new database and fill it with sample data.
```sh
server -load-sample-data -db=./base.db
```
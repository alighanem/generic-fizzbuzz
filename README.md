# Generic - Fizzbuzz

## Description

It is a rest API which can generate strings representations of integer values
and gives the most frequent request.

## Endpoints

### Generate

Returns an array containing strings representation of integer values generated 
from the request parameters.

The parameters are:
 * two integers `int1`, `int2`: value that much be replaced by a string. 
 * one integer `limit`: maximum integer to generate.
 * two strings `str1` (string representation of `int1`) and `str2` (string representation of `int2`). 

> Request

```http
GET /generate?int1=3&int2=5&limit=20&str1=fizz&str2=buzz HTTP/1.1
Content-Type: application/json
``` 

If the generation succeeded: 

> Response

```http
HTTP/1.1 200 OK
Content-Type: application/json

{
    "data": [
        "1",
        "2",
        "fiz",
        "4",
        "buzz",
        "fiz",
        "7",
        "8",
        "fiz",
        "buzz",
        "11",
        "fiz",
        "13",
        "14",
        "fizbuzz",
        "16",
        "17",
        "fiz",
        "19",
        "buzz"
    ]
}
```

If parameters are invalid:

```http
HTTP/1.1 400 BadRequest
Content-Type: application/json

{
    "errors": [
        {
            "code": "internal_error",
            "detail": "cannot read int1: param int1 empty"
        },
        {
            "code": "internal_error",
            "detail": "cannot read str1: param str1 not found"
        }
    ]
}
```

### Statistics

Returns the most frequent request and the number of hits.


> Request

```http
GET /statistics HTTP/1.1
Content-Type: application/json
``` 

If the statistics exists: 

> Response

```http
HTTP/1.1 200 OK
Content-Type: application/json

{
    "data": {
        "converter": {
            "value1": 3,
            "word1": "fiz",
            "value2": 5,
            "word2": "buzz",
            "limit": 20
        },
        "hits": 1
    }
}
```

If there are no statistics:

```http
HTTP/1.1 404 NotFound
Content-Type: application/json

{
    "errors": [
        {
            "code": "not_found",
            "detail": "max number of hits not found"
        }
    ]
}
```

## Configuration

Before launching the API, you need to configure it
by using the environments variables.

Two variables must be configured:
 * `API_LOG_PATH`: file path of the log.
 * `API_PORT`: to set the port.

## Future improvements

For the need of the exercise, the metrics store have been simplified.
It is only in memory.

It is possible to use a simple database to store metrics to persist previous executions.
But it is better to have a real metric tool like `prometheus`.

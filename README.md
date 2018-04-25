# Graphql tester - Tool for testing GraphQL

[![Build Status](https://travis-ci.org/iMega/graphql-tester.svg?branch=master)](https://travis-ci.org/iMega/graphql-tester)

### Motivation
in progress

### Usage

```
docker run -v `pwd`/github_api:/data imega/graphql-tester:0.0.1 -H 'Authorization: Bearer <token>' -u https://api.github.com/graphql /data
```

token - [Creating a personal access token](https://help.github.com/articles/creating-a-personal-access-token-for-the-command-line/)

### Referense

#### The test case

The test case specification in the data part is composed by a series of test blocks.
Each test block usually corresponds to a single test case, which has a title,
an optional description, and a series of data sections. The structure of a test block
is described by the following template.


```
# Title suite test

=== Test 1.

--- query
query {
    viewer {
        id
        login
    }
}

--- query_var
{
    "variables": {
        "input": {
            "login": "blah-blah-blah"
        }
    },
    "operationName": "GetLogin"
}

--- expected_response
{
    "data": {
        "viewer": {
            "login": "iMega"
        }
    }
}

--- assert_not_empty
data.viewer.login

--- assert_empty
data.viewer.id

--- set_vars_from_response
@login = data.viewer.login

```

### Examples

See folder github_api

# Lab 3

The objectives of this lab are as follows:

- Learn about and implement a simple RESTful API.
- Implement user identification and authentication mechanisms.
- Implement confidentiality mechanisms using HTTPS.

## Statement

We will create a prototype of a [database as a service][1] consisting of an HTTP RESTful service capable of storing documents in JSON format. The database will allow different user accounts, as well as basic user management.

For this prototype, assume the server runs on port 5000 on the machine `myserver.local`; therefore, the base endpoint will be `https://myserver.local:5000`. `myserver.local` must resolve to IP `127.0.0.1`.

### Document Storage

Each user will have their own space to store JSON documents. This space will be accessible via `/<user>`, and documents can be created and retrieved via `/<user>/<id>`, where `<id>` is a string that identifies the user's document. For example, `https://myserver.local:5000/luffy/plan-para-ser-el-rey-de-los-piratas` will access the document with identifier `plan-para-ser-el-rey-de-los-piratas` belonging to user `luffy`.

Documents will be stored as files on the server itself. Following the previous example, there would be a file for that document at `<root_path>/luffy/plan-para-ser-el-rey-de-los-piratas.json`, with `<root_path>` being configurable.

### User Authentication

Service users will have a username (`username`) and an associated password (`password`). To register, users will use the `/signup` endpoint.

To use the service, users need a temporary token issued by the server when calling `/login`. This token expires 5 minutes after its creation, so users must call `/login` periodically to obtain a new token. All operations performed while the token is valid will be allowed.

A user may have several active tokens at the same time.

### API

The objects to be implemented, along with their HTTP verbs, are explained below. If an HTTP verb is not specified for an object, that method is not accepted and its call must not be allowed.

### `/version`

- `GET`: returns the program version.

  If successful, returns a program version number in the format `X.Y.Z`. Example:

  ```json
  {"version": "0.1.0"}
  ```

### `/signup`

- `POST`: creates a new user. Requires the following arguments:

  - `username`: the username.
  - `password`: the user's password.

  If the user cannot be created, the response will include a JSON object with the reason for rejection:

  ```json
  {"reason": "Username already used"}
  ```

  ```json
  {"reason": "Too simple password"}
  ```

### `/login`

- `POST`: authenticates an existing user. Requires the following arguments:

  - `username`: the username.
  - `password`: the user's password.

  If successful, returns the user's access token in the format:

  ```json
  {"access_token": "123456789"}
  ```

All operations described below require a valid authentication token to be processed correctly. The token must be sent in the HTTP `Authorization` header of the request. The header value must be `token <user-token>`. Example:

```
Authorization: token <user-token>
```

### `/<string:username>/<string:doc_id>`

- `GET`: retrieves the content of document `doc_id` for user `username`.

  If successful, returns the full content of the document in JSON format.

- `POST`: creates a new document with identifier `doc_id` for user `username`. Requires the following argument:

  - `doc_content`: the content of the document to create, in JSON format.

  For example, if we wanted to upload the following document:

  ```json
  {
    "1": "Encontrar un espadachín para la tripulación",
    "2": "Encontrar una navegante"
  }
  ```

  The request should be:

  ```json
  {
    "doc_content": {
      "1": "Encontrar un espadachín para la tripulación",
      "2": "Encontrar una navegante"
    }
  }
  ```

  **Note that a JSON document may be just an array or even a single value.**

  If successful, returns the number of bytes written to disk in the format:

  ```json
  {"size": <total_bytes>}
  ```

  Remember to store only the inner JSON document (the value of the `doc_content` field), not the outer wrapper object sent in the request. This will affect the total number of bytes written.

- `PUT`: updates the content of document `doc_id` for user `username`. Requires the following argument:

  - `doc_content`: new content of the document in JSON format (see the `POST` description).

  If successful, returns the number of bytes written to disk in the format:

  ```json
  {"size": <total_bytes>}
  ```

- `DELETE`: deletes the document `doc_id` for user `username`.

  If successful, returns an empty response:

  ```json
  {}
  ```

### `/<string:username>/_all_docs`

- `GET`: retrieves the content of all documents for user `username`.

  If successful, a JSON object will be returned with the following structure:

  ```json
  {
    "document_name1": {
      ...content of document_name1...
    },
    "document_name2": {
      ...content of document_name2...
    },
    ...
  }
  ```

## Requirements

1. Implement the service in the Go programming language. During the lab sessions, a recommended framework may be suggested to implement the service.

   - For each input argument, its validity and acceptance must be specified and checked. It is important to ensure that the values provided by the user are **correct and secure**.

   - For each error situation, an appropriate HTTP error code must be returned. For example, if a document does not exist, code `404` is appropriate, while code `403` is not.

   - It is recommended to use a user authentication mechanism discussed in class (such as `shadow` files, use of _salt_ and _digest_ for the password...).

   - In case the service is restarted, tokens do not need to be persistent. It is fine if the user's program has to request the token again.

   - Obviously, documents and user authentication must be persistent.

   - Authentication with the token must be performed in each request that requires it. Token generation must be good enough to avoid collisions between new tokens.

2. `README` and documentation: include a `README.md` (or `README.pdf`) that briefly explains the project, how to compile, install, and run it. It must contain the necessary instructions so that any professional, without knowing how it is implemented, can run it locally without problems.

   - The document must explain how to obtain SSL certificates to deploy the service and how to make a standard GNU/Linux client [2] (for example, `curl`) access the service without certificate problems.

   - The document must also describe aspects that ensure service availability, such as per-user storage quota limitation or limiting user requests per time interval. These measures do not need to be implemented but should be documented for future work.

3. Scripts or tools to generate and install certificates may be included.

Deliverables will be tested with `curl` and the Python `requests` library for automated evaluation, so it is very important to implement the API exactly as specified. The evaluation will also attempt to exploit security flaws inherent to this specification.

## Evaluation Criteria
 
1. Implementation of unit and integration tests.
1. Properly checking and validating input parameters and anticipating possible security flaws that users could exploit.
1. Creating a container image for the service.

[1]: https://en.wikipedia.org/wiki/Cloud_database
[2]: https://wiki.debian.org/Firefox/PrivateCertificateAuthority

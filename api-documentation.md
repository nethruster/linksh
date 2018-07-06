# API methods
## Users
**Get Users**
----
  Returns the list of registered users in the server. <br />
  Requires admin permissions
* **URL**

  `/api/users`

* **Method:**

  `GET` 
  
*  **URL Params**

   **Optional:**
 
   `email=[text]`
   
   * Filter users by its email, it does not need a perfect match

   `offset=[uint]`
   
   * It specifies an offset for the results

   * default: 0

   `limit=[uint]`
   
   * It limits the number of users that will be returned per request.

   * default: 0 (unlimited)

* **Success Response:**

  * **Code:** 200 OK<br />
    **Content:** 
    ```javascript
    {
        [
            {
                "id": "FfXBOYCcTU9riSMg0UWn4qqX6Kt2~bpXSPID",
                "username": "Example user",
                "email": "example@user.com",
                "apikey": "7FrCYXDROFLoJ4EnFf6L5irmKS5L1eb64Q17",
                "isAdmin": false,
                "links": null,
                "sessions": null
                "createdAt": "2018-06-21T19:43:20.553Z"
                "updatedAt": "2018-06-21T19:43:20.553Z"
            }
        ]
    }
    ```
 
* **Error Response:**

  * **Code:** 401 UNAUTHORIZED <br />
    **Content:** `{"error": "UNAUTHORIZED"}`

  * **Code:** 403 FORBIDDEN <br />
    **Content:** `{ error : "FORBIDDEN" }`
  * **Code:** 500 INTERNAL SERVER ERROR <br />
    **Content:** `{ error : "Internal server error" }`

**Get User**
----
  Returns the specified user. <br />
  The requested user should be the same as logged user or the logged user 
* **URL**

  `/api/users/:id`

* **Method:**

  `GET` 
  
*  **URL Params**

   **Optional:**
 
   `includeSessions=[bool]`
   
   * Whether or not to include the sesions data of the user

   * default: false

   `sessionsOffset=[uint]`
   
   * It specifies an offset for the sessions that will be included for each user.
   * This parameter is only used when includeSessions is set to true, otherwise it will be ignored.
   * default: 0

   `sessionsLimit=[uint]`
   
   * It limits the number of sessions that will be included per user.
   * This parameter is only used when includeSessions is set to true, otherwise it will be ignored.
   * default: 0 (unlimited)
   
   `includeLinks=[bool]`
   
   * Whether or not to include the links data of the user

   * default: false

   `linksOffset=[uint]`
   
   * It specifies an offset for the links that will be included for each user.
   * This parameter is only used when includeLinks is set to true, otherwise it will be ignored.
   * default: 0

   `linksLimit=[uint]`
   
   * It limits the number of links that will be included per user.
   * This parameter is only used when includeLinks is set to true, otherwise it will be ignored.
   * default: 0 (unlimited)

* **Success Response:**

  * **Code:** 200 OK<br />
    **Content:** 
    ```javascript
        {
            "id": "FfXBOYCcTU9riSMg0UWn4qqX6Kt2~bpXSPID",
            "username": "Example user",
            "email": "example@user.com",
            "apikey": "7FrCYXDROFLoJ4EnFf6L5irmKS5L1eb64Q17",
            "isAdmin": true,
            "links": [
                {
                    "id": "RRqltG",
                    "content": "https://example.org",
                    "Hits": 0,
                    "UserId": "FfXBOYCcTU9riSMg0UWn4qqX6Kt2~bpXSPID",
                    "CreatedAt": "2018-06-22T01:25:34+01:00",
                    "UpdatedAt": "2018-06-22T01:25:34+01:00"
                }
            ],
            "sessions": [
                {
                    "id": "3bxBaFtRBz0fKHMThVJd48spaeHsWTsEUCRP",
                    "UserId": "FfXBOYCcTU9riSMg0UWn4qqX6Kt2~bpXSPID",
                    "createdAt": "2018-06-17T19:49:41+01:00"
                }
            ],
            "createdAt": "2018-06-17T19:49:41+01:00",
            "updatedAt": "2018-06-17T19:49:41+01:00"
        }
    ```
* **Error Response:**

  * **Code:** 401 UNAUTHORIZED <br />
    **Content:** `{"error": "UNAUTHORIZED"}`
  * **Code:** 403 FORBIDDEN <br />
    **Content:** `{ error : "FORBIDDEN" }`
  * **Code:** 404 NOT FOUND <br />
    **Content:** `{ error : "User not found" }`
  * **Code:** 500 INTERNAL SERVER ERROR <br />
    **Content:** `{ error : "Internal server error" }`


**Create User**
----
  Creates a new user <br />
  The current user should be an admin.
* **URL**

  `/api/users`

* **Method:**

  `POST` 
  
* **Data Params**

    Data params should be provided in JSON format.

    **Required:**

    `username=[string(max=255)]`

    `email=[string(max=255)]`

    `password=[string]`

    **Optional:**

  `isAdmin=[bool]` //In order to modify this parameter, the current user should be an admin
* **Success Response:**

  * **Code:** 201 CREATED<br />
    **Content:** 
    ```javascript
        {
            "id": "HDwvc5zV37_u4VToo9ahJkFwqBv3ugvoDY~L",
            "username": "Example user",
            "email": "example@user.com",
            "apikey": "bRTiLUrxPd5jS5DgVwCRq91XQ7X1UY8wYf5x",
            "isAdmin": false,
            "links": null,
            "sessions": null,
            "createdAt": "2018-06-22T01:40:49.1475591+01:00",
            "updatedAt": "2018-06-22T02:00:00.1475591+01:00"
        }
    ```
* **Error Response:**
    * **Code:** 400 BAD REQUEST <br />
        **Content:** `{"error": [Possible errors]}`

        **Possible errors:** 
        * "Missing username"
        * "Username must not be longer than 255 characters"
        * "Missing email"
        * "Email is to long"
        * "Invalid email format"
        * "Missing password"
        
    * **Code:** 401 UNAUTHORIZED <br />
    **Content:** `{"error": "UNAUTHORIZED"}`
    * **Code:** 403 FORBIDDEN <br />
    **Content:** `{ error : "FORBIDDEN" }`
    * **Code:** 500 INTERNAL SERVER ERROR <br />
    **Content:** `{ error : "Internal server error" }`

**Edit User**
----
  Edits a user <br />
  The current user should be the edited user or an admin.
* **URL**

  `/api/users/:id`

* **Method:**

  `PUT` 
  
* **Data Params**

    Data params should be provided in JSON format.

    **Optional:**

    `username=[string(max=255)]`

    `email=[string(max=255)]`

    `password=[string]`

    `originalPassword=[string]` //originalPassword should be provided in order to change the user's password if modifing current user.

    `isAdmin=[bool]` //In order to modify this parameter, the current user should be an admin

    `apikey=[bool]` //This would make the apikey to regenerate.
* **Success Response:**

  * **Code:** 200 OK<br />
    **Content:** 
    ```javascript
        {
            "id": "HDwvc5zV37_u4VToo9ahJkFwqBv3ugvoDY~L",
            "username": "Example user",
            "email": "example@user.com",
            "apikey": "bRTiLUrxPd5jS5DgVwCRq91XQ7X1UY8wYf5x",
            "isAdmin": false,
            "links": null,
            "sessions": null,
            "createdAt": "2018-06-22T01:40:49.1475591+01:00",
            "updatedAt": "2018-06-22T01:40:49.1475591+01:00"
        }
    ```
* **Error Response:**
    * **Code:** 400 BAD REQUEST <br />
      **Content:** `{"error": "Missing or invalid originalPassword"}`
    * **Code:** 400 BAD REQUEST <br />
      **Content:** `{"error": [Possible errors]}`

        **Possible errors:** 
        * "Missing username"
        * "Username must not be longer than 255 characters"
        * "Missing email"
        * "Email is to long"
        * "Invalid email format"
        * "Missing password"
        
    * **Code:** 401 UNAUTHORIZED <br />
    **Content:** `{"error": "UNAUTHORIZED"}`
    * **Code:** 403 FORBIDDEN <br />
    **Content:** `{ error : "FORBIDDEN" }`
    * **Code:** 404 NOT FOUND <br />
    **Content:** `{ error : "User not found" }`
    * **Code:** 500 INTERNAL SERVER ERROR <br />
    **Content:** `{ error : "Internal server error" }`

**Delete User**
----
  Deletes the specified user. <br />
  The requested user should be the same as logged user or the logged user 
* **URL**

  `/api/users/:id`

* **Method:**

  `DELETE` 

* **Success Response:**

  * **Code:** 204 NO CONTENT<br />

* **Error Response:**

  * **Code:** 401 UNAUTHORIZED <br />
    **Content:** `{"error": "UNAUTHORIZED"}`
  * **Code:** 403 FORBIDDEN <br />
    **Content:** `{ error : "FORBIDDEN" }`
  * **Code:** 404 NOT FOUND <br />
    **Content:** `{ error : "User not found" }`
  * **Code:** 500 INTERNAL SERVER ERROR <br />
    **Content:** `{ error : "Internal server error" }`
## Links
**Get Links**
----
  Returns the list of links, it can be limited to one user or to the entire server <br />
  Requires admin permissions in order to view links which doesn't belong to the current user
* **URL**

  `/api/links`

* **Method:**

  `GET` 
  
*  **URL Params**

   **Optional:**
 
   `ownerId=[text]`
   
   * Limits links listed to one specific user

   `offset=[uint]`
   
   * It specifies an offset for the results

   * default: 0

   `limit=[uint]`
   
   * It limits the number of links that will be returned per request.

   * default: 0 (unlimited)

* **Success Response:**

  * **Code:** 200 OK<br />
    **Content:** 
    ```javascript
    [
        {
            "id": "RRqltG",
            "content": "google.es",
            "Hits": 0,
            "UserId": "HDwvc5zV37_u4VToo9ahJkFwqBv3ugvoDY~L",
            "CreatedAt": "2018-06-22T01:25:34+01:00",
            "UpdatedAt": "2018-06-22T01:25:34+01:00"
        }
    ]
    ```
 
* **Error Response:**

  * **Code:** 401 UNAUTHORIZED <br />
    **Content:** `{"error": "UNAUTHORIZED"}`

  * **Code:** 403 FORBIDDEN <br />
    **Content:** `{ error : "FORBIDDEN" }`
  * **Code:** 500 INTERNAL SERVER ERROR <br />
    **Content:** `{ error : "Internal server error" }`

**Get Link**
----
  Returns the specified link. <br />
  The requested link should be owned by the current user or the current user should be an admin.
* **URL**

  `/api/users/:id`

* **Method:**

  `GET` 

* **Success Response:**

  * **Code:** 200 OK<br />
    **Content:** 
    ```javascript
        {
            "id": "RRqltG",
            "content": "https://example.org",
            "Hits": 0,
            "UserId": "HDwvc5zV37_u4VToo9ahJkFwqBv3ugvoDY~L",
            "CreatedAt": "2018-06-22T01:25:34+01:00",
            "UpdatedAt": "2018-06-22T01:25:34+01:00"
        }
    ```
* **Error Response:**

  * **Code:** 401 UNAUTHORIZED <br />
    **Content:** `{"error": "UNAUTHORIZED"}`
  * **Code:** 403 FORBIDDEN <br />
    **Content:** `{ error : "FORBIDDEN" }`
  * **Code:** 404 NOT FOUND <br />
    **Content:** `{ error : "Link not found" }`
  * **Code:** 500 INTERNAL SERVER ERROR <br />
    **Content:** `{ error : "Internal server error" }`

**Create Link**
----
  Creates a new link

* **URL**

  `/api/links`

* **Method:**

  `POST` 
  
* **Data Params**

    Data params should be provided in JSON format.

    **Required:**

    `Content=[string(max=2000)]`

    **Optional:**

    `CustomId=[string(max=255)]` 

* **Success Response:**

  * **Code:** 201 CREATED<br />
    **Content:** 
    ```javascript
        {
            "id": "WcnB_N",
            "content": "https://example.org",
            "Hits": 0,
            "UserId": "FfXBOYCcTU9riSMg0UWn4qqX6Kt2~bpXSPID",
            "CreatedAt": "2018-06-22T17:46:33.2726225+01:00",
            "UpdatedAt": "2018-06-22T17:46:33.2726225+01:00"
        }
    ```
* **Error Response:**
    * **Code:** 400 BAD REQUEST <br />
      **Content:** `{"error": "Link already exists"}`
    * **Code:** 401 UNAUTHORIZED <br />
    **Content:** `{"error": "UNAUTHORIZED"}`
    * **Code:** 403 FORBIDDEN <br />
    **Content:** `{ error : "FORBIDDEN" }`
    * **Code:** 500 INTERNAL SERVER ERROR <br />
    **Content:** `{ error : "Internal server error" }`

**Edit Link**
----
  Edits a Link <br />
  The current user should be the link owner or an admin.
* **URL**

  `/api/links/:id`

* **Method:**

  `PUT` 
  
* **Data Params**

    Data params should be provided in JSON format.

    **Required:**

    `Content=[string(max=2000)]`

    **Optional:**

    `CustomId=[string(max=255)]` 

* **Success Response:**

  * **Code:** 200 OK<br />
    **Content:** 
    ```javascript
        {
            "id": "WcnB_N",
            "content": "https://example.com",
            "Hits": 0,
            "UserId": "FfXBOYCcTU9riSMg0UWn4qqX6Kt2~bpXSPID",
            "CreatedAt": "2018-06-22T17:46:33+01:00",
            "UpdatedAt": "2018-06-22T17:53:21.4461048+01:00"
        }
    ```
* **Error Response:**
    * **Code:** 400 BAD REQUEST <br />
      **Content:** `{"error": "Link already exists"}`
    * **Code:** 401 UNAUTHORIZED <br />
    **Content:** `{"error": "UNAUTHORIZED"}`
    * **Code:** 403 FORBIDDEN <br />
    **Content:** `{ error : "FORBIDDEN" }`
    * **Code:** 404 NOT FOUND <br />
    **Content:** `{ error : "Link not found" }`
    * **Code:** 500 INTERNAL SERVER ERROR <br />
    **Content:** `{ error : "Internal server error" }`

**Delete Link**
----
  Deletes the specified link. <br />
  The requested link should be owned by the current user or the current user is an admin
* **URL**

  `/api/users/:id`

* **Method:**

  `DELETE` 

* **Success Response:**

  * **Code:** 204 NO CONTENT<br />

* **Error Response:**

  * **Code:** 401 UNAUTHORIZED <br />
    **Content:** `{"error": "UNAUTHORIZED"}`
  * **Code:** 403 FORBIDDEN <br />
    **Content:** `{ error : "FORBIDDEN" }`
  * **Code:** 404 NOT FOUND <br />
    **Content:** `{ error : "Link not found" }`
  * **Code:** 500 INTERNAL SERVER ERROR <br />
    **Content:** `{ error : "Internal server error" }`

## Sessions
**Get Sessions**
----
  Returns the list of sessions, it can be limited to one user or to the entire server <br />
  Requires admin permissions in order to view sessions which doesn't belong to the current user
* **URL**

  `/api/sessions`

* **Method:**

  `GET` 
  
*  **URL Params**

   **Optional:**
 
   `ownerId=[text]`
   
   * Limits sessions listed to one specific user

   `offset=[uint]`
   
   * It specifies an offset for the results

   * default: 0

   `limit=[uint]`
   
   * It limits the number of sessions that will be returned per request.

   * default: 0 (unlimited)

* **Success Response:**

  * **Code:** 200 OK<br />
    **Content:** 
    ```javascript
    [
        {
            "id": "3bxBaFtRBz0fKHMThVJd48spaeHsWTsEUCRP",
            "UserId": "FfXBOYCcTU9riSMg0UWn4qqX6Kt2~bpXSPID",
            "createdAt": "2018-06-17T19:49:41+01:00"
        }
    ]
    ```
 
* **Error Response:**

  * **Code:** 401 UNAUTHORIZED <br />
    **Content:** `{"error": "UNAUTHORIZED"}`

  * **Code:** 403 FORBIDDEN <br />
    **Content:** `{ error : "FORBIDDEN" }`
  * **Code:** 500 INTERNAL SERVER ERROR <br />
    **Content:** `{ error : "Internal server error" }`

**Get Session**
----
  Returns the specified session. <br />
  The requested session should be owned by the current user or the current user should be an admin.
* **URL**

  `/api/sessions/:id`

* **Method:**

  `GET` 

* **Success Response:**

  * **Code:** 200 OK<br />
    **Content:** 
    ```javascript
        {
            "id": "3bxBaFtRBz0fKHMThVJd48spaeHsWTsEUCRP",
            "UserId": "FfXBOYCcTU9riSMg0UWn4qqX6Kt2~bpXSPID",
            "createdAt": "2018-06-17T19:49:41+01:00"
        }
    ```
* **Error Response:**

  * **Code:** 401 UNAUTHORIZED <br />
    **Content:** `{"error": "UNAUTHORIZED"}`
  * **Code:** 403 FORBIDDEN <br />
    **Content:** `{ error : "FORBIDDEN" }`
  * **Code:** 404 NOT FOUND <br />
    **Content:** `{ error : "Link not found" }`
  * **Code:** 500 INTERNAL SERVER ERROR <br />
    **Content:** `{ error : "Internal server error" }`

**Create Session**
----
  Creates a new session for the current user

* **URL**

  `/api/session`

* **Method:**

  `POST` 

* **Success Response:**

  * **Code:** 201 CREATED<br />
    **Content:** 
    ```javascript
        {
            "id": "Yx34RbdddkyYQNbiv696J1ismF6r7nazn_fr",
            "UserId": "FfXBOYCcTU9riSMg0UWn4qqX6Kt2~bpXSPID",
            "createdAt": "2018-06-22T18:24:34.3053876+01:00"
        }
    ```
* **Error Response:**
    * **Code:** 400 BAD REQUEST <br />
      **Content:** `{"error": "Link already exists"}`
    * **Code:** 401 UNAUTHORIZED <br />
    **Content:** `{"error": "UNAUTHORIZED"}`
    * **Code:** 403 FORBIDDEN <br />
    **Content:** `{ error : "FORBIDDEN" }`
    * **Code:** 500 INTERNAL SERVER ERROR <br />
    **Content:** `{ error : "Internal server error" }`


**Delete Session**
----
  Deletes the specified session. <br />
  The requested session should be owned by the current user or the current user is an admin
* **URL**

  `/api/sessions/:id`

* **Method:**

  `DELETE` 

* **Success Response:**

  * **Code:** 204 NO CONTENT<br />

* **Error Response:**

  * **Code:** 401 UNAUTHORIZED <br />
    **Content:** `{"error": "UNAUTHORIZED"}`
  * **Code:** 403 FORBIDDEN <br />
    **Content:** `{ error : "FORBIDDEN" }`
  * **Code:** 404 NOT FOUND <br />
    **Content:** `{ error : "Link not found" }`
  * **Code:** 500 INTERNAL SERVER ERROR <br />
    **Content:** `{ error : "Internal server error" }`

## Auth

**Login**
----
 Log in

* **URL**

  `/api/auth/login`

* **Method:**

  `POST` 
  
* **Data Params**

    Data params should be provided in JSON format.

    **Required:**

    `email=[string(max=255)]`

    `password=[string]`

    **Optional:**

    `noExpire=[bool]` 

    `useCookie=[bool]`

* **Success Response:**

  * **Code:** 200 OK<br />
    **Content:** 
    ```javascript
        {
            {
                "sessionId": "WAlnF8zqUmORo6h44nfbRXcbLMmJLwzt5Gyu",
                "userId": "FfXBOYCcTU9riSMg0UWn4qqX6Kt2~bpXSPID",
                "expiresAt": "2018-06-24 18:03:47.0576086 +0100 BST"
            }
        }
    ```
* **Error Response:**
    * **Code:** 400 BAD REQUEST <br />
      **Content:** `{"error": "Missing email or password"}`
    * **Code:** 400 BAD REQUEST <br />
      **Content:** `{"error": "The email or the password are invalid"}`
    * **Code:** 500 INTERNAL SERVER ERROR <br />
    **Content:** `{ error : "Internal server error" }`

**Register**
----
 Register

* **URL**

  `/api/auth/register`

* **Method:**

  `POST` 
  
* **Data Params**

    Data params should be provided in JSON format.

    **Required:**

    `username=[string(max=255)]`

    `email=[string(max=255)]`

    `password=[string]`

    **Optional:**

    `noExpire=[bool]` 

    `useCookie=[bool]`

* **Success Response:**

  * **Code:** 200 OK<br />
    **Content:** 
    ```javascript
        {
            {
                "sessionId": "WAlnF8zqUmORo6h44nfbRXcbLMmJLwzt5Gyu",
                "userId": "FfXBOYCcTU9riSMg0UWn4qqX6Kt2~bpXSPID",
                "expiresAt": "2018-06-24 18:03:47.0576086 +0100 BST"
            }
        }
    ```
* **Error Response:**
    * **Code:** 400 BAD REQUEST <br />
        **Content:** `{"error": [Possible errors]}`

        **Possible errors:** 
        * "Missing username"
        * "Username must not be longer than 255 characters"
        * "Missing email"
        * "Email is to long"
        * "Invalid email format"
        * "Missing password"

    * **Code:** 500 INTERNAL SERVER ERROR <br />
    **Content:** `{ error : "Internal server error" }`

**Log out**
----
  Log out from the current session. <br />
  Using this method authenticated via api-key will cause an 400 Bad Request error.
* **URL**

  `/api/auth/logout`

* **Method:**

  `DELETE` 

* **Success Response:**

  * **Code:** 204 NO CONTENT<br />

* **Error Response:**

  * **Code:** 400 BAD REQUEST <br />
    **Content:** `{"error": "Bad request"}`
  * **Code:** 401 UNAUTHORIZED <br />
    **Content:** `{"error": "UNAUTHORIZED"}`
  * **Code:** 403 FORBIDDEN <br />
    **Content:** `{ error : "FORBIDDEN" }`
  * **Code:** 500 INTERNAL SERVER ERROR <br />
    **Content:** `{ error : "Internal server error" }`

# Authentication methods

There are 3 different ways to authenticate yourself against the API, using a session via cookie, using a session via header or using the api-key. For the first two you will need a valid session for you user, you can use the methods Login, Register and Create Session in order to get one.

## Session via Cookie
This option is the recommended way for webapps.
In order to use this method, you need to get your session via the Login or Register methods. You need to ensure that you set the `useCookie` parameter is set to true, it will create cookie called *linksh-auth* with all the necessary data in the browser, this cookie will have the flag HTTP-Only for security reasons, so it cannot  be seen by the front-end. Once you have the cookie you only need to ensure that the browser sends it with every request.
## Session via Header
In order to use this method, you need the *userId* and a valid *sessionId* for your user (you can get it from the Login, Register and Create Session API methods). Once you have it, you must add a header named `X-Auth` with a value consisting of the *sessionId* followed by the *userId* without spaces and in that order.

Example code:
```javascript
{
    var userId = "AAA"
    var sessionId = "XXX"

    var myHeaders = new Headers({
        "X-Auth": sessionId + userId
    });

    fetch('/api/users', {
        method: 'GET',
        headers: myHeaders,
        mode: 'cors',
        cache: 'default'
    })
}
```
## Api-key
In order to use this method, you need the user's *email* and *api-key*. <br />
You must add a header named `X-Auth-Key` with the *api-key* as value, and a header named `X-Auth-Email` with the email. Notice that both headers are required.

Example code:
```javascript
{
    var userId = "AAA"
    var userEmail = "example@mail.org"

    var myHeaders = new Headers({
        "X-Auth-Key": userId,
        "X-Auth-Email": userEmail
    });

    fetch('/api/users', {
        method: 'GET',
        headers: myHeaders,
        mode: 'cors',
        cache: 'default'
    })
}
```

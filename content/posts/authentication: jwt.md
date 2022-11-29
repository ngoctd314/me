# Authentication

## Json Web Token

In its compact form, JSON Web Tokens consist of three parts separated by dots (.) which are:

- Header
- Payload
- Signature

**Header**
The header typically consists of two parts: the type of the token, which is JWT, and the signing algorithm being used, such as HMAC SHA256 or RSA

```json
{
    "alg": "HS256",
    "typ": "JWT"
}
```
Then, this JSON is Base64Url encoded to form the first part of the JWT.

**Payload**

The second part of the token is the payload, which contains the claims. Claims are statements about an entity (typically, the user) and additional data. There are three types of claims: registered, public, and private claims.

- Registered claims: There are a set of predefined claims which are not mandatory but recommended, to provide a set of useful
iss(issuer), exp(expiration time), sub(subject), aud(audience)
- Public claims:
- Private claims: These are the custom claims created to share information between parties that agree on using them and are neither registered or public claims.

```json
{
    "sub": "123456789",
    "name": "John",
    "admin": true
}
```

The payload is then Base64Url encoded to form the second part of the JSON Web Token.

**Signature**

To create the signature part you have to take the encoded header, the encoded payload, a secret, the algorithm specified in the header, and sign that.
The signature is used to verify the message wasn't changed along the way, and, in the case of tokens signed with a private key, it can also verify that the sender of the JWT is who it says it is.

**How do JSON Web Tokens work?**

In authentication, when the user successfully logs in using their credentials, a JSON Web Token will be returned. Since tokens are credentials,, greate care must be taken to prevent security issues. In general, you should not keep tokens longer than required.

You also should not store sensitive session data in browser storage due to lack of security.

Whenever the user wants to access a protected route or resource, the user agent should send the JWT, typically in the Authorization header using the Bearer schema.

The content os the header should look like the following:
```txt
Authorization: Bearer <token>
```
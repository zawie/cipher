# End-to-End Encryption Message Application

## About End-to-End Encryption
### What is End-to-End Encryption?
Suppose two users Alice and Bob want to communicate, but they do not trust any other party to process or store their plain-text messages. 
More specifically, they do not want any messaging service to have any computationally feasible way to read their messages. End-to-end encryption solves for this problem by utilizing asymmetric encryption.

Here is an illustrative example of the procedure. Suppose Alice wants to send plain-text message `M` to Bob via a service ran by Sally. At a high level, the procedure works as follows:
- Bob generates a private and public key pair.
- Bob shares his public key with Sally.
- Alice retrieves Bob's public key from Sally.
- Alice encrypts her message `M` with Bob's public key, producing cipher text `C`
- Alice sends cipher text `C` to Sally.
- Sally passes `C` to Bob
- Bob uses his private key to unencrypt the cipher text `C` to retrieve Alice's plain text message `M`

Throughout this process Sally never recieved enough information to be able to retrieve `M`, as sally only ever recieved Bob's public key and the cipher text. Since Sally never recieved the private key, she is unable to decrypt `C`. 
Therefore, Alice and Bob are able to succesfully communicate through Sally without exposing any information about their message contents.

### Why should I care about End-to-End Encryption?

There are various reasons parties should prefer End-to-End encryption. It is possible that will Sally maliciously attempts to read private communications between Alice and Bob.
Morever, even if you trust Sally, it is possible an unrelated party compromises Sally and reads private communicate between Alice and Bob. 

# About Our Implementation

TODO: write

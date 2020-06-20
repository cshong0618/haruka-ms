# Haruka
No special meaning with the name. 

## Introduction
This is a monorepo of a simple service done with the microservice architecture style.

## Components
Directory|Usage
---|---
`user`|Handles user data
`gateway`|REST entry point

### `gateway`
The Gateway service is what the outside world should be hitting. It does
not do any domain logic but to only proxy message to the inner services 
through `nats.io` as message broker.

```

[Outside world] --> [Gateway] ---> [nats.io]
                                   /   |    \
                                  /    |     \
                                 /     |      \
                                /      |       \
                             api1     api2     ...
```

### `user`
The User service houses the CRUD of user data. It has both `nats` and `REST` 
executables available to be built. The underlying database is MongoDB.


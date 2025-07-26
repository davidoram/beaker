# API Layer

As discussed in the architecture document our API exposes a JSON based API. It can be accessed over two protocols, NATS or HTTP.

- NATS is our preferred protocol. It's high performance, and ideal for system integration because it offers great security, low latency, network failure resiliance, and good scaleability. On the downside, you need to use one of the client libraries, it's less well known than HTTP, and it might be blocked on some corporate firewalls.
- HTTP is the defacto standard protocol for API delivery. It's compatible with standard web infrastructure and has no special client library requirements. However its performance isn't nearly as good as the NATS protocol.  

Our system will be deployed on [Synadia Cloud](https://www.synadia.com/cloud), which offers support for both NATS and HTTP protocols simultaneously.

The access looks something like this:

```mermaid
flowchart TD
    subgraph Client Side
        A[NATS Client<br/>Go, JS, etc.]
        B[HTTP Client<br/>Curl, Browser, etc.]
    end

    subgraph Synadia Cloud
        A --> C[NATS Server]
        B --> D[HTTP Gateway]

        D --> C

        style C fill:#eef,stroke:#88c
        style D fill:#efe,stroke:#8c8
    end

    subgraph Codespace
        C --> E[Service API<br/>NATS Microservice]
    end

    subgraph Authentication
        F1[NATS JWT<br/>+ TLS]
        F2[Bearer Token<br/>+ TLS]

        A -->|Connect with| F1
        B -->|Send token with| F2
    end

    click C "https://docs.nats.io/using-nats/developer/services" _blank
    click D "https://docs.synadia.com/nats-api-gateway/" _blank

```

- A NATS **user** is created for each client/caller.
- NATS Clients authenticate using JWTs issued by the NATS account system, and TLS secures the connection.
- HTTP Clients authenticate by passing a Bearer token, issued by Synadid Cloud, which the HTTP Gateway validates.
- Both paths end up delivering the message to the same NATS microservice, with identity and auth context injected.
- Both connections are **Authorized** the same way.

In turn our microservice itself must make an authenticated connection to NATS, so that API requests can be routed to the service, and responses returned.  Thats configured like this diagram:

```mermaid
flowchart TD

    subgraph Synadia Cloud
        C[NATS Server]
    end

    subgraph Codespace
        E[Service API<br/>NATS Microservice]
        D[credentials file]
        D --> |uses| E
        E --> |authenticates| C
    end

    subgraph Authentication
        F1[NATS JWT<br/>+ TLS]

        C -->|Verifies| F1
        F1 -->|Issues| D
    end

    click C "https://docs.nats.io/using-nats/developer/services" _blank
```

- A NATS **user** is created for our Microservice on Synadia Cloud
- When we configure our Microservice,  we copy the **user** credentials file down to our codespace.
- When the Microservice connects to Synadia Cloud (NATS), it passes the credentials which NATS used to **Authenticate** and **Authorize**, sometimes called AuthN/Z.
    - **Authentication** confirms who we are
    - **Authorization** tells the system what we can do. In NATS that means what **subjects** we are allowed to subscribe & publish to.

So for our purposes we need to signup to [Synadia Cloud](https://cloud.synadia.com) and create some **users**.

# Signup to Synadia Cloud

First signup to [Synadia Cloud](https://cloud.synadia.com).  They have a free offering which is all we will need for now.

Once you accept the Ts&Cs you are ready to go.  Note there is no need to install the `nats` cli tool becuase we already have that installed inside our codespace.

Our goal is to create two separate NATS **users** and download their credentials.  Credentials are a kind of "secret", and as such, we never want to share them with anyone else, so we can't add them to our git projects. Luckily GitHub codespaces has a feature called Github secrets, that are stored securely against your own personal profile. We have a way of "injecting" Github secrets into our projects.

The process is as follows:

1. Create our secret. In our case its a 'credentials' file created via Synadia Cloud.
2. Encode our secret in [base64](https://en.wikipedia.org/wiki/Base64)
3. Add the base64 encoded secret to [Github codespace user secrets](https://github.com/settings/codespaces)
4. When our codespace starts up, transform the Github codespace user secrets back into a 'credentials' file.

We will repreat this process three times for each of the following NATS user credentials files:

- `CLI` the user that has full access to everything, useful for troubleshooting. Github codespace user secret name `NATS_CREDS_CLI`
- `APP` for the application. Github codespace user secret name `NATS_CREDS_APP`
- `CALLER` for a caller to use. Github codespace user secret name `NATS_CREDS_CALLER`

Use the Synadia website to create a new user. Then click 'Get Connected' and download the credentials file.  Next upload that to your codespace by dragging and dropping the file (or copy and pasing the content). Encode the file using the following command `base64 <path to creds file>` eg: `base64 /workspaces/beaker/NGS-Default-CLI.creds`.  Copy exactly the output, and paste it in as a new [Github codespace user secrets](https://github.com/settings/codespaces) withe the correct name as listed above.

> [!NOTE]
> This project demonstrates a session-based microservices architecture using JavaScript, Redis, Express, and NGINX. Sessions are stored in Redis via `express-session` and `connect-redis`, allowing stateless Node.js services like `profile-service` to persist user sessions. An NGINX API Gateway proxies requests to individual services (e.g., `/auth/` and `/profile/`), forwarding client headers and maintaining a consistent interface. This setup enables centralized session management, scalability, and a clean separation between services behind a single gateway.

```mermaid
%%{init: {"themeVariables": {"primaryColor": "#4a90e2", "edgeLabelBackground":"#ffffff"}}}%%
graph TD
    %% Title
    classDef title fill:#222,stroke:#fff,color:#fff,font-size:20px,font-weight:bold;
    titleNode["Session Storage Architecture\n(JavaScript, Redis, NGINX, API Gateway)"] 
    class titleNode title
    titleNode:::title

    %% Nodes
    Client[Client Browser / Frontend App]
    NGINX["API Gateway (NGINX)"]
    Auth["Auth Service (Express)"]
    Profile["Profile Service (Express)"]
    Redis["(Redis Session Store)"]

    %% Edges
    Client -->|HTTP /auth/*| NGINX
    Client -->|HTTP /profile/*| NGINX

    NGINX -->|Proxy /auth/*| Auth
    NGINX -->|Proxy /profile/*| Profile

    Auth -->|Session read/write| Redis
    Profile -->|Session read/write| Redis

    %% Clickable links (optional - add URLs or remove if unwanted)
    click Client "https://example.com/frontend" "Go to Frontend"
    click NGINX "https://nginx.org/en/docs/" "NGINX Docs"
    click Auth "https://expressjs.com/" "Express Docs"
    click Profile "https://expressjs.com/" "Express Docs"
    click Redis "https://redis.io/docs/" "Redis Docs"
```
## Description

Rating Service is service that handle user to give rating.

Rating direction consist of:

- Customer -> Technician
- Customer -> Service Result
- Customer -> App
- Technician -> Customer
- Technician -> App

RabbitMQ Diagram:

```mermaid
flowchart LR
    %% Styling mirip dengan gambar referensi
    classDef producer fill:#6ba4ff,stroke:#333,stroke-width:2px,color:#000,rx:10px,ry:10px;
    classDef exchange fill:#95d166,stroke:#333,stroke-width:2px,color:#000,rx:10px,ry:10px;
    classDef queue fill:#ffdd66,stroke:#333,stroke-width:2px,color:#000,rx:10px,ry:10px;
    classDef consumer fill:#6ba4ff,stroke:#333,stroke-width:2px,color:#000,rx:10px,ry:10px;

    %% Nodes
    P["Producer\n(NestJS Monolith)"]:::producer

    subgraph RabbitMQ ["RabbitMQ"]
        direction LR
        E["Exchange\n(sidingin.events)"]:::exchange
        Q1["Queue\n(rating_service.order_completed.queue)"]:::queue
        E -- "order.completed" --> Q1
    end
    style RabbitMQ fill:#a6c0db,stroke:#333,stroke-width:2px,color:#000,rx:15px,ry:15px;

    C1["Consumer 1\n(Golang rating-service)"]:::consumer

    %% Connections
    P -- "Publish JSON" --> E
    Q1 -- "Consume JSON" --> C1

```

- Exchange: sidingin.events
- Queue: rating_service.order_completed.queue
- Binding Routing Key: order.completed

# Room

## Sequence Diagram
### Assigned Agent with Webhook CAA
Here is a sequence process for assigned flow agent when get webhook from CAA:
```mermaid
sequenceDiagram
    participant O as Omnichannel
    participant S as CAA Service
    participant D as CAA Database
    O->>S: Webhook CAA
    S->>S: Check match APP ID
    S->>D: Added Queue
    D-->>S: Succesfully added queue
    S->>D: Get room detail
    D-->>S: Succesfully get room detail
    S->>S: Check queue isResolved
    S->>S: Check agent id is empty
    S->>O: Get available agent in room
    O-->>S: Succesfully get available agent in room
    S->>S: Check max handle customer
    S->>O: Assign agent in room
    O-->>S: Succesfully assign agent in room
    S->>D: Updated queue
    D-->>S: Succesfully updated queue
    S-->>O: Succesfully send webhook
```
### Assigned Agent with Webhook Mark As Resolved
Here is a sequence process for assigned flow agent when get webhook resolved:
```mermaid
sequenceDiagram
    participant O as Omnichannel
    participant S as CAA Service
    participant D as CAA Database
    O->>S: Webhook Resolved
    S->>D: Update queue
    D-->>S: Succesfully update queue
    S->>D: Get first room unasigned queue
    D-->>S: Succesfully get first room unasigned queue
    S->>O: Get available agent in room
    O-->>S: Succesfully get available agent in room
    S->>S: Check max handle customer
    S->>O: Assign agent in room
    O-->>S: Succesfully assign agent in room
    S->>D: Updated queue
    D-->>S: Succesfully updated queue
    S-->>O: Succesfully send webhook
```

### Assigned Agent with Cronjob
Here is a sequence process for assigned flow with cronjob in delay 60 seconds:
```mermaid
sequenceDiagram
    participant O as Omnichannel
    participant S as CAA Service
    participant D as CAA Database
    note over S: This cronjob running every 60 seconds
    S->>S: Periodic run
    S->>D: Get first room unasigned queue
    D-->>S: Succesfully get first room unasigned queue
    S->>O: Get available agent in room
    O-->>S: Succesfully get available agent in room
    S->>S: Check max handle customer
    S->>O: Assign agent in room
    O-->>S: Succesfully assign agent in room
    S->>D: Updated queue
    D-->>S: Succesfully updated queue
```

## FlowChart
### Assigned Agent with Webhook CAA
Here is a flowchart for assigned flow agent when get webhook from CAA:
```mermaid
flowchart TD
    A[Start] --> B{App id valid?}
    B -->|No| Z[Finish]
    B -->|Yes| C[Added Queue]
    C --> D[Get Room Detail]
    D --> E{Room is resolved?}
    E -->|No| F{AgentID is empty?}
    E -->|Yes| Z[Finish]
    F -->|No| Z[Finish]
    F -->|Yes| G[Get available agent]
    G --> H{Have criteria agent?}
    H -->|No| Z[Finish]
    H -->|Yes| I[Assign Agent]
    I -->J[Update Queue]
    J -->Z[Finish]
```

### Assigned Agent Webhook Mark As Resolved
Here is a flowchart for assigned flow agent when get webhook resolved:

```mermaid
flowchart TD
    A[Start] --> B[Update Room Resolved]
    B --> C[Get First Queue Today]
    C --> D[Get available agent]
    D --> E{Have criteria agent?}
    E -->|No| Z[Finish]
    E -->|Yes| F[Assign Agent]
    F -->H[Update Queue]
    H -->Z[Finish]
```

## ERD
Here is a ERD for Queue
```mermaid
erDiagram
    agent_allocation_queues {
        uint id PK
        string room_id
        string agent_id
        bool is_resolved
        datetime created_at
        datetime update_at
    }



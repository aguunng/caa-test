# Room 

## FlowChart
### Webhook CAA
Here is a flowchart for Webhhok CAA:
```mermaid
flowchart TD
    A[Start] --> B[Add Queue]
    B --> C[Get First Queue Today]
    C --> D{Available Agent?}
    D -->|No| Z[Stop]
    D -->|Yes| E[Assign Agent]
    E -->G[Update Queue]
    G -->Z[Stop]
```

### Webhook Mark As Resolved
Here is a flowchart for Webhook Mark As Resolved:

```mermaid
flowchart TD
    A[Start] --> B[Get First Queue Today]
    B --> C{Available Agent?}
    C -->|No| Z[Stop]
    C -->|Yes| D[Assign Agent]
    D -->E[Update Queue]
    E -->Z[Stop]
```

## ERD
Here is a ERD for Queue
```mermaid
erDiagram
    agent_allocation_queues {
        uint ID PK
        string RoomID
        string AgentID
        datetime CreatedAt
    }



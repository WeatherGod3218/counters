# counters
Days Since Last Useless House Service was made
```mermaid
erDiagram

    RESET {
        string username
        string description
        time timestamp
    }
    COUNTER {
        string id
        string createdBy
        string title
        string description
        RESET lastReset
        list[RESET] history
    }


```
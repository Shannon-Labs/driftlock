# FastAPI Integration

Integrate Driftlock anomaly detection into your FastAPI application using dependency injection and middleware.

## Installation

```bash
pip install driftlock
```

## Middleware Integration

The easiest way to monitor all requests is using the `DriftlockMiddleware`.

```python
from fastapi import FastAPI
from driftlock.fastapi import DriftlockMiddleware

app = FastAPI()

app.add_middleware(
    DriftlockMiddleware,
    api_key="your-api-key",
    stream_id="fastapi-app",
    exclude_paths=["/docs", "/openapi.json"]
)

@app.get("/")
async def root():
    return {"message": "Hello World"}
```

## Dependency Injection

For more control, you can inject the `DriftlockClient` into your path operations.

```python
from fastapi import FastAPI, Depends, HTTPException
from driftlock import DriftlockClient
from driftlock.fastapi import get_driftlock_client

app = FastAPI()

@app.post("/items/")
async def create_item(
    item: dict, 
    client: DriftlockClient = Depends(get_driftlock_client)
):
    # Detect anomalies in the request body
    result = await client.detect(
        stream_id="items-creation",
        events=[{"type": "create_item", "body": item}]
    )
    
    if result.anomalies:
        raise HTTPException(status_code=400, detail="Anomaly detected")
        
    return {"item": item, "status": "created"}
```

## Background Tasks

To avoid adding latency to your endpoints, use FastAPI's `BackgroundTasks`.

```python
from fastapi import BackgroundTasks

@app.post("/logs/")
async def log_event(
    event: dict, 
    background_tasks: BackgroundTasks,
    client: DriftlockClient = Depends(get_driftlock_client)
):
    # Schedule detection to run after the response is sent
    background_tasks.add_task(
        client.detect, 
        stream_id="logs", 
        events=[{"type": "log", "body": event}]
    )
    
    return {"status": "accepted"}
```

## Pydantic Integration

Driftlock works seamlessly with Pydantic models.

```python
from pydantic import BaseModel

class User(BaseModel):
    username: str
    email: str
    age: int

@app.post("/users/")
async def create_user(user: User, client: DriftlockClient = Depends(get_driftlock_client)):
    # Pydantic model is automatically serialized
    await client.detect(
        stream_id="user-signups",
        events=[{"type": "signup", "body": user.dict()}]
    )
    return user
```

## Configuration

You can configure the client globally or per-request.

```python
# main.py
from driftlock import DriftlockClient

# Global configuration
client = DriftlockClient(
    api_key="...",
    timeout=5.0
)

def get_client():
    return client
```

## Error Handling

Handle Driftlock errors gracefully in your endpoints.

```python
from driftlock.exceptions import RateLimitError

@app.get("/data")
async def get_data(client: DriftlockClient = Depends(get_driftlock_client)):
    try:
        await client.detect(...)
    except RateLimitError:
        # Log warning but proceed
        pass
    return {"data": "..."}
```

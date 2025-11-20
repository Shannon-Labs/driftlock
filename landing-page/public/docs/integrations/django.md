# Django Integration

Integrate Driftlock anomaly detection into your Django project using our middleware.

## Installation

Install the Python SDK:

```bash
pip install driftlock
```

## Middleware Configuration

Add the Driftlock middleware to your `MIDDLEWARE` setting in `settings.py`.

```python
# settings.py

MIDDLEWARE = [
    # ... other middleware
    'driftlock.django.DriftlockMiddleware',
]

DRIFTLOCK = {
    'API_KEY': 'your-api-key',
    'STREAM_ID': 'django-app',
    'EXCLUDE_PATHS': ['/admin/', '/health/'],
    'ASYNC_MODE': True,  # Use async processing to avoid blocking response
}
```

## Async Support

Driftlock fully supports Django's async views and middleware (Django 3.1+).

```python
# views.py
from django.http import JsonResponse
import asyncio

async def my_async_view(request):
    # Driftlock middleware automatically captures this request
    await asyncio.sleep(0.1)
    return JsonResponse({'status': 'ok'})
```

## Manual Detection in Views

You can also use the client directly within your views for more granular control.

```python
from driftlock import DriftlockClient
from django.conf import settings

client = DriftlockClient(api_key=settings.DRIFTLOCK['API_KEY'])

def process_payment(request):
    amount = request.POST.get('amount')
    
    # Check for anomalies before processing
    result = client.detect_sync(
        stream_id='payments',
        events=[{'type': 'payment', 'body': {'amount': amount}}]
    )
    
    if result.anomalies:
        return JsonResponse({'error': 'Suspicious activity detected'}, status=403)
        
    # Process payment...
    return JsonResponse({'status': 'processed'})
```

## Customizing Events

You can attach extra context to the request object, which the middleware will include in the event.

```python
def my_view(request):
    request.driftlock_context = {
        'user_tier': request.user.profile.tier,
        'ip_address': request.META.get('REMOTE_ADDR')
    }
    # ...
```

## Celery Integration

For heavy workloads, you might want to offload detection to a Celery task.

```python
# tasks.py
from celery import shared_task
from driftlock import DriftlockSyncClient

@shared_task
def detect_anomaly_task(event_data):
    client = DriftlockSyncClient(api_key='...')
    client.detect(stream_id='background-jobs', events=[event_data])

# views.py
def my_view(request):
    # ... logic ...
    detect_anomaly_task.delay({
        'timestamp': datetime.now().isoformat(),
        'type': 'job_metric',
        'body': {'duration': duration}
    })
```

## Troubleshooting

### Middleware Order
Place `DriftlockMiddleware` after `AuthenticationMiddleware` if you want to include user information in the anomaly detection context.

### Performance
Ensure `ASYNC_MODE` is set to `True` (default) in production to prevent the API call from blocking your response time.

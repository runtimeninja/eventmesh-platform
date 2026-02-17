# EventMesh Platform

Event-driven automation SaaS:
- ingest events
- evaluate rules
- trigger actions (webhook/email/slack)
- reliable delivery (retry + DLQ)
- audit + replay

## Local Dev
Gateway: http://localhost:8080

### Run
docker compose -f deploy/docker-compose.yml up --build

### Test
cd services/gateway-service && go test ./...

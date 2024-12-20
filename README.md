# Service Backend API Sample
This is a sample backend API service that can easily be used in demos &amp; tests.

## Operations
### GET /sleep
Call the /sleep endpoint to have the service sleep for 500ms before returning a positive 200 response. This is useful to a have a stable backend service call for load tests of API gateways and middleware, to see how much latency is added to the 500ms processing time. You can also set the sleep time with the **ms** query parameter.

#### Sample calls
```sh
# sleep for 500ms
curl https://service-backend-api-323709580283.europe-west1.run.app/sleep

# sleep for 1500ms
curl https://service-backend-api-323709580283.europe-west1.run.app/sleep?ms=1500
```
### POST /payload
Gets a 3kb payload similar to the output when calling an LLM like Gemini or Mistral.

```sh
# get 3kb llm chat output payload
curl -X POST https://service-backend-api-323709580283.europe-west1.run.app/payload
```

### POST /payload/:sizeInMb
Gets a payload of a specified size **sizeInMb**, with a maximum size of 20mb.

```sh
# get 10mb llm chat output payload
curl -X POST https://service-backend-api-323709580283.europe-west1.run.app/payload/10
```
## Deploy to Cloud Run
```sh
# set your project id and region
PROJECT_ID=
REGION=

# deploy
gcloud run deploy service-backend-api --source . --project $PROJECT_ID --region $REGION --allow-unauthenticated

# allow public traffic
gcloud run services add-iam-policy-binding service-backend-api \
  --member="allUsers" \
  --role="roles/run.invoker"
´´´
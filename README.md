# mystic-case-online-quest
Online quest for Mystic Case

# Attribute repository should point to the repository from what we're deploying
gcloud iam service-accounts add-iam-policy-binding <SERVICE_ACCOUNT> --project=<PROJECT_ID> --role=roles/iam.serviceAccountTokenCreator --member=<ATTRIBUTE_REPOSITORY>

# Assign corresponding role to the service account to allow upload built images to the Artifact repository
gcloud artifacts repositories add-iam-policy-binding <REPOSITORY_NAME> --location=<LOCATION> --member=<SERVICE_ACCOUNT> --role=roles/artifactregistry.writer

# Maybe this is step is not required
gcloud run services update <SERIVCE_NAME> --service-account=<SERVICE_ACCOUNT>

# Assign developer role on the Cloud Run service to the service account that is going to deploy it
gcloud run services add-iam-policy-binding <SERVICE_NAME> --member=<SERVICE_ACCOUNT> --role=roles/run.developer
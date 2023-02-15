# cloud-platform-go-get-module

cloud-platform-go-get-module is an API which provides an interface to GET and POST github repository latest release tags. The aim is to enforce the use of the latest terraform modules on our users by checking their module use when they raise a PR. This is done through a github action in the [environments repo](https://github.com/ministryofjustice/cloud-platform-environments/tree/main/cmd/check-terraform-modules-are-latest)

The API stores the latest release tag for all "cloud-platform-terraform-\*" repos. When the api starts it searches for the latest release tags for the source repos ("cloud-platform-terraform-\*") and then updates redis. The API also has functionality to receive version updates, these can be updated from a github action run from the source repos [for example](./.github/workflows/push-terraform-module-version.yaml).

## Usage

Head to the [Makefile](./Makefile) for basic commands

## Deployment strategy

Once you have merged your changes into main github actions will deploy the app to dev (which is hosted on the live-2 cluster). Once your changes have been deploy you can find them in the `cloud-platform-go-get-module` namespace and the API can be found on [dev url](https://modules.apps.live-2.cloud-platform.service.justice.gov.uk/) When you are happy that there are no issues with this deployment you can release to production (which is hosted on the live cluster) by tagging a github release using semver eg. 1.2.3. Once the action has deployed the change you will find the API in the `cloud-platform-got-get-module-prod` namespace and the deployed to the [production url](https://modules.apps.live-2.cloud-platform.service.justice.gov.uk/)

### Usage

```bash get all route
curl -i https://modules.apps.live-2.cloud-platform.service.justice.gov.uk/
```

```bash get one repo route
curl -i https://modules.apps.live-2.cloud-platform.service.justice.gov.uk/cloud-platform-terraform-$REPO_SUFFIX
```

```bash update one repo route
curl -X POST -H 'X-API-Key: $API_KEY' -i https://modules.apps.live-2.cloud-platform.service.justice.gov.uk/update/$REPO_NAME/$NEW_VERSION  
```

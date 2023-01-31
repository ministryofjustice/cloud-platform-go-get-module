# cloud-platform-go-get-module

cloud-platform-go-get-module is an API which provides an interface to GET and POST github repository latest release tags. The aim is to enforce the use of the latest terraform modules on our users.

The API stores the latest release tag for all "cloud-platform-terraform-*" repos. When the api starts it searches for the latest release tags for the source repos ("cloud-platform-terraform-*") and updates the redis store. The API also has functionality to receive version updates, these can be updated perhaps from a github action from the source repos.

# Bored API Pipeline
Application to fetch random activities from the [Bored API](https://www.boredapi.com/) and store them in files.

### Setup
Maksure all below tools are installed inorder to run this pipeline
1. [Docker](https://www.docker.com/products/docker-desktop/)
1. [Helm](https://helm.sh/docs/intro/install/)
1. [kubectl](https://kubernetes.io/docs/tasks/tools/#kubectl)
1. [minikube](https://minikube.sigs.k8s.io/docs/start/)
1. [Make](https://www.gnu.org/software/make/)

- run `bash scripts/install-dev-tools.sh` to check any of above packages missing
- run `AUTO_INSTALL='true' bash scripts/install-dev-tools.sh` to auto install above packages. Only mac brew and linux apt-get manager are supported for auto install

### Development Tools
The following tools are used for development(you dont need these if you just want to run the pipeline):
1. [Go](https://go.dev/doc/install)


### What does this app do?
- This app contains a Golang application that fetches random activities from the Bored API and stores them in files. 
- It uses Kubernetes and Helm for deployment and orchestration.
- Continously fetches data from activity endpoint for a `maxPollTime` configured in [dev.yaml](config/dev.yaml)
- All the responses are written to a file and a file is rotated(i.e closed and new file is opened) when one of the below conditions is met
  - current file size exceeds `rotate.size` configured in [dev.yaml](config/dev.yaml)
  - current file is opened for more than `rotate.interval` configured in [dev.yaml](config/dev.yaml)
- All files are stored in `outputDir` configured in [dev.yaml](config/dev.yaml)
- each file run outputs can be found under `${outputDir}/{run_id}` where run_id is cronjob run name

### How to run?
- `git clone git@github.com:sappusaketh/bored-api-pipeline.git` 
- `cd bored-api-pipeline`
- Once you have installed above tools
- Start docker desktop. Image is not published to any central docker registry we will be using our local docker registry for storing images we need this step to build and store image
- run `minikube start` to start the Kubernetes cluster on your local machine
- run `eval $(minikube -p minikube docker-env)` this helps to use our local images on k8s
- run `make docker-build` this builds our pipeline image `bored-api-pipeline:latest`
- run `make deploy` this will create kubernetes resources below 
  - Namespace - to pull all these resources under one namespace
  - CronJob - This will run the data pipeline on schedule(default 1 min) to fetch activities and store them in files
  - PersistentVolume - To persist data we copied in our datapipeline job
  - PersistentVolumeClaim - Claim volume for our job
- By default persistence volume is enabled which means all files are stored on node running in kubernetes cluster. 
- You can check the outputs by running `minikube ssh -n minikube`(if your node name is not minikube please replace it) and checking files under `/data/boredapi-pipeline-runs` folder.
- If you want to mount your local volume you can set `pvEnabled` to false in [values.yaml](helm/values.yaml)
- run `minikube mount /tmp/minikube:/tmp` to mount your local dir
- run make deploy to disable persistence volume and use local mount
- once job run is done you can find all outputs on your mounted volume eg `ls /tmp/minikube` in this case. 
- Its better to create new dir and use it for mounting because minkiube [can't do 9P Mounts with large folders](https://minikube.sigs.k8s.io/docs/handbook/mount/#9p-mounts) 
- Once testing is done makesure to delete infra either by running `make delete` or `minikube delete` to clean up resources


### Configuration
- Edit jobSchedule in [values.yaml](helm/values.yaml) to change schedule and run `make deploy`
- Any value edited in [dev.yaml](config/dev.yaml) will take effect only if you rebuild image `make docker-build`
  | Config | Description |
  |-|-|  
  | endpoint | API endpoint to fetch activities |
  | maxPollTime | Total time to keep fetching new activities |
  | outputDir | Directory to write output files |
  | rotate.interval | Interval to rotate output files |
  | rotate.size | Max size of each output file |
- If you want to create new config file in [config](config) copy current one and make changes to it and then run `make docker-build` and then update env in [values.yaml](helm/values.yaml) and run `make deploy`
  | Config | Description |
  |-|-|
  | namespace | Kubernetes namespace for resources |  
  | appName | Name given to k8s resources like cronjob, applabels etc |
  | pvcName | Persistent volume claim name |
  | jobSchedule | Cron schedule for pipeline job |
  | env | Config file name without extension |
  | pvEnabled | Whether to use persistent volume or local mount |

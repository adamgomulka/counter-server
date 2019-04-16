# Readme

This program starts a simple API server that responds to calls for a series of predefined routes. It is written in Golang, and produces a completely self-contained binary when built. Also included in this project are Terraform descriptors for a GKE cluster on which to deploy the application, a Helm chart to deploy the application on Kubernetes, and a Makefile to automate the entire process.

### Prerequisites

To build this project, you must have access to a Docker client and server (either locally or remotely), and to the GNU make utility. Deployment on GKE requires the Google Cloud SDK, kubectl, Helm, and Terraform which are downloaded using the first job defined in the Makefile. Downloading the SDK this way assumes the existence of the wget, tar, unzip, and gzip utilities on the system where the build is being run from, and if these tools are not available, the SDK can be obtained independently, and placed in the tools/ subdirectory of the project. For compatibility, this project is tested with google-cloud-sdk-241.0.0. Other versions may or may not work as expected. Additionally, the Docker image is built using Docker version 18.09.4. Other versions of Docker should work, as long as the same images are used, but again, compatibility cannot be guaranteed. On other operating systems, Terraform v0.11.13, Helm v2.13.1, and kubectl v1.14.0 can be downloaded and placed in the tools/ directory, and the project should work as expected.

### Building

To build this project on Linux, simply navigate into the projects root directory and run 'make all'. This will initialize your local machine for the Google Cloud SDK, retrieve credentials (see section on this later), build the Docker image, create a GKE cluster, and deploy the application to it. If you only want to perform part of this process, you can run one of the following:

- `make get_cloud_sdk`: Downloads the Google Cloud SDK to your local machine, and unpacks it in the projects tools/ directory.
- `make auth`: Configures Google Cloud authentication for your local machine
- `make build`: Builds the Docker image defined in the project's Dockerfile and pushes it to GCR.
- `make create`: creates the GKE cluster using Terraform.
- `make config_kubernetes`: Configures an existing Kubernetes cluster to allow deployment of Helm charts.
- `make deploy`: Deploys the Helm chart onto the Kubernetes cluster.

### Credentials

There is a JSON file included in this archive that contains the Service Account credentials used for the Google Cloud API. While this file is in the archive, it is not checked into the archive, as it would not be on a production system. In an ideal case, this should be checked into a system like Parameter Store or Hashicorp Vault, and distributed securely, but it is included here to provide required access in the absence of infrastructure for proper credential distribution.

### Using and Testing

Running `make all` will deploy the application. Once it is on GKE, the service will create a LoadBalancer. Once that is ready, simply navigate to the IP address it creates, where the application will be exposed on port 80. The following routes are exposed:

- GET http://service/hello/:name => "Hello, name!"
- GET http://service/counts => {"name":1}
- DELETE http://service/counts => {}
- GET http://service/health => {SYSINFO}

Once deployed, the application can be destroyed by running `make destroy`. Local installation artifacts are removed upon running `make clean`.

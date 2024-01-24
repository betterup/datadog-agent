# Datadog Agent

This repository is dedicated to the Datadog Agent, a software that helps in collecting metrics, traces, and logs from your environment for visualization and monitoring in Datadog's platform.

## Build and Signing for GOV

The process involves downloading the latest version of the Datadog Agent, building a container with it, and then signing the container. The signed container is then stored in GitHub Container Registry (GHCR).

### Downloading the Datadog Agent

The first step involves downloading the latest version of the Datadog Agent. This ensures that we are always using the most up-to-date and secure version of the agent.

### Building the Container

Once the Datadog Agent is downloaded, the next step is to build a container with it. This involves creating a Dockerfile that specifies the Datadog Agent as the base image and then building the container from this Dockerfile.

### Signing the Container

After the container is built, it is then signed. This involves using Cosign to sign the container image, which adds a layer of security and authenticity to the image.

### Storing in GHCR

Finally, the signed container image is pushed and stored in the GitHub Container Registry (GHCR). This allows for easy distribution and deployment of the image.
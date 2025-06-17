#####
# Create Dev stage
#####

FROM mcr.microsoft.com/devcontainers/go
ENV TZ=America/New_York

# Install golang-ebitengine dependency packages
RUN \
	apt-get update && apt-get install -y \
	sudo \
    libc6-dev \
    libgl1-mesa-dev \
    libxcursor-dev \
    libxi-dev \
    libxinerama-dev \
    libxrandr-dev \
    libxxf86vm-dev \
    libasound2-dev \
    pkg-config \
	&& apt-get autoremove -y

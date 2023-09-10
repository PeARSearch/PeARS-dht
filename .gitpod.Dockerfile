FROM gitpod/workspace-full:latest

RUN sudo apt -y update

RUN sudo sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d -b ~/.local/bin

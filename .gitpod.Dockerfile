FROM gitpod/workspace-full:latest

RUN sudo apt -y update
RUN sudo apt install -y tmux pre-commit

RUN git clone -b fruitfly-api https://github.com/PeARSearch/PeARS-orchard.git /workspace/PeARS-orchard

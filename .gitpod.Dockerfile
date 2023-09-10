FROM gitpod/workspace-full:latest

RUN sudo apt -y update

RUN sudo sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d -b ~/.local/bin

RUN sudo apt install -y protobuf-compiler tmux

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

#!/bin/bash

protoc --go_out=./ --go-grpc_out=./ users/userpb/user.proto
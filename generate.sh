#!/bin/bash

protoc greet/greetpb/greet.proto --go_out=plugins=grpc:.

protoc calculator/calculatorpb/calculator.proto --go_out=plugins=grpc:.

protoc blog-with-grpc/blog-protobuf/blog.proto --go_out=plugins=grpc:.
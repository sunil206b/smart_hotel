#!/bin/bash

go build -o smartbooking cmd/web/*.go
./smartbooking -cache=false -production=false
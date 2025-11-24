#!/bin/sh

# Write incoming code to a file
echo "$CODE" > code.cpp

# Compile the C++ program
g++ code.cpp -o code.out

# If compilation fails, exit immediately
if [ $? -ne 0 ]; then
    exit 1
fi

# Run with 1-second time limit
timeout 1 ./code.out

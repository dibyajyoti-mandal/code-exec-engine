#!/bin/sh

# Write the incoming code to a file
echo "$CODE" > code.cpp

# Compile the C++ program
g++ code.cpp -o code.out

# If compilation fails, exit immediately
if [ $? -ne 0 ]; then
    exit 1
fi

# Run the compiled program
./code.out

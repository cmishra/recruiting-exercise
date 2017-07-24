#!/bin/bash

ExecuteAndPrint ()
{
    echo "~~~~ Running command $1"
    echo "$($1)"
    echo "~~~~ "
    echo " "
}

# Test cases in README
ExecuteAndPrint "curl localhost:9000/rates?base=USD"

ExecuteAndPrint "curl localhost:9000/rates?base=USD&target=CAD"

ExecuteAndPrint "curl localhost:9000/rates?base=USD&timestamp=2016-05-01T14:34:46Z"


# Multiple targets
ExecuteAndPrint "curl localhost:9000/rates?base=USD&target=CAD&target=INR"

# Default base is USD
ExecuteAndPrint "curl localhost:9000/rates?target=CAD"

# Default target is all currencies
ExecuteAndPrint "curl localhost:9000/rates"

# Catches incorrect currency requests
ExecuteAndPrint "curl localhost:9000/rates?target=ABC"
ExecuteAndPrint "curl localhost:9000/rates?base=ABC&target=CAD"
ExecuteAndPrint "curl localhost:9000/rates?target=CAD&target=ABC"

# Catches future dates
ExecuteAndPrint "curl localhost:9000/rates?base=USD&timestamp=2018-05-01T14:34:46Z"

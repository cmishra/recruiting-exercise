# Instructions on how to run

This homework assignment is containerized and can be build and run with:

```
docker build -f Dockerfile . -t recruiting-exercise 
docker run -it -p 9000:9000 recruiting-exercise
```

I would recommend then trying out some of the examples in `example_commands.sh`. 
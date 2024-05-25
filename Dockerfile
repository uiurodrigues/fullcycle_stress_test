 # Use an official Golang runtime as a parent image
 FROM golang:latest

 # Set the working directory inside the container
 WORKDIR /app

 # Copy the local package files to the container's workspace
 COPY . /app

 # Build the Go application inside the container
 RUN go build -o fullcycle_stress_test

 # Define the command to run your application
 ENTRYPOINT ["./fullcycle_stress_test"]

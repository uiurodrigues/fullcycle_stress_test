# Stress TEst

Welcome to Stress Test Application! This application its an example of an stress test application.

## Usage

1. Run the command [docker build -t fullcycle_stress_test .] to build the docker image
2. Run the command [ docker run fullcycle_stress_test stressTest -u -r -c ] to execute the test
    - -u = url that will be tested
    - -r = count of request that will be made
    - -c = count of requests that will be made at the same time
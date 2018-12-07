# GO

# Inputs

You can either create the inputs manually (copying into /aoc/[year]/inputs/data/[day]/input.txt) or automatically by running (I've saved my configuration into a .env file for docker-compose, you can specify the variables in your own shell, at prompt, or passing literals); this will only make a request if the input exists:

    docker-compose run go sh -c 'cd /aoc/go && go run cmd/inputs.go --base-url $AOC_BASEURL $AOC_SESSIONID'

# Running (using docker)

Individual commands:

    docker-compose run go sh -c 'cd /aoc/go && go run [day]/[part]/main.go'

Or just open an interactive shell (year mount available at /aoc):

    docker-compose run go sh
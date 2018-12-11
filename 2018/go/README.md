# GO

# Inputs

Inputs are created automagically from the function common.AOCInputFile([day]) (returns *os.File), call this at the start of the main function for the aoc part.

# Running (using docker)

Individual commands:

    docker-compose run go sh -c 'cd /aoc/go && go run [day]/[part]/main.go'

Or just open an interactive shell (year mount available at /aoc):

    docker-compose run go sh
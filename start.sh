#! /usr/bin/env bash
usage() {
    echo -e "\033[1;37mUsage:\033[0m"
    echo -e "  bash $0 [options]"
    echo

    echo -e "\033[1;37mDescription:\033[0m"
    echo -e "  start server, run unit test"
    echo

    echo -e "\033[1;37mOptions:\033[0m"
    echo -e "\033[1;37m  -h, --help\033[0m       display this help message"
    echo -e "\033[1;37m  -s, --start\033[0m      start server"
    echo -e "\033[1;37m  -t, --test\033[0m       run unit tests"
    echo

    echo -e "\033[1;37mExample:\033[0m"
    echo -e "\033[1;37m  ./start.sh --start\033[0m"
    echo -e "      Start web server."
    echo
    echo -e "\033[1;37m  ./start.sh --test\033[0m"
    echo -e "      Run unit test."
    echo
}

start() {
    ./app
}

run_test() {
    go test -v ./...
}

if [ "$1" = "" ]; then
    usage
    exit 0
fi

while true; do
    case "$1" in
        -h|--help)       usage; exit 0;;
        -s|--start)      start;;
        -t|--test)       run_test;;
        --)              usage; exit 0;;
    esac
    shift

    if [ "$1" = "" ]; then
        break
    fi
done

package main

import "flag"

var idOpt = flag.Int64("id", 0, "id")
var usernameOpt = flag.String("username", "", "username")
var configPathOpt = flag.String("config", "", "path to config file")
var deleteOpt = flag.Bool("delete", false, "delete user after db operation")
var loglevelOpt = flag.String("loglevel", "info", "logger level")

package main

import "flag"

var usersServiceIDOpt = flag.Int64("userID", 1, "users service id")
var authServiceIDOpt = flag.Int64("authID", 1, "auth service id")
var usernameOpt = flag.String("username", "", "username")
var configPathOpt = flag.String("config", "", "path to config file")
var deleteOpt = flag.Bool("delete", true, "delete user after db operation")
var loglevelOpt = flag.String("loglevel", "info", "logger level")

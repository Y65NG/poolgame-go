package main

type Server struct {
	Members map[*Client]bool
}

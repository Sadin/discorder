package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
	"github.com/bwmarrin/discord.go"
)

// system variables
var (
	Token string
	buffer = make([][]byte, 0)
)
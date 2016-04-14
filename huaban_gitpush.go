package main

import (
	"log"
	"os/exec"
)

func main() {
	add := exec.Command("git", "add", ".")
	cmt := exec.Command("git", "commit", "-m", "init")
	push := exec.Command("git", "push")

	var e error
	if e = add.Run(); e == nil {
		if e = cmt.Run(); e == nil {
			if e = push.Run(); e == nil {
				log.Println("push YES")
			}
		}
	}
	if e != nil {
		log.Println(e)
	}
}

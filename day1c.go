package main

import "fmt"

type Player struct {
	Name   string
	Health int
}

func TakeDamage(player *Player, damage int) {
	player.Health -= damage

	if player.Health < 0 {
		player.Health = 0
	}
}

func (p *Player) TakeDamage(damage int) {
	p.Health -= damage

	if p.Health < 0 {
		p.Health = 0
	}
}

func (p Player) Heal(value int) Player {
	p.Health += value

	return p
}

func (p *Player) Heal2(value int) {
	p.Health += value
}

type user struct {
	Name string
	Age  *int
}

func day1c() {
	// ex 16
	player := Player{
		Name:   "Robert",
		Health: 100,
	}

	fmt.Println(player)

	TakeDamage(&player, 20)

	fmt.Println(player)

	player.TakeDamage(20)

	fmt.Println(player) // 60

	// ex 17
	var healedCopy = player.Heal(1)
	player.Heal2(20)

	fmt.Println(healedCopy) // 61 health
	fmt.Println(player)     // 80 health

	// ex 18
	newPlayer := &Player{Name: "New", Health: 100}

	newPlayer.Health += 20
	(*newPlayer).Health += 1

	fmt.Println(newPlayer)

	// ex 19

	user1 := user{Name: "John"}
	age := 40
	user2 := user{Name: "Sophia", Age: &age}

	fmt.Println(user1, user2)

	// ex 20

	slice := []Player{{}, {}}
	fmt.Println(slice)

	for i := 0; i <= len(slice); i++ {
		if i == len(slice) {
			// fmt.Println(slice[i]) - out of range
		} else {
			fmt.Println(slice[i])
		}
	}

	for _, p := range slice {
		p.Health = 10 // works on copy, not reference
	}

	fmt.Println(slice)
}

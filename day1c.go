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
}

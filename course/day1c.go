package course

import "fmt"

type player struct {
	Name   string
	Health int
}

func takeDamage(p *player, damage int) {
	p.Health -= damage

	if p.Health < 0 {
		p.Health = 0
	}
}

func (p *player) TakeDamage(damage int) {
	p.Health -= damage

	if p.Health < 0 {
		p.Health = 0
	}
}

func (p player) Heal(value int) player {
	p.Health += value

	return p
}

func (p *player) Heal2(value int) {
	p.Health += value
}

type optionalUser struct {
	Name string
	Age  *int
}

func Day1c() {
	// ex 16
	p := player{
		Name:   "Robert",
		Health: 100,
	}

	fmt.Println(p)

	takeDamage(&p, 20)

	fmt.Println(p)

	p.TakeDamage(20)

	fmt.Println(p) // 60

	// ex 17
	var healedCopy = p.Heal(1)
	p.Heal2(20)

	fmt.Println(healedCopy) // 61 health
	fmt.Println(p)          // 80 health

	// ex 18
	newPlayer := &player{Name: "New", Health: 100}

	newPlayer.Health += 20
	(*newPlayer).Health += 1

	fmt.Println(newPlayer)

	// ex 19

	user1 := optionalUser{Name: "John"}
	age := 40
	user2 := optionalUser{Name: "Sophia", Age: &age}

	fmt.Println(user1, user2)

	// ex 20

	slice := []player{{}, {}}
	fmt.Println(slice)

	for i := 0; i <= len(slice); i++ {
		if i == len(slice) {
			// fmt.Println(slice[i]) - out of range
		} else {
			fmt.Println(slice[i])
		}
	}

	for _, item := range slice {
		item.Health = 10 // works on copy, not reference
	}

	fmt.Println(slice)
}

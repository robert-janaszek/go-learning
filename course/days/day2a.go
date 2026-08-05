package days

import "fmt"

func Day2a() {
	// ex 1
	var age1 int = 10
	var age2 = 20
	age3 := 30

	fmt.Println(age1, age2, age3)

	// ex 2

	var a int = 10
	var b float64 = 20.5
	var c = float64(a) + b

	fmt.Println(c)

	// ex 3
	var a1 int
	var b1 float64
	var c1 string
	var d1 *int
	var e1 bool

	fmt.Println(a1, b1, c1, d1, e1)
	fmt.Println(c1 == "")
	fmt.Printf("%q\n", c1)

	// ex 4
	const Mon = 1
	const Tue = 2
	const Wed = 3
	const Thu = 4
	const Fri = 5
	const Sat = 6
	const Sun = 7

	const (
		Mon1 = iota + 1
		Tue1 = iota + 1
		Wed1 = iota + 1
		Thu1 = iota + 1
		Fri1 = iota + 1
		Sat1 = iota + 1
		Sun1 = iota + 1
	)

	const (
		Mon2 = iota + 1
		Tue2
		Wed2
		Thu2
		Fri2
		Sat2
		Sun2
	)

	fmt.Println(Mon, Tue, Wed, Thu, Fri, Sat, Sun)
	fmt.Println(Mon1, Tue1, Wed1, Thu1, Fri1, Sat1, Sun1)
	fmt.Println(Mon2, Tue2, Wed2, Thu2, Fri2, Sat2, Sun2)

	// ex 5
	const a2 = 2
	{
		const a2 = 3
		fmt.Println(a2)
	}
	fmt.Println(a2)

	// ex 6
	// const score = 100 - failure, can't take reference from constant
	score := 100
	ptr := &score

	// ptr+1 - failure, can't move through pointers
	fmt.Println(ptr, *ptr, *ptr+1)

	// ex 7
	fmt.Println("ex7")
	*ptr++
	fmt.Println(score) // 101

	*ptr = *ptr + 100 - 1
	fmt.Println(*ptr) // 200

	secondScore := 300
	*ptr = secondScore
	*ptr = *ptr + 1

	fmt.Println(score)       // 301
	fmt.Println(secondScore) // 300

	ptr = &secondScore
	*ptr = *ptr + 1
	fmt.Println(*ptr) // 301

	// ex 8
	fmt.Println("ex 8")
	x := 100
	pointer1 := &x
	pointer2 := &x

	*pointer2 += 2

	fmt.Println(*pointer1, *pointer2, x) // 102 102 102

	// ex 9
	var p *int
	fmt.Println(p == nil)
	fmt.Println(p) // nil
	// fmt.Println(*p) -- panic
	// *p = 1 -- panic
	a3 := 1
	p = &a3
	fmt.Println(*p) // 1

	// ex 10
	a4 := 10
	var ptr1 *int
	var pptr1 **int
	var pptr2 ***int

	ptr1 = &a4
	pptr1 = &ptr1
	pptr2 = &pptr1

	fmt.Println(pptr2)
	fmt.Println(*pptr2)
	fmt.Println(**pptr2)
	fmt.Println(***pptr2) // here we have value 10
}

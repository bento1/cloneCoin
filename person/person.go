package person

type Person struct {
	name string
	age  int
}

func (p *Person) SetDetails(name_ string, age_ int) {
	p.name = name_
	p.age = age_
}
func (p Person) SeeName() string {
	return p.name
}

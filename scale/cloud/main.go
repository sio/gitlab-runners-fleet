package cloud

import "fmt"

func Run() {
	fmt.Println("hello world")
	fleet := Fleet{
		Entrypoint: "hello",
		Hosts: []Host{
			{Name: "foo"},
			{Name: "bar"},
		},
	}
	fmt.Println(fleet)
	fleet.Save("fleet.json")
}
